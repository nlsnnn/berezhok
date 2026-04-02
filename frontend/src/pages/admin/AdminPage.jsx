import { useEffect, useState } from 'react'
import { observer } from 'mobx-react-lite'
import { toast } from 'sonner'
import { Building2, CheckCircle2, Eye, Mail, MapPin, Phone, RefreshCw, XCircle } from 'lucide-react'
import { BUSINESS_CATEGORIES } from '@/lib/constants'
import { formatDateTime, getErrorMessage } from '@/lib/utils'
import AdminNav from '@/components/AdminNav'
import Spinner from '@/components/ui/feedback/Spinner'
import StatusBadge from '@/components/ui/feedback/StatusBadge'
import Button from '@/components/ui/actions/Button'
import Modal from '@/components/ui/overlay/Modal'
import { useStores } from '@/context/StoresContext'

const STATUS_FILTERS = [
  { value: 'all', label: 'Все' },
  { value: 'pending', label: 'На рассмотрении' },
  { value: 'approved', label: 'Одобрены' },
  { value: 'rejected', label: 'Отклонены' },
]

function getCategoryLabel(code) {
  return BUSINESS_CATEGORIES.find((c) => c.code === code)?.label ?? code
}

function ApplicationDetailModal({ application, onClose, onApprove, onReject, loading }) {
  const [reason, setReason] = useState('')
  const [showRejectForm, setShowRejectForm] = useState(false)
  if (!application) return null

  const canAct = application.status === 'pending'

  const handleReject = () => {
    if (!reason.trim()) {
      toast.error('Укажите причину отказа')
      return
    }
    onReject(application.id, reason)
  }

  return (
    <Modal open={Boolean(application)} onClose={onClose} title="Заявка на партнерство" className="max-w-xl">
      <div className="space-y-5">
        <div className="flex items-center justify-between">
          <StatusBadge status={application.status} />
          <span className="text-xs text-cream-500">{formatDateTime(application.created_at)}</span>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <DetailRow icon={Building2} label="Заведение" value={application.business_name} />
          <DetailRow label="Категория" value={getCategoryLabel(application.category_code)} />
          <DetailRow icon={Mail} label="Email" value={application.contact_email} />
          <DetailRow icon={Phone} label="Телефон" value={application.contact_phone} />
          <div className="col-span-2">
            <DetailRow icon={MapPin} label="Адрес" value={application.address} />
          </div>
        </div>

        {canAct && !showRejectForm && (
          <div className="flex gap-3 pt-2">
            <Button variant="primary" className="flex-1" onClick={() => onApprove(application.id)} disabled={loading}>
              <CheckCircle2 size={16} /> Одобрить
            </Button>
            <Button variant="danger" className="flex-1" onClick={() => setShowRejectForm(true)}>
              <XCircle size={16} /> Отклонить
            </Button>
          </div>
        )}

        {canAct && showRejectForm && (
          <div className="space-y-3 pt-2">
            <textarea
              rows={3}
              className="input-base resize-none"
              placeholder="Укажите причину отказа..."
              value={reason}
              onChange={(e) => setReason(e.target.value)}
            />
            <div className="flex gap-3">
              <Button variant="danger" className="flex-1" onClick={handleReject} disabled={loading}>Подтвердить отказ</Button>
              <Button variant="secondary" onClick={() => setShowRejectForm(false)}>Отмена</Button>
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

function AdminPageBase() {
  const [filter, setFilter] = useState('all')
  const [selected, setSelected] = useState(null)
  const { adminApplicationsStore } = useStores()

  useEffect(() => {
    adminApplicationsStore.load()
  }, [adminApplicationsStore])

  const applications = adminApplicationsStore.items
  const filtered = filter === 'all' ? applications : applications.filter((item) => item.status === filter)

  const counts = {
    all: applications.length,
    pending: applications.filter((a) => a.status === 'pending').length,
    approved: applications.filter((a) => a.status === 'approved').length,
    rejected: applications.filter((a) => a.status === 'rejected').length,
  }

  const handleApprove = async (id) => {
    try {
      await adminApplicationsStore.approve(id)
      toast.success('Заявка одобрена')
      setSelected(null)
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const handleReject = async (id, reason) => {
    try {
      await adminApplicationsStore.reject(id, reason)
      toast.success('Заявка отклонена')
      setSelected(null)
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  return (
    <div className="min-h-screen flex flex-col bg-cream-50">
      <AdminNav />

      <main className="flex-1 max-w-6xl mx-auto w-full px-4 sm:px-6 py-8">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-2xl font-bold text-brand-900">Заявки на партнерство</h1>
            <p className="text-sm text-brand-500 mt-1">{counts.pending} ожидают рассмотрения</p>
          </div>
          <button onClick={() => adminApplicationsStore.load()} className="btn-ghost gap-1.5">
            <RefreshCw size={16} /> Обновить
          </button>
        </div>

        <div className="flex gap-2 flex-wrap mb-6">
          {STATUS_FILTERS.map((item) => (
            <button
              key={item.value}
              onClick={() => setFilter(item.value)}
              className={`px-4 py-2 rounded-xl text-sm font-medium transition-colors ${
                filter === item.value
                  ? 'bg-brand-500 text-white'
                  : 'bg-white border border-cream-200 text-brand-600 hover:border-brand-300'
              }`}
            >
              {item.label}
              <span className={`ml-1.5 text-xs ${filter === item.value ? 'opacity-80' : 'text-cream-400'}`}>
                {counts[item.value]}
              </span>
            </button>
          ))}
        </div>

        {adminApplicationsStore.loading && (
          <div className="flex justify-center py-20"><Spinner size={32} /></div>
        )}

        {adminApplicationsStore.error && (
          <div className="card text-center py-12 text-red-500">Не удалось загрузить заявки</div>
        )}

        {!adminApplicationsStore.loading && !adminApplicationsStore.error && filtered.length === 0 && (
          <div className="card text-center py-12 text-cream-500">Нет заявок в этом статусе</div>
        )}

        {!adminApplicationsStore.loading && filtered.length > 0 && (
          <div className="space-y-3">
            {filtered.map((app) => (
              <div key={app.id} className="card hover:shadow-md transition-shadow cursor-pointer" onClick={() => setSelected(app)}>
                <div className="flex items-start justify-between gap-4">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-3 mb-2">
                      <StatusBadge status={app.status} />
                      <span className="text-xs text-cream-400">{formatDateTime(app.created_at)}</span>
                    </div>
                    <h3 className="font-semibold text-brand-900 text-base truncate">{app.business_name}</h3>
                    <div className="flex flex-wrap gap-x-4 gap-y-1 mt-2 text-sm text-brand-500">
                      <span className="flex items-center gap-1"><Building2 size={13} /> {getCategoryLabel(app.category_code)}</span>
                      <span className="flex items-center gap-1"><Mail size={13} /> {app.contact_email}</span>
                      <span className="flex items-center gap-1"><MapPin size={13} /> {app.address}</span>
                    </div>
                  </div>
                  <button className="btn-ghost py-1.5 px-2" onClick={(e) => { e.stopPropagation(); setSelected(app) }}>
                    <Eye size={16} />
                  </button>
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
        loading={adminApplicationsStore.actionLoading}
      />
    </div>
  )
}

export default observer(AdminPageBase)
