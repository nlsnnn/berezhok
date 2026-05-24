import { useEffect } from 'react'
import { useState } from 'react'
import { observer } from 'mobx-react-lite'
import { toast } from 'sonner'
import { RefreshCw } from 'lucide-react'
import AdminBulkActions from '@/components/admin/AdminBulkActions'
import AdminDataState from '@/components/admin/AdminDataState'
import AdminLayout from '@/components/admin/layout/AdminLayout'
import AdminPagination from '@/components/admin/AdminPagination'
import Button from '@/components/ui/actions/Button'
import Input from '@/components/ui/form/Input'
import { useStores } from '@/context/StoresContext'
import { toggleAllPageIds, toggleSelectedId } from '@/lib/adminSelection'
import { formatDateTime, getErrorMessage } from '@/lib/utils'

function AdminAuditPageBase() {
  const { adminAuditStore } = useStores()
  const [selectedIds, setSelectedIds] = useState([])

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

  const getEventId = (event) => event.id || `${event.created_at}-${event.action}-${event.entity_id}`
  const pageIds = adminAuditStore.items.map(getEventId)
  const selectedPageIds = selectedIds.filter((id) => pageIds.includes(id))
  const allPageSelected = pageIds.length > 0 && pageIds.every((id) => selectedPageIds.includes(id))

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
            <AdminBulkActions
              selectedIds={selectedPageIds}
              pageIds={pageIds}
              allPageSelected={allPageSelected}
              onToggleAll={() => setSelectedIds((current) => toggleAllPageIds(current, pageIds))}
              onClear={() => setSelectedIds((current) => current.filter((id) => !pageIds.includes(id)))}
              loading={adminAuditStore.loading}
            />

            {adminAuditStore.items.map((event) => (
              <article key={getEventId(event)} className="card">
                <div className="flex flex-wrap items-start justify-between gap-4">
                  <label className="flex items-center pt-1">
                    <input
                      type="checkbox"
                      className="h-4 w-4 rounded border-cream-300 text-brand-600"
                      checked={selectedPageIds.includes(getEventId(event))}
                      onChange={() => setSelectedIds((current) => toggleSelectedId(current, getEventId(event)))}
                      aria-label="Выбрать событие аудита"
                    />
                  </label>
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
