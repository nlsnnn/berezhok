import { useState } from 'react'
import Input from '@/components/ui/form/Input'
import Label from '@/components/ui/form/Label'
import Select from '@/components/ui/form/Select'
import Button from '@/components/ui/actions/Button'
import TimeInput from '@/components/partner/forms/TimeInput'
import ImageUpload from '@/components/partner/forms/ImageUpload'
import LocationSelect from '@/components/partner/forms/LocationSelect'

export default function BoxForm({ initialData, locations, onSubmit, isLoading }) {
  const [formData, setFormData] = useState({
    location_id: initialData?.location_id || '',
    name: initialData?.name || '',
    description: initialData?.description || '',
    original_price: initialData?.original_price || '',
    discount_price: initialData?.discount_price || '',
    pickup_time_start: initialData?.pickup_time?.start || initialData?.pickup_time_start || '',
    pickup_time_end: initialData?.pickup_time?.end || initialData?.pickup_time_end || '',
    quantity_available: initialData?.quantity_available || initialData?.quantity || '',
    image_url: initialData?.image_url || '',
    status: initialData?.status || 'inactive',
  })

  const [errors, setErrors] = useState({})

  const handleChange = (field, value) => {
    setFormData((prev) => ({ ...prev, [field]: value }))
    if (errors[field]) {
      setErrors((prev) => ({ ...prev, [field]: '' }))
    }
  }

  const validate = () => {
    const nextErrors = {}
    if (!formData.location_id) nextErrors.location_id = 'Выберите локацию'
    if (!formData.name || formData.name.length < 2) nextErrors.name = 'Минимум 2 символа'
    if (!formData.description) nextErrors.description = 'Обязательное поле'
    if (!formData.discount_price || parseFloat(formData.discount_price) <= 0) {
      nextErrors.discount_price = 'Цена должна быть больше 0'
    }
    if (formData.original_price && parseFloat(formData.original_price) <= parseFloat(formData.discount_price || 0)) {
      nextErrors.original_price = 'Должна быть больше цены со скидкой'
    }
    if (!formData.pickup_time_start) nextErrors.pickup_time_start = 'Обязательное поле'
    if (!formData.pickup_time_end) nextErrors.pickup_time_end = 'Обязательное поле'
    if (formData.pickup_time_start && formData.pickup_time_end && formData.pickup_time_end <= formData.pickup_time_start) {
      nextErrors.pickup_time_end = 'Время окончания должно быть позже времени начала'
    }
    if (!formData.quantity_available || parseInt(formData.quantity_available, 10) < 1) {
      nextErrors.quantity_available = 'Количество должно быть больше 0'
    }

    setErrors(nextErrors)
    return Object.keys(nextErrors).length === 0
  }

  const handleSubmit = (e) => {
    e.preventDefault()
    if (!validate()) return
    onSubmit(formData)
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <LocationSelect
        locations={locations}
        value={formData.location_id}
        onChange={(e) => handleChange('location_id', e.target.value)}
        error={errors.location_id}
      />

      <div>
        <Label htmlFor="name">Название бокса *</Label>
        <Input
          id="name"
          value={formData.name}
          onChange={(e) => handleChange('name', e.target.value)}
          placeholder="Например: Утренний бокс"
          error={errors.name}
          maxLength={100}
        />
      </div>

      <div>
        <Label htmlFor="description">Описание *</Label>
        <textarea
          id="description"
          value={formData.description}
          onChange={(e) => handleChange('description', e.target.value)}
          className="input-base min-h-[100px] resize-none"
          placeholder="Опишите содержимое бокса"
          rows={4}
        />
        {errors.description && <p className="mt-1 text-xs text-red-500">{errors.description}</p>}
      </div>

      <div className="grid sm:grid-cols-2 gap-4">
        <div>
          <Label htmlFor="discount_price">Цена со скидкой (₽) *</Label>
          <Input
            id="discount_price"
            type="number"
            step="0.01"
            min="0"
            value={formData.discount_price}
            onChange={(e) => handleChange('discount_price', e.target.value)}
            error={errors.discount_price}
          />
        </div>
        <div>
          <Label htmlFor="original_price">Обычная цена (₽)</Label>
          <Input
            id="original_price"
            type="number"
            step="0.01"
            min="0"
            value={formData.original_price}
            onChange={(e) => handleChange('original_price', e.target.value)}
            error={errors.original_price}
          />
        </div>
      </div>

      <div className="grid sm:grid-cols-2 gap-4">
        <TimeInput
          label="Окно выдачи: с *"
          value={formData.pickup_time_start}
          onChange={(e) => handleChange('pickup_time_start', e.target.value)}
          error={errors.pickup_time_start}
        />
        <TimeInput
          label="Окно выдачи: до *"
          value={formData.pickup_time_end}
          onChange={(e) => handleChange('pickup_time_end', e.target.value)}
          error={errors.pickup_time_end}
        />
      </div>

      <div>
        <Label htmlFor="quantity_available">Количество *</Label>
        <Input
          id="quantity_available"
          type="number"
          min="1"
          value={formData.quantity_available}
          onChange={(e) => handleChange('quantity_available', e.target.value)}
          error={errors.quantity_available}
        />
      </div>

      <ImageUpload value={formData.image_url} onChange={(url) => handleChange('image_url', url)} error={errors.image_url} />

      <div>
        <Label htmlFor="status">Статус *</Label>
        <Select id="status" value={formData.status} onChange={(e) => handleChange('status', e.target.value)}>
          <option value="inactive">Неактивен</option>
          <option value="active">Активен</option>
          <option value="draft">Черновик</option>
        </Select>
      </div>

      <Button type="submit" disabled={isLoading} className="w-full">
        {isLoading ? 'Сохраняем...' : initialData ? 'Сохранить изменения' : 'Создать бокс'}
      </Button>
    </form>
  )
}
