import { useEffect, useState } from 'react'
import { observer } from 'mobx-react-lite'
import { toast } from 'sonner'
import { Building2, CheckCircle2, Eye, Mail, MapPin, Phone, RefreshCw, Trash2, XCircle } from 'lucide-react'
import AdminBulkActions from '@/components/admin/AdminBulkActions'
import AdminDataState from '@/components/admin/AdminDataState'
import AdminDetailModal, { DetailGrid, DetailItem } from '@/components/admin/AdminDetailModal'
import AdminLayout from '@/components/admin/layout/AdminLayout'
import AdminPagination from '@/components/admin/AdminPagination'
import AdminToolbar from '@/components/admin/AdminToolbar'
import Button from '@/components/ui/actions/Button'
import StatusBadge from '@/components/ui/feedback/StatusBadge'
import { useStores } from '@/context/StoresContext'
import { BUSINESS_CATEGORIES } from '@/lib/constants'
import { canMutateOperations } from '@/lib/admin'
import { toggleAllPageIds, toggleSelectedId } from '@/lib/adminSelection'
import { formatDateTime, getErrorMessage } from '@/lib/utils'

const STATUS_FILTERS = [
  { value: 'all', label: 'Все' },
  { value: 'pending', label: 'На рассмотрении' },
  { value: 'approved', label: 'Одобрены' },
  { value: 'rejected', label: 'Отклонены' },
]

function getCategoryLabel(code) {
  return BUSINESS_CATEGORIES.find((category) => category.code === code)?.label ?? code
}

function AdminApplicationsPageBase() {
  const { adminApplicationsStore, adminAuthStore } = useStores()
  const [selectedId, setSelectedId] = useState(null)
  const [selectedIds, setSelectedIds] = useState([])
  const [reason, setReason] = useState('')
  const canAct = canMutateOperations(adminAuthStore.user)

  useEffect(() => {
    adminApplicationsStore.load().catch(() => {})
  }, [adminApplicationsStore])

  const pageIds = adminApplicationsStore.items.map((item) => item.id)
  const selectedPageIds = selectedIds.filter((id) => pageIds.includes(id))
  const allPageSelected = pageIds.length > 0 && pageIds.every((id) => selectedPageIds.includes(id))

  const openDetail = async (id) => {
    setSelectedId(id)
    setReason('')
    try {
      await adminApplicationsStore.loadDetail(id)
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const closeDetail = () => {
    setSelectedId(null)
    setReason('')
    adminApplicationsStore.clearCurrent()
  }

  const handleApplyFilters = (event) => {
    event.preventDefault()
    adminApplicationsStore.load().catch((error) => toast.error(getErrorMessage(error)))
  }

  const handleResetFilters = () => {
    adminApplicationsStore.resetFilters()
    adminApplicationsStore.load().catch((error) => toast.error(getErrorMessage(error)))
  }

  const handleApprove = async (id) => {
    try {
      await adminApplicationsStore.approve(id)
      toast.success('Заявка одобрена')
      closeDetail()
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const handleReject = async (id) => {
    if (!reason.trim()) {
      toast.error('Укажите причину отказа')
      return
    }

    try {
      await adminApplicationsStore.reject(id, reason.trim())
      toast.success('Заявка отклонена')
      closeDetail()
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const handleDelete = async (id) => {
    if (!window.confirm('Удалить заявку?')) return

    try {
      await adminApplicationsStore.delete(id)
      toast.success('Заявка удалена')
      closeDetail()
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const handleBulkApprove = async (ids) => {
    try {
      await adminApplicationsStore.approveMany(ids)
      toast.success('Выбранные заявки одобрены')
      setSelectedIds([])
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const handleBulkReject = async (ids) => {
    const bulkReason = window.prompt('Причина отказа для выбранных заявок')
    if (!bulkReason?.trim()) return

    try {
      await adminApplicationsStore.rejectMany(ids, bulkReason.trim())
      toast.success('Выбранные заявки отклонены')
      setSelectedIds([])
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const handleBulkDelete = async (ids) => {
    if (!window.confirm(`Удалить выбранные заявки: ${ids.length}?`)) return

    try {
      await adminApplicationsStore.deleteMany(ids)
      toast.success('Выбранные заявки удалены')
      setSelectedIds([])
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const current = adminApplicationsStore.current
  const canReviewCurrent = canAct && current?.status === 'pending'

  return (
    <AdminLayout
      title="Заявки на партнёрство"
      subtitle="Рассмотрение новых партнёров и история решений"
      actions={
        <Button variant="secondary" onClick={() => adminApplicationsStore.load()} disabled={adminApplicationsStore.loading}>
          <RefreshCw size={16} />
          Обновить
        </Button>
      }
    >
      <div className="space-y-6">
        <AdminToolbar
          search={adminApplicationsStore.filters.search}
          onSearchChange={(value) => adminApplicationsStore.setFilter('search', value)}
          status={adminApplicationsStore.filters.status}
          onStatusChange={(value) => adminApplicationsStore.setFilter('status', value)}
          statusOptions={STATUS_FILTERS}
          onApply={handleApplyFilters}
          onReset={handleResetFilters}
          loading={adminApplicationsStore.loading}
        />

        <AdminDataState
          loading={adminApplicationsStore.loading}
          error={adminApplicationsStore.error}
          empty={adminApplicationsStore.items.length === 0}
          emptyText="Заявок по текущим фильтрам нет"
          onRetry={() => adminApplicationsStore.load()}
        >
          <section className="space-y-4">
            <AdminBulkActions
              selectedIds={selectedPageIds}
              pageIds={pageIds}
              allPageSelected={allPageSelected}
              onToggleAll={() => setSelectedIds((current) => toggleAllPageIds(current, pageIds))}
              onClear={() => setSelectedIds((current) => current.filter((id) => !pageIds.includes(id)))}
              actions={
                canAct
                  ? [
                      { label: 'Одобрить', icon: <CheckCircle2 size={16} />, onClick: handleBulkApprove },
                      { label: 'Отклонить', icon: <XCircle size={16} />, variant: 'danger', onClick: handleBulkReject },
                      { label: 'Удалить', icon: <Trash2 size={16} />, variant: 'secondary', onClick: handleBulkDelete },
                    ]
                  : []
              }
              loading={adminApplicationsStore.loading || adminApplicationsStore.actionLoading}
            />

            {adminApplicationsStore.items.map((application) => (
              <article key={application.id} className="card cursor-pointer transition-shadow hover:shadow-md" onClick={() => openDetail(application.id)}>
                <div className="flex flex-wrap items-start justify-between gap-4">
                  <label className="flex items-center pt-1" onClick={(event) => event.stopPropagation()}>
                    <input
                      type="checkbox"
                      className="h-4 w-4 rounded border-cream-300 text-brand-600"
                      checked={selectedPageIds.includes(application.id)}
                      onChange={() => setSelectedIds((current) => toggleSelectedId(current, application.id))}
                      aria-label="Выбрать заявку"
                    />
                  </label>
                  <div className="min-w-0 flex-1">
                    <div className="mb-2 flex flex-wrap items-center gap-3">
                      <StatusBadge status={application.status} />
                      <span className="text-xs text-cream-500">{formatDateTime(application.created_at)}</span>
                    </div>
                    <h2 className="truncate text-base font-semibold text-brand-900">{application.business_name}</h2>
                    <div className="mt-2 flex flex-wrap gap-x-4 gap-y-1 text-sm text-brand-600">
                      <span className="flex items-center gap-1">
                        <Building2 size={13} /> {getCategoryLabel(application.category_code)}
                      </span>
                      <span className="flex items-center gap-1">
                        <Mail size={13} /> {application.contact_email || '—'}
                      </span>
                      <span className="flex items-center gap-1">
                        <MapPin size={13} /> {application.address || '—'}
                      </span>
                    </div>
                  </div>
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={(event) => {
                      event.stopPropagation()
                      openDetail(application.id)
                    }}
                  >
                    <Eye size={16} />
                    Открыть
                  </Button>
                </div>
              </article>
            ))}

            <AdminPagination
              pagination={adminApplicationsStore.pagination}
              loading={adminApplicationsStore.loading}
              onPrev={() => adminApplicationsStore.prevPage()}
              onNext={() => adminApplicationsStore.nextPage()}
            />
          </section>
        </AdminDataState>
      </div>

      <AdminDetailModal
        open={Boolean(selectedId)}
        title="Заявка на партнёрство"
        loading={adminApplicationsStore.detailLoading}
        onClose={closeDetail}
        actions={
          current && (
            <>
              {canReviewCurrent && (
                <>
                  <Button onClick={() => handleApprove(current.id)} disabled={adminApplicationsStore.actionLoading}>
                    <CheckCircle2 size={16} />
                    Одобрить
                  </Button>
                  <Button variant="danger" onClick={() => handleReject(current.id)} disabled={adminApplicationsStore.actionLoading}>
                    <XCircle size={16} />
                    Отклонить
                  </Button>
                </>
              )}
              {canAct && (
                <Button variant="secondary" onClick={() => handleDelete(current.id)} disabled={adminApplicationsStore.actionLoading}>
                  <Trash2 size={16} />
                  Удалить
                </Button>
              )}
            </>
          )
        }
      >
        {current && (
          <>
            <div className="flex flex-wrap items-center justify-between gap-3">
              <StatusBadge status={current.status} />
              <span className="text-sm text-brand-600">{formatDateTime(current.created_at)}</span>
            </div>
            <DetailGrid>
              <DetailItem label="Заведение" value={current.business_name} />
              <DetailItem label="Категория" value={getCategoryLabel(current.category_code)} />
              <DetailItem label="Email" value={current.contact_email} />
              <DetailItem label="Телефон" value={current.contact_phone} />
              <DetailItem label="Адрес" value={current.address} wide />
              <DetailItem label="Причина отказа" value={current.rejection_reason} wide />
            </DetailGrid>
            {canReviewCurrent && (
              <label className="block space-y-2">
                <span className="text-sm font-medium text-brand-700">Причина отказа</span>
                <textarea
                  className="input-base min-h-24 resize-none"
                  value={reason}
                  onChange={(event) => setReason(event.target.value)}
                  placeholder="Укажите причину перед отклонением"
                />
              </label>
            )}
          </>
        )}
      </AdminDetailModal>
    </AdminLayout>
  )
}

export default observer(AdminApplicationsPageBase)
