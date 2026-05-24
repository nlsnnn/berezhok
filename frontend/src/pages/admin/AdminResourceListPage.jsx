import { useEffect, useState } from 'react'
import { observer } from 'mobx-react-lite'
import { toast } from 'sonner'
import { Eye, RefreshCw } from 'lucide-react'
import AdminBulkActions from '@/components/admin/AdminBulkActions'
import AdminDataState from '@/components/admin/AdminDataState'
import AdminDetailModal, { DetailGrid, DetailItem } from '@/components/admin/AdminDetailModal'
import AdminLayout from '@/components/admin/layout/AdminLayout'
import AdminPagination from '@/components/admin/AdminPagination'
import AdminToolbar from '@/components/admin/AdminToolbar'
import Button from '@/components/ui/actions/Button'
import StatusBadge from '@/components/ui/feedback/StatusBadge'
import { getStatusMeta } from '@/lib/adminResources'
import { toggleAllPageIds, toggleSelectedId } from '@/lib/adminSelection'
import { formatDateTime, getErrorMessage } from '@/lib/utils'

function AdminResourceListPageBase({
  store,
  title,
  subtitle,
  emptyText,
  statusOptions = [],
  statusMap = {},
  extraFilters = [],
  getTitle,
  getDescription,
  getStatus,
  getMeta = () => [],
  detailFields = [],
  detailTitle = title,
  renderDetailExtra,
  renderActions,
  bulkActions = [],
  search = true,
}) {
  const [selectedId, setSelectedId] = useState(null)
  const [selectedIds, setSelectedIds] = useState([])

  useEffect(() => {
    store.load().catch(() => {})
  }, [store])

  const pageIds = store.items.map((item) => item.id)
  const selectedPageIds = selectedIds.filter((id) => pageIds.includes(id))
  const allPageSelected = pageIds.length > 0 && pageIds.every((id) => selectedPageIds.includes(id))

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
            extraFilters={extraFilters.map((filter) => ({
              ...filter,
              value: store.filters[filter.name] || filter.fallback || 'all',
              onChange: (value) => store.setFilter(filter.name, value),
            }))}
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
            <AdminBulkActions
              selectedIds={selectedPageIds}
              pageIds={pageIds}
              allPageSelected={allPageSelected}
              onToggleAll={() => setSelectedIds((current) => toggleAllPageIds(current, pageIds))}
              onClear={() => setSelectedIds((current) => current.filter((id) => !pageIds.includes(id)))}
              actions={bulkActions.map((action) => ({
                ...action,
                onClick: async (ids) => {
                  await action.onClick(ids)
                  setSelectedIds([])
                },
              }))}
              loading={store.loading || store.actionLoading}
            />

            {store.items.map((item) => {
              const status = getStatus?.(item)
              const statusMeta = getStatusMeta(status, statusMap)

              return (
                <article key={item.id} className="card cursor-pointer transition-shadow hover:shadow-md" onClick={() => openDetail(item.id)}>
                  <div className="flex flex-wrap items-start justify-between gap-4">
                    <label className="flex items-center pt-1" onClick={(event) => event.stopPropagation()}>
                      <input
                        type="checkbox"
                        className="h-4 w-4 rounded border-cream-300 text-brand-600"
                        checked={selectedPageIds.includes(item.id)}
                        onChange={() => setSelectedIds((current) => toggleSelectedId(current, item.id))}
                        aria-label="Выбрать строку"
                      />
                    </label>
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
