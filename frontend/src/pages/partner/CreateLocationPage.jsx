import { useState } from 'react'
import { observer } from 'mobx-react-lite'
import { useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { Plus, Phone } from 'lucide-react'
import { BUSINESS_CATEGORIES } from '@/lib/constants'
import { getErrorMessage } from '@/lib/utils'
import PartnerLayout from '@/components/partner/layout/PartnerLayout'
import AddressAutocomplete from '@/components/AddressAutocomplete'
import Input from '@/components/ui/form/Input'
import Select from '@/components/ui/form/Select'
import Label from '@/components/ui/form/Label'
import Button from '@/components/ui/actions/Button'
import { useStores } from '@/context/StoresContext'

const INITIAL = {
  name: '',
  address: '',
  latitude: null,
  longitude: null,
  category_code: '',
  phone: '',
}

function CreateLocationPageBase() {
  const [form, setForm] = useState(INITIAL)
  const [errors, setErrors] = useState({})
  const navigate = useNavigate()
  const { locationsStore } = useStores()

  const setField = (field) => (e) => setForm((f) => ({ ...f, [field]: e.target.value }))

  const validate = () => {
    const nextErrors = {}
    if (!form.name.trim()) nextErrors.name = 'Введите название'
    if (!form.category_code) nextErrors.category_code = 'Выберите категорию'
    if (!form.address || !form.latitude) nextErrors.address = 'Выберите адрес из списка'
    return nextErrors
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    const nextErrors = validate()
    if (Object.keys(nextErrors).length) {
      setErrors(nextErrors)
      return
    }
    setErrors({})

    const payload = {
      name: form.name,
      address: form.address,
      latitude: form.latitude,
      longitude: form.longitude,
      category_code: form.category_code,
    }
    if (form.phone.trim()) payload.phone = form.phone

    try {
      await locationsStore.create(payload)
      toast.success('Локация создана')
      navigate('/partner/locations')
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
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
    <PartnerLayout title="Добавить локацию" subtitle="Создайте новую точку продаж">
      <div className="max-w-3xl">
        <form onSubmit={handleSubmit} className="card space-y-5" noValidate>
          <div className="grid sm:grid-cols-2 gap-5">
            <div>
              <Label required>Название точки</Label>
              <Input value={form.name} onChange={setField('name')} error={errors.name} />
            </div>
            <div>
              <Label required>Категория</Label>
              <Select value={form.category_code} onChange={setField('category_code')} error={errors.category_code}>
                <option value="">Выберите категорию</option>
                {BUSINESS_CATEGORIES.map((category) => (
                  <option key={category.code} value={category.code}>{category.label}</option>
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
            <Label>Телефон</Label>
            <div className="relative">
              <Phone size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-cream-400" />
              <Input type="tel" value={form.phone} onChange={setField('phone')} className="pl-9" />
            </div>
          </div>

          <div className="flex gap-3 pt-1">
            <Button type="submit" className="flex-1" disabled={locationsStore.submitting}>
              {locationsStore.submitting ? 'Создаем...' : (<><Plus size={16} /> Создать локацию</>)}
            </Button>
            <Button type="button" variant="secondary" onClick={() => navigate('/partner/locations')}>
              Отмена
            </Button>
          </div>
        </form>
      </div>
    </PartnerLayout>
  )
}

export default observer(CreateLocationPageBase)
