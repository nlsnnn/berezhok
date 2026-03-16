import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { toast } from 'sonner'
import { Plus, MapPin, Phone } from 'lucide-react'
import { createLocation } from '@/api/partner'
import { BUSINESS_CATEGORIES } from '@/lib/constants'
import { getErrorMessage } from '@/lib/utils'
import PartnerNav from '@/components/PartnerNav'
import Input from '@/components/ui/Input'
import Select from '@/components/ui/Select'
import Label from '@/components/ui/Label'
import Button from '@/components/ui/Button'
import AddressAutocomplete from '@/components/AddressAutocomplete'

const INITIAL = {
  name: '',
  address: '',
  latitude: null,
  longitude: null,
  category_code: '',
  phone: '',
}

export default function CreateLocationPage() {
  const [form, setForm] = useState(INITIAL)
  const [errors, setErrors] = useState({})
  const navigate = useNavigate()

  const mutation = useMutation({
    mutationFn: createLocation,
    onSuccess: (data) => {
      toast.success('Локация создана')
      navigate('/partner/dashboard')
    },
    onError: (err) => toast.error(getErrorMessage(err)),
  })

  const set = (field) => (e) => setForm((f) => ({ ...f, [field]: e.target.value }))

  const validate = () => {
    const e = {}
    if (!form.name.trim()) e.name = 'Введите название'
    if (!form.category_code) e.category_code = 'Выберите категорию'
    if (!form.address || !form.latitude) e.address = 'Выберите адрес из списка'
    return e
  }

  const handleSubmit = (e) => {
    e.preventDefault()
    const errs = validate()
    if (Object.keys(errs).length) { setErrors(errs); return }
    setErrors({})

    const payload = {
      name: form.name,
      address: form.address,
      latitude: form.latitude,
      longitude: form.longitude,
      category_code: form.category_code,
    }
    if (form.phone.trim()) payload.phone = form.phone

    mutation.mutate(payload)
  }

  const handleAddressSelect = (suggestion) => {
    if (!suggestion) {
      setForm((f) => ({ ...f, address: '', latitude: null, longitude: null }))
      return
    }
    setForm((f) => ({
      ...f,
      address: suggestion.display_name,
      latitude: suggestion.latitude,
      longitude: suggestion.longitude,
    }))
    setErrors((e) => ({ ...e, address: undefined }))
  }

  return (
    <div className="min-h-screen flex flex-col bg-cream-50">
      <PartnerNav />
      <main className="flex-1 max-w-2xl mx-auto w-full px-4 sm:px-6 py-8">
        <div className="mb-8">
          <h1 className="text-2xl font-bold text-brand-900">Добавить локацию</h1>
          <p className="text-sm text-brand-500 mt-1">Укажите данные вашей точки продаж</p>
        </div>

        <form onSubmit={handleSubmit} className="card space-y-5" noValidate>
          <div className="grid sm:grid-cols-2 gap-5">
            <div>
              <Label required>Название точки</Label>
              <Input
                placeholder="Пекарня на Ленина"
                value={form.name}
                onChange={set('name')}
                error={errors.name}
              />
            </div>
            <div>
              <Label required>Категория</Label>
              <Select
                value={form.category_code}
                onChange={set('category_code')}
                error={errors.category_code}
              >
                <option value="">Выберите категорию</option>
                {BUSINESS_CATEGORIES.map((c) => (
                  <option key={c.code} value={c.code}>{c.label}</option>
                ))}
              </Select>
            </div>
          </div>

          <div>
            <Label required>Адрес</Label>
            <AddressAutocomplete
              value={form.address}
              onChange={handleAddressSelect}
              placeholder="Начните вводить адрес..."
              error={errors.address}
            />
          </div>

          <div>
            <Label>Телефон <span className="text-cream-400 font-normal">(необязательно)</span></Label>
            <div className="relative">
              <Phone size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-cream-400 pointer-events-none" />
              <Input
                type="tel"
                placeholder="+74951234567"
                value={form.phone}
                onChange={set('phone')}
                error={errors.phone}
                className="pl-9"
              />
            </div>
          </div>

          <div className="flex gap-3 pt-2">
            <Button type="submit" className="flex-1" disabled={mutation.isPending} size="lg">
              {mutation.isPending ? 'Создаём...' : (
                <><Plus size={16} /> Создать локацию</>
              )}
            </Button>
            <Button
              type="button"
              variant="secondary"
              onClick={() => navigate('/partner/dashboard')}
            >
              Отмена
            </Button>
          </div>
        </form>
      </main>
    </div>
  )
}
