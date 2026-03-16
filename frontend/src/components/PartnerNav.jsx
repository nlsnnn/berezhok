import { Link, useNavigate } from 'react-router-dom'
import { Leaf, LogOut, User } from 'lucide-react'
import { useAuth } from '@/context/AuthContext'

export default function PartnerNav() {
  const { logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/partner/login')
  }

  return (
    <nav className="bg-white border-b border-cream-200 sticky top-0 z-40">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 h-14 flex items-center justify-between">
        <Link to="/partner/dashboard" className="flex items-center gap-2 font-semibold text-brand-700">
          <Leaf size={20} className="text-brand-500" />
          <span>Бережок</span>
          <span className="text-cream-400 font-normal text-sm ml-1">Партнёр</span>
        </Link>
        <button onClick={handleLogout} className="btn-ghost gap-1.5 text-sm">
          <LogOut size={16} />
          <span className="hidden sm:inline">Выйти</span>
        </button>
      </div>
    </nav>
  )
}
