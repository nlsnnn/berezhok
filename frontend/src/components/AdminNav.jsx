import { Link } from 'react-router-dom'
import { Leaf, ShieldCheck } from 'lucide-react'

export default function AdminNav() {
  return (
    <nav className="bg-brand-800 text-white sticky top-0 z-40">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 h-14 flex items-center justify-between">
        <Link to="/admin" className="flex items-center gap-2 font-semibold text-white">
          <Leaf size={20} className="text-brand-300" />
          <span>Бережок</span>
          <span className="text-brand-300 font-normal text-sm ml-1">Администратор</span>
        </Link>
        <div className="flex items-center gap-1.5 text-sm text-brand-300">
          <ShieldCheck size={16} />
          <span className="hidden sm:inline">Панель управления</span>
        </div>
      </div>
    </nav>
  )
}
