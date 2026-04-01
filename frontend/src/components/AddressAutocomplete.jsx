import { useState, useRef, useEffect } from 'react'
import { MapPin, Loader2, X } from 'lucide-react'
import { useAddressSearch } from '@/hooks/useAddressSearch'
import { cn } from '@/lib/utils'

/**
 * AddressAutocomplete
 * props:
 *   value       – current display value (string)
 *   onChange    – called with { address, display_name, latitude, longitude }
 *   placeholder
 *   error
 *   disabled
 */
export default function AddressAutocomplete({ value, onChange, placeholder = 'Начните вводить адрес...', error, disabled }) {
  const { query, setQuery, suggestions, loading, clear } = useAddressSearch(400)
  const [open, setOpen] = useState(false)
  const [selected, setSelected] = useState(false)
  const wrapperRef = useRef(null)

  // If external value changes (e.g. form reset)
  useEffect(() => {
    if (!value) {
      clear()
    }
  }, [value, clear])

  // Close dropdown on outside click
  useEffect(() => {
    const handler = (e) => {
      if (wrapperRef.current && !wrapperRef.current.contains(e.target)) {
        setOpen(false)
      }
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [])

  const handleInput = (e) => {
    setSelected(false)
    setQuery(e.target.value)
    setOpen(true)
  }

  const handleSelect = (suggestion) => {
    setSelected(true)
    setQuery(suggestion.display_name)
    setOpen(false)
    onChange(suggestion)
  }

  const handleClear = () => {
    clear()
    setSelected(false)
    setOpen(false)
    onChange(null)
  }

  const displayValue = selected ? query : (query || value || '')

  return (
    <div ref={wrapperRef} className="relative w-full">
      <div className="relative">
        <MapPin size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-cream-400 pointer-events-none" />
        <input
          type="text"
          value={displayValue}
          onChange={handleInput}
          onFocus={() => query.length >= 3 && setOpen(true)}
          placeholder={placeholder}
          disabled={disabled}
          className={cn(
            'input-base pl-9 pr-8',
            error && 'border-red-400 focus:ring-red-400'
          )}
        />
        <div className="absolute right-3 top-1/2 -translate-y-1/2 flex items-center">
          {loading && <Loader2 size={14} className="animate-spin text-cream-400" />}
          {!loading && displayValue && (
            <button type="button" onClick={handleClear} className="text-cream-400 hover:text-brand-500">
              <X size={14} />
            </button>
          )}
        </div>
      </div>

      {error && <p className="mt-1 text-xs text-red-500">{error}</p>}

      {open && suggestions.length > 0 && (
        <ul className="absolute z-20 mt-1 w-full bg-white rounded-xl border border-cream-200 shadow-lg overflow-hidden">
          {suggestions.map((s, i) => (
            <li key={i}>
              <button
                type="button"
                className="w-full px-4 py-3 text-left text-sm hover:bg-cream-100 transition-colors flex gap-2 items-start"
                onClick={() => handleSelect(s)}
              >
                <MapPin size={14} className="text-brand-400 mt-0.5 shrink-0" />
                <span className="text-brand-800">{s.display_name}</span>
              </button>
            </li>
          ))}
        </ul>
      )}

      {open && !loading && query.length >= 3 && suggestions.length === 0 && (
        <div className="absolute z-20 mt-1 w-full bg-white rounded-xl border border-cream-200 shadow-lg p-4 text-sm text-cream-500 text-center">
          Адресов не найдено
        </div>
      )}
    </div>
  )
}
