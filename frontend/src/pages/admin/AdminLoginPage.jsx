import { useState } from 'react'
import { Navigate, useLocation, useNavigate } from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import { ShieldCheck } from 'lucide-react'
import { toast } from 'sonner'
import Button from '@/components/ui/actions/Button'
import Input from '@/components/ui/form/Input'
import { useStores } from '@/context/StoresContext'
import { getErrorMessage } from '@/lib/utils'

function AdminLoginPageBase() {
  const { adminAuthStore } = useStores()
  const navigate = useNavigate()
  const location = useLocation()
  const [email, setEmail] = useState('admin@berezhok.local')
  const [password, setPassword] = useState('test12345')

  if (adminAuthStore.isAuthenticated) {
    return <Navigate to="/admin/applications" replace />
  }

  const handleSubmit = async (event) => {
    event.preventDefault()
    try {
      await adminAuthStore.login(email, password)
      const nextPath = location.state?.from?.pathname || '/admin/applications'
      navigate(nextPath, { replace: true })
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  return (
    <main className="min-h-screen bg-gradient-to-br from-cream-50 via-white to-brand-50/60 px-4 py-10">
      <section className="mx-auto flex min-h-[calc(100svh-5rem)] w-full max-w-md items-center">
        <form className="card w-full space-y-6" onSubmit={handleSubmit}>
          <div className="flex items-center gap-3">
            <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-brand-100">
              <ShieldCheck size={24} className="text-brand-700" />
            </div>
            <div>
              <h1 className="text-2xl font-bold text-brand-900">Вход администратора</h1>
              <p className="text-sm text-brand-600">Панель управления Бережок</p>
            </div>
          </div>

          <label className="block space-y-2">
            <span className="text-sm font-medium text-brand-700">Email</span>
            <Input
              type="email"
              autoComplete="email"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              required
            />
          </label>

          <label className="block space-y-2">
            <span className="text-sm font-medium text-brand-700">Пароль</span>
            <Input
              type="password"
              autoComplete="current-password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              required
            />
          </label>

          <Button type="submit" className="w-full" disabled={adminAuthStore.loading}>
            {adminAuthStore.loading ? 'Входим...' : 'Войти'}
          </Button>
        </form>
      </section>
    </main>
  )
}

export default observer(AdminLoginPageBase)
