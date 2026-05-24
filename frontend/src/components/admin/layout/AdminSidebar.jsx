import { NavLink, useLocation, useNavigate } from 'react-router-dom'
import { ChevronLeft, ChevronRight, LogOut, X } from 'lucide-react'
import { useEffect, useState } from 'react'
import { useStores } from '@/context/StoresContext'
import { getAdminRoleLabel } from '@/lib/admin'
import { cn } from '@/lib/utils'
import {
  adminMoreItem,
  getAllowedAdminLinks,
  getDrawerAdminLinks,
  getPrimaryAdminLinks,
} from '@/components/admin/layout/adminNav'

const SIDEBAR_KEY = 'admin_sidebar_collapsed'

function SidebarLink({ to, label, icon: Icon, collapsed, onClick }) {
  return (
    <NavLink
      to={to}
      onClick={() => onClick?.()}
      className={({ isActive }) =>
        cn(
          'flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-colors',
          collapsed && 'justify-center px-3',
          isActive ? 'bg-brand-500 text-white shadow-sm' : 'text-brand-700 hover:bg-cream-100'
        )
      }
    >
      <Icon size={17} className="shrink-0" />
      {!collapsed && <span>{label}</span>}
    </NavLink>
  )
}

function SidebarContent({ links, user, onClose, collapsed = false, onToggle, showCollapse = false }) {
  const { adminAuthStore } = useStores()
  const navigate = useNavigate()

  const handleLogout = () => {
    adminAuthStore.logout()
    onClose?.()
    navigate('/admin/login')
  }

  return (
    <div className="flex h-full flex-col overflow-hidden border-r border-cream-200 bg-white">
      <div className={cn('flex h-20 items-center border-b border-cream-200', collapsed ? 'justify-center px-4' : 'justify-between px-5')}>
        <NavLink to="/admin/applications" className={cn('flex items-center gap-3', collapsed && 'justify-center')} onClick={onClose}>
          <img src="/logo.png" alt="Бережок" className="h-10 w-10 shrink-0 rounded-xl object-cover" />
          {!collapsed && (
            <div>
              <p className="text-base font-bold leading-none text-brand-800">Бережок</p>
              <p className="mt-1 text-xs text-brand-500">Админ-панель</p>
            </div>
          )}
        </NavLink>

        {!collapsed && onClose && (
          <button className="btn-ghost p-2" onClick={onClose} aria-label="Закрыть меню">
            <X size={18} />
          </button>
        )}
      </div>

      {!collapsed && (
        <div className="border-b border-cream-200 px-5 py-4">
          <p className="truncate text-sm font-semibold text-brand-900">{user?.name || 'Администратор'}</p>
          <p className="mt-1 truncate text-xs text-brand-500">{getAdminRoleLabel(user?.role)}</p>
        </div>
      )}

      <div className="flex-1 space-y-1 overflow-y-auto px-3 py-4">
        {links.map((link) => (
          <SidebarLink key={link.to} {...link} collapsed={collapsed} onClick={onClose} />
        ))}
      </div>

      {showCollapse && (
        <div className="border-t border-cream-200 p-3">
          <button
            onClick={onToggle}
            onMouseDown={(event) => event.preventDefault()}
            className={cn('btn-ghost w-full text-sm focus:outline-none focus:ring-0', collapsed ? 'justify-center' : 'justify-start gap-3')}
          >
            {collapsed ? <ChevronRight size={16} className="shrink-0" /> : <ChevronLeft size={16} className="shrink-0" />}
            {!collapsed && <span>Свернуть</span>}
          </button>
        </div>
      )}

      <div className={cn('p-3', showCollapse ? 'pt-0' : 'border-t border-cream-200')}>
        <button onClick={handleLogout} className={cn('btn-ghost w-full text-sm', collapsed ? 'justify-center' : 'justify-start gap-3')}>
          <LogOut size={16} className="shrink-0" />
          {!collapsed && <span>Выйти</span>}
        </button>
      </div>
    </div>
  )
}

function MobileBottomNav({ links, pathname, onOpenMenu, onNavigate }) {
  const MoreIcon = adminMoreItem.icon

  return (
    <nav
      className="fixed inset-x-0 bottom-0 z-40 border-t border-cream-200 bg-white/95 shadow-[0_-10px_30px_rgba(15,23,42,0.08)] backdrop-blur-xl md:hidden"
      style={{ paddingBottom: 'max(env(safe-area-inset-bottom), 0.75rem)' }}
    >
      <div className="grid grid-cols-4 gap-2 px-3 pt-2">
        {links.map(({ to, label, icon: Icon }) => {
          const isActive = pathname === to || pathname.startsWith(`${to}/`)

          return (
            <NavLink
              key={to}
              to={to}
              onClick={onNavigate}
              className={cn(
                'flex min-h-16 flex-col items-center justify-center gap-1 rounded-2xl px-2 py-2 text-[11px] font-semibold transition-colors',
                isActive ? 'bg-brand-500 text-white shadow-sm' : 'text-brand-700 hover:bg-cream-100'
              )}
            >
              <Icon size={18} className="shrink-0" />
              <span>{label}</span>
            </NavLink>
          )
        })}

        <button
          type="button"
          onClick={onOpenMenu}
          className="flex min-h-16 flex-col items-center justify-center gap-1 rounded-2xl px-2 py-2 text-[11px] font-semibold text-brand-700 transition-colors hover:bg-cream-100"
        >
          <MoreIcon size={18} className="shrink-0" />
          <span>{adminMoreItem.label}</span>
        </button>
      </div>
    </nav>
  )
}

export default function AdminSidebar() {
  const { adminAuthStore } = useStores()
  const { pathname } = useLocation()
  const [open, setOpen] = useState(false)
  const [collapsed, setCollapsed] = useState(() => localStorage.getItem(SIDEBAR_KEY) === 'true')
  const role = adminAuthStore.user?.role || 'support'

  useEffect(() => {
    localStorage.setItem(SIDEBAR_KEY, String(collapsed))
  }, [collapsed])

  const desktopLinks = getAllowedAdminLinks(role)
  const mobilePrimaryLinks = getPrimaryAdminLinks(role)
  const mobileDrawerLinks = getDrawerAdminLinks(role)

  return (
    <>
      <aside className={cn('hidden shrink-0 transition-all duration-300 ease-in-out md:block', collapsed ? 'w-16' : 'w-72')}>
        <SidebarContent
          links={desktopLinks}
          user={adminAuthStore.user}
          collapsed={collapsed}
          onToggle={() => setCollapsed(!collapsed)}
          showCollapse
        />
      </aside>

      <MobileBottomNav
        links={mobilePrimaryLinks}
        pathname={pathname}
        onOpenMenu={() => setOpen(true)}
        onNavigate={() => setOpen(false)}
      />

      {open && (
        <div className="fixed inset-0 z-50 bg-black/40 md:hidden" onClick={() => setOpen(false)}>
          <div className="ml-auto h-full w-full max-w-sm" onClick={(event) => event.stopPropagation()}>
            <SidebarContent links={mobileDrawerLinks} user={adminAuthStore.user} onClose={() => setOpen(false)} />
          </div>
        </div>
      )}
    </>
  )
}
