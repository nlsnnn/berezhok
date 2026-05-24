import { useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { BarChart3, RefreshCw } from 'lucide-react'
import AdminDataState from '@/components/admin/AdminDataState'
import AdminLayout from '@/components/admin/layout/AdminLayout'
import Button from '@/components/ui/actions/Button'
import { useStores } from '@/context/StoresContext'

function valueLabel(value) {
  if (value === null || value === undefined || value === '') return '—'
  if (typeof value === 'number') return value.toLocaleString('ru-RU')
  return String(value)
}

function AdminStatsPageBase() {
  const { adminStatsStore } = useStores()

  useEffect(() => {
    adminStatsStore.load().catch(() => {})
  }, [adminStatsStore])

  const data = adminStatsStore.data || {}
  const summary = data.summary || data
  const entries = Object.entries(summary).filter(([, value]) => typeof value !== 'object' || value === null)

  return (
    <AdminLayout
      title="Статистика"
      subtitle="Сводные показатели платформы"
      actions={
        <Button variant="secondary" onClick={() => adminStatsStore.load()} disabled={adminStatsStore.loading}>
          <RefreshCw size={16} />
          Обновить
        </Button>
      }
    >
      <AdminDataState
        loading={adminStatsStore.loading && !adminStatsStore.data}
        error={adminStatsStore.error && !adminStatsStore.data ? adminStatsStore.error : null}
        empty={entries.length === 0}
        emptyText="Статистика пока недоступна"
        onRetry={() => adminStatsStore.load()}
      >
        <section className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
          {entries.map(([key, value]) => (
            <article key={key} className="rounded-2xl border border-cream-200 bg-white p-4 shadow-sm">
              <div className="mb-3 flex h-10 w-10 items-center justify-center rounded-xl bg-brand-100">
                <BarChart3 size={18} className="text-brand-600" />
              </div>
              <p className="text-sm text-brand-600">{key.replaceAll('_', ' ')}</p>
              <p className="mt-1 text-2xl font-bold text-brand-950">{valueLabel(value)}</p>
            </article>
          ))}
        </section>
      </AdminDataState>
    </AdminLayout>
  )
}

export default observer(AdminStatsPageBase)
