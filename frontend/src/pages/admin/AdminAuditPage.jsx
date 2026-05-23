import { useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { toast } from 'sonner'
import { RefreshCw } from 'lucide-react'
import AdminDataState from '@/components/admin/AdminDataState'
import AdminLayout from '@/components/admin/layout/AdminLayout'
import AdminPagination from '@/components/admin/AdminPagination'
import Button from '@/components/ui/actions/Button'
import Input from '@/components/ui/form/Input'
import { useStores } from '@/context/StoresContext'
import { formatDateTime, getErrorMessage } from '@/lib/utils'

function AdminAuditPageBase() {
  const { adminAuditStore } = useStores()

  useEffect(() => {
    adminAuditStore.load().catch(() => {})
  }, [adminAuditStore])

  const handleApply = (event) => {
    event.preventDefault()
    adminAuditStore.load().catch((error) => toast.error(getErrorMessage(error)))
  }

  const handleReset = () => {
    adminAuditStore.resetFilters()
    adminAuditStore.load().catch((error) => toast.error(getErrorMessage(error)))
  }

  return (
    <AdminLayout
      title="Аудит"
      subtitle="История административных действий"
      actions={
        <Button variant="secondary" onClick={() => adminAuditStore.load()} disabled={adminAuditStore.loading}>
          <RefreshCw size={16} />
          Обновить
        </Button>
      }
    >
      <div className="space-y-6">
        <form className="card grid gap-4 md:grid-cols-[1fr_1fr_auto_auto]" onSubmit={handleApply}>
          <label className="space-y-2">
            <span className="text-sm font-medium text-brand-700">Действие</span>
            <Input value={adminAuditStore.filters.action} onChange={(event) => adminAuditStore.setFilter('action', event.target.value)} placeholder="approve, update, delete" />
          </label>
          <label className="space-y-2">
            <span className="text-sm font-medium text-brand-700">Сущность</span>
            <Input value={adminAuditStore.filters.entity_type} onChange={(event) => adminAuditStore.setFilter('entity_type', event.target.value)} placeholder="partner, box, order" />
          </label>
          <div className="flex items-end">
            <Button type="submit" disabled={adminAuditStore.loading}>Применить</Button>
          </div>
          <div className="flex items-end">
            <Button type="button" variant="secondary" onClick={handleReset} disabled={adminAuditStore.loading}>Сбросить</Button>
          </div>
        </form>

        <AdminDataState
          loading={adminAuditStore.loading}
          error={adminAuditStore.error}
          empty={adminAuditStore.items.length === 0}
          emptyText="Событий аудита по текущим фильтрам нет"
          onRetry={() => adminAuditStore.load()}
        >
          <section className="space-y-4">
            {adminAuditStore.items.map((event) => (
              <article key={event.id || `${event.created_at}-${event.action}`} className="card">
                <div className="flex flex-wrap items-start justify-between gap-4">
                  <div>
                    <h2 className="text-base font-semibold text-brand-900">{event.action || 'Действие'}</h2>
                    <p className="mt-1 text-sm text-brand-600">{event.entity_type || 'Сущность'} · {event.entity_id || '—'}</p>
                    <p className="mt-2 text-sm text-brand-700">{event.admin_email || event.admin_id || 'Администратор не указан'}</p>
                  </div>
                  <span className="text-sm text-brand-600">{formatDateTime(event.created_at)}</span>
                </div>
              </article>
            ))}
            <AdminPagination pagination={adminAuditStore.pagination} loading={adminAuditStore.loading} onPrev={() => adminAuditStore.prevPage()} onNext={() => adminAuditStore.nextPage()} />
          </section>
        </AdminDataState>
      </div>
    </AdminLayout>
  )
}

export default observer(AdminAuditPageBase)
