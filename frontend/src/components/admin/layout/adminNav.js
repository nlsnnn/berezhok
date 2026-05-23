import {
  BadgeRussianRuble,
  BarChart3,
  Boxes,
  ClipboardList,
  FileClock,
  FileText,
  MapPin,
  ShieldCheck,
  ShoppingBag,
  Store,
  Users,
} from 'lucide-react'
import { canManageAdmins } from '../../../lib/admin.js'

export const adminLinks = [
  { to: '/admin/applications', label: 'Заявки', icon: FileText },
  { to: '/admin/partners', label: 'Партнёры', icon: Store },
  { to: '/admin/locations', label: 'Точки', icon: MapPin },
  { to: '/admin/boxes', label: 'Боксы', icon: Boxes },
  { to: '/admin/customers', label: 'Клиенты', icon: Users },
  { to: '/admin/orders', label: 'Заказы', icon: ShoppingBag },
  { to: '/admin/payments', label: 'Платежи', icon: BadgeRussianRuble },
  { to: '/admin/stats', label: 'Статистика', icon: BarChart3 },
  { to: '/admin/audit', label: 'Аудит', icon: FileClock },
  { to: '/admin/admins', label: 'Админы', icon: ShieldCheck, superAdminOnly: true },
]

export const adminMoreItem = { label: 'Ещё', icon: ClipboardList }

export function getAllowedAdminLinks(role) {
  const user = { role }
  return adminLinks.filter((link) => !link.superAdminOnly || canManageAdmins(user))
}

export function getPrimaryAdminLinks(role) {
  return getAllowedAdminLinks(role).filter((link) =>
    ['/admin/applications', '/admin/partners', '/admin/orders'].includes(link.to)
  )
}

export function getDrawerAdminLinks(role) {
  const primary = new Set(getPrimaryAdminLinks(role).map((link) => link.to))
  return getAllowedAdminLinks(role).filter((link) => !primary.has(link.to))
}
