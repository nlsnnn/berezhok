import { useState } from 'react'
import { observer } from 'mobx-react-lite'
import { toast } from 'sonner'
import { CheckCircle2, Send } from 'lucide-react'
import { BUSINESS_CATEGORIES } from '@/lib/constants'
import { getErrorMessage } from '@/lib/utils'
import AddressAutocomplete from '@/components/AddressAutocomplete'
import Input from '@/components/ui/form/Input'
import Select from '@/components/ui/form/Select'
import Label from '@/components/ui/form/Label'
import Button from '@/components/ui/actions/Button'
import { useStores } from '@/context/StoresContext'

const INITIAL = {
  contact_name: '',
  contact_email: '',
  contact_phone: '',
  business_name: '',
  category_code: '',
  address: '',
  latitude: null,
  longitude: null,
  description: '',
}

function ApplicationFormBase() {
  const [form, setForm] = useState(INITIAL)
  const [errors, setErrors] = useState({})
  const [submitted, setSubmitted] = useState(false)
  const { applicationStore } = useStores()

  const setField = (field) => (e) => setForm((prev) => ({ ...prev, [field]: e.target.value }))

  const validate = () => {
    const nextErrors = {}
    if (!form.contact_name.trim()) nextErrors.contact_name = 'Введите имя'
    if (!form.contact_email.trim()) nextErrors.contact_email = 'Введите email'
    else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.contact_email)) nextErrors.contact_email = 'Некорректный email'
    if (!form.contact_phone.trim()) nextErrors.contact_phone = 'Введите телефон'
    else if (!/^\+7\d{10}$/.test(form.contact_phone)) nextErrors.contact_phone = 'Формат: +7XXXXXXXXXX'
    if (!form.business_name.trim()) nextErrors.business_name = 'Введите название'
    if (!form.category_code) nextErrors.category_code = 'Выберите категорию'
    if (!form.address) nextErrors.address = 'Выберите адрес из списка'
    if (!form.latitude || !form.longitude) nextErrors.address = 'Выберите адрес из выпадающего списка'
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
    try {
      await applicationStore.create(form)
      setSubmitted(true)
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const handleAddressSelect = (suggestion) => {
    if (!suggestion) {
      setForm((prev) => ({ ...prev, address: '', latitude: null, longitude: null }))
      return
    }
    setForm((prev) => ({
      ...prev,
      address: suggestion.display_name,
      latitude: suggestion.latitude,
      longitude: suggestion.longitude,
    }))
    setErrors((prev) => ({ ...prev, address: undefined }))
  }

  if (submitted) {
    return (
      <div className="card text-center py-14">
        <div className="w-16 h-16 rounded-full bg-brand-100 flex items-center justify-center mx-auto mb-5">
          <CheckCircle2 size={32} className="text-brand-500" />
        </div>
        <h3 className="text-2xl font-semibold text-brand-900 mb-3">Заявка отправлена!</h3>
        <p className="text-brand-600 max-w-sm mx-auto">
          Мы получили вашу заявку и рассмотрим ее в течение 1-2 рабочих дней.
          Ожидайте ответа на почту <strong>{form.contact_email}</strong>.
        </p>
        <button className="mt-8 btn-secondary" onClick={() => { setSubmitted(false); setForm(INITIAL) }}>
          Подать еще одну заявку
        </button>
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit} className="card space-y-5" noValidate>
      <div className="grid sm:grid-cols-2 gap-5">
        <div>
          <Label required>Ваше имя</Label>
          <Input value={form.contact_name} onChange={setField('contact_name')} error={errors.contact_name} />
        </div>
        <div>
          <Label required>Телефон</Label>
          <Input type="tel" value={form.contact_phone} onChange={setField('contact_phone')} error={errors.contact_phone} />
        </div>
      </div>

      <div>
        <Label required>Email</Label>
        <Input type="email" value={form.contact_email} onChange={setField('contact_email')} error={errors.contact_email} />
      </div>

      <div className="grid sm:grid-cols-2 gap-5">
        <div>
          <Label required>Название заведения</Label>
          <Input value={form.business_name} onChange={setField('business_name')} error={errors.business_name} />
        </div>
        <div>
          <Label required>Категория</Label>
          <Select value={form.category_code} onChange={setField('category_code')} error={errors.category_code}>
            <option value="">Выберите категорию</option>
            {BUSINESS_CATEGORIES.map((c) => <option key={c.code} value={c.code}>{c.label}</option>)}
          </Select>
        </div>
      </div>

      <div>
        <Label required>Адрес</Label>
        <AddressAutocomplete value={form.address} onChange={handleAddressSelect} error={errors.address} />
      </div>

      <div>
        <Label>Описание</Label>
        <textarea
          rows={3}
          value={form.description}
          onChange={setField('description')}
          className="input-base resize-none"
          placeholder="Расскажите о заведении"
        />
      </div>

      <Button type="submit" size="lg" className="w-full" disabled={applicationStore.submitting}>
        {applicationStore.submitting ? 'Отправляем...' : (<><Send size={16} /> Отправить заявку</>)}
      </Button>
    </form>
  )
}

export default observer(ApplicationFormBase)
