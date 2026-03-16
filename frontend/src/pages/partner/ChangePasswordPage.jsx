import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { toast } from 'sonner'
import { KeyRound, Eye, EyeOff, CheckCircle2 } from 'lucide-react'
import { changePassword } from '@/api/partner'
import { useAuth } from '@/context/AuthContext'
import { getErrorMessage } from '@/lib/utils'
import Input from '@/components/ui/Input'
import Button from '@/components/ui/Button'
import Label from '@/components/ui/Label'
import PartnerNav from '@/components/PartnerNav'

export default function ChangePasswordPage() {
  const [form, setForm] = useState({ current_password: '', new_password: '', confirm: '' })
  const [show, setShow] = useState({ current: false, new_: false, confirm: false })
  const [errors, setErrors] = useState({})
  const { user, login: refreshUser } = useAuth()
  const navigate = useNavigate()

  // Update must_change_password flag in localStorage after success
  const handleSuccess = () => {
    const stored = JSON.parse(localStorage.getItem('partner_user') || '{}')
    stored.must_change_password = false
    localStorage.setItem('partner_user', JSON.stringify(stored))
    toast.success('Пароль успешно изменён')
    navigate('/partner/dashboard', { replace: true })
  }

  const mutation = useMutation({
    mutationFn: () => changePassword(form.current_password, form.new_password),
    onSuccess: handleSuccess,
    onError: (err) => toast.error(getErrorMessage(err)),
  })

  const set = (field) => (e) => setForm((f) => ({ ...f, [field]: e.target.value }))
  const toggle = (field) => () => setShow((s) => ({ ...s, [field]: !s[field] }))

  const validate = () => {
    const e = {}
    if (!form.current_password) e.current_password = 'Введите текущий пароль'
    if (!form.new_password) e.new_password = 'Введите новый пароль'
    else if (form.new_password.length < 8) e.new_password = 'Минимум 8 символов'
    if (form.new_password !== form.confirm) e.confirm = 'Пароли не совпадают'
    return e
  }

  const handleSubmit = (e) => {
    e.preventDefault()
    const errs = validate()
    if (Object.keys(errs).length) { setErrors(errs); return }
    setErrors({})
    mutation.mutate()
  }

  const isMustChange = user?.must_change_password

  return (
    <div className="min-h-screen flex flex-col bg-cream-50">
      <PartnerNav />
      <main className="flex-1 flex items-center justify-center p-4">
        <div className="w-full max-w-md">
          {isMustChange && (
            <div className="bg-yellow-50 border border-yellow-200 rounded-xl p-4 mb-6 text-sm text-yellow-800">
              <strong>Требуется смена пароля.</strong> Для продолжения работы установите новый пароль.
            </div>
          )}

          <div className="text-center mb-8">
            <div className="inline-flex items-center justify-center w-14 h-14 rounded-2xl bg-brand-100 mb-4">
              <KeyRound size={26} className="text-brand-500" />
            </div>
            <h1 className="text-2xl font-bold text-brand-900">
              {isMustChange ? 'Установите новый пароль' : 'Смена пароля'}
            </h1>
          </div>

          <form onSubmit={handleSubmit} className="card space-y-4" noValidate>
            <PasswordField
              label="Текущий пароль"
              value={form.current_password}
              onChange={set('current_password')}
              show={show.current}
              onToggle={toggle('current')}
              error={errors.current_password}
              autoComplete="current-password"
            />
            <PasswordField
              label="Новый пароль"
              value={form.new_password}
              onChange={set('new_password')}
              show={show.new_}
              onToggle={toggle('new_')}
              error={errors.new_password}
              autoComplete="new-password"
              hint="Не менее 8 символов"
            />
            <PasswordField
              label="Повторите новый пароль"
              value={form.confirm}
              onChange={set('confirm')}
              show={show.confirm}
              onToggle={toggle('confirm')}
              error={errors.confirm}
              autoComplete="new-password"
            />

            <Button type="submit" className="w-full" disabled={mutation.isPending} size="lg">
              {mutation.isPending ? 'Сохраняем...' : (
                <><CheckCircle2 size={16} /> Сохранить пароль</>
              )}
            </Button>
          </form>
        </div>
      </main>
    </div>
  )
}

function PasswordField({ label, value, onChange, show, onToggle, error, autoComplete, hint }) {
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
      {hint && !error && <p className="mt-1 text-xs text-cream-500">{hint}</p>}
    </div>
  )
}
