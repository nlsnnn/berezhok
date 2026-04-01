import { NavLink, useNavigate } from 'react-router-dom'
import { LayoutDashboard, LogOut, MapPin, Menu, Package, QrCode, X } from 'lucide-react'
import { useState } from 'react'
import { useAuth } from '@/context/AuthContext'
import { cn } from '@/lib/utils'

const links = [
  { to: '/partner/dashboard', label: 'Дашборд', icon: LayoutDashboard },
  { to: '/partner/locations', label: 'Локации', icon: MapPin },
  { to: '/partner/boxes', label: 'Боксы', icon: Package },
  { to: '/partner/orders/pickup', label: 'Выдача', icon: QrCode },
]

function SidebarContent({ onClose }) {
  const { logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/partner/login')
  }

  return (
    <div className="h-full flex flex-col bg-white border-r border-cream-200">
      <div className="h-20 px-5 border-b border-cream-200 flex items-center justify-between">
        <NavLink to="/partner/dashboard" className="flex items-center gap-3" onClick={onClose}>
          <img src="/logo.png" alt="Бережок" className="w-10 h-10 rounded-xl object-cover" />
          <div>
            <p className="font-bold text-brand-800 text-base leading-none">Бережок</p>
            <p className="text-xs text-brand-500 mt-1">Панель партнера</p>
          </div>
        </NavLink>
        <button className="md:hidden btn-ghost p-2" onClick={onClose}>
          <X size={18} />
        </button>
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
                isActive
                  ? 'bg-brand-500 text-white shadow-sm'
                  : 'text-brand-700 hover:bg-cream-100'
              )
            }
          >
            <Icon size={17} />
            {label}
          </NavLink>
        ))}
      </div>

      <div className="p-3 border-t border-cream-200">
        <button onClick={handleLogout} className="btn-ghost w-full justify-start gap-3 text-sm">
          <LogOut size={16} />
          Выйти
        </button>
      </div>
    </div>
  )
}

export default function PartnerSidebar() {
  const [open, setOpen] = useState(false)

  return (
    <>
      <aside className="hidden md:block w-72 shrink-0">
        <SidebarContent />
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
