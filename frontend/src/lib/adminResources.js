import { BOX_STATUS, PARTNER_STATUS } from '@/lib/constants'
import { getOrderStatusMeta } from '@/lib/orderStatus'

export const ALL_STATUS_OPTION = { value: 'all', label: 'Все' }

export const PARTNER_STATUS_OPTIONS = [
  ALL_STATUS_OPTION,
  { value: 'pending_documents', label: 'Ожидает документов' },
  { value: 'active', label: 'Активен' },
  { value: 'suspended', label: 'Приостановлен' },
  { value: 'blocked', label: 'Заблокирован' },
]

export const LEGAL_VERIFICATION_STATUS_OPTIONS = [
  ALL_STATUS_OPTION,
  { value: 'pending', label: 'На проверке' },
  { value: 'verified', label: 'Подтверждены' },
  { value: 'failed', label: 'Отклонены' },
]

export const LOCATION_STATUS_OPTIONS = [
  ALL_STATUS_OPTION,
  { value: 'pending_review', label: 'На модерации' },
  { value: 'active', label: 'Активна' },
  { value: 'inactive', label: 'Неактивна' },
  { value: 'blocked', label: 'Заблокирована' },
]

export const BOX_STATUS_OPTIONS = [
  ALL_STATUS_OPTION,
  { value: 'active', label: 'Активен' },
  { value: 'inactive', label: 'Неактивен' },
  { value: 'draft', label: 'Черновик' },
  { value: 'sold_out', label: 'Распродан' },
]

export const ORDER_STATUS_OPTIONS = [
  ALL_STATUS_OPTION,
  { value: 'created', label: 'Создан' },
  { value: 'paid', label: 'Оплачен' },
  { value: 'confirmed', label: 'Подтверждён' },
  { value: 'picked_up', label: 'Выдан' },
  { value: 'completed', label: 'Завершён' },
  { value: 'cancelled', label: 'Отменён' },
]

export const PAYMENT_STATUS_OPTIONS = [
  ALL_STATUS_OPTION,
  { value: 'pending', label: 'Ожидает' },
  { value: 'succeeded', label: 'Успешен' },
  { value: 'canceled', label: 'Отменён' },
  { value: 'failed', label: 'Ошибка' },
]

export const ADMIN_ROLE_OPTIONS = [
  { value: 'admin', label: 'Администратор' },
  { value: 'support', label: 'Поддержка' },
]

export const ADMIN_EDIT_ROLE_OPTIONS = [
  { value: 'super_admin', label: 'Суперадминистратор' },
  ...ADMIN_ROLE_OPTIONS,
]

export const PARTNER_STATUS_MAP = {
  ...PARTNER_STATUS,
  suspended: { label: 'Приостановлен', color: 'bg-orange-100 text-orange-800' },
}

export const LEGAL_VERIFICATION_STATUS_MAP = {
  pending: { label: 'Юр. данные на проверке', color: 'bg-yellow-100 text-yellow-800' },
  verified: { label: 'Юр. данные подтверждены', color: 'bg-green-100 text-green-800' },
  failed: { label: 'Юр. данные отклонены', color: 'bg-red-100 text-red-800' },
}

export const LOCATION_STATUS_MAP = {
  pending_review: { label: 'На модерации', color: 'bg-yellow-100 text-yellow-800' },
  active: { label: 'Активна', color: 'bg-green-100 text-green-800' },
  inactive: { label: 'Неактивна', color: 'bg-gray-100 text-gray-800' },
  blocked: { label: 'Заблокирована', color: 'bg-red-100 text-red-800' },
}

export const BOX_STATUS_MAP = BOX_STATUS

export const PAYMENT_STATUS_MAP = {
  pending: { label: 'Ожидает', color: 'bg-yellow-100 text-yellow-800' },
  succeeded: { label: 'Успешен', color: 'bg-green-100 text-green-800' },
  canceled: { label: 'Отменён', color: 'bg-gray-100 text-gray-800' },
  failed: { label: 'Ошибка', color: 'bg-red-100 text-red-800' },
}

export const ADMIN_ACTIVE_STATUS_MAP = {
  active: { label: 'Активен', color: 'bg-green-100 text-green-800' },
  inactive: { label: 'Отключён', color: 'bg-gray-100 text-gray-800' },
}

export function getStatusMeta(status, map = {}) {
  if (map === 'order') return getOrderStatusMeta(status)
  return map[status] || { label: status || '—', color: 'bg-gray-100 text-gray-700' }
}

export function money(value) {
  if (value === null || value === undefined || value === '') return '—'
  return `${Number(value).toLocaleString('ru-RU')} ₽`
}

export function percent(value) {
  if (value === null || value === undefined || value === '') return '—'
  return `${Math.round(Number(value) * 100)}%`
}
