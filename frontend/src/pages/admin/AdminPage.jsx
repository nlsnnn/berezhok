import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { listApplications, approveApplication, rejectApplication } from '@/api/applications'
import { BUSINESS_CATEGORIES, APPLICATION_STATUS } from '@/lib/constants'
import { formatDateTime, getErrorMessage } from '@/lib/utils'
import AdminNav from '@/components/AdminNav'
import Spinner from '@/components/ui/Spinner'
import StatusBadge from '@/components/ui/StatusBadge'
import Button from '@/components/ui/Button'
import Modal from '@/components/ui/Modal'
import { CheckCircle2, XCircle, Eye, RefreshCw, Phone, Mail, MapPin, Building2 } from 'lucide-react'

const STATUS_FILTERS = [
  { value: 'all', label: 'Все' },
  { value: 'pending', label: 'На рассмотрении' },
  { value: 'approved', label: 'Одобрены' },
  { value: 'rejected', label: 'Отклонены' },
]

function getCategoryLabel(code) {
  return BUSINESS_CATEGORIES.find((c) => c.code === code)?.label ?? code
}

function ApplicationDetailModal({ application, onClose, onApprove, onReject }) {
  const [rejectReason, setRejectReason] = useState('')
  const [showRejectForm, setShowRejectForm] = useState(false)

  if (!application) return null

  const canAct = application.status === 'pending'

  const handleReject = () => {
    if (!rejectReason.trim()) {
      toast.error('Укажите причину отказа')
      return
    }
    onReject(application.id, rejectReason)
  }

  return (
    <Modal open={!!application} onClose={onClose} title="Заявка на партнёрство" className="max-w-xl">
      <div className="space-y-5">
        <div className="flex items-center justify-between">
          <StatusBadge status={application.status} />
          <span className="text-xs text-cream-500">{formatDateTime(application.created_at)}</span>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <DetailRow icon={Building2} label="Заведение" value={application.business_name} />
          <DetailRow icon={null} label="Категория" value={getCategoryLabel(application.category_code)} />
          <DetailRow icon={Mail} label="Email" value={application.contact_email} />
          <DetailRow icon={Phone} label="Телефон" value={application.contact_phone} />
          <div className="col-span-2">
            <DetailRow icon={MapPin} label="Адрес" value={application.address} />
          </div>
        </div>

        {application.description && (
          <div>
            <p className="text-xs font-medium text-cream-500 uppercase tracking-wider mb-1.5">Описание</p>
            <p className="text-sm text-brand-700 bg-cream-50 rounded-xl p-3">{application.description}</p>
          </div>
        )}

        {application.rejection_reason && (
          <div className="bg-red-50 rounded-xl p-3">
            <p className="text-xs font-medium text-red-500 mb-1">Причина отказа</p>
            <p className="text-sm text-red-700">{application.rejection_reason}</p>
          </div>
        )}

        {canAct && !showRejectForm && (
          <div className="flex gap-3 pt-2">
            <Button variant="primary" className="flex-1" onClick={() => onApprove(application.id)}>
              <CheckCircle2 size={16} />
              Одобрить
            </Button>
            <Button variant="danger" className="flex-1" onClick={() => setShowRejectForm(true)}>
              <XCircle size={16} />
              Отклонить
            </Button>
          </div>
        )}

        {canAct && showRejectForm && (
          <div className="space-y-3 pt-2">
            <div>
              <label className="block text-sm font-medium text-brand-700 mb-1">Причина отказа</label>
              <textarea
                rows={3}
                className="input-base resize-none"
                placeholder="Укажите причину отказа..."
                value={rejectReason}
                onChange={(e) => setRejectReason(e.target.value)}
              />
            </div>
            <div className="flex gap-3">
              <Button variant="danger" className="flex-1" onClick={handleReject}>
                Подтвердить отказ
              </Button>
              <Button variant="secondary" onClick={() => setShowRejectForm(false)}>
                Отмена
              </Button>
            </div>
          </div>
        )}
      </div>
    </Modal>
  )
}

function DetailRow({ icon: Icon, label, value }) {
  return (
    <div>
      <p className="text-xs font-medium text-cream-500 uppercase tracking-wider mb-1">{label}</p>
      <p className="text-sm text-brand-800 font-medium flex items-center gap-1.5">
        {Icon && <Icon size={13} className="text-brand-400 shrink-0" />}
        {value || '—'}
      </p>
    </div>
  )
}

export default function AdminPage() {
  const [filter, setFilter] = useState('all')
  const [selected, setSelected] = useState(null)
  const qc = useQueryClient()

  const { data: applications = [], isLoading, isError, refetch } = useQuery({
    queryKey: ['admin', 'applications'],
    queryFn: listApplications,
  })

  const approveMutation = useMutation({
    mutationFn: (id) => approveApplication(id),
    onSuccess: () => {
      toast.success('Заявка одобрена. Партнёр создан, пароль отправлен на email.')
      qc.invalidateQueries(['admin', 'applications'])
      setSelected(null)
    },
    onError: (err) => toast.error(getErrorMessage(err)),
  })

  const rejectMutation = useMutation({
    mutationFn: ({ id, reason }) => rejectApplication(id, reason),
    onSuccess: () => {
      toast.success('Заявка отклонена.')
      qc.invalidateQueries(['admin', 'applications'])
      setSelected(null)
    },
    onError: (err) => toast.error(getErrorMessage(err)),
  })

  const handleApprove = (id) => approveMutation.mutate(id)
  const handleReject = (id, reason) => rejectMutation.mutate({ id, reason })

  const filtered = filter === 'all' ? applications : applications.filter((a) => a.status === filter)

  const counts = {
    all: applications.length,
    pending: applications.filter((a) => a.status === 'pending').length,
    approved: applications.filter((a) => a.status === 'approved').length,
    rejected: applications.filter((a) => a.status === 'rejected').length,
  }

  return (
    <div className="min-h-screen flex flex-col bg-cream-50">
      <AdminNav />

      <main className="flex-1 max-w-6xl mx-auto w-full px-4 sm:px-6 py-8">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-2xl font-bold text-brand-900">Заявки на партнёрство</h1>
            <p className="text-sm text-brand-500 mt-1">{counts.pending} ожидают рассмотрения</p>
          </div>
          <button onClick={() => refetch()} className="btn-ghost gap-1.5">
            <RefreshCw size={16} />
            Обновить
          </button>
        </div>

        {/* Filter tabs */}
        <div className="flex gap-2 flex-wrap mb-6">
          {STATUS_FILTERS.map((f) => (
            <button
              key={f.value}
              onClick={() => setFilter(f.value)}
              className={`px-4 py-2 rounded-xl text-sm font-medium transition-colors ${
                filter === f.value
                  ? 'bg-brand-500 text-white'
                  : 'bg-white border border-cream-200 text-brand-600 hover:border-brand-300'
              }`}
            >
              {f.label}
              <span className={`ml-1.5 text-xs ${filter === f.value ? 'opacity-80' : 'text-cream-400'}`}>
                {counts[f.value]}
              </span>
            </button>
          ))}
        </div>

        {/* Content */}
        {isLoading && (
          <div className="flex justify-center py-20">
            <Spinner size={32} />
          </div>
        )}

        {isError && (
          <div className="card text-center py-12 text-red-500">
            Не удалось загрузить заявки. <button onClick={() => refetch()} className="underline">Попробовать снова</button>
          </div>
        )}

        {!isLoading && !isError && filtered.length === 0 && (
          <div className="card text-center py-12 text-cream-500">
            {filter === 'all' ? 'Заявок пока нет' : 'Нет заявок в этом статусе'}
          </div>
        )}

        {!isLoading && !isError && filtered.length > 0 && (
          <div className="space-y-3">
            {filtered.map((app) => (
              <div
                key={app.id}
                className="card hover:shadow-md transition-shadow cursor-pointer"
                onClick={() => setSelected(app)}
              >
                <div className="flex items-start justify-between gap-4">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-3 mb-2">
                      <StatusBadge status={app.status} />
                      <span className="text-xs text-cream-400">{formatDateTime(app.created_at)}</span>
                    </div>
                    <h3 className="font-semibold text-brand-900 text-base truncate">{app.business_name}</h3>
                    <div className="flex flex-wrap gap-x-4 gap-y-1 mt-2 text-sm text-brand-500">
                      <span className="flex items-center gap-1">
                        <Building2 size={13} />
                        {getCategoryLabel(app.category_code)}
                      </span>
                      <span className="flex items-center gap-1">
                        <Mail size={13} />
                        {app.contact_email}
                      </span>
                      <span className="flex items-center gap-1">
                        <MapPin size={13} className="shrink-0" />
                        <span className="truncate max-w-xs">{app.address}</span>
                      </span>
                    </div>
                  </div>
                  <div className="flex items-center gap-2 shrink-0">
                    {app.status === 'pending' && (
                      <>
                        <button
                          className="btn-primary py-1.5 px-3 text-xs"
                          onClick={(e) => { e.stopPropagation(); handleApprove(app.id) }}
                          disabled={approveMutation.isPending}
                        >
                          <CheckCircle2 size={13} />
                          Одобрить
                        </button>
                        <button
                          className="btn-danger py-1.5 px-3 text-xs"
                          onClick={(e) => { e.stopPropagation(); setSelected(app) }}
                        >
                          <XCircle size={13} />
                          Отклонить
                        </button>
                      </>
                    )}
                    <button className="btn-ghost py-1.5 px-2" onClick={(e) => { e.stopPropagation(); setSelected(app) }}>
                      <Eye size={16} />
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </main>

      <ApplicationDetailModal
        application={selected}
        onClose={() => setSelected(null)}
        onApprove={handleApprove}
        onReject={handleReject}
      />
    </div>
  )
}
