import { useState } from 'react'
import { observer } from 'mobx-react-lite'
import { toast } from 'sonner'
import { BadgeRussianRuble, Mail, MapPin, Package, Phone, ShoppingBag, Store, User } from 'lucide-react'
import AdminResourceListPage from '@/pages/admin/AdminResourceListPage'
import Button from '@/components/ui/actions/Button'
import Input from '@/components/ui/form/Input'
import Select from '@/components/ui/form/Select'
import { DetailGrid, DetailItem } from '@/components/admin/AdminDetailModal'
import { useStores } from '@/context/StoresContext'
import { canMutateOperations } from '@/lib/admin'
import {
  BOX_STATUS_MAP,
  BOX_STATUS_OPTIONS,
  LOCATION_STATUS_MAP,
  LOCATION_STATUS_OPTIONS,
  ORDER_STATUS_OPTIONS,
  PARTNER_STATUS_MAP,
  PARTNER_STATUS_OPTIONS,
  PAYMENT_STATUS_MAP,
  PAYMENT_STATUS_OPTIONS,
  money,
  percent,
} from '@/lib/adminResources'
import { formatDateTime, getErrorMessage } from '@/lib/utils'

function PartnerActions({ item, store, canAct }) {
  const [status, setStatus] = useState(item.status || 'active')
  const [commissionRate, setCommissionRate] = useState(item.commission_rate ?? '')
  const [promoCommissionRate, setPromoCommissionRate] = useState(item.promo_commission_rate ?? '')
  const [promoCommissionUntil, setPromoCommissionUntil] = useState(item.promo_commission_until?.slice?.(0, 10) || '')

  if (!canAct) return null

  const submit = async () => {
    try {
      await store.update(item.id, {
        status,
        commission_rate: commissionRate === '' ? null : Number(commissionRate),
        promo_commission_rate: promoCommissionRate === '' ? null : Number(promoCommissionRate),
        promo_commission_until: promoCommissionUntil || null,
      })
      toast.success('Партнёр обновлён')
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  return (
    <div className="grid w-full gap-3 md:grid-cols-4">
      <Select value={status} onChange={(event) => setStatus(event.target.value)}>
        {PARTNER_STATUS_OPTIONS.filter((option) => option.value !== 'all').map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </Select>
      <Input type="number" step="0.01" value={commissionRate} onChange={(event) => setCommissionRate(event.target.value)} placeholder="Комиссия" />
      <Input type="number" step="0.01" value={promoCommissionRate} onChange={(event) => setPromoCommissionRate(event.target.value)} placeholder="Промо" />
      <Input type="date" value={promoCommissionUntil} onChange={(event) => setPromoCommissionUntil(event.target.value)} />
      <div className="md:col-span-4">
        <Button onClick={submit} disabled={store.actionLoading}>Сохранить изменения</Button>
      </div>
    </div>
  )
}

function StatusActions({ item, store, options, canAct, successText }) {
  const [status, setStatus] = useState(item.status || options[0]?.value)

  if (!canAct) return null

  const submit = async () => {
    try {
      await store.updateStatus(item.id, status)
      toast.success(successText)
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  return (
    <div className="flex w-full flex-wrap gap-3">
      <div className="w-full max-w-xs">
        <Select value={status} onChange={(event) => setStatus(event.target.value)}>
          {options.filter((option) => option.value !== 'all').map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </Select>
      </div>
      <Button onClick={submit} disabled={store.actionLoading}>Сменить статус</Button>
    </div>
  )
}

export const AdminPartnersPage = observer(function AdminPartnersPage() {
  const { adminPartnersStore, adminAuthStore } = useStores()
  const canAct = canMutateOperations(adminAuthStore.user)
  const bulkStatus = async (ids, status, message) => {
    try {
      await adminPartnersStore.updateStatusMany(ids, status)
      toast.success(message)
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  return (
    <AdminResourceListPage
      store={adminPartnersStore}
      title="Партнёры"
      subtitle="Статусы партнёров, комиссии и промо-условия"
      emptyText="Партнёров по текущим фильтрам нет"
      statusOptions={PARTNER_STATUS_OPTIONS}
      statusMap={PARTNER_STATUS_MAP}
      getTitle={(item) => item.business_name || item.name || item.email || item.id}
      getDescription={(item) => item.email || item.contact_email}
      getStatus={(item) => item.status}
      getMeta={(item) => [
        { icon: Mail, label: 'Email', value: item.email || item.contact_email },
        { icon: Phone, label: 'Телефон', value: item.phone || item.contact_phone },
        { icon: BadgeRussianRuble, label: 'Комиссия', value: percent(item.commission_rate) },
      ]}
      detailFields={[
        { label: 'Название', value: (item) => item.business_name || item.name },
        { label: 'Email', value: (item) => item.email || item.contact_email },
        { label: 'Телефон', value: (item) => item.phone || item.contact_phone },
        { label: 'Статус', value: (item) => PARTNER_STATUS_MAP[item.status]?.label || item.status },
        { label: 'Комиссия', value: (item) => percent(item.commission_rate) },
        { label: 'Промо-комиссия', value: (item) => percent(item.promo_commission_rate) },
        { label: 'Промо до', value: (item) => item.promo_commission_until },
        { label: 'Создан', value: (item) => formatDateTime(item.created_at) },
      ]}
      bulkActions={
        canAct
          ? [
              { label: 'Активировать', onClick: (ids) => bulkStatus(ids, 'active', 'Партнёры активированы') },
              { label: 'Приостановить', variant: 'secondary', onClick: (ids) => bulkStatus(ids, 'suspended', 'Партнёры приостановлены') },
              { label: 'Заблокировать', variant: 'danger', onClick: (ids) => bulkStatus(ids, 'blocked', 'Партнёры заблокированы') },
            ]
          : []
      }
      renderActions={({ item }) => <PartnerActions item={item} store={adminPartnersStore} canAct={canAct} />}
    />
  )
})

export const AdminLocationsPage = observer(function AdminLocationsPage() {
  const { adminLocationsStore, adminAuthStore } = useStores()
  const canAct = canMutateOperations(adminAuthStore.user)
  const bulkStatus = async (ids, status, message) => {
    try {
      await adminLocationsStore.updateStatusMany(ids, status)
      toast.success(message)
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  return (
    <AdminResourceListPage
      store={adminLocationsStore}
      title="Точки"
      subtitle="Контроль статусов и привязки точек партнёров"
      emptyText="Точек по текущим фильтрам нет"
      statusOptions={LOCATION_STATUS_OPTIONS}
      statusMap={LOCATION_STATUS_MAP}
      getTitle={(item) => item.name || item.address || item.id}
      getDescription={(item) => item.address}
      getStatus={(item) => item.status}
      getMeta={(item) => [
        { icon: Store, label: 'Партнёр', value: item.partner_name || item.partner?.name || item.partner_id },
        { icon: MapPin, label: 'Адрес', value: item.address },
      ]}
      detailFields={[
        { label: 'Название', value: (item) => item.name },
        { label: 'Адрес', value: (item) => item.address, wide: true },
        { label: 'Партнёр', value: (item) => item.partner_name || item.partner?.name || item.partner_id },
        { label: 'Статус', value: (item) => LOCATION_STATUS_MAP[item.status]?.label || item.status },
        { label: 'Создана', value: (item) => formatDateTime(item.created_at) },
      ]}
      bulkActions={
        canAct
          ? [
              { label: 'Активировать', onClick: (ids) => bulkStatus(ids, 'active', 'Точки активированы') },
              { label: 'Выключить', variant: 'secondary', onClick: (ids) => bulkStatus(ids, 'inactive', 'Точки выключены') },
              { label: 'Заблокировать', variant: 'danger', onClick: (ids) => bulkStatus(ids, 'blocked', 'Точки заблокированы') },
            ]
          : []
      }
      renderActions={({ item }) => (
        <StatusActions item={item} store={adminLocationsStore} options={LOCATION_STATUS_OPTIONS} canAct={canAct} successText="Статус точки обновлён" />
      )}
    />
  )
})

export const AdminBoxesPage = observer(function AdminBoxesPage() {
  const { adminBoxesStore, adminAuthStore } = useStores()
  const canAct = canMutateOperations(adminAuthStore.user)
  const bulkStatus = async (ids, status, message) => {
    try {
      await adminBoxesStore.updateStatusMany(ids, status)
      toast.success(message)
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  return (
    <AdminResourceListPage
      store={adminBoxesStore}
      title="Боксы"
      subtitle="Каталог боксов и смена статусов"
      emptyText="Боксов по текущим фильтрам нет"
      statusOptions={BOX_STATUS_OPTIONS}
      statusMap={BOX_STATUS_MAP}
      getTitle={(item) => item.name || item.title || item.id}
      getDescription={(item) => item.description}
      getStatus={(item) => item.status}
      getMeta={(item) => [
        { icon: Store, label: 'Точка', value: item.location_name || item.location?.name || item.location_id },
        { icon: BadgeRussianRuble, label: 'Цена', value: money(item.price) },
      ]}
      detailFields={[
        { label: 'Название', value: (item) => item.name || item.title },
        { label: 'Описание', value: (item) => item.description, wide: true },
        { label: 'Точка', value: (item) => item.location_name || item.location?.name || item.location_id },
        { label: 'Статус', value: (item) => BOX_STATUS_MAP[item.status]?.label || item.status },
        { label: 'Цена', value: (item) => money(item.price) },
        { label: 'Создан', value: (item) => formatDateTime(item.created_at) },
      ]}
      bulkActions={
        canAct
          ? [
              { label: 'Активировать', onClick: (ids) => bulkStatus(ids, 'active', 'Боксы активированы') },
              { label: 'Снять с продажи', variant: 'secondary', onClick: (ids) => bulkStatus(ids, 'inactive', 'Боксы выключены') },
              { label: 'Распроданы', variant: 'danger', onClick: (ids) => bulkStatus(ids, 'sold_out', 'Боксы отмечены распроданными') },
            ]
          : []
      }
      renderActions={({ item }) => (
        <StatusActions item={item} store={adminBoxesStore} options={BOX_STATUS_OPTIONS} canAct={canAct} successText="Статус бокса обновлён" />
      )}
    />
  )
})

export const AdminCustomersPage = observer(function AdminCustomersPage() {
  const { adminCustomersStore } = useStores()
  return (
    <AdminResourceListPage
      store={adminCustomersStore}
      title="Клиенты"
      subtitle="Просмотр клиентских профилей"
      emptyText="Клиентов по текущему поиску нет"
      getTitle={(item) => item.name || item.phone || item.email || item.id}
      getDescription={(item) => item.email}
      getMeta={(item) => [
        { icon: Phone, label: 'Телефон', value: item.phone },
        { icon: ShoppingBag, label: 'Заказы', value: item.orders_count },
      ]}
      detailFields={[
        { label: 'Имя', value: (item) => item.name },
        { label: 'Телефон', value: (item) => item.phone },
        { label: 'Email', value: (item) => item.email },
        { label: 'Заказов', value: (item) => item.orders_count },
        { label: 'Создан', value: (item) => formatDateTime(item.created_at) },
      ]}
    />
  )
})

export const AdminOrdersPage = observer(function AdminOrdersPage() {
  const { adminOrdersStore } = useStores()
  return (
    <AdminResourceListPage
      store={adminOrdersStore}
      title="Заказы"
      subtitle="Read-only просмотр заказов"
      emptyText="Заказов по текущим фильтрам нет"
      statusOptions={ORDER_STATUS_OPTIONS}
      statusMap="order"
      getTitle={(item) => item.pickup_code || item.id}
      getDescription={(item) => item.box_name || item.box?.name}
      getStatus={(item) => item.status}
      getMeta={(item) => [
        { icon: User, label: 'Клиент', value: item.customer_name || item.customer?.name || item.customer_phone },
        { icon: Store, label: 'Точка', value: item.location_name || item.location?.name },
        { icon: BadgeRussianRuble, label: 'Сумма', value: money(item.amount) },
      ]}
      detailFields={[
        { label: 'Код выдачи', value: (item) => item.pickup_code },
        { label: 'Статус', value: (item) => item.status },
        { label: 'Клиент', value: (item) => item.customer_name || item.customer?.name || item.customer_phone },
        { label: 'Бокс', value: (item) => item.box_name || item.box?.name },
        { label: 'Точка', value: (item) => item.location_name || item.location?.name },
        { label: 'Сумма', value: (item) => money(item.amount) },
        { label: 'Создан', value: (item) => formatDateTime(item.created_at) },
        { label: 'Выдача', value: (item) => `${formatDateTime(item.pickup_time_start || item.pickup_time?.start)} - ${formatDateTime(item.pickup_time_end || item.pickup_time?.end)}`, wide: true },
      ]}
    />
  )
})

export const AdminPaymentsPage = observer(function AdminPaymentsPage() {
  const { adminPaymentsStore } = useStores()
  return (
    <AdminResourceListPage
      store={adminPaymentsStore}
      title="Платежи"
      subtitle="Read-only просмотр платежей и событий"
      emptyText="Платежей по текущим фильтрам нет"
      statusOptions={PAYMENT_STATUS_OPTIONS}
      statusMap={PAYMENT_STATUS_MAP}
      getTitle={(item) => item.provider_payment_id || item.id}
      getDescription={(item) => item.order_id}
      getStatus={(item) => item.status}
      getMeta={(item) => [
        { icon: ShoppingBag, label: 'Заказ', value: item.order_id },
        { icon: BadgeRussianRuble, label: 'Сумма', value: money(item.amount) },
      ]}
      detailFields={[
        { label: 'ID платежа', value: (item) => item.id, wide: true },
        { label: 'Провайдер', value: (item) => item.provider_payment_id, wide: true },
        { label: 'Заказ', value: (item) => item.order_id },
        { label: 'Статус', value: (item) => PAYMENT_STATUS_MAP[item.status]?.label || item.status },
        { label: 'Сумма', value: (item) => money(item.amount) },
        { label: 'Создан', value: (item) => formatDateTime(item.created_at) },
      ]}
      renderDetailExtra={() => (
        <section>
          <h3 className="mb-3 text-base font-semibold text-brand-900">События платежа</h3>
          <div className="space-y-2">
            {adminPaymentsStore.events.map((event) => (
              <div key={event.id || event.created_at} className="rounded-xl border border-cream-200 bg-cream-50 p-3 text-sm text-brand-700">
                <p className="font-medium text-brand-900">{event.event_type || event.type || 'Событие'}</p>
                <p className="mt-1 text-xs text-brand-500">{formatDateTime(event.created_at)}</p>
              </div>
            ))}
            {adminPaymentsStore.events.length === 0 && <p className="text-sm text-brand-600">Событий пока нет</p>}
          </div>
        </section>
      )}
    />
  )
})
