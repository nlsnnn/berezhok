import {
  BarChart3,
  Banknote,
  ClipboardList,
  KeyRound,
  LayoutDashboard,
  MapPin,
  Menu,
  Package,
  QrCode,
  UserCircle,
  Users,
  MessageSquare,
  Gift,
} from 'lucide-react'

export const partnerNavLinks = [
  {
    to: '/partner/dashboard',
    label: 'Дашборд',
    icon: LayoutDashboard,
    roles: ['owner'],
    mobilePrimary: true,
    matchPaths: ['/partner/dashboard'],
  },

  {
    to: '/partner/locations',
    label: 'Локации',
    icon: MapPin,
    roles: ['owner'],
    matchPaths: ['/partner/locations'],
  },
  {
    to: '/partner/boxes',
    label: 'Боксы',
    icon: Package,
    roles: ['owner'],
    matchPaths: ['/partner/boxes'],
  },
  {
    to: '/partner/orders',
    label: 'Заказы',
    icon: ClipboardList,
    roles: ['owner', 'employee'],
    mobilePrimary: true,
    matchPaths: ['/partner/orders'],
  },
  {
    to: '/partner/orders/pickup',
    label: 'Выдача',
    icon: QrCode,
    roles: ['owner', 'employee'],
    mobilePrimary: true,
    matchPaths: ['/partner/orders/pickup'],
  },
  {
    to: '/partner/chat',
    label: 'Чаты',
    icon: MessageSquare,
    roles: ['owner', 'employee'],
  },
  {
    to: '/partner/employees',
    label: 'Сотрудники',
    icon: Users,
    roles: ['owner'],
    matchPaths: ['/partner/employees'],
  },
  {
    to: '/partner/stats',
    label: 'Статистика',
    icon: BarChart3,
    roles: ['owner'],
    matchPaths: ['/partner/stats'],
  },
  {
    to: '/partner/payouts',
    label: 'Выплаты',
    icon: Banknote,
    roles: ['owner'],
    matchPaths: ['/partner/payouts'],
  },
  {
    to: '/partner/profile',
    label: 'Профиль',
    icon: UserCircle,
    roles: ['owner'],
    matchPaths: ['/partner/profile', '/partner/legal-info'],
  },
  {
    to: '/partner/change-password',
    label: 'Пароль',
    icon: KeyRound,
    roles: ['owner', 'employee'],
    matchPaths: ['/partner/change-password'],
  },
]

export const mobileMoreItem = {
  key: 'more',
  label: 'Еще',
  icon: Menu,
}

export function getAllowedPartnerLinks(role) {
  const resolvedRole = role || 'owner'
  return partnerNavLinks.filter((link) => link.roles.includes(resolvedRole))
}

export function getMobilePrimaryLinks(role) {
  return getAllowedPartnerLinks(role).filter((link) => link.mobilePrimary)
}

export function getMobileDrawerLinks(role) {
  return getAllowedPartnerLinks(role).filter((link) => !link.mobilePrimary)
}

export function getActiveMobileLink(pathname, links) {
  return (
    links
      .filter((link) => link.matchPaths?.some((path) => pathname.startsWith(path)))
      .sort((a, b) => {
        const aLen = Math.max(...a.matchPaths.filter((path) => pathname.startsWith(path)).map((path) => path.length))
        const bLen = Math.max(...b.matchPaths.filter((path) => pathname.startsWith(path)).map((path) => path.length))
        return bLen - aLen
      })[0] || null
  )
}
