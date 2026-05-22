import { useEffect, useRef, useState } from 'react'
import { observer } from 'mobx-react-lite'
import { AlertCircle, CheckCircle2, Clock3, MessageSquare, RefreshCcw, Send, Store, User } from 'lucide-react'
import { toast } from 'sonner'
import PartnerLayout from '@/components/partner/layout/PartnerLayout'
import Button from '@/components/ui/actions/Button'
import Spinner from '@/components/ui/feedback/Spinner'
import Input from '@/components/ui/form/Input'
import { useStores } from '@/context/StoresContext'
import { getOrderStatusMeta } from '@/lib/orderStatus'
import { cn, formatDateTime, getErrorMessage } from '@/lib/utils'

function PartnerChatPageBase() {
  const { chatStore } = useStores()
  const [draft, setDraft] = useState('')
  const messagesEndRef = useRef(null)

  useEffect(() => {
    chatStore.loadOrders().catch((error) => {
      toast.error(getErrorMessage(error))
    })

    return () => {
      chatStore.closeSocket()
    }
  }, [chatStore])

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ block: 'end' })
  }, [chatStore.messages.length, chatStore.activeOrderId])

  const handleSelectOrder = async (orderId) => {
    try {
      await chatStore.selectOrder(orderId)
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const handleRetryOrders = async () => {
    try {
      await chatStore.loadOrders()
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const handleRetryMessages = async () => {
    try {
      await chatStore.selectOrder(chatStore.activeOrderId)
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const handleLoadMore = async () => {
    try {
      await chatStore.loadMoreOrders()
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const handleSubmit = async (event) => {
    event.preventDefault()

    const sent = await chatStore.sendMessage(draft)
    if (sent) {
      setDraft('')
      return
    }

    if (chatStore.isActiveChatClosed) {
      toast.error('Чат закрыт для этого заказа')
      return
    }

    if (chatStore.socketStatus !== 'open') {
      toast.error('Чат ещё подключается')
    }
  }

  const inputDisabled = chatStore.isActiveChatClosed || chatStore.sending || chatStore.socketStatus !== 'open'

  return (
    <PartnerLayout title="Чаты" subtitle="Сообщения с клиентами по заказам">
      <div className="overflow-hidden rounded-2xl border border-cream-200 bg-white shadow-sm">
        <div className="grid min-h-[calc(100svh-230px)] lg:grid-cols-[360px_minmax(0,1fr)]">
          <aside className="border-b border-cream-200 bg-cream-50 lg:border-b-0 lg:border-r">
            <div className="flex items-center justify-between gap-3 border-b border-cream-200 bg-white px-4 py-4">
              <div>
                <h2 className="text-base font-semibold text-brand-900">Диалоги</h2>
                <p className="text-xs text-brand-500">По заказам партнёра</p>
              </div>
              {chatStore.loadingOrders && <Spinner size={20} />}
            </div>

            <OrderList
              orders={chatStore.orders}
              activeOrderId={chatStore.activeOrderId}
              loading={chatStore.loadingOrders}
              error={chatStore.ordersError}
              hasMore={chatStore.pagination.has_more}
              onSelect={handleSelectOrder}
              onRetry={handleRetryOrders}
              onLoadMore={handleLoadMore}
            />
          </aside>

          <section className="flex min-h-[620px] flex-col bg-gray-50">
            {chatStore.activeOrder ? (
              <>
                <ChatHeader
                  order={chatStore.activeOrder}
                  socketStatus={chatStore.socketStatus}
                  socketError={chatStore.socketError}
                  isClosed={chatStore.isActiveChatClosed}
                />

                <div className="flex-1 overflow-y-auto px-4 py-5">
                  {chatStore.loadingMessages && chatStore.messages.length === 0 && (
                    <div className="flex h-full items-center justify-center">
                      <Spinner size={34} />
                    </div>
                  )}

                  {chatStore.messagesError && chatStore.messages.length === 0 && (
                    <EmptyState
                      icon={RefreshCcw}
                      title="Не удалось загрузить сообщения"
                      text={getErrorMessage(chatStore.messagesError)}
                      action={
                        <Button variant="secondary" onClick={handleRetryMessages}>
                          Повторить
                        </Button>
                      }
                    />
                  )}

                  {!chatStore.loadingMessages && !chatStore.messagesError && chatStore.messages.length === 0 && (
                    <EmptyState
                      icon={MessageSquare}
                      title="Сообщений пока нет"
                      text="Начните диалог с клиентом по этому заказу."
                    />
                  )}

                  {chatStore.messages.length > 0 && (
                    <div className="space-y-3">
                      {chatStore.messages.map((message) => (
                        <MessageBubble key={message.id} message={message} />
                      ))}
                      <div ref={messagesEndRef} />
                    </div>
                  )}
                </div>

                {chatStore.socketError && (
                  <div className="border-t border-amber-200 bg-amber-50 px-4 py-2 text-sm text-amber-700">
                    {chatStore.socketError.message}
                  </div>
                )}

                {chatStore.isActiveChatClosed && (
                  <div className="border-t border-cream-200 bg-white px-4 py-3 text-sm text-brand-600">
                    Чат закрыт: заказ больше не находится в активном статусе.
                  </div>
                )}

                <form onSubmit={handleSubmit} className="border-t border-cream-200 bg-white p-4">
                  <div className="flex items-center gap-2">
                    <Input
                      value={draft}
                      onChange={(event) => setDraft(event.target.value)}
                      placeholder={chatStore.isActiveChatClosed ? 'Чат закрыт' : 'Введите сообщение...'}
                      disabled={inputDisabled}
                    />
                    <Button type="submit" disabled={!draft.trim() || inputDisabled} className="shrink-0 gap-2">
                      <span className="hidden sm:inline">{chatStore.sending ? 'Отправляем...' : 'Отправить'}</span>
                      <Send size={18} />
                    </Button>
                  </div>
                </form>
              </>
            ) : (
              <EmptyState
                icon={MessageSquare}
                title="Выберите диалог"
                text="Здесь появится переписка по выбранному заказу."
              />
            )}
          </section>
        </div>
      </div>
    </PartnerLayout>
  )
}

function OrderList({ orders, activeOrderId, loading, error, hasMore, onSelect, onRetry, onLoadMore }) {
  if (loading && orders.length === 0) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Spinner size={30} />
      </div>
    )
  }

  if (error && orders.length === 0) {
    return (
      <div className="p-4">
        <EmptyState
          icon={RefreshCcw}
          title="Заказы не загрузились"
          text={getErrorMessage(error)}
          action={
            <Button variant="secondary" onClick={onRetry}>
              Повторить
            </Button>
          }
        />
      </div>
    )
  }

  if (!loading && orders.length === 0) {
    return (
      <div className="p-4">
        <EmptyState icon={MessageSquare} title="Диалогов пока нет" text="Когда появятся заказы, здесь будут чаты с клиентами." />
      </div>
    )
  }

  return (
    <div className="max-h-[420px] overflow-y-auto lg:max-h-[calc(100svh-310px)]">
      {orders.map((order) => (
        <OrderChatItem key={order.id} order={order} active={activeOrderId === order.id} onSelect={onSelect} />
      ))}

      {hasMore && (
        <div className="p-4">
          <Button variant="secondary" className="w-full" onClick={onLoadMore} disabled={loading}>
            {loading ? 'Загружаем...' : 'Показать ещё'}
          </Button>
        </div>
      )}
    </div>
  )
}

function OrderChatItem({ order, active, onSelect }) {
  const statusMeta = getOrderStatusMeta(order.status)
  const closed = !['confirmed', 'picked_up'].includes(order.status)

  return (
    <button
      type="button"
      onClick={() => onSelect(order.id)}
      className={cn(
        'flex w-full gap-3 border-b border-cream-100 px-4 py-4 text-left transition-colors hover:bg-white',
        active && 'bg-white shadow-inner'
      )}
    >
      <div className={cn('mt-1 flex h-10 w-10 shrink-0 items-center justify-center rounded-xl', closed ? 'bg-cream-200' : 'bg-brand-500')}>
        <User size={20} className={closed ? 'text-brand-500' : 'text-white'} />
      </div>
      <div className="min-w-0 flex-1">
        <div className="flex items-start justify-between gap-2">
          <p className="truncate text-sm font-semibold text-brand-900">{getCustomerLabel(order)}</p>
          <span className={cn('badge shrink-0 text-[11px]', statusMeta.className)}>{statusMeta.label}</span>
        </div>
        <p className="mt-1 truncate text-sm text-brand-600">Заказ {order.pickup_code || order.id}</p>
        <div className="mt-2 flex flex-wrap gap-x-3 gap-y-1 text-xs text-brand-500">
          <span className="inline-flex items-center gap-1">
            <Store size={12} />
            {order.location?.name || 'Локация'}
          </span>
          <span className="inline-flex items-center gap-1">
            <Clock3 size={12} />
            {formatPickupWindow(order)}
          </span>
        </div>
      </div>
    </button>
  )
}

function ChatHeader({ order, socketStatus, socketError, isClosed }) {
  const statusMeta = getOrderStatusMeta(order.status)

  return (
    <div className="flex flex-wrap items-center justify-between gap-3 border-b border-cream-200 bg-white px-4 py-4">
      <div className="min-w-0">
        <div className="flex flex-wrap items-center gap-2">
          <h3 className="truncate text-lg font-semibold text-brand-900">{getCustomerLabel(order)}</h3>
          <span className={cn('badge', statusMeta.className)}>{statusMeta.label}</span>
        </div>
        <p className="mt-1 text-sm text-brand-500">
          {order.location?.name || 'Локация'} · код {order.pickup_code || '—'}
        </p>
      </div>
      <ConnectionStatus socketStatus={socketStatus} socketError={socketError} isClosed={isClosed} />
    </div>
  )
}

function ConnectionStatus({ socketStatus, socketError, isClosed }) {
  if (isClosed) {
    return <span className="inline-flex items-center gap-2 text-sm text-brand-500">Чат закрыт</span>
  }

  if (socketError) {
    return (
      <span className="inline-flex items-center gap-2 text-sm text-amber-700">
        <AlertCircle size={16} />
        Нет соединения
      </span>
    )
  }

  if (socketStatus === 'open') {
    return (
      <span className="inline-flex items-center gap-2 text-sm text-emerald-700">
        <CheckCircle2 size={16} />
        Онлайн
      </span>
    )
  }

  return <span className="text-sm text-brand-500">Подключаемся...</span>
}

function MessageBubble({ message }) {
  const isMine = message.sender_type === 'partner'

  return (
    <div className={cn('flex', isMine ? 'justify-end' : 'justify-start')}>
      <div
        className={cn(
          'max-w-[85%] rounded-2xl px-4 py-2 text-sm shadow-sm md:max-w-[70%]',
          isMine ? 'rounded-br-sm bg-brand-500 text-white' : 'rounded-bl-sm border border-gray-200 bg-white text-gray-800'
        )}
      >
        <p className="whitespace-pre-wrap break-words">{message.message}</p>
        <p className={cn('mt-1 text-right text-xs opacity-75', isMine ? 'text-brand-100' : 'text-gray-400')}>
          {formatMessageTime(message.created_at)}
        </p>
      </div>
    </div>
  )
}

function EmptyState({ icon: Icon, title, text, action }) {
  return (
    <div className="flex h-full min-h-64 flex-col items-center justify-center px-6 py-10 text-center">
      <div className="mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-brand-100">
        <Icon size={24} className="text-brand-600" />
      </div>
      <h3 className="text-lg font-semibold text-brand-900">{title}</h3>
      <p className="mt-2 max-w-sm text-sm text-brand-600">{text}</p>
      {action && <div className="mt-4">{action}</div>}
    </div>
  )
}

function getCustomerLabel(order) {
  if (order.customer?.name) return `${order.customer.name} (${order.customer.phone || 'телефон не указан'})`
  return order.customer?.phone || 'Клиент'
}

function formatPickupWindow(order) {
  if (!order.pickup_time?.start || !order.pickup_time?.end) return 'Время не указано'
  return `${formatDateTime(order.pickup_time.start)} - ${formatDateTime(order.pickup_time.end)}`
}

function formatMessageTime(date) {
  if (!date) return ''
  return new Date(date).toLocaleTimeString('ru-RU', {
    hour: '2-digit',
    minute: '2-digit',
  })
}

export default observer(PartnerChatPageBase)
