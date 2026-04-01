import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import { toast } from 'sonner'
import { CheckCircle2, Eye, EyeOff, KeyRound } from 'lucide-react'
import { changePassword } from '@/api/partner'
import { getErrorMessage } from '@/lib/utils'
import { useAuth } from '@/context/AuthContext'
import PartnerLayout from '@/components/partner/layout/PartnerLayout'
import Input from '@/components/ui/form/Input'
import Label from '@/components/ui/form/Label'
import Button from '@/components/ui/actions/Button'

function ChangePasswordPageBase() {
  const [form, setForm] = useState({ current_password: '', new_password: '', confirm: '' })
  const [show, setShow] = useState({ current: false, next: false, confirm: false })
  const [errors, setErrors] = useState({})
  const [loading, setLoading] = useState(false)
  const { user, markPasswordChanged } = useAuth()
  const navigate = useNavigate()

  const setField = (field) => (e) => setForm((prev) => ({ ...prev, [field]: e.target.value }))

  const validate = () => {
    const nextErrors = {}
    if (!form.current_password) nextErrors.current_password = 'Введите текущий пароль'
    if (!form.new_password) nextErrors.new_password = 'Введите новый пароль'
    else if (form.new_password.length < 8) nextErrors.new_password = 'Минимум 8 символов'
    if (form.confirm !== form.new_password) nextErrors.confirm = 'Пароли не совпадают'
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
    setLoading(true)
    try {
      await changePassword(form.current_password, form.new_password)
      markPasswordChanged()
      toast.success('Пароль обновлен')
      navigate('/partner/dashboard', { replace: true })
    } catch (error) {
      toast.error(getErrorMessage(error))
    } finally {
      setLoading(false)
    }
  }

  return (
    <PartnerLayout
      title={user?.must_change_password ? 'Установите новый пароль' : 'Смена пароля'}
      subtitle="Пароль от партнерского кабинета"
    >
      <div className="max-w-xl">
        {user?.must_change_password && (
          <div className="mb-5 rounded-xl bg-yellow-50 border border-yellow-200 p-3 text-sm text-yellow-800">
            Для продолжения работы необходимо сменить временный пароль.
          </div>
        )}

        <form onSubmit={handleSubmit} className="card space-y-4" noValidate>
          <PasswordField
            label="Текущий пароль"
            value={form.current_password}
            onChange={setField('current_password')}
            error={errors.current_password}
            show={show.current}
            onToggle={() => setShow((prev) => ({ ...prev, current: !prev.current }))}
            autoComplete="current-password"
          />
          <PasswordField
            label="Новый пароль"
            value={form.new_password}
            onChange={setField('new_password')}
            error={errors.new_password}
            show={show.next}
            onToggle={() => setShow((prev) => ({ ...prev, next: !prev.next }))}
            autoComplete="new-password"
          />
          <PasswordField
            label="Повторите новый пароль"
            value={form.confirm}
            onChange={setField('confirm')}
            error={errors.confirm}
            show={show.confirm}
            onToggle={() => setShow((prev) => ({ ...prev, confirm: !prev.confirm }))}
            autoComplete="new-password"
          />

          <Button type="submit" disabled={loading} className="w-full">
            {loading ? 'Сохраняем...' : (<><CheckCircle2 size={16} /> Сохранить пароль</>)}
          </Button>
        </form>
      </div>
    </PartnerLayout>
  )
}

function PasswordField({ label, value, onChange, error, show, onToggle, autoComplete }) {
  return (
    <div>
      <Label required>{label}</Label>
      <div className="relative">
        <Input
          type={show ? 'text' : 'password'}
          value={value}
          onChange={onChange}
          error={error}
          autoComplete={autoComplete}
          className="pr-10"
        />
        <button type="button" onClick={onToggle} className="absolute right-3 top-2.5 text-cream-400 hover:text-brand-500">
          {show ? <EyeOff size={16} /> : <Eye size={16} />}
        </button>
      </div>
    </div>
  )
}

export default observer(ChangePasswordPageBase)
