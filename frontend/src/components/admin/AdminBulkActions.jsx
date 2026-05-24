import { CheckSquare, Copy, Square } from 'lucide-react'
import { toast } from 'sonner'
import Button from '@/components/ui/actions/Button'

async function copyIds(ids) {
  const text = ids.join('\n')
  if (navigator.clipboard?.writeText) {
    await navigator.clipboard.writeText(text)
    return
  }
  window.prompt('Скопируйте ID выбранных сущностей', text)
}

export default function AdminBulkActions({
  selectedIds,
  pageIds,
  allPageSelected,
  onToggleAll,
  onClear,
  actions = [],
  loading,
}) {
  if (pageIds.length === 0) return null

  const handleCopy = async () => {
    try {
      await copyIds(selectedIds)
      toast.success('ID выбранных сущностей скопированы')
    } catch {
      toast.error('Не удалось скопировать ID')
    }
  }

  return (
    <div className="card flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div className="flex flex-wrap items-center gap-3">
        <button type="button" className="btn-secondary px-4" onClick={onToggleAll} disabled={loading}>
          {allPageSelected ? <CheckSquare size={16} /> : <Square size={16} />}
          {allPageSelected ? 'Снять выбор на странице' : 'Выбрать страницу'}
        </button>
        <span className="text-sm font-medium text-brand-700">Выбрано: {selectedIds.length}</span>
      </div>

      {selectedIds.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {actions.map((action) => (
            <Button
              key={action.label}
              type="button"
              variant={action.variant || 'secondary'}
              size="sm"
              onClick={() => action.onClick(selectedIds)}
              disabled={loading || action.disabled}
            >
              {action.icon}
              {action.label}
            </Button>
          ))}
          <Button type="button" variant="secondary" size="sm" onClick={handleCopy} disabled={loading}>
            <Copy size={16} />
            Скопировать ID
          </Button>
          <Button type="button" variant="ghost" size="sm" onClick={onClear} disabled={loading}>
            Снять выбор
          </Button>
        </div>
      )}
    </div>
  )
}
