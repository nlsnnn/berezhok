import { Link } from 'react-router-dom'

export default function LandingNav() {
  return (
    <nav className="fixed top-0 inset-x-0 z-40 bg-white/90 backdrop-blur-md border-b border-cream-200">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 h-16 flex items-center justify-between">
        <Link to="/" className="flex items-center gap-2 font-semibold text-brand-700 text-lg">
          <img src="/logo.png" alt="Бережок" className="w-7 h-7 rounded-lg object-cover" />
          Бережок
        </Link>
        <div className="hidden sm:flex items-center gap-6 text-sm text-brand-600">
          <a href="#how-it-works" className="hover:text-brand-500 transition-colors">Как это работает</a>
          <a href="#apply" className="btn-primary text-sm px-4 py-2">Стать партнером</a>
        </div>
        <a href="#apply" className="sm:hidden btn-primary text-xs px-3 py-2">Стать партнером</a>
      </div>
    </nav>
  )
}
