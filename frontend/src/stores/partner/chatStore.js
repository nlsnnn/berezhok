import { makeAutoObservable, runInAction } from 'mobx'
import { createOrderChatSocket, listOrderMessages, markOrderMessagesRead } from '@/api/chat'
import { listPartnerOrders } from '@/api/partner'
import { getLatestMessageId, isOrderChatClosed, mergeMessages } from '@/lib/chat'

class ChatStore {
  orders = []
  pagination = {
    total: 0,
    limit: 50,
    offset: 0,
    has_more: false,
  }
  activeOrderId = null
  messagesByOrder = new Map()
  closedOrderIds = new Set()
  loadingOrders = false
  loadingMessages = false
  sending = false
  ordersError = null
  messagesError = null
  socketStatus = 'idle'
  socketError = null
  socket = null

  constructor() {
    makeAutoObservable(this, { socket: false }, { autoBind: true })
  }

  get activeOrder() {
    return this.orders.find((order) => order.id === this.activeOrderId) || null
  }

  get messages() {
    return this.messagesByOrder.get(this.activeOrderId) || []
  }

  get isActiveChatClosed() {
    if (!this.activeOrder) return true
    return this.closedOrderIds.has(this.activeOrderId) || isOrderChatClosed(this.activeOrder.status)
  }

  get canSend() {
    return Boolean(this.activeOrderId) && !this.isActiveChatClosed && !this.sending
  }

  async loadOrders({ limit = this.pagination.limit, offset = 0, append = false } = {}) {
    this.loadingOrders = true
    this.ordersError = null

    try {
      const data = await listPartnerOrders({ limit, offset })
      const nextItems = Array.isArray(data?.items) ? data.items : []

      runInAction(() => {
        this.orders = append ? [...this.orders, ...nextItems] : nextItems
        this.pagination = {
          total: data?.pagination?.total ?? this.orders.length,
          limit: data?.pagination?.limit ?? limit,
          offset: data?.pagination?.offset ?? offset,
          has_more: Boolean(data?.pagination?.has_more),
        }

        if (!this.activeOrderId && this.orders.length > 0) {
          this.activeOrderId = this.orders[0].id
        }
      })

      if (this.activeOrderId && !append) {
        await this.selectOrder(this.activeOrderId).catch(() => {})
      }

      return data
    } catch (error) {
      runInAction(() => {
        if (!append) {
          this.orders = []
          this.activeOrderId = null
        }
        this.ordersError = error
      })
      throw error
    } finally {
      runInAction(() => {
        this.loadingOrders = false
      })
    }
  }

  async loadMoreOrders() {
    if (this.loadingOrders || !this.pagination.has_more) return null

    return this.loadOrders({
      limit: this.pagination.limit,
      offset: this.orders.length,
      append: true,
    })
  }

  async selectOrder(orderId) {
    if (!orderId) return null

    this.activeOrderId = orderId
    this.messagesError = null
    this.socketError = null
    this.closeSocket()

    try {
      await this.loadMessages(orderId)
      this.openSocket(orderId)
    } catch (error) {
      runInAction(() => {
        this.messagesError = error
      })
      throw error
    }

    return this.messagesByOrder.get(orderId) || []
  }

  async loadMessages(orderId = this.activeOrderId) {
    if (!orderId) return []

    this.loadingMessages = true
    this.messagesError = null

    try {
      const data = await listOrderMessages(orderId, { limit: 50 })
      const items = Array.isArray(data?.items) ? data.items : []

      runInAction(() => {
        this.messagesByOrder.set(orderId, mergeMessages([], items))
      })

      await this.markLatestRead(orderId)
      return items
    } catch (error) {
      runInAction(() => {
        this.messagesError = error
      })
      throw error
    } finally {
      runInAction(() => {
        this.loadingMessages = false
      })
    }
  }

  openSocket(orderId = this.activeOrderId) {
    if (!orderId || this.closedOrderIds.has(orderId) || this.isActiveChatClosed) return

    let socket
    try {
      socket = createOrderChatSocket(orderId)
    } catch (error) {
      this.socketStatus = 'closed'
      this.socketError = error
      return
    }

    this.socket = socket
    this.socketStatus = 'connecting'

    socket.onopen = () => {
      runInAction(() => {
        if (this.socket === socket) {
          this.socketStatus = 'open'
          this.socketError = null
        }
      })
    }

    socket.onmessage = (event) => {
      this.handleSocketMessage(orderId, event)
    }

    socket.onerror = () => {
      runInAction(() => {
        if (this.socket === socket) {
          this.socketError = new Error('Не удалось подключиться к чату')
        }
      })
    }

    socket.onclose = () => {
      runInAction(() => {
        if (this.socket === socket) {
          this.socketStatus = 'closed'
          this.socket = null
        }
      })
    }
  }

  handleSocketMessage(orderId, event) {
    let payload
    try {
      payload = JSON.parse(event.data)
    } catch {
      return
    }

    if (payload.type === 'message.created' && payload.message) {
      runInAction(() => {
        const existing = this.messagesByOrder.get(orderId) || []
        this.messagesByOrder.set(orderId, mergeMessages(existing, [payload.message]))
      })

      this.markLatestRead(orderId).catch(() => {})
      return
    }

    if (payload.type === 'chat.closed') {
      runInAction(() => {
        this.closedOrderIds.add(payload.order_id || orderId)
      })
      return
    }

    if (payload.type === 'error') {
      runInAction(() => {
        this.socketError = new Error(payload.message || 'Ошибка чата')
        if (payload.code === 'chat_closed') {
          this.closedOrderIds.add(orderId)
        }
      })
    }
  }

  async sendMessage(text) {
    const trimmed = text.trim()
    if (!trimmed || !this.canSend || !this.socket || this.socket.readyState !== WebSocket.OPEN) {
      return false
    }

    this.sending = true
    this.socketError = null

    try {
      this.socket.send(
        JSON.stringify({
          type: 'message.send',
          request_id: crypto.randomUUID(),
          message: trimmed,
        })
      )
      return true
    } finally {
      runInAction(() => {
        this.sending = false
      })
    }
  }

  async markLatestRead(orderId = this.activeOrderId) {
    const messages = this.messagesByOrder.get(orderId) || []
    const messageId = getLatestMessageId(messages)
    if (!orderId || !messageId) return

    await markOrderMessagesRead(orderId, messageId)
  }

  closeSocket() {
    if (this.socket) {
      this.socket.onopen = null
      this.socket.onmessage = null
      this.socket.onerror = null
      this.socket.onclose = null
      this.socket.close()
      this.socket = null
    }
    this.socketStatus = 'idle'
  }

  reset() {
    this.closeSocket()
    this.orders = []
    this.activeOrderId = null
    this.messagesByOrder.clear()
    this.closedOrderIds.clear()
    this.loadingOrders = false
    this.loadingMessages = false
    this.sending = false
    this.ordersError = null
    this.messagesError = null
    this.socketError = null
  }
}

export const chatStore = new ChatStore()
