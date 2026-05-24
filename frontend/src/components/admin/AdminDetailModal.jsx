import Modal from '@/components/ui/overlay/Modal'
import Spinner from '@/components/ui/feedback/Spinner'

export function DetailGrid({ children }) {
  return <div className="grid gap-4 sm:grid-cols-2">{children}</div>
}

export function DetailItem({ label, value, wide = false }) {
  return (
    <div className={wide ? 'sm:col-span-2' : ''}>
      <p className="mb-1 text-xs font-medium uppercase tracking-wider text-cream-500">{label}</p>
      <p className="break-words text-sm font-medium text-brand-800">{value || '—'}</p>
    </div>
  )
}

export default function AdminDetailModal({ open, title, loading, children, actions, onClose }) {
  return (
    <Modal open={open} onClose={onClose} title={title} className="max-w-2xl">
      {loading ? (
        <div className="flex justify-center py-16">
          <Spinner size={32} />
        </div>
      ) : (
        <div className="space-y-6">
          {children}
          {actions && <div className="flex flex-wrap gap-3 border-t border-cream-200 pt-4">{actions}</div>}
        </div>
      )}
    </Modal>
  )
}
