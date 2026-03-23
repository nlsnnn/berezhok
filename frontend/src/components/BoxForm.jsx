import { useState, useEffect } from 'react'
import Input from '@/components/ui/Input'
import Label from '@/components/ui/Label'
import Select from '@/components/ui/Select'
import Button from '@/components/ui/Button'
import TimeInput from '@/components/ui/TimeInput'
import ImageUpload from '@/components/ui/ImageUpload'
import LocationSelect from '@/components/ui/LocationSelect'

export default function BoxForm({ initialData, locations, onSubmit, isLoading }) {
  const [formData, setFormData] = useState({
    location_id: initialData?.location_id || '',
    name: initialData?.name || '',
    description: initialData?.description || '',
    original_price: initialData?.original_price || '',
    discount_price: initialData?.discount_price || '',
    pickup_time_start: initialData?.pickup_time?.start || '',
    pickup_time_end: initialData?.pickup_time?.end || '',
    quantity: initialData?.quantity || '',
    image_url: initialData?.image_url || '',
    status: initialData?.status || 'draft',
  })

  const [errors, setErrors] = useState({})

  const handleChange = (field, value) => {
    setFormData(prev => ({ ...prev, [field]: value }))
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: '' }))
    }
  }

  const validate = () => {
    const newErrors = {}

    if (!formData.location_id) newErrors.location_id = 'Выберите локацию'
    if (!formData.name || formData.name.length < 2) newErrors.name = 'Минимум 2 символа'
    if (formData.name && formData.name.length > 100) newErrors.name = 'Максимум 100 символов'
    if (!formData.description) newErrors.description = 'Обязательное поле'
    if (!formData.discount_price || parseFloat(formData.discount_price) <= 0) {
      newErrors.discount_price = 'Цена должна быть больше 0'
    }
    if (formData.original_price && parseFloat(formData.original_price) <= parseFloat(formData.discount_price)) {
      newErrors.original_price = 'Должна быть больше цены со скидкой'
    }
    if (!formData.pickup_time_start) newErrors.pickup_time_start = 'Обязательное поле'
    if (!formData.pickup_time_end) newErrors.pickup_time_end = 'Обязательное поле'
    if (formData.pickup_time_start && formData.pickup_time_end) {
      if (formData.pickup_time_end <= formData.pickup_time_start) {
        newErrors.pickup_time_end = 'Время окончания должно быть позже времени начала'
      }
    }
    if (!formData.quantity || parseInt(formData.quantity) < 1) {
      newErrors.quantity = 'Количество должно быть больше 0'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = (e) => {
    e.preventDefault()
    if (validate()) {
      onSubmit(formData)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Location */}
      <LocationSelect
        locations={locations}
        value={formData.location_id}
        onChange={(e) => handleChange('location_id', e.target.value)}
        error={errors.location_id}
      />

      {/* Name */}
      <div>
        <Label htmlFor="name">Название бокса *</Label>
        <Input
          id="name"
          value={formData.name}
          onChange={(e) => handleChange('name', e.target.value)}
          placeholder="Например: Утренний завтрак"
          error={errors.name}
          maxLength={100}
        />
      </div>

      {/* Description */}
      <div>
        <Label htmlFor="description">Описание *</Label>
        <textarea
          id="description"
          value={formData.description}
          onChange={(e) => handleChange('description', e.target.value)}
          className="input-base min-h-[100px] resize-none"
          placeholder="Расскажите, что входит в бокс..."
          rows={4}
        />
        {errors.description && (
          <p className="mt-1 text-xs text-red-500">{errors.description}</p>
        )}
      </div>

      {/* Prices */}
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
            placeholder="1500"
            error={errors.discount_price}
          />
        </div>
        <div>
          <Label htmlFor="original_price">Оригинальная цена (₽)</Label>
          <Input
            id="original_price"
            type="number"
            step="0.01"
            min="0"
            value={formData.original_price}
            onChange={(e) => handleChange('original_price', e.target.value)}
            placeholder="2500"
            error={errors.original_price}
          />
          <p className="mt-1 text-xs text-brand-500">Для отображения скидки</p>
        </div>
      </div>

      {/* Pickup time */}
      <div className="grid sm:grid-cols-2 gap-4">
        <TimeInput
          label="Время получения с *"
          value={formData.pickup_time_start}
          onChange={(e) => handleChange('pickup_time_start', e.target.value)}
          error={errors.pickup_time_start}
        />
        <TimeInput
          label="Время получения до *"
          value={formData.pickup_time_end}
          onChange={(e) => handleChange('pickup_time_end', e.target.value)}
          error={errors.pickup_time_end}
        />
      </div>

      {/* Quantity */}
      <div>
        <Label htmlFor="quantity">Количество *</Label>
        <Input
          id="quantity"
          type="number"
          min="1"
          value={formData.quantity}
          onChange={(e) => handleChange('quantity', e.target.value)}
          placeholder="10"
          error={errors.quantity}
        />
      </div>

      {/* Image */}
      <ImageUpload
        value={formData.image_url}
        onChange={(url) => handleChange('image_url', url)}
        error={errors.image_url}
      />

      {/* Status */}
      <div>
        <Label htmlFor="status">Статус *</Label>
        <Select
          id="status"
          value={formData.status}
          onChange={(e) => handleChange('status', e.target.value)}
        >
          <option value="draft">Черновик</option>
          <option value="active">Активен</option>
          <option value="inactive">Неактивен</option>
        </Select>
        <p className="mt-1 text-xs text-brand-500">
          Активные боксы отображаются в приложении
        </p>
      </div>

      {/* Submit */}
      <div className="flex gap-3 pt-4">
        <Button type="submit" disabled={isLoading} className="flex-1">
          {isLoading ? 'Сохранение...' : initialData ? 'Сохранить изменения' : 'Создать бокс'}
        </Button>
      </div>
    </form>
  )
}
