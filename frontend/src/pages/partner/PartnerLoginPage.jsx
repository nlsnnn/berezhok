import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { Eye, EyeOff, LogIn } from 'lucide-react'
import { useAuth } from '@/context/AuthContext'
import { getErrorMessage } from '@/lib/utils'
import Input from '@/components/ui/form/Input'
import Label from '@/components/ui/form/Label'
import Button from '@/components/ui/actions/Button'

export default function PartnerLoginPage() {
  const [form, setForm] = useState({ email: '', password: '' })
  const [showPassword, setShowPassword] = useState(false)
  const [errors, setErrors] = useState({})
  const [loading, setLoading] = useState(false)
  const { login } = useAuth()
  const navigate = useNavigate()

  const setField = (field) => (e) => setForm((prev) => ({ ...prev, [field]: e.target.value }))

  const validate = () => {
    const nextErrors = {}
    if (!form.email.trim()) nextErrors.email = 'Введите email'
    if (!form.password) nextErrors.password = 'Введите пароль'
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
      const data = await login(form.email, form.password)
      if (data.must_change_password) {
        navigate('/partner/change-password', { replace: true })
      } else {
        navigate('/partner/dashboard', { replace: true })
      }
    } catch (error) {
      toast.error(getErrorMessage(error))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-cream-50 via-white to-brand-50/70 flex items-center justify-center px-4 py-10">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <img src="/logo.png" alt="Бережок" className="w-16 h-16 rounded-2xl object-cover mx-auto mb-4 shadow-sm" />
          <h1 className="text-3xl font-bold text-brand-900">Вход для партнеров</h1>
          <p className="text-sm text-brand-600 mt-2">Управляйте локациями, боксами и выдачей заказов</p>
        </div>

        <form onSubmit={handleSubmit} className="card space-y-4" noValidate>
          <div>
            <Label required>Email</Label>
            <Input
              type="email"
              placeholder="owner@bakery.ru"
              value={form.email}
              onChange={setField('email')}
              error={errors.email}
              autoComplete="email"
            />
          </div>

          <div>
            <Label required>Пароль</Label>
            <div className="relative">
              <Input
                type={showPassword ? 'text' : 'password'}
                placeholder="Введите пароль"
                value={form.password}
                onChange={setField('password')}
                error={errors.password}
                autoComplete="current-password"
                className="pr-10"
              />
              <button
                type="button"
                onClick={() => setShowPassword((prev) => !prev)}
                className="absolute right-3 top-2.5 text-cream-400 hover:text-brand-500"
              >
                {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
              </button>
            </div>
          </div>

          <Button type="submit" className="w-full" disabled={loading} size="lg">
            {loading ? 'Входим...' : (<><LogIn size={16} /> Войти</>)}
          </Button>
        </form>
      </div>
    </div>
  )
}
