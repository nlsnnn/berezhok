import { NavLink, useNavigate } from 'react-router-dom'
import { ChevronLeft, ChevronRight, LayoutDashboard, LogOut, MapPin, Menu, Package, QrCode, X } from 'lucide-react'
import { useEffect, useState } from 'react'
import { useAuth } from '@/context/AuthContext'
import { cn } from '@/lib/utils'

const SIDEBAR_KEY = 'partner_sidebar_collapsed'

const links = [
  { to: '/partner/dashboard', label: 'Дашборд', icon: LayoutDashboard },
  { to: '/partner/locations', label: 'Локации', icon: MapPin },
  { to: '/partner/boxes', label: 'Боксы', icon: Package },
  { to: '/partner/orders/pickup', label: 'Выдача', icon: QrCode },
]

function SidebarContent({ onClose, collapsed, onToggle }) {
  const { logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/partner/login')
  }

  return (
    <div className="h-full flex flex-col bg-white border-r border-cream-200 overflow-hidden">
      <div className={cn('h-20 border-b border-cream-200 flex items-center', collapsed ? 'px-4 justify-center' : 'px-5 justify-between')}>
        <NavLink to="/partner/dashboard" className={cn('flex items-center gap-3', collapsed && 'justify-center')} onClick={onClose}>
          <img src="/logo.png" alt="Бережок" className="w-10 h-10 rounded-xl object-cover shrink-0" />
          {!collapsed && (
            <div>
              <p className="font-bold text-brand-800 text-base leading-none">Бережок</p>
              <p className="text-xs text-brand-500 mt-1">Панель партнера</p>
            </div>
          )}
        </NavLink>
        {!collapsed && (
          <button className="md:hidden btn-ghost p-2" onClick={onClose}>
            <X size={18} />
          </button>
        )}
      </div>

      <div className="flex-1 py-4 px-3 space-y-1">
        {links.map(({ to, label, icon: Icon }) => (
          <NavLink
            key={to}
            to={to}
            onClick={onClose}
            className={({ isActive }) =>
              cn(
                'flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-colors',
                collapsed && 'justify-center px-3',
                isActive
                  ? 'bg-brand-500 text-white shadow-sm'
                  : 'text-brand-700 hover:bg-cream-100'
              )
            }
          >
            <Icon size={17} className="shrink-0" />
            {!collapsed && <span>{label}</span>}
          </NavLink>
        ))}
      </div>

      <div className="p-3 border-t border-cream-200">
        <button
          onClick={onToggle}
          onMouseDown={(e) => e.preventDefault()}
          className={cn('btn-ghost w-full text-sm transition-colors focus:outline-none focus:ring-0', collapsed ? 'justify-center' : 'justify-start gap-3')}
        >
          {collapsed ? <ChevronRight size={16} className="shrink-0" /> : <ChevronLeft size={16} className="shrink-0" />}
          {!collapsed && <span>Свернуть</span>}
        </button>
      </div>

      <div className="p-3 pt-0">
        <button onClick={handleLogout} className={cn('btn-ghost w-full text-sm', collapsed ? 'justify-center' : 'justify-start gap-3')}>
          <LogOut size={16} className="shrink-0" />
          {!collapsed && <span>Выйти</span>}
        </button>
      </div>
    </div>
  )
}

export default function PartnerSidebar() {
  const [open, setOpen] = useState(false)
  const [collapsed, setCollapsed] = useState(() => localStorage.getItem(SIDEBAR_KEY) === 'true')

  useEffect(() => {
    localStorage.setItem(SIDEBAR_KEY, String(collapsed))
  }, [collapsed])

  return (
    <>
      <aside
        className={cn(
          'hidden md:block shrink-0 transition-all duration-300 ease-in-out',
          collapsed ? 'w-16' : 'w-72'
        )}
      >
        <SidebarContent collapsed={collapsed} onToggle={() => setCollapsed(!collapsed)} />
      </aside>

      <button className="md:hidden fixed top-4 left-4 z-40 btn-secondary p-2" onClick={() => setOpen(true)}>
        <Menu size={18} />
      </button>

      {open && (
        <div className="md:hidden fixed inset-0 z-50 bg-black/40">
          <div className="w-full h-full">
            <SidebarContent onClose={() => setOpen(false)} />
          </div>
        </div>
      )}
    </>
  )
}
