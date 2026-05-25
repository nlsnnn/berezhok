import { Download, X } from 'lucide-react'
import { useInstallPrompt } from '@/hooks/useInstallPrompt'

export default function InstallPrompt() {
  const { canInstall, install, dismiss } = useInstallPrompt()

  if (!canInstall) return null

  return (
    <div className="fixed bottom-4 left-4 right-4 z-50 flex items-center gap-3 rounded-2xl bg-white px-4 py-3 shadow-xl ring-1 ring-black/5 md:left-auto md:right-4 md:w-80">
      <img src="/logo-sm.png" alt="Berezhok" className="h-10 w-10 flex-shrink-0 rounded-xl object-cover" />
      <div className="flex-1 min-w-0">
        <p className="text-sm font-semibold text-gray-900">Установить приложение</p>
        <p className="text-xs text-gray-500">Быстрый доступ с главного экрана</p>
      </div>
      <button
        onClick={install}
        className="flex items-center gap-1.5 rounded-xl bg-brand-500 px-3 py-1.5 text-xs font-semibold text-white hover:bg-brand-600 transition-colors"
      >
        <Download className="h-3.5 w-3.5" />
        Установить
      </button>
      <button
        onClick={dismiss}
        className="flex-shrink-0 text-gray-400 hover:text-gray-600 transition-colors"
        aria-label="Закрыть"
      >
        <X className="h-4 w-4" />
      </button>
    </div>
  )
}
