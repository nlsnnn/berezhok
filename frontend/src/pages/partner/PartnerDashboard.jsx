import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { getPartnerProfile } from '@/api/partner'
import { PARTNER_STATUS, BUSINESS_CATEGORIES } from '@/lib/constants'
import { formatDate, cn } from '@/lib/utils'
import PartnerNav from '@/components/PartnerNav'
import Spinner from '@/components/ui/Spinner'
import Button from '@/components/ui/Button'
import StatusBadge from '@/components/ui/StatusBadge'
import { Building2, MapPin, Mail, User, KeyRound, Plus, AlertTriangle } from 'lucide-react'

export default function PartnerDashboard() {
  const { data, isLoading, isError, refetch } = useQuery({
    queryKey: ['partner', 'profile'],
    queryFn: getPartnerProfile,
  })

  return (
    <div className="min-h-screen flex flex-col bg-cream-50">
      <PartnerNav />

      <main className="flex-1 max-w-5xl mx-auto w-full px-4 sm:px-6 py-8">
        {isLoading && (
          <div className="flex justify-center py-20">
            <Spinner size={32} />
          </div>
        )}

        {isError && (
          <div className="card text-center py-12 text-red-500">
            Не удалось загрузить профиль.{' '}
            <button onClick={() => refetch()} className="underline">Попробовать снова</button>
          </div>
        )}

        {data && <DashboardContent data={data} />}
      </main>
    </div>
  )
}

function DashboardContent({ data }) {
  const { partner, employee, location } = data

  const partnerStatus = PARTNER_STATUS[partner.status]
  const commissionPct = Math.round((partner.commission_rate ?? 0.2) * 100)
  const isPromo = !!partner.promo_until && new Date(partner.promo_until) > new Date()

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-start justify-between flex-wrap gap-4">
        <div>
          <h1 className="text-2xl font-bold text-brand-900">{partner.brand_name}</h1>
          <div className="flex items-center gap-3 mt-2">
            <span className={cn('badge', partnerStatus?.color ?? 'bg-gray-100 text-gray-700')}>
              {partnerStatus?.label ?? partner.status}
            </span>
            <span className="text-sm text-brand-500">
              Партнёр с {formatDate(partner.created_at)}
            </span>
          </div>
        </div>
        <Link to="/partner/change-password" className="btn-secondary gap-2 text-sm">
          <KeyRound size={15} />
          Сменить пароль
        </Link>
      </div>

      {/* Status alert */}
      {partner.status === 'pending_documents' && (
        <div className="bg-yellow-50 border border-yellow-200 rounded-xl p-4 flex gap-3">
          <AlertTriangle size={18} className="text-yellow-500 shrink-0 mt-0.5" />
          <div className="text-sm text-yellow-800">
            <strong>Ожидает проверки документов.</strong> Наш менеджер скоро свяжется с вами для завершения регистрации.
          </div>
        </div>
      )}

      <div className="grid md:grid-cols-2 gap-6">
        {/* Partner card */}
        <div className="card space-y-4">
          <h2 className="text-base font-semibold text-brand-900 flex items-center gap-2">
            <Building2 size={18} className="text-brand-400" />
            Партнёр
          </h2>
          <InfoRow label="Название бренда" value={partner.brand_name} />
          <InfoRow label="Комиссия" value={
            <span className="flex items-center gap-2">
              {commissionPct}%
              {isPromo && (
                <span className="badge bg-brand-100 text-brand-700">Промо до {formatDate(partner.promo_until)}</span>
              )}
            </span>
          } />
        </div>

        {/* Employee card */}
        <div className="card space-y-4">
          <h2 className="text-base font-semibold text-brand-900 flex items-center gap-2">
            <User size={18} className="text-brand-400" />
            Сотрудник
          </h2>
          <InfoRow label="Имя" value={employee.name || '—'} />
          <InfoRow label="Email" value={employee.email} icon={Mail} />
          <InfoRow label="Роль" value={ROLE_LABELS[employee.role] ?? employee.role} />
        </div>
      </div>

      {/* Location card */}
      {location ? (
        <div className="card space-y-4">
          <h2 className="text-base font-semibold text-brand-900 flex items-center gap-2">
            <MapPin size={18} className="text-brand-400" />
            Локация
          </h2>
          <div className="grid sm:grid-cols-2 gap-4">
            <InfoRow label="Название" value={location.name} />
            <InfoRow label="Адрес" value={location.address} />
          </div>
        </div>
      ) : (
        <div className="card border-dashed border-cream-300 text-center py-10">
          <MapPin size={32} className="text-cream-400 mx-auto mb-3" />
          <h3 className="font-semibold text-brand-700 mb-2">Нет активной локации</h3>
          <p className="text-sm text-brand-500 mb-5">Добавьте точку продаж для начала работы</p>
          <Link to="/partner/locations/new" className="btn-primary">
            <Plus size={16} />
            Добавить локацию
          </Link>
        </div>
      )}

      {/* Quick actions */}
      <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <QuickAction
          icon={Plus}
          title="Новая локация"
          desc="Добавить точку продаж"
          to="/partner/locations/new"
        />
        <QuickAction
          icon={KeyRound}
          title="Сменить пароль"
          desc="Обновить пароль входа"
          to="/partner/change-password"
        />
      </div>
    </div>
  )
}

const ROLE_LABELS = {
  owner: 'Владелец',
  manager: 'Менеджер',
  employee: 'Сотрудник',
}

function InfoRow({ label, value, icon: Icon }) {
  return (
    <div>
      <p className="text-xs text-cream-500 uppercase tracking-wider font-medium mb-1">{label}</p>
      <p className="text-sm font-medium text-brand-800 flex items-center gap-1.5">
        {Icon && <Icon size={13} className="text-brand-400" />}
        {value ?? '—'}
      </p>
    </div>
  )
}

function QuickAction({ icon: Icon, title, desc, to }) {
  return (
    <Link to={to} className="card hover:shadow-md transition-shadow flex items-center gap-4 group">
      <div className="w-11 h-11 rounded-xl bg-brand-100 flex items-center justify-center shrink-0 group-hover:bg-brand-500 transition-colors">
        <Icon size={20} className="text-brand-500 group-hover:text-white transition-colors" />
      </div>
      <div>
        <p className="font-semibold text-brand-800 text-sm">{title}</p>
        <p className="text-xs text-brand-500">{desc}</p>
      </div>
    </Link>
  )
}
