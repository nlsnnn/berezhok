import { makeAutoObservable, runInAction } from 'mobx'
import { getOrderByPickupCode, listPartnerOrders, listPendingConfirmationOrders, pickupOrder } from '@/api/partner'
import { normalizePartnerOrderStatus } from '@/lib/orderStatus'

class OrdersStore {
  pending = []
  items = []
  pagination = {
    total: 0,
    limit: 20,
    offset: 0,
    has_more: false,
  }
  statusFilter = 'all'
  current = null
  loading = false
  lookupLoading = false
  pickupLoading = false
  pickupOrderId = null
  error = null

  constructor() {
    makeAutoObservable(this, {}, { autoBind: true })
  }

  setStatusFilter(status) {
    this.statusFilter = status || 'all'
  }

  async loadPending() {
    this.loading = true
    this.error = null
    try {
      const data = await listPendingConfirmationOrders()
      const items = Array.isArray(data?.items) ? data.items : data
      runInAction(() => {
        this.pending = items || []
      })
    } catch (error) {
      runInAction(() => {
        this.error = error
      })
    } finally {
      runInAction(() => {
        this.loading = false
      })
    }
  }

  async loadList({ status = this.statusFilter, limit = 20, offset = 0, append = false } = {}) {
    const normalizedStatus = normalizePartnerOrderStatus(status)

    this.loading = true
    this.error = null

    if (!append) {
      this.statusFilter = status || 'all'
    }

    try {
      const data = await listPartnerOrders({
        status: normalizedStatus,
        limit,
        offset,
      })

      runInAction(() => {
        const nextItems = Array.isArray(data?.items) ? data.items : []
        this.items = append ? [...this.items, ...nextItems] : nextItems
        this.pagination = {
          total: data?.pagination?.total ?? this.items.length,
          limit: data?.pagination?.limit ?? limit,
          offset: data?.pagination?.offset ?? offset,
          has_more: Boolean(data?.pagination?.has_more),
        }
        this.statusFilter = status || 'all'
      })

      return data
    } catch (error) {
      runInAction(() => {
        if (!append) {
          this.items = []
          this.pagination = {
            total: 0,
            limit,
            offset: 0,
            has_more: false,
          }
        }
        this.error = error
      })
      throw error
    } finally {
      runInAction(() => {
        this.loading = false
      })
    }
  }

  async loadMore() {
    if (this.loading || !this.pagination.has_more) return null

    return this.loadList({
      status: this.statusFilter,
      limit: this.pagination.limit,
      offset: this.items.length,
      append: true,
    })
  }

  async lookupByCode(code) {
    this.lookupLoading = true
    this.error = null
    try {
      const data = await getOrderByPickupCode(code)
      runInAction(() => {
        this.current = data
      })
      return data
    } catch (error) {
      runInAction(() => {
        this.current = null
        this.error = error
      })
      throw error
    } finally {
      runInAction(() => {
        this.lookupLoading = false
      })
    }
  }

  async pickup(orderId) {
    this.pickupLoading = true
    this.pickupOrderId = orderId
    try {
      const data = await pickupOrder(orderId)
      runInAction(() => {
        const nextStatus = data?.status ?? 'completed'
        if (this.current) {
          this.current = {
            ...this.current,
            status: nextStatus,
          }
        }

        this.items = this.items.map((item) =>
          item.id === orderId
            ? {
                ...item,
                status: nextStatus,
                can_pickup: false,
              }
            : item
        )
      })
      return data
    } finally {
      runInAction(() => {
        this.pickupLoading = false
        this.pickupOrderId = null
      })
    }
  }
}

export const ordersStore = new OrdersStore()
