import { useState, useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { toast } from 'sonner'
import { CheckCircle2, Eye, EyeOff, Pencil, X, FileText, Building2, User } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { getPartnerProfile, updatePartnerProfile, changePassword } from '@/api/partner'
import { getErrorMessage } from '@/lib/utils'
import { useAuth } from '@/context/AuthContext'
import PartnerLayout from '@/components/partner/layout/PartnerLayout'
import Input from '@/components/ui/form/Input'
import Label from '@/components/ui/form/Label'
import Button from '@/components/ui/actions/Button'

function ProfilePageBase() {
  const { markPasswordChanged } = useAuth()
  const navigate = useNavigate()

  const [profile, setProfile] = useState(null)
  const [loadingProfile, setLoadingProfile] = useState(true)

  const [editingName, setEditingName] = useState(false)
  const [nameValue, setNameValue] = useState('')
  const [savingName, setSavingName] = useState(false)

  const [pwForm, setPwForm] = useState({ current_password: '', new_password: '', confirm: '' })
  const [pwShow, setPwShow] = useState({ current: false, next: false, confirm: false })
  const [pwErrors, setPwErrors] = useState({})
  const [savingPw, setSavingPw] = useState(false)

  useEffect(() => {
    getPartnerProfile()
      .then((data) => {
        setProfile(data)
        setNameValue(data.employee?.name || '')
      })
      .catch(() => toast.error('Не удалось загрузить профиль'))
      .finally(() => setLoadingProfile(false))
  }, [])

  const handleSaveName = async () => {
    if (!nameValue.trim() || nameValue.trim().length < 2) {
      toast.error('Имя должно содержать минимум 2 символа')
      return
    }
    setSavingName(true)
    try {
      await updatePartnerProfile({ name: nameValue.trim() })
      setProfile((prev) => ({ ...prev, employee: { ...prev.employee, name: nameValue.trim() } }))
      setEditingName(false)
      toast.success('Имя обновлено')
    } catch (error) {
      toast.error(getErrorMessage(error))
    } finally {
      setSavingName(false)
    }
  }

  const handleCancelName = () => {
    setNameValue(profile?.employee?.name || '')
    setEditingName(false)
  }

  const setPwField = (field) => (e) => setPwForm((prev) => ({ ...prev, [field]: e.target.value }))

  const validatePw = () => {
    const errs = {}
    if (!pwForm.current_password) errs.current_password = 'Введите текущий пароль'
    if (!pwForm.new_password) errs.new_password = 'Введите новый пароль'
    else if (pwForm.new_password.length < 8) errs.new_password = 'Минимум 8 символов'
    if (pwForm.confirm !== pwForm.new_password) errs.confirm = 'Пароли не совпадают'
    return errs
  }

  const handleChangePw = async (e) => {
    e.preventDefault()
    const errs = validatePw()
    if (Object.keys(errs).length) {
      setPwErrors(errs)
      return
    }
    setPwErrors({})
    setSavingPw(true)
    try {
      await changePassword(pwForm.current_password, pwForm.new_password)
      markPasswordChanged()
      setPwForm({ current_password: '', new_password: '', confirm: '' })
      toast.success('Пароль обновлён')
    } catch (error) {
      toast.error(getErrorMessage(error))
    } finally {
      setSavingPw(false)
    }
  }

  const roleLabel = { owner: 'Владелец', manager: 'Менеджер', employee: 'Сотрудник' }

  return (
    <PartnerLayout title="Профиль" subtitle="Настройки аккаунта и данные партнёра">
      <div className="grid gap-5 lg:grid-cols-2">
        {/* Profile info */}
        <div className="card space-y-5">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-brand-100 flex items-center justify-center shrink-0">
              <User size={18} className="text-brand-600" />
            </div>
            <h3 className="text-base font-semibold text-brand-900">Данные аккаунта</h3>
          </div>

          {loadingProfile ? (
            <div className="h-32 flex items-center justify-center text-sm text-brand-400">Загрузка...</div>
          ) : (
            <div className="space-y-4">
              <div>
                <Label>Имя</Label>
                {editingName ? (
                  <div className="flex gap-2 mt-1">
                    <Input
                      value={nameValue}
                      onChange={(e) => setNameValue(e.target.value)}
                      autoFocus
                      className="flex-1"
                    />
                    <button
                      onClick={handleSaveName}
                      disabled={savingName}
                      className="btn-primary px-3 py-2 text-sm rounded-xl"
                    >
                      <CheckCircle2 size={16} />
                    </button>
                    <button onClick={handleCancelName} className="btn-ghost px-3 py-2 rounded-xl">
                      <X size={16} />
                    </button>
                  </div>
                ) : (
                  <div className="flex items-center justify-between mt-1">
                    <span className="text-sm text-brand-800">{profile?.employee?.name || '—'}</span>
                    <button
                      onClick={() => setEditingName(true)}
                      className="btn-ghost p-1.5 text-brand-400 hover:text-brand-600"
                      title="Редактировать"
                    >
                      <Pencil size={14} />
                    </button>
                  </div>
                )}
              </div>

              <div>
                <Label>Email</Label>
                <p className="mt-1 text-sm text-brand-800">{profile?.employee?.email || '—'}</p>
              </div>

              <div>
                <Label>Роль</Label>
                <p className="mt-1 text-sm text-brand-800">
                  {roleLabel[profile?.employee?.role] || profile?.employee?.role || '—'}
                </p>
              </div>

              <div>
                <Label>Название бренда</Label>
                <p className="mt-1 text-sm text-brand-800">{profile?.partner?.brand_name || '—'}</p>
              </div>
            </div>
          )}
        </div>

        {/* Legal info */}
        <div className="card flex flex-col">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-10 h-10 rounded-xl bg-brand-100 flex items-center justify-center shrink-0">
              <Building2 size={18} className="text-brand-600" />
            </div>
            <h3 className="text-base font-semibold text-brand-900">Юридические данные</h3>
          </div>
          <p className="text-sm text-brand-600 flex-1">
            Заполните реквизиты компании, чтобы активировать боксы и начать принимать заказы без ограничений.
          </p>
          <div className="mt-5">
            <Button onClick={() => navigate('/partner/legal-info')}>
              <FileText size={15} />
              Перейти к заполнению
            </Button>
          </div>
        </div>

        {/* Change password */}
        <div className="card lg:col-span-2 max-w-xl">
          <div className="flex items-center gap-3 mb-5">
            <div className="w-10 h-10 rounded-xl bg-brand-100 flex items-center justify-center shrink-0">
              <Eye size={18} className="text-brand-600" />
            </div>
            <h3 className="text-base font-semibold text-brand-900">Смена пароля</h3>
          </div>

          <form onSubmit={handleChangePw} className="space-y-4" noValidate>
            <PasswordField
              label="Текущий пароль"
              value={pwForm.current_password}
              onChange={setPwField('current_password')}
              error={pwErrors.current_password}
              show={pwShow.current}
              onToggle={() => setPwShow((prev) => ({ ...prev, current: !prev.current }))}
              autoComplete="current-password"
            />
            <PasswordField
              label="Новый пароль"
              value={pwForm.new_password}
              onChange={setPwField('new_password')}
              error={pwErrors.new_password}
              show={pwShow.next}
              onToggle={() => setPwShow((prev) => ({ ...prev, next: !prev.next }))}
              autoComplete="new-password"
            />
            <PasswordField
              label="Повторите новый пароль"
              value={pwForm.confirm}
              onChange={setPwField('confirm')}
              error={pwErrors.confirm}
              show={pwShow.confirm}
              onToggle={() => setPwShow((prev) => ({ ...prev, confirm: !prev.confirm }))}
              autoComplete="new-password"
            />
            <Button type="submit" disabled={savingPw}>
              {savingPw ? 'Сохраняем...' : <><CheckCircle2 size={16} /> Сохранить пароль</>}
            </Button>
          </form>
        </div>
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
        <button
          type="button"
          onClick={onToggle}
          className="absolute right-3 top-2.5 text-cream-400 hover:text-brand-500"
        >
          {show ? <EyeOff size={16} /> : <Eye size={16} />}
        </button>
      </div>
    </div>
  )
}

export default observer(ProfilePageBase)
