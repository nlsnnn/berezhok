import TimeInput from './TimeInput'

const DAYS = [
  { key: 'monday', label: 'Понедельник' },
  { key: 'tuesday', label: 'Вторник' },
  { key: 'wednesday', label: 'Среда' },
  { key: 'thursday', label: 'Четверг' },
  { key: 'friday', label: 'Пятница' },
  { key: 'saturday', label: 'Суббота' },
  { key: 'sunday', label: 'Воскресенье' },
]

const DEFAULT_OPEN = '09:00'
const DEFAULT_CLOSE = '21:00'

function parseHours(value) {
  if (!value || typeof value !== 'string') return { open: '', close: '', closed: true }
  if (value === 'closed') return { open: '', close: '', closed: true }
  const [open = '', close = ''] = value.split('-')
  return { open, close, closed: false }
}

function formatHours({ open, close }) {
  return `${open}-${close}`
}

export default function WorkingHoursEditor({ value, onChange, label = 'График работы' }) {
  const data = value || {}

  const updateDay = (dayKey, next) => {
    const updated = { ...data }
    if (next.closed) {
      delete updated[dayKey]
    } else {
      updated[dayKey] = formatHours(next)
    }
    onChange(updated)
  }

  return (
    <div>
      {label && <label className="block text-sm font-medium text-brand-700 mb-2">{label}</label>}
      <div className="space-y-2">
        {DAYS.map(({ key, label: dayLabel }) => {
          const parsed = parseHours(data[key])
          return (
            <div key={key} className="flex items-center gap-3 p-3 rounded-xl bg-cream-50 border border-cream-200">
              <label className="flex items-center gap-2 w-44 shrink-0 cursor-pointer">
                <input
                  type="checkbox"
                  checked={!parsed.closed}
                  onChange={(e) => {
                    const open = parsed.open || DEFAULT_OPEN
                    const close = parsed.close || DEFAULT_CLOSE
                    updateDay(key, { open, close, closed: !e.target.checked })
                  }}
                  className="h-4 w-4 accent-brand-600"
                />
                <span className="text-sm text-brand-900">{dayLabel}</span>
              </label>

              {parsed.closed ? (
                <span className="text-sm text-brand-500">Выходной</span>
              ) : (
                <div className="flex items-center gap-2 flex-1">
                  <TimeInput
                    value={parsed.open}
                    onChange={(e) => updateDay(key, { open: e.target.value, close: parsed.close || DEFAULT_CLOSE, closed: false })}
                  />
                  <span className="text-brand-500">—</span>
                  <TimeInput
                    value={parsed.close}
                    onChange={(e) => updateDay(key, { open: parsed.open || DEFAULT_OPEN, close: e.target.value, closed: false })}
                  />
                </div>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}
