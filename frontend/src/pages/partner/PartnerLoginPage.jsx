import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { Leaf, Eye, EyeOff, LogIn } from 'lucide-react'
import { useAuth } from '@/context/AuthContext'
import { getErrorMessage } from '@/lib/utils'
import Input from '@/components/ui/Input'
import Button from '@/components/ui/Button'
import Label from '@/components/ui/Label'

export default function PartnerLoginPage() {
  const [form, setForm] = useState({ email: '', password: '' })
  const [showPw, setShowPw] = useState(false)
  const [errors, setErrors] = useState({})
  const [loading, setLoading] = useState(false)
  const { login } = useAuth()
  const navigate = useNavigate()

  const set = (field) => (e) => setForm((f) => ({ ...f, [field]: e.target.value }))

  const validate = () => {
    const e = {}
    if (!form.email.trim()) e.email = 'Введите email'
    if (!form.password) e.password = 'Введите пароль'
    return e
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    const errs = validate()
    if (Object.keys(errs).length) { setErrors(errs); return }
    setErrors({})
    setLoading(true)
    try {
      const data = await login(form.email, form.password)
      if (data.must_change_password) {
        navigate('/partner/change-password', { replace: true })
      } else {
        navigate('/partner/dashboard', { replace: true })
      }
    } catch (err) {
      toast.error(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-cream-50 flex items-center justify-center p-4">
      <div className="w-full max-w-sm">
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-14 h-14 rounded-2xl bg-brand-500 mb-4">
            <Leaf size={28} className="text-white" />
          </div>
          <h1 className="text-2xl font-bold text-brand-900">Вход для партнёров</h1>
          <p className="text-sm text-brand-500 mt-2">Войдите в личный кабинет партнёра</p>
        </div>

        <form onSubmit={handleSubmit} className="card space-y-4" noValidate>
          <div>
            <Label required>Email</Label>
            <Input
              type="email"
              placeholder="owner@bakery.ru"
              value={form.email}
              onChange={set('email')}
              error={errors.email}
              autoComplete="email"
            />
          </div>

          <div>
            <Label required>Пароль</Label>
            <div className="relative">
              <Input
                type={showPw ? 'text' : 'password'}
                placeholder="Введите пароль"
                value={form.password}
                onChange={set('password')}
                error={errors.password}
                autoComplete="current-password"
                className="pr-10"
              />
              <button
                type="button"
                onClick={() => setShowPw((v) => !v)}
                className="absolute right-3 top-2.5 text-cream-400 hover:text-brand-500"
              >
                {showPw ? <EyeOff size={16} /> : <Eye size={16} />}
              </button>
            </div>
          </div>

          <Button type="submit" className="w-full" disabled={loading} size="lg">
            {loading ? 'Входим...' : (
              <><LogIn size={16} /> Войти</>
            )}
          </Button>
        </form>

        <p className="text-center text-sm text-brand-500 mt-6">
          Нет аккаунта?{' '}
          <a href="/#apply" className="text-brand-500 font-medium hover:underline">
            Подать заявку
          </a>
        </p>
      </div>
    </div>
  )
}
