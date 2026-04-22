import { useEffect } from 'react'
import { X } from 'lucide-react'
import { cn } from '@/lib/utils'

export default function Modal({
  open,
  onClose,
  title,
  children,
  className,
  closeOnBackdrop = true,
  closeOnEscape = true,
  showCloseButton = true,
}) {
  useEffect(() => {
    if (!open) return
    const handler = (e) => {
      if (e.key === 'Escape' && closeOnEscape) onClose?.()
    }
    document.addEventListener('keydown', handler)
    return () => document.removeEventListener('keydown', handler)
  }, [closeOnEscape, onClose, open])

  useEffect(() => {
    document.body.style.overflow = open ? 'hidden' : ''
    return () => {
      document.body.style.overflow = ''
    }
  }, [open])

  if (!open) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        className="absolute inset-0 bg-black/40 backdrop-blur-sm"
        onClick={closeOnBackdrop ? onClose : undefined}
      />
      <div className={cn('relative bg-white rounded-2xl shadow-xl w-full max-w-lg max-h-[90vh] overflow-y-auto', className)}>
        <div className="flex items-center justify-between p-6 border-b border-cream-200">
          <h3 className="text-lg font-semibold text-brand-900">{title}</h3>
          {showCloseButton && (
            <button onClick={onClose} className="btn-ghost p-1.5 -mr-1.5">
              <X size={18} />
            </button>
          )}
        </div>
        <div className="p-6">{children}</div>
      </div>
    </div>
  )
}
