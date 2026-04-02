import Select from '@/components/ui/form/Select'

export default function LocationSelect({ locations, value, onChange, error }) {
  if (!locations?.length) {
    return (
      <div>
        <label className="block text-sm font-medium text-brand-700 mb-1">Локация</label>
        <div className="input-base bg-gray-50 text-gray-500 cursor-not-allowed">Нет доступных локаций</div>
      </div>
    )
  }

  return (
    <div>
      <label className="block text-sm font-medium text-brand-700 mb-1">Локация *</label>
      <Select value={value} onChange={onChange} error={error}>
        <option value="">Выберите локацию</option>
        {locations.map((location) => (
          <option key={location.id} value={location.id}>
            {location.name} - {location.address}
          </option>
        ))}
      </Select>
    </div>
  )
}
