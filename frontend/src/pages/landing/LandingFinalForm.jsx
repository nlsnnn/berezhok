import { useState } from 'react'
import { observer } from 'mobx-react-lite'
import { toast } from 'sonner'
import { useStores } from '@/context/StoresContext'
import { getErrorMessage } from '@/lib/utils'
import { BUSINESS_CATEGORIES } from '@/lib/constants'
import { IconArrow } from './LandingStickers'
import { HedgehogMoney } from './LandingMascot'
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

function LandingFinalFormBase() {
  const [form, setForm] = useState(INITIAL)
  const [errors, setErrors] = useState({})
  const [submitted, setSubmitted] = useState(false)
  const { applicationStore } = useStores()

  const set = (field) => (e) => setForm((prev) => ({ ...prev, [field]: e.target.value }))

  const validate = () => {
    const e = {}
    if (!form.contact_name.trim()) e.contact_name = 'Введите имя'
    if (!form.contact_email.trim()) e.contact_email = 'Введите email'
    else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.contact_email)) e.contact_email = 'Некорректный email'
    if (!form.contact_phone.trim()) e.contact_phone = 'Введите телефон'
    else if (!/^\+7\d{10}$/.test(form.contact_phone)) e.contact_phone = 'Формат: +7XXXXXXXXXX'
    if (!form.business_name.trim()) e.business_name = 'Введите название'
    if (!form.category_code) e.category_code = 'Выберите тип'
    if (!form.address) e.address = 'Выберите адрес из списка'
    if (!form.latitude || !form.longitude) e.address = 'Выберите адрес из выпадающего списка'
    return e
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    const errs = validate()
    if (Object.keys(errs).length) { setErrors(errs); return }
    setErrors({})
    try {
      await applicationStore.create(form)
      setSubmitted(true)
    } catch (err) {
      toast.error(getErrorMessage(err))
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
      <div className="form-success">
        <div style={{ margin: '0 auto 16px', width: 80, height: 80 }}>
          <HedgehogMoney />
        </div>
        <h3>Заявка отправлена!</h3>
        <p>Менеджер позвонит в течение пары часов в рабочее время. Готовьте остатки 🥐</p>
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit} noValidate>
      <div className="form-row">
        <div className="form-field">
          <label>Ваше имя</label>
          <input type="text" placeholder="Анна" required value={form.contact_name} onChange={set('contact_name')} />
          {errors.contact_name && <span style={{ fontSize: 12, color: 'var(--c-warn, #c53030)' }}>{errors.contact_name}</span>}
        </div>
        <div className="form-field">
          <label>Телефон</label>
          <input type="tel" placeholder="+7XXXXXXXXXX" required value={form.contact_phone} onChange={set('contact_phone')} />
          {errors.contact_phone && <span style={{ fontSize: 12, color: 'var(--c-warn, #c53030)' }}>{errors.contact_phone}</span>}
        </div>
      </div>

      <div className="form-row">
        <div className="form-field full">
          <label>Email</label>
          <input type="email" placeholder="email@example.com" required value={form.contact_email} onChange={set('contact_email')} />
          {errors.contact_email && <span style={{ fontSize: 12, color: 'var(--c-warn, #c53030)' }}>{errors.contact_email}</span>}
        </div>
      </div>

      <div className="form-row">
        <div className="form-field">
          <label>Название заведения</label>
          <input type="text" placeholder='Пекарня «Корица»' required value={form.business_name} onChange={set('business_name')} />
          {errors.business_name && <span style={{ fontSize: 12, color: 'var(--c-warn, #c53030)' }}>{errors.business_name}</span>}
        </div>
        <div className="form-field">
          <label>Тип заведения</label>
          <select required value={form.category_code} onChange={set('category_code')}>
            <option value="">Выберите тип</option>
            {BUSINESS_CATEGORIES.map((c) => (
              <option key={c.code} value={c.code}>{c.label}</option>
            ))}
          </select>
          {errors.category_code && <span style={{ fontSize: 12, color: 'var(--c-warn, #c53030)' }}>{errors.category_code}</span>}
        </div>
      </div>

      <div className="form-row">
        <div className="form-field full">
          <label>Адрес</label>
          <AddressAutocomplete
            value={form.address}
            onChange={handleAddressSelect}
            error={errors.address}
          />
        </div>
      </div>

      <div className="form-row">
        <div className="form-field full">
          <label>Сколько примерно остатков в день?</label>
          <input type="text" placeholder="2–3 коробки выпечки, литр супа..." value={form.description} onChange={set('description')} />
        </div>
      </div>

      <button
        type="submit"
        className="btn btn-accent btn-lg form-submit"
        disabled={applicationStore.submitting}
      >
        {applicationStore.submitting ? 'Отправляем...' : (
          <>Отправить заявку <IconArrow size={18} /></>
        )}
      </button>

      <p className="form-disclaimer">
        Нажимая, соглашаетесь с обработкой персональных данных. Без спама, обещаем.
      </p>
    </form>
  )
}

export default observer(LandingFinalFormBase)
