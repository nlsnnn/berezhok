export const ORDER_STATUS_META = {
  pending: { label: 'Ожидает оплаты', className: 'bg-yellow-100 text-yellow-800' },
  paid: { label: 'Оплачен', className: 'bg-emerald-100 text-emerald-800' },
  confirmed: { label: 'Подтверждён', className: 'bg-green-100 text-green-800' },
  completed: { label: 'Выдан', className: 'bg-blue-100 text-blue-800' },
  picked_up: { label: 'Получен клиентом', className: 'bg-brand-100 text-brand-800' },
  cancelled: { label: 'Отменён', className: 'bg-red-100 text-red-800' },
  disputed: { label: 'Спор', className: 'bg-orange-100 text-orange-800' },
  refunded: { label: 'Возврат', className: 'bg-slate-200 text-slate-800' },
}

export const PARTNER_ORDER_FILTERS = [
  { value: 'all', label: 'Все' },
  { value: 'pending', label: 'Ожидают оплаты' },
  { value: 'paid', label: 'Оплачены' },
  { value: 'confirmed', label: 'Подтверждены' },
  { value: 'completed', label: 'Выданы' },
  { value: 'picked_up', label: 'Получены клиентом' },
  { value: 'cancelled', label: 'Отменены' },
  { value: 'disputed', label: 'Спорные' },
  { value: 'refunded', label: 'Возвраты' },
]

export function getOrderStatusMeta(status) {
  return ORDER_STATUS_META[status] || {
    label: status || '—',
    className: 'bg-gray-100 text-gray-800',
  }
}

export function normalizePartnerOrderStatus(status) {
  if (!status || status === 'all') return ''
  return ORDER_STATUS_META[status] ? status : ''
}
