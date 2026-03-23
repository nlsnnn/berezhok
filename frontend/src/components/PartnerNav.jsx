import { Link, NavLink, useNavigate, useLocation } from 'react-router-dom'
import { Leaf, LogOut, LayoutDashboard, MapPin, Package } from 'lucide-react'
import { useAuth } from '@/context/AuthContext'
import { cn } from '@/lib/utils'

export default function PartnerNav() {
  const { logout } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()

  const handleLogout = () => {
    logout()
    navigate('/partner/login')
  }

  const navLinks = [
    { to: '/partner/dashboard', label: 'Дашборд', icon: LayoutDashboard },
    { to: '/partner/locations', label: 'Локации', icon: MapPin },
    { to: '/partner/boxes', label: 'Боксы', icon: Package },
  ]

  return (
    <nav className="bg-white border-b border-cream-200 sticky top-0 z-40">
      <div className="max-w-7xl mx-auto px-4 sm:px-6">
        <div className="flex items-center justify-between h-14">
          {/* Logo */}
          <Link to="/partner/dashboard" className="flex items-center gap-2 font-semibold text-brand-700">
            <Leaf size={20} className="text-brand-500" />
            <span>Бережок</span>
            <span className="text-cream-400 font-normal text-sm ml-1">Партнёр</span>
          </Link>

          {/* Nav links */}
          <div className="hidden md:flex items-center gap-1">
            {navLinks.map(({ to, label, icon: Icon }) => (
              <NavLink
                key={to}
                to={to}
                className={({ isActive }) =>
                  cn(
                    'px-3 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-2',
                    isActive
                      ? 'bg-brand-50 text-brand-700'
                      : 'text-brand-600 hover:bg-cream-100 hover:text-brand-800'
                  )
                }
              >
                <Icon size={16} />
                {label}
              </NavLink>
            ))}
          </div>

          {/* Logout */}
          <button onClick={handleLogout} className="btn-ghost gap-1.5 text-sm">
            <LogOut size={16} />
            <span className="hidden sm:inline">Выйти</span>
          </button>
        </div>
      </div>
    </nav>
  )
}
