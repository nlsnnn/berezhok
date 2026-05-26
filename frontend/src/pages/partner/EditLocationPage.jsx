import { useEffect, useState } from 'react'
import { observer } from 'mobx-react-lite'
import { useNavigate, useParams } from 'react-router-dom'
import { toast } from 'sonner'
import { Phone, Save } from 'lucide-react'
import { getErrorMessage } from '@/lib/utils'
import PartnerLayout from '@/components/partner/layout/PartnerLayout'
import Input from '@/components/ui/form/Input'
import Label from '@/components/ui/form/Label'
import Button from '@/components/ui/actions/Button'
import Spinner from '@/components/ui/feedback/Spinner'
import ImageUpload from '@/components/partner/forms/ImageUpload'
import WorkingHoursEditor from '@/components/partner/forms/WorkingHoursEditor'
import { useStores } from '@/context/StoresContext'

const INITIAL = {
  phone: '',
  logo_url: '',
  cover_image_url: '',
  working_hours: {},
}

function EditLocationPageBase() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { locationsStore } = useStores()

  const [form, setForm] = useState(INITIAL)
  const [errors, setErrors] = useState({})

  useEffect(() => {
    locationsStore
      .loadById(id)
      .then((location) => {
        if (!location) return
        setForm({
          phone: location.phone || '',
          logo_url: location.logo_url || '',
          cover_image_url: location.cover_image_url || '',
          working_hours: location.working_hours || {},
        })
      })
      .catch((error) => toast.error(getErrorMessage(error)))
  }, [id, locationsStore])

  const setField = (field) => (value) => setForm((f) => ({ ...f, [field]: value }))

  const validate = () => {
    const next = {}
    if (form.phone && !/^\+[1-9]\d{1,14}$/.test(form.phone.trim())) {
      next.phone = 'Формат: +71234567890'
    }
    return next
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    const next = validate()
    if (Object.keys(next).length) {
      setErrors(next)
      return
    }
    setErrors({})

    const payload = {}
    if (form.phone.trim()) payload.phone = form.phone.trim()
    if (form.logo_url) payload.logo_url = form.logo_url
    if (form.cover_image_url) payload.cover_image_url = form.cover_image_url
    if (form.working_hours && Object.keys(form.working_hours).length > 0) {
      payload.working_hours = form.working_hours
    }

    try {
      await locationsStore.update(id, payload)
      toast.success('Настройки сохранены')
      navigate('/partner/locations')
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  if (locationsStore.loading && !locationsStore.current) {
    return (
      <PartnerLayout title="Настройка заведения">
        <div className="flex justify-center py-10">
          <Spinner size={28} />
        </div>
      </PartnerLayout>
    )
  }

  const location = locationsStore.current

  return (
    <PartnerLayout
      title="Настройка заведения"
      subtitle={location?.name || 'Загрузка...'}
    >
      <div className="max-w-3xl">
        <form onSubmit={handleSubmit} className="space-y-5" noValidate>
          <div className="card space-y-5">
            <h2 className="text-base font-semibold text-brand-900">Изображения</h2>
            <ImageUpload
              value={form.logo_url}
              onChange={setField('logo_url')}
              label="Логотип"
              hint="Квадратное изображение, до 5 МБ"
            />
            <ImageUpload
              value={form.cover_image_url}
              onChange={setField('cover_image_url')}
              label="Обложка"
              hint="Широкое изображение, до 5 МБ"
            />
          </div>

          <div className="card space-y-5">
            <h2 className="text-base font-semibold text-brand-900">Контакты</h2>
            <div>
              <Label>Телефон</Label>
              <div className="relative">
                <Phone size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-cream-400" />
                <Input
                  type="tel"
                  value={form.phone}
                  onChange={(e) => setForm((f) => ({ ...f, phone: e.target.value }))}
                  placeholder="+71234567890"
                  className="pl-9"
                  error={errors.phone}
                />
              </div>
            </div>
          </div>

          <div className="card space-y-3">
            <h2 className="text-base font-semibold text-brand-900">График работы</h2>
            <WorkingHoursEditor
              value={form.working_hours}
              onChange={setField('working_hours')}
              label=""
            />
          </div>

          <div className="flex gap-3">
            <Button type="submit" className="flex-1" disabled={locationsStore.submitting}>
              {locationsStore.submitting ? 'Сохраняем...' : (<><Save size={16} /> Сохранить</>)}
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

export default observer(EditLocationPageBase)
