import { useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { Link } from 'react-router-dom'
import { Clock3, ClipboardList, PackageCheck, QrCode, RefreshCcw, Store, User } from 'lucide-react'
import { toast } from 'sonner'
import PartnerLayout from '@/components/partner/layout/PartnerLayout'
import Button from '@/components/ui/actions/Button'
import Spinner from '@/components/ui/feedback/Spinner'
import { useStores } from '@/context/StoresContext'
import { getOrderStatusMeta, PARTNER_ORDER_FILTERS } from '@/lib/orderStatus'
import { formatDateTime, getErrorMessage } from '@/lib/utils'

function OrdersPageBase() {
  const { ordersStore } = useStores()

  useEffect(() => {
    ordersStore.loadList().catch(() => {})
  }, [ordersStore])

  const handleFilterChange = async (status) => {
    try {
      await ordersStore.loadList({ status })
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const handleRetry = async () => {
    try {
      await ordersStore.loadList({ status: ordersStore.statusFilter })
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const handleLoadMore = async () => {
    try {
      await ordersStore.loadMore()
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const handlePickup = async (orderId) => {
    try {
      const data = await ordersStore.pickup(orderId)
      toast.success(data?.message || 'Заказ отмечен как выданный')
    } catch (error) {
      if (error?.response?.status === 409) {
        toast.error('Заказ нельзя выдать в текущем статусе')
        return
      }

      toast.error(getErrorMessage(error))
    }
  }

  return (
    <PartnerLayout
      title="Заказы"
      subtitle="Все заказы по вашим локациям с фильтрами, статусами и быстрым переходом к выдаче"
      actions={
        <Link to="/partner/orders/pickup">
          <Button className="gap-2">
            <QrCode size={16} />
            Сканер выдачи
          </Button>
        </Link>
      }
    >
      <div className="space-y-6">
        <section className="card space-y-4">
          <div className="flex flex-wrap gap-2">
            {PARTNER_ORDER_FILTERS.map((filter) => {
              const isActive = ordersStore.statusFilter === filter.value

              return (
                <button
                  key={filter.value}
                  type="button"
                  onClick={() => handleFilterChange(filter.value)}
                  disabled={ordersStore.loading && isActive}
                  className={[
                    'rounded-full border px-4 py-2 text-sm font-medium transition-colors',
                    isActive
                      ? 'border-brand-500 bg-brand-500 text-white shadow-sm'
                      : 'border-cream-300 bg-white text-brand-700 hover:border-brand-300 hover:text-brand-900',
                  ].join(' ')}
                >
                  {filter.label}
                </button>
              )
            })}
          </div>
        </section>

        {ordersStore.loading && ordersStore.items.length === 0 && (
          <div className="card flex justify-center py-20">
            <Spinner size={34} />
          </div>
        )}

        {ordersStore.error && ordersStore.items.length === 0 && (
          <div className="card flex flex-col items-center gap-4 py-14 text-center">
            <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-red-50">
              <RefreshCcw size={24} className="text-red-500" />
            </div>
            <div>
              <h3 className="text-lg font-semibold text-brand-900">Не удалось загрузить заказы</h3>
              <p className="mt-2 max-w-md text-sm text-brand-600">{getErrorMessage(ordersStore.error)}</p>
            </div>
            <Button variant="secondary" onClick={handleRetry}>
              Повторить
            </Button>
          </div>
        )}

        {!ordersStore.loading && !ordersStore.error && ordersStore.items.length === 0 && (
          <div className="card flex flex-col items-center justify-center py-16 text-center">
            <div className="mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-brand-100">
              <ClipboardList size={24} className="text-brand-600" />
            </div>
            <h3 className="mb-2 text-lg font-semibold text-brand-900">Заказов пока нет</h3>
            <p className="max-w-sm text-sm text-brand-600">
              Здесь появятся все заказы по вашим точкам. Когда клиенты начнут оформлять боксы, вы увидите их статусы и сможете быстро перейти к выдаче.
            </p>
          </div>
        )}

        {ordersStore.items.length > 0 && (
          <section className="space-y-4">
            {ordersStore.items.map((order) => (
              <OrderCard
                key={order.id}
                order={order}
                pickupLoading={ordersStore.pickupLoading && ordersStore.pickupOrderId === order.id}
                onPickup={handlePickup}
              />
            ))}

            {ordersStore.pagination.has_more && (
              <div className="flex justify-center pt-2">
                <Button variant="secondary" onClick={handleLoadMore} disabled={ordersStore.loading}>
                  {ordersStore.loading ? 'Загружаем...' : 'Показать ещё'}
                </Button>
              </div>
            )}
          </section>
        )}
      </div>
    </PartnerLayout>
  )
}

function OrderCard({ order, pickupLoading, onPickup }) {
  const statusMeta = getOrderStatusMeta(order.status)

  return (
    <article className="card space-y-5">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <p className="text-xs uppercase tracking-[0.22em] text-cream-500">Код получения</p>
          <p className="mt-1 text-xl font-semibold tracking-[0.18em] text-brand-950">{order.pickup_code}</p>
        </div>
        <span className={`badge ${statusMeta.className}`}>{statusMeta.label}</span>
      </div>

      <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-4">
        <InfoRow
          icon={PackageCheck}
          label="Бокс"
          value={order.box?.name || '—'}
          hint={order.box?.image_url ? 'Фото бокса добавлено' : 'Фото бокса отсутствует'}
        />
        <InfoRow
          icon={User}
          label="Клиент"
          value={order.customer?.name ? `${order.customer.name} (${order.customer.phone})` : order.customer?.phone || '—'}
        />
        <InfoRow icon={Store} label="Локация" value={order.location?.name || '—'} hint={order.location?.address || '—'} />
        <InfoRow
          icon={Clock3}
          label="Окно выдачи"
          value={`${formatDateTime(order.pickup_time?.start)} - ${formatDateTime(order.pickup_time?.end)}`}
          hint={`Создан ${formatDateTime(order.created_at)}`}
        />
      </div>

      <div className="flex flex-wrap gap-3">
        {order.can_pickup && (
          <Button onClick={() => onPickup(order.id)} disabled={pickupLoading} className="gap-2">
            <PackageCheck size={16} />
            {pickupLoading ? 'Выдаём заказ...' : 'Выдать заказ'}
          </Button>
        )}

        <Link to="/partner/orders/pickup">
          <Button variant="secondary" className="gap-2">
            <QrCode size={16} />
            Открыть сканер
          </Button>
        </Link>
      </div>
    </article>
  )
}

function InfoRow({ icon: Icon, label, value, hint }) {
  return (
    <div className="rounded-xl border border-cream-200 bg-cream-50 p-4">
      <p className="mb-1 text-xs uppercase tracking-wider text-cream-500">{label}</p>
      <p className="flex items-center gap-2 text-sm font-medium text-brand-800">
        <Icon size={14} className="text-brand-500" />
        {value || '—'}
      </p>
      {hint && <p className="mt-2 text-xs text-brand-500">{hint}</p>}
    </div>
  )
}

export default observer(OrdersPageBase)
