const GEOCODER_API_KEY = import.meta.env.VITE_YANDEX_GEOCODER_API_KEY

/**
 * Search addresses via Yandex Geocoder API
 * Returns array of { display_name, address, latitude, longitude }
 */
export const searchAddress = async (query) => {
  if (!query || query.trim().length < 3) return []

  const params = new URLSearchParams({
    apikey: GEOCODER_API_KEY,
    geocode: query,
    format: 'json',
    results: 5,
    lang: 'ru_RU',
  })

  const res = await fetch(
    `https://geocode-maps.yandex.ru/v1/?${params}`
  )

  if (!res.ok) throw new Error('Geocoder request failed')

  const json = await res.json()
  const featureMembers =
    json?.response?.GeoObjectCollection?.featureMember ?? []

  return featureMembers.map((f) => {
    const geo = f.GeoObject
    const pos = geo.Point.pos.split(' ') // "lon lat"
    return {
      display_name: geo.metaDataProperty.GeocoderMetaData.text,
      address: geo.name,
      longitude: parseFloat(pos[0]),
      latitude: parseFloat(pos[1]),
    }
  })
}
