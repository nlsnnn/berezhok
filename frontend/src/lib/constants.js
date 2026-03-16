export const BUSINESS_CATEGORIES = [
  { code: 'bakery',     label: 'Пекарня' },
  { code: 'cafe',       label: 'Кафе' },
  { code: 'restaurant', label: 'Ресторан' },
  { code: 'grocery',   label: 'Продуктовый магазин' },
  { code: 'hotel',     label: 'Отель' },
  { code: 'confectionery', label: 'Кондитерская' },
  { code: 'sushi',     label: 'Суши / японская кухня' },
  { code: 'pizza',     label: 'Пиццерия' },
  { code: 'other',     label: 'Другое' },
]

export const APPLICATION_STATUS = {
  pending:  { label: 'На рассмотрении', color: 'bg-yellow-100 text-yellow-800' },
  approved: { label: 'Одобрена',        color: 'bg-green-100 text-green-800' },
  rejected: { label: 'Отклонена',       color: 'bg-red-100 text-red-800' },
}

export const PARTNER_STATUS = {
  pending_documents: { label: 'Ожидает документов', color: 'bg-yellow-100 text-yellow-800' },
  active:            { label: 'Активен',             color: 'bg-green-100 text-green-800' },
  blocked:           { label: 'Заблокирован',        color: 'bg-red-100 text-red-800' },
}
