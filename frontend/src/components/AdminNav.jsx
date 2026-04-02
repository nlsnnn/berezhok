import { Link, useLocation } from 'react-router-dom'
import { ShieldCheck, FileText, Store, AlertCircle, BarChart3 } from 'lucide-react'
import { cn } from '@/lib/utils'

const links = [
  { to: '/admin', label: 'Заявки', icon: FileText },
  { to: '/admin/partners', label: 'Партнёры', icon: Store },
  { to: '/admin/disputes', label: 'Споры', icon: AlertCircle },
  { to: '/admin/stats', label: 'Статистика', icon: BarChart3 },
]

export default function AdminNav() {
  const location = useLocation()

  return (
    <nav className="bg-brand-800 text-white sticky top-0 z-40">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 h-14 flex items-center justify-between">
        <Link to="/admin" className="flex items-center gap-2 font-semibold text-white">
          <img src="/logo.png" alt="Бережок" className="w-6 h-6 rounded-md object-cover" />
          <span>Бережок</span>
          <span className="text-brand-300 font-normal text-sm ml-1">Администратор</span>
        </Link>
        <div className="hidden sm:flex items-center gap-1">
          {links.map(({ to, label, icon: Icon }) => {
            const isActive = location.pathname === to
            return (
              <Link
                key={to}
                to={to}
                className={cn(
                  'flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm transition-colors',
                  isActive
                    ? 'bg-white/15 text-white'
                    : 'text-brand-300 hover:text-white hover:bg-white/10'
                )}
              >
                <Icon size={14} />
                <span>{label}</span>
              </Link>
            )
          })}
        </div>
        <div className="sm:hidden flex items-center gap-1.5 text-sm text-brand-300">
          <ShieldCheck size={16} />
          <span>Панель управления</span>
        </div>
      </div>
    </nav>
  )
}
