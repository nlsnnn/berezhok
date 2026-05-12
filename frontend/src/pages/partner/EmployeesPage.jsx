import { useEffect, useState } from 'react'
import { observer } from 'mobx-react-lite'
import { Mail, Trash2, UserPlus, Users } from 'lucide-react'
import { toast } from 'sonner'
import PartnerLayout from '@/components/partner/layout/PartnerLayout'
import Input from '@/components/ui/form/Input'
import Select from '@/components/ui/form/Select'
import Label from '@/components/ui/form/Label'
import Button from '@/components/ui/actions/Button'
import Spinner from '@/components/ui/feedback/Spinner'
import { useStores } from '@/context/StoresContext'
import { formatDateTime, getErrorMessage } from '@/lib/utils'

const INITIAL_FORM = {
  location_id: '',
  name: '',
  email: '',
}

function EmployeesPageBase() {
  const { employeesStore } = useStores()
  const [form, setForm] = useState(INITIAL_FORM)
  const [errors, setErrors] = useState({})

  useEffect(() => {
    employeesStore.load().catch(() => {})
  }, [employeesStore])

  const setField = (field) => (event) => {
    setForm((current) => ({ ...current, [field]: event.target.value }))
  }

  const validate = () => {
    const nextErrors = {}

    if (!form.location_id) nextErrors.location_id = 'Выберите заведение'
    if (!form.name.trim()) nextErrors.name = 'Введите имя сотрудника'
    if (!form.email.trim()) nextErrors.email = 'Введите email'

    return nextErrors
  }

  const handleSubmit = async (event) => {
    event.preventDefault()
    const nextErrors = validate()
    if (Object.keys(nextErrors).length > 0) {
      setErrors(nextErrors)
      return
    }

    setErrors({})

    try {
      await employeesStore.create({
        location_id: form.location_id,
        name: form.name.trim(),
        email: form.email.trim(),
      })
      toast.success('Сотрудник добавлен. Временный пароль отправлен на email')
      setForm(INITIAL_FORM)
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const handleDelete = async (employee) => {
    const confirmed = window.confirm(`Удалить сотрудника ${employee.name || employee.email}?`)
    if (!confirmed) return

    try {
      await employeesStore.remove(employee.id)
      toast.success('Сотрудник удален')
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  return (
    <PartnerLayout
      title="Сотрудники"
      subtitle="Добавляйте сотрудников для выдачи заказов и управляйте доступом к вашим заведениям"
    >
      <div className="grid gap-6 xl:grid-cols-[420px,minmax(0,1fr)]">
        <section className="card">
          <div className="mb-5 flex items-center gap-3">
            <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-brand-100">
              <UserPlus size={22} className="text-brand-600" />
            </div>
            <div>
              <h2 className="text-lg font-semibold text-brand-900">Новый сотрудник</h2>
              <p className="text-sm text-brand-600">Сотрудник получит временный пароль на email и сможет работать только со своей локацией.</p>
            </div>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4" noValidate>
            <div>
              <Label required>Заведение</Label>
              <Select value={form.location_id} onChange={setField('location_id')} error={errors.location_id}>
                <option value="">Выберите локацию</option>
                {employeesStore.locations.map((location) => (
                  <option key={location.id} value={location.id}>
                    {location.name}
                  </option>
                ))}
              </Select>
            </div>

            <div>
              <Label required>Имя сотрудника</Label>
              <Input value={form.name} onChange={setField('name')} error={errors.name} placeholder="Например, Иван" />
            </div>

            <div>
              <Label required>Email</Label>
              <div className="relative">
                <Mail size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-cream-400" />
                <Input
                  type="email"
                  value={form.email}
                  onChange={setField('email')}
                  error={errors.email}
                  className="pl-9"
                  placeholder="employee@coffee.ru"
                />
              </div>
            </div>

            <Button type="submit" className="w-full" disabled={employeesStore.submitting || employeesStore.locations.length === 0}>
              {employeesStore.submitting ? 'Добавляем...' : (<><UserPlus size={16} /> Добавить сотрудника</>)}
            </Button>
          </form>
        </section>

        <section className="space-y-4">
          <div className="card flex items-center justify-between gap-4">
            <div className="flex items-center gap-3">
              <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-cream-100">
                <Users size={22} className="text-brand-700" />
              </div>
              <div>
                <h2 className="text-lg font-semibold text-brand-900">Команда</h2>
                <p className="text-sm text-brand-600">Всего аккаунтов: {employeesStore.items.length}</p>
              </div>
            </div>
          </div>

          {employeesStore.loading && (
            <div className="card flex justify-center py-16">
              <Spinner size={32} />
            </div>
          )}

          {!employeesStore.loading && employeesStore.items.length === 0 && (
            <div className="card flex flex-col items-center justify-center py-16 text-center">
              <div className="mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-brand-100">
                <Users size={24} className="text-brand-600" />
              </div>
              <h3 className="mb-2 text-lg font-semibold text-brand-900">Сотрудников пока нет</h3>
              <p className="max-w-sm text-sm text-brand-600">
                Создайте первый аккаунт для сотрудника, чтобы он мог просматривать активные заказы своей локации и выдавать их клиентам.
              </p>
            </div>
          )}

          {!employeesStore.loading && employeesStore.items.length > 0 && (
            <div className="space-y-4">
              {employeesStore.items.map((employee) => (
                <article key={employee.id} className="card space-y-4">
                  <div className="flex flex-wrap items-start justify-between gap-3">
                    <div>
                      <p className="text-lg font-semibold text-brand-900">{employee.name}</p>
                      <p className="text-sm text-brand-600">{employee.email}</p>
                    </div>
                    <div className="flex items-center gap-2">
                      <span className={`badge ${employee.role === 'owner' ? 'bg-brand-100 text-brand-800' : 'bg-cream-100 text-brand-700'}`}>
                        {employee.role === 'owner' ? 'Владелец' : 'Сотрудник'}
                      </span>
                      {employee.must_change_password && (
                        <span className="badge bg-yellow-100 text-yellow-800">Ожидает входа</span>
                      )}
                    </div>
                  </div>

                  <div className="grid gap-3 md:grid-cols-3">
                    <InfoRow label="Заведение" value={employee.location_name || '—'} />
                    <InfoRow label="Роль" value={employee.role === 'owner' ? 'Владелец' : 'Сотрудник'} />
                    <InfoRow label="Создан" value={formatDateTime(employee.created_at)} />
                  </div>

                  <div className="flex justify-end">
                    <Button
                      type="button"
                      variant="danger"
                      onClick={() => handleDelete(employee)}
                      disabled={employeesStore.deletingId === employee.id || employee.role === 'owner'}
                      className="gap-2"
                    >
                      <Trash2 size={16} />
                      {employeesStore.deletingId === employee.id ? 'Удаляем...' : 'Удалить'}
                    </Button>
                  </div>
                </article>
              ))}
            </div>
          )}
        </section>
      </div>
    </PartnerLayout>
  )
}

function InfoRow({ label, value }) {
  return (
    <div className="rounded-xl border border-cream-200 bg-cream-50 p-4">
      <p className="mb-1 text-xs uppercase tracking-wider text-cream-500">{label}</p>
      <p className="text-sm font-medium text-brand-800">{value || '—'}</p>
    </div>
  )
}

export default observer(EmployeesPageBase)
