import { cn } from '@/lib/utils'

export default function LocationPinsSelector({ availablePins = [], selectedCodes = [], onChange }) {
  const toggle = (code) => {
    if (selectedCodes.includes(code)) {
      onChange(selectedCodes.filter((c) => c !== code))
    } else {
      onChange([...selectedCodes, code])
    }
  }

  return (
    <div className="flex flex-wrap gap-2">
      {availablePins.map((pin) => {
        const selected = selectedCodes.includes(pin.code)
        return (
          <button
            key={pin.code}
            type="button"
            onClick={() => toggle(pin.code)}
            className={cn(
              'px-3 py-1.5 rounded-full text-sm font-medium border transition-colors',
              selected
                ? 'bg-brand-600 text-white border-brand-600'
                : 'bg-white text-brand-700 border-brand-200 hover:border-brand-400'
            )}
          >
            {pin.name_ru}
          </button>
        )
      })}
    </div>
  )
}
