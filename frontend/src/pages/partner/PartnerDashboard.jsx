import { useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { Link } from 'react-router-dom'
import { AlertTriangle, CalendarDays, CheckCircle2, Clock3, Coins, Package, Plus, Star, Store } from 'lucide-react'
import PartnerLayout from '@/components/partner/layout/PartnerLayout'
import Spinner from '@/components/ui/feedback/Spinner'
import Button from '@/components/ui/actions/Button'
import { formatDate } from '@/lib/utils'
import { useStores } from '@/context/StoresContext'

function PartnerDashboardBase() {
  const { dashboardStore } = useStores()

  useEffect(() => {
    dashboardStore.load()
  }, [dashboardStore])

  const data = dashboardStore.data

  return (
    <PartnerLayout
      title="Дашборд"
      subtitle="Ключевые метрики по локациям, заказам и выплатам"
      actions={
        <Link to="/partner/boxes/new">
          <Button className="gap-2">
            <Plus size={16} />
            Создать бокс
          </Button>
        </Link>
      }
    >
      {dashboardStore.loading && (
        <div className="flex justify-center py-24">
          <Spinner size={34} />
        </div>
      )}

      {dashboardStore.error && (
        <div className="card text-center py-12 text-red-600">
          Не удалось загрузить дашборд.
          <button className="underline ml-2" onClick={() => dashboardStore.load()}>Попробовать снова</button>
        </div>
      )}

      {data && (
        <div className="space-y-6">
          <section className="grid md:grid-cols-4 gap-4">
            <StatCard icon={Clock3} label="Ожидают подтверждения" value={data?.today?.pending_confirmation ?? 0} />
            <StatCard icon={CheckCircle2} label="Подтверждены сегодня" value={data?.today?.confirmed ?? 0} />
            <StatCard icon={Coins} label="Выручка за неделю" value={`${data?.week?.gross_revenue ?? 0} ₽`} />
            <StatCard icon={Star} label="Рейтинг" value={data?.week?.avg_rating ?? '0.0'} />
          </section>

          <section className="grid xl:grid-cols-3 gap-5">
            <article className="card xl:col-span-2">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-lg font-semibold text-brand-900">Локации</h2>
                <Link to="/partner/locations" className="text-sm text-brand-600 hover:text-brand-800">Все локации</Link>
              </div>
              <div className="space-y-3">
                {(data.locations || []).map((location) => (
                  <div key={location.id} className="rounded-xl border border-cream-200 p-4 bg-white">
                    <div className="flex items-center justify-between gap-3">
                      <div>
                        <p className="font-semibold text-brand-900">{location.name}</p>
                        <p className="text-sm text-brand-600 mt-1">{location.address}</p>
                      </div>
                      <span className="badge bg-brand-100 text-brand-700">{location.active_boxes_count ?? 0} активных боксов</span>
                    </div>
                  </div>
                ))}
              </div>
            </article>

            <article className="card">
              <h2 className="text-lg font-semibold text-brand-900 mb-4">Финансы</h2>
              <div className="space-y-3 text-sm">
                <InfoRow icon={Store} label="Бренд" value={data?.partner?.brand_name} />
                <InfoRow icon={Package} label="Завершено за неделю" value={data?.week?.orders_completed ?? 0} />
                <InfoRow icon={Coins} label="Ожидает выплаты" value={`${data?.finance?.balance_pending ?? 0} ₽`} />
                <InfoRow icon={CalendarDays} label="Следующая выплата" value={formatDate(data?.finance?.next_payout_date)} />
              </div>
            </article>
          </section>

          {data?.partner?.status === 'pending_documents' && (
            <div className="rounded-xl border border-yellow-300 bg-yellow-50 p-4 text-yellow-900 text-sm flex items-start gap-3">
              <AlertTriangle size={18} className="shrink-0 mt-0.5" />
              <p>Профиль ожидает проверку документов. После верификации все функции будут доступны без ограничений.</p>
            </div>
          )}
        </div>
      )}
    </PartnerLayout>
  )
}

function StatCard({ icon: Icon, label, value }) {
  return (
    <article className="bg-white rounded-2xl border border-cream-200 p-4 shadow-sm">
      <div className="w-10 h-10 rounded-xl bg-brand-100 flex items-center justify-center mb-3">
        <Icon size={18} className="text-brand-600" />
      </div>
      <p className="text-sm text-brand-600">{label}</p>
      <p className="text-2xl font-bold text-brand-900 mt-1">{value}</p>
    </article>
  )
}

function InfoRow({ icon: Icon, label, value }) {
  return (
    <div className="rounded-xl border border-cream-200 px-3 py-2.5 bg-cream-50">
      <p className="text-xs uppercase tracking-wider text-cream-500 mb-1">{label}</p>
      <p className="text-brand-800 font-medium flex items-center gap-2">
        <Icon size={14} className="text-brand-500" />
        {value || '—'}
      </p>
    </div>
  )
}

export default observer(PartnerDashboardBase)
