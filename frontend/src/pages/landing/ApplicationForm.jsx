import { useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { toast } from 'sonner'
import { CheckCircle2, Send } from 'lucide-react'
import { createApplication } from '@/api/applications'
import { BUSINESS_CATEGORIES } from '@/lib/constants'
import { getErrorMessage } from '@/lib/utils'
import Input from '@/components/ui/Input'
import Select from '@/components/ui/Select'
import Label from '@/components/ui/Label'
import Button from '@/components/ui/Button'
import AddressAutocomplete from '@/components/AddressAutocomplete'

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

export default function ApplicationForm() {
  const [form, setForm] = useState(INITIAL)
  const [errors, setErrors] = useState({})
  const [submitted, setSubmitted] = useState(false)

  const mutation = useMutation({
    mutationFn: createApplication,
    onSuccess: () => {
      setSubmitted(true)
    },
    onError: (err) => {
      toast.error(getErrorMessage(err))
    },
  })

  const set = (field) => (e) =>
    setForm((f) => ({ ...f, [field]: e.target.value }))

  const validate = () => {
    const e = {}
    if (!form.contact_name.trim()) e.contact_name = 'Введите имя'
    if (!form.contact_email.trim()) e.contact_email = 'Введите email'
    else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.contact_email)) e.contact_email = 'Некорректный email'
    if (!form.contact_phone.trim()) e.contact_phone = 'Введите телефон'
    else if (!/^\+7\d{10}$/.test(form.contact_phone)) e.contact_phone = 'Формат: +7XXXXXXXXXX'
    if (!form.business_name.trim()) e.business_name = 'Введите название'
    if (!form.category_code) e.category_code = 'Выберите категорию'
    if (!form.address) e.address = 'Выберите адрес из списка'
    if (!form.latitude || !form.longitude) e.address = 'Выберите адрес из выпадающего списка'
    return e
  }

  const handleSubmit = (e) => {
    e.preventDefault()
    const e_ = validate()
    if (Object.keys(e_).length) {
      setErrors(e_)
      return
    }
    setErrors({})
    mutation.mutate(form)
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

  if (submitted) {
    return (
      <div className="card text-center py-14">
        <div className="w-16 h-16 rounded-full bg-brand-100 flex items-center justify-center mx-auto mb-5">
          <CheckCircle2 size={32} className="text-brand-500" />
        </div>
        <h3 className="text-2xl font-semibold text-brand-900 mb-3">Заявка отправлена!</h3>
        <p className="text-brand-600 max-w-sm mx-auto">
          Мы получили вашу заявку и рассмотрим её в течение 1–2 рабочих дней.
          Ожидайте ответа на почту <strong>{form.contact_email}</strong>.
        </p>
        <button
          className="mt-8 btn-secondary"
          onClick={() => { setSubmitted(false); setForm(INITIAL) }}
        >
          Подать ещё одну заявку
        </button>
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit} className="card space-y-5" noValidate>
      <div className="grid sm:grid-cols-2 gap-5">
        <div>
          <Label required>Ваше имя</Label>
          <Input
            placeholder="Иван Петров"
            value={form.contact_name}
            onChange={set('contact_name')}
            error={errors.contact_name}
          />
        </div>
        <div>
          <Label required>Телефон</Label>
          <Input
            type="tel"
            placeholder="+79001234567"
            value={form.contact_phone}
            onChange={set('contact_phone')}
            error={errors.contact_phone}
          />
        </div>
      </div>

      <div>
        <Label required>Email</Label>
        <Input
          type="email"
          placeholder="ivan@example.ru"
          value={form.contact_email}
          onChange={set('contact_email')}
          error={errors.contact_email}
        />
      </div>

      <div className="grid sm:grid-cols-2 gap-5">
        <div>
          <Label required>Название заведения</Label>
          <Input
            placeholder="Пекарня «Ромашка»"
            value={form.business_name}
            onChange={set('business_name')}
            error={errors.business_name}
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
          placeholder="Начните вводить адрес заведения..."
          error={errors.address}
        />
        <p className="mt-1 text-xs text-cream-500">Начните вводить адрес и выберите из предложенных вариантов</p>
      </div>

      <div>
        <Label>Описание <span className="text-cream-400 font-normal">(необязательно)</span></Label>
        <textarea
          rows={3}
          placeholder="Расскажите о вашем заведении, что планируете продавать..."
          value={form.description}
          onChange={set('description')}
          className="input-base resize-none"
        />
      </div>

      <Button
        type="submit"
        variant="primary"
        size="lg"
        className="w-full"
        disabled={mutation.isPending}
      >
        {mutation.isPending ? (
          <>Отправляем...</>
        ) : (
          <>
            <Send size={16} />
            Отправить заявку
          </>
        )}
      </Button>
    </form>
  )
}
