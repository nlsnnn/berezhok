import { useState } from 'react'
import Input from '@/components/ui/form/Input'
import Label from '@/components/ui/form/Label'
import Toggle from '@/components/ui/form/Toggle'
import Button from '@/components/ui/actions/Button'
import TimeInput from '@/components/partner/forms/TimeInput'
import ImageUpload from '@/components/partner/forms/ImageUpload'
import LocationSelect from '@/components/partner/forms/LocationSelect'

const isInitiallyFree = (initialData) => {
  if (!initialData) return false
  const value = initialData.discount_price
  if (value === undefined || value === null || value === '') return false
  return Number(value) === 0
}

export default function BoxForm({
  initialData,
  locations,
  onSubmit,
  isLoading,
  canActivateBoxes = true,
  mode = 'create',
}) {
  const [formData, setFormData] = useState({
    location_id: initialData?.location_id || '',
    name: initialData?.name || '',
    description: initialData?.description || '',
    original_price: initialData?.original_price || '',
    discount_price: initialData?.discount_price ?? '',
    pickup_time_start: initialData?.pickup_time?.start || initialData?.pickup_time_start || '',
    pickup_time_end: initialData?.pickup_time?.end || initialData?.pickup_time_end || '',
    quantity_available: initialData?.quantity_available || initialData?.quantity || '',
    image_url: initialData?.image_url || '',
  })

  const initialStatus = initialData?.status || 'draft'
  const [isFree, setIsFree] = useState(isInitiallyFree(initialData))
  const [isPublished, setIsPublished] = useState(initialStatus === 'active')
  const [errors, setErrors] = useState({})

  const handleChange = (field, value) => {
    setFormData((prev) => ({ ...prev, [field]: value }))
    if (errors[field]) {
      setErrors((prev) => ({ ...prev, [field]: '' }))
    }
  }

  const handleFreeToggle = (next) => {
    setIsFree(next)
    if (next) {
      setFormData((prev) => ({ ...prev, discount_price: '', original_price: '' }))
      setErrors((prev) => ({ ...prev, discount_price: '', original_price: '' }))
    }
  }

  const validate = () => {
    const nextErrors = {}
    if (!formData.location_id) nextErrors.location_id = 'Выберите локацию'
    if (!formData.name || formData.name.length < 2) nextErrors.name = 'Минимум 2 символа'
    if (!formData.description) nextErrors.description = 'Обязательное поле'

    if (!isFree) {
      if (formData.discount_price === '' || parseFloat(formData.discount_price) < 0) {
        nextErrors.discount_price = 'Укажите цену'
      }
      if (
        formData.original_price &&
        parseFloat(formData.original_price) <= parseFloat(formData.discount_price || 0)
      ) {
        nextErrors.original_price = 'Должна быть больше цены со скидкой'
      }
    }

    if (!formData.pickup_time_start) nextErrors.pickup_time_start = 'Обязательное поле'
    if (!formData.pickup_time_end) nextErrors.pickup_time_end = 'Обязательное поле'
    if (
      formData.pickup_time_start &&
      formData.pickup_time_end &&
      formData.pickup_time_end <= formData.pickup_time_start
    ) {
      nextErrors.pickup_time_end = 'Время окончания должно быть позже времени начала'
    }
    if (!formData.quantity_available || parseInt(formData.quantity_available, 10) < 1) {
      nextErrors.quantity_available = 'Количество должно быть больше 0'
    }

    setErrors(nextErrors)
    return Object.keys(nextErrors).length === 0
  }

  const buildPayload = (statusOverride) => {
    return {
      ...formData,
      discount_price: isFree ? '0' : formData.discount_price,
      original_price: isFree ? '' : formData.original_price,
      status: statusOverride,
    }
  }

  const submit = (status) => {
    if (!validate()) return
    onSubmit(buildPayload(status))
  }

  const handleFormSubmit = (e) => {
    e.preventDefault()
  }

  const showDraftActions = mode === 'create' || initialStatus === 'draft'
  const isSoldOut = initialStatus === 'sold_out'

  return (
    <form onSubmit={handleFormSubmit} className="space-y-6">
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

      <div className="space-y-3">
        <Toggle
          id="is_free"
          checked={isFree}
          onChange={handleFreeToggle}
          label="Бесплатный бокс"
          description="Отдаёте еду бесплатно — поля цены скрыты"
        />
        {!isFree && (
          <div className="grid sm:grid-cols-2 gap-4">
            <div>
              <Label htmlFor="discount_price">Цена (₽) *</Label>
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
              <Label htmlFor="original_price">Цена до скидки (₽)</Label>
              <Input
                id="original_price"
                type="number"
                step="0.01"
                min="0"
                value={formData.original_price}
                onChange={(e) => handleChange('original_price', e.target.value)}
                error={errors.original_price}
                placeholder="Необязательно"
              />
            </div>
          </div>
        )}
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

      {mode === 'edit' && !showDraftActions && !isSoldOut && (
        <div className="rounded-xl border border-cream-200 bg-cream-50 px-4 py-3">
          <Toggle
            id="is_published"
            checked={isPublished}
            onChange={(next) => setIsPublished(next)}
            disabled={!isPublished && !canActivateBoxes}
            label="Опубликовано"
            description={
              !isPublished && !canActivateBoxes
                ? 'Заполните юридические данные партнёра, чтобы публиковать боксы'
                : isPublished
                  ? 'Бокс виден покупателям'
                  : 'Бокс скрыт от покупателей'
            }
          />
        </div>
      )}

      {isSoldOut && (
        <div className="rounded-xl border border-cream-200 bg-cream-50 px-4 py-3 text-sm text-brand-700">
          Бокс распродан. Чтобы возобновить продажи, увеличьте количество и сохраните изменения.
        </div>
      )}

      <div className="flex flex-col sm:flex-row gap-3 pt-2">
        {showDraftActions ? (
          <>
            <Button
              type="button"
              variant="secondary"
              disabled={isLoading}
              onClick={() => submit('draft')}
              className="sm:flex-1"
            >
              {isLoading ? 'Сохраняем...' : 'Сохранить как черновик'}
            </Button>
            <Button
              type="button"
              disabled={isLoading || !canActivateBoxes}
              onClick={() => submit('active')}
              className="sm:flex-1"
              title={!canActivateBoxes ? 'Заполните юридические данные, чтобы публиковать боксы' : undefined}
            >
              {isLoading ? 'Сохраняем...' : 'Опубликовать'}
            </Button>
          </>
        ) : (
          <Button
            type="button"
            disabled={isLoading}
            onClick={() =>
              submit(isSoldOut ? 'sold_out' : isPublished ? 'active' : 'inactive')
            }
            className="w-full"
          >
            {isLoading ? 'Сохраняем...' : 'Сохранить изменения'}
          </Button>
        )}
      </div>
    </form>
  )
}
