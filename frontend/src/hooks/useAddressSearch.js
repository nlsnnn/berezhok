import { useState, useEffect, useRef, useCallback } from 'react'
import { searchAddress } from '@/api/geocoder'

/**
 * useAddressSearch — debounced Yandex Geocoder address search
 * Returns { suggestions, loading, error, search, clear }
 */
export function useAddressSearch(delay = 400) {
  const [query, setQuery] = useState('')
  const [suggestions, setSuggestions] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const timerRef = useRef(null)

  useEffect(() => {
    if (!query || query.trim().length < 3) {
      setSuggestions([])
      return
    }

    clearTimeout(timerRef.current)
    timerRef.current = setTimeout(async () => {
      setLoading(true)
      setError(null)
      try {
        const results = await searchAddress(query)
        setSuggestions(results)
      } catch {
        setError('Не удалось загрузить подсказки адресов')
        setSuggestions([])
      } finally {
        setLoading(false)
      }
    }, delay)

    return () => clearTimeout(timerRef.current)
  }, [query, delay])

  const clear = useCallback(() => {
    setQuery('')
    setSuggestions([])
    setError(null)
  }, [])

  return { query, setQuery, suggestions, loading, error, clear }
}
