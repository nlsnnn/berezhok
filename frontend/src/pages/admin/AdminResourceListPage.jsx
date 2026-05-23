import { useEffect, useState } from 'react'
import { observer } from 'mobx-react-lite'
import { toast } from 'sonner'
import { Eye, RefreshCw } from 'lucide-react'
import AdminDataState from '@/components/admin/AdminDataState'
import AdminDetailModal, { DetailGrid, DetailItem } from '@/components/admin/AdminDetailModal'
import AdminLayout from '@/components/admin/layout/AdminLayout'
import AdminPagination from '@/components/admin/AdminPagination'
import AdminToolbar from '@/components/admin/AdminToolbar'
import Button from '@/components/ui/actions/Button'
import StatusBadge from '@/components/ui/feedback/StatusBadge'
import { getStatusMeta } from '@/lib/adminResources'
import { formatDateTime, getErrorMessage } from '@/lib/utils'

function AdminResourceListPageBase({
  store,
  title,
  subtitle,
  emptyText,
  statusOptions = [],
  statusMap = {},
  getTitle,
  getDescription,
  getStatus,
  getMeta = () => [],
  detailFields = [],
  detailTitle = title,
  renderDetailExtra,
  renderActions,
  search = true,
}) {
  const [selectedId, setSelectedId] = useState(null)

  useEffect(() => {
    store.load().catch(() => {})
  }, [store])

  const openDetail = async (id) => {
    setSelectedId(id)
    try {
      await store.loadDetail(id)
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const closeDetail = () => {
    setSelectedId(null)
    store.clearCurrent()
  }

  const handleApplyFilters = (event) => {
    event.preventDefault()
    store.load().catch((error) => toast.error(getErrorMessage(error)))
  }

  const handleResetFilters = () => {
    store.resetFilters()
    store.load().catch((error) => toast.error(getErrorMessage(error)))
  }

  return (
    <AdminLayout
      title={title}
      subtitle={subtitle}
      actions={
        <Button variant="secondary" onClick={() => store.load()} disabled={store.loading}>
          <RefreshCw size={16} />
          Обновить
        </Button>
      }
    >
      <div className="space-y-6">
        {(search || statusOptions.length > 0) && (
          <AdminToolbar
            search={store.filters.search || ''}
            onSearchChange={(value) => store.setFilter('search', value)}
            status={store.filters.status || 'all'}
            onStatusChange={(value) => store.setFilter('status', value)}
            statusOptions={statusOptions}
            onApply={handleApplyFilters}
            onReset={handleResetFilters}
            loading={store.loading}
          />
        )}

        <AdminDataState
          loading={store.loading}
          error={store.error}
          empty={store.items.length === 0}
          emptyText={emptyText}
          onRetry={() => store.load()}
        >
          <section className="space-y-4">
            {store.items.map((item) => {
              const status = getStatus?.(item)
              const statusMeta = getStatusMeta(status, statusMap)

              return (
                <article key={item.id} className="card cursor-pointer transition-shadow hover:shadow-md" onClick={() => openDetail(item.id)}>
                  <div className="flex flex-wrap items-start justify-between gap-4">
                    <div className="min-w-0 flex-1">
                      <div className="mb-2 flex flex-wrap items-center gap-3">
                        {status && <StatusBadge status={status} customLabel={statusMeta.label} customColor={statusMeta.color || statusMeta.className} />}
                        {item.created_at && <span className="text-xs text-cream-500">{formatDateTime(item.created_at)}</span>}
                      </div>
                      <h2 className="truncate text-base font-semibold text-brand-900">{getTitle(item)}</h2>
                      {getDescription?.(item) && <p className="mt-1 text-sm text-brand-600">{getDescription(item)}</p>}
                      <div className="mt-3 flex flex-wrap gap-x-4 gap-y-1 text-sm text-brand-600">
                        {getMeta(item).map((meta) => (
                          <span key={meta.label} className="flex items-center gap-1">
                            {meta.icon && <meta.icon size={13} />}
                            {meta.label}: {meta.value || '—'}
                          </span>
                        ))}
                      </div>
                    </div>
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={(event) => {
                        event.stopPropagation()
                        openDetail(item.id)
                      }}
                    >
                      <Eye size={16} />
                      Открыть
                    </Button>
                  </div>
                </article>
              )
            })}

            <AdminPagination pagination={store.pagination} loading={store.loading} onPrev={() => store.prevPage()} onNext={() => store.nextPage()} />
          </section>
        </AdminDataState>
      </div>

      <AdminDetailModal
        open={Boolean(selectedId)}
        title={detailTitle}
        loading={store.detailLoading}
        onClose={closeDetail}
        actions={store.current && renderActions?.({ item: store.current, close: closeDetail })}
      >
        {store.current && (
          <>
            <DetailGrid>
              {detailFields.map((field) => (
                <DetailItem key={field.label} label={field.label} value={field.value(store.current)} wide={field.wide} />
              ))}
            </DetailGrid>
            {renderDetailExtra?.(store.current)}
          </>
        )}
      </AdminDetailModal>
    </AdminLayout>
  )
}

export default observer(AdminResourceListPageBase)
