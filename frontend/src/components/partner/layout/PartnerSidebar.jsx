import { NavLink, useLocation, useNavigate } from 'react-router-dom'
import { ChevronLeft, ChevronRight, LogOut, X } from 'lucide-react'
import { useEffect, useState } from 'react'
import { useAuth } from '@/context/AuthContext'
import { cn } from '@/lib/utils'
import {
  getActiveMobileLink,
  getAllowedPartnerLinks,
  getMobileDrawerLinks,
  getMobilePrimaryLinks,
  mobileMoreItem,
} from '@/components/partner/layout/partnerNav'

const SIDEBAR_KEY = 'partner_sidebar_collapsed'

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

function SidebarContent({ links, homeTo, onClose, collapsed = false, onToggle, showCollapse = false }) {
  const { logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    onClose?.()
    navigate('/partner/login')
  }

  return (
    <div className="h-full flex flex-col bg-white border-r border-cream-200 overflow-hidden">
      <div className={cn('h-20 border-b border-cream-200 flex items-center', collapsed ? 'px-4 justify-center' : 'px-5 justify-between')}>
        <NavLink to={homeTo} className={cn('flex items-center gap-3', collapsed && 'justify-center')} onClick={onClose}>
          <img src="/logo.png" alt="Бережок" className="w-10 h-10 rounded-xl object-cover shrink-0" />
          {!collapsed && (
            <div>
              <p className="font-bold text-brand-800 text-base leading-none">Бережок</p>
              <p className="text-xs text-brand-500 mt-1">Панель партнера</p>
            </div>
          )}
        </NavLink>

        {!collapsed && onClose && (
          <button className="btn-ghost p-2" onClick={onClose} aria-label="Закрыть меню">
            <X size={18} />
          </button>
        )}
      </div>

      <div className="flex-1 py-4 px-3 space-y-1">
        {links.map((link) => (
          <SidebarLink key={link.to} {...link} collapsed={collapsed} onClick={onClose} />
        ))}
      </div>

      {showCollapse && (
        <div className="p-3 border-t border-cream-200">
          <button
            onClick={onToggle}
            onMouseDown={(e) => e.preventDefault()}
            className={cn(
              'btn-ghost w-full text-sm transition-colors focus:outline-none focus:ring-0',
              collapsed ? 'justify-center' : 'justify-start gap-3'
            )}
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
  const activeLink = getActiveMobileLink(pathname, links)
  const MoreIcon = mobileMoreItem.icon

  return (
    <nav
      className="md:hidden fixed inset-x-0 bottom-0 z-40 border-t border-cream-200 bg-white/95 backdrop-blur-xl shadow-[0_-10px_30px_rgba(15,23,42,0.08)]"
      style={{ paddingBottom: 'max(env(safe-area-inset-bottom), 0.75rem)' }}
    >
      <div className={cn('grid gap-2 px-3 pt-2', links.length === 2 ? 'grid-cols-3' : 'grid-cols-4')}>
        {links.map(({ to, label, icon: Icon }) => {
          const isActive = activeLink?.to === to

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
              <span>{label === 'Выдача' ? 'Сканер' : label}</span>
            </NavLink>
          )
        })}

        <button
          type="button"
          onClick={onOpenMenu}
          className="flex min-h-16 flex-col items-center justify-center gap-1 rounded-2xl px-2 py-2 text-[11px] font-semibold text-brand-700 transition-colors hover:bg-cream-100"
        >
          <MoreIcon size={18} className="shrink-0" />
          <span>{mobileMoreItem.label}</span>
        </button>
      </div>
    </nav>
  )
}

export default function PartnerSidebar() {
  const { user } = useAuth()
  const { pathname } = useLocation()
  const [open, setOpen] = useState(false)
  const [collapsed, setCollapsed] = useState(() => localStorage.getItem(SIDEBAR_KEY) === 'true')

  useEffect(() => {
    localStorage.setItem(SIDEBAR_KEY, String(collapsed))
  }, [collapsed])

  const role = user?.role || 'owner'
  const desktopLinks = getAllowedPartnerLinks(role)
  const mobilePrimaryLinks = getMobilePrimaryLinks(role)
  const mobileDrawerLinks = getMobileDrawerLinks(role)
  const homeTo = mobilePrimaryLinks[0]?.to || desktopLinks[0]?.to || '/partner/dashboard'

  return (
    <>
      <aside
        className={cn(
          'hidden md:block shrink-0 transition-all duration-300 ease-in-out',
          collapsed ? 'w-16' : 'w-72'
        )}
      >
        <SidebarContent
          links={desktopLinks}
          homeTo={homeTo}
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
        <div className="md:hidden fixed inset-0 z-50 bg-black/40" onClick={() => setOpen(false)}>
          <div className="ml-auto h-full w-full max-w-sm" onClick={(e) => e.stopPropagation()}>
            <SidebarContent links={mobileDrawerLinks} homeTo={homeTo} onClose={() => setOpen(false)} />
          </div>
        </div>
      )}
    </>
  )
}
