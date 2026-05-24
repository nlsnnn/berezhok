import { Search } from 'lucide-react'
import Button from '@/components/ui/actions/Button'
import Input from '@/components/ui/form/Input'
import Select from '@/components/ui/form/Select'

export default function AdminToolbar({
  search,
  onSearchChange,
  status,
  onStatusChange,
  statusOptions = [],
  extraFilters = [],
  onApply,
  onReset,
  loading,
}) {
  return (
    <form className="card grid gap-4 lg:grid-cols-[minmax(0,1fr)_220px_220px_auto_auto]" onSubmit={onApply}>
      <label className="space-y-2">
        <span className="text-sm font-medium text-brand-700">Поиск</span>
        <div className="relative">
          <Search size={16} className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-cream-500" />
          <Input className="pl-9" value={search} onChange={(event) => onSearchChange(event.target.value)} placeholder="Название, email, телефон или ID" />
        </div>
      </label>

      {statusOptions.length > 0 && (
        <label className="space-y-2">
          <span className="text-sm font-medium text-brand-700">Статус</span>
          <Select value={status} onChange={(event) => onStatusChange(event.target.value)}>
            {statusOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </Select>
        </label>
      )}

      {extraFilters.map((filter) => (
        <label key={filter.name} className="space-y-2">
          <span className="text-sm font-medium text-brand-700">{filter.label}</span>
          <Select value={filter.value} onChange={(event) => filter.onChange(event.target.value)}>
            {filter.options.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </Select>
        </label>
      ))}

      <div className="flex items-end">
        <Button type="submit" className="w-full" disabled={loading}>
          Применить
        </Button>
      </div>

      <div className="flex items-end">
        <Button type="button" variant="secondary" className="w-full" onClick={onReset} disabled={loading}>
          Сбросить
        </Button>
      </div>
    </form>
  )
}
