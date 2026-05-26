import { useEffect, useState } from 'react'
import { observer } from 'mobx-react-lite'
import { toast } from 'sonner'
import { Banknote, ChevronRight, Pencil, Save, X } from 'lucide-react'

import PartnerLayout from '@/components/partner/layout/PartnerLayout'
import Input from '@/components/ui/form/Input'
import PhoneInput from '@/components/ui/form/PhoneInput'
import Label from '@/components/ui/form/Label'
import Select from '@/components/ui/form/Select'
import Button from '@/components/ui/actions/Button'
import StatusBadge from '@/components/ui/feedback/StatusBadge'
import { payoutsStore } from '@/stores/partner/payoutsStore'

const STATUS_LABELS = {
  pending: 'Ожидает',
  processing: 'Обрабатывается',
  completed: 'Выплачено',
  failed: 'Ошибка',
}

function formatDate(isoStr) {
  if (!isoStr) return '—'
  return new Date(isoStr).toLocaleDateString('ru-RU')
}

function formatAmount(str) {
  if (!str) return '—'
  return new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'RUB' }).format(str)
}

function friendlyError(msg) {
  if (!msg) return null
  if (msg.startsWith('yookassa [') || msg.startsWith('error:')) {
    return 'Ошибка провайдера. Обратитесь в поддержку.'
  }
  if (msg === 'payout destination not configured') return 'Реквизиты не настроены.'
  if (msg === 'payout canceled by provider') return 'Выплата отклонена провайдером.'
  return msg
}

const PartnerPayoutsPage = observer(() => {
  const store = payoutsStore

  const [modalOpen, setModalOpen] = useState(false)
  const [form, setForm] = useState({ type: 'sbp', sbp_phone: '', sbp_bank_id: '', recipient_name: '' })
  const [errors, setErrors] = useState({})
  const [detailPayout, setDetailPayout] = useState(null)

  useEffect(() => {
    store.loadDestination()
    store.loadHistory({ limit: 20, offset: 0 })
  }, [store])

  useEffect(() => {
    if (store.destination) {
      setForm({
        type: store.destination.type || 'sbp',
        sbp_phone: store.destination.sbp_phone || '',
        sbp_bank_id: store.destination.sbp_bank_id || '',
        recipient_name: store.destination.recipient_name || '',
      })
    }
  }, [store.destination])

  function openModal() {
    store.loadBanks()
    setModalOpen(true)
  }

  function handleChange(field, value) {
    setForm((f) => ({ ...f, [field]: value }))
    setErrors((e) => ({ ...e, [field]: undefined }))
  }

  function validate() {
    const errs = {}
    if (!form.sbp_phone.match(/^\+7\d{10}$/)) {
      errs.sbp_phone = 'Введите полный номер телефона'
    }
    if (!form.sbp_bank_id) {
      errs.sbp_bank_id = 'Выберите банк'
    }
    if (form.recipient_name.trim().length < 2) {
      errs.recipient_name = 'Укажите ФИО получателя'
    }
    return errs
  }

  async function handleSave(e) {
    e.preventDefault()
    const errs = validate()
    if (Object.keys(errs).length > 0) {
      setErrors(errs)
      return
    }
    try {
      await store.saveDestination(form)
      toast.success('Реквизиты сохранены')
      setModalOpen(false)
    } catch {
      toast.error('Не удалось сохранить реквизиты')
    }
  }

  async function handleShowDetail(id) {
    try {
      const data = await store.loadById(id)
      setDetailPayout(data)
    } catch {
      toast.error('Не удалось загрузить детали выплаты')
    }
  }

  const bankName = store.banks.find((b) => b.bank_id === store.destination?.sbp_bank_id)?.name

  return (
    <PartnerLayout>
      <div className="max-w-4xl mx-auto space-y-6 p-4 sm:p-6">

        {/* Page header */}
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Выплаты</h1>
          <p className="text-sm text-gray-500 mt-1">Получайте вознаграждение за выкупленные боксы</p>
        </div>

        {/* Requisites card */}
        <section className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6">
          {store.destination ? (
            <div className="flex items-center justify-between gap-4">
              <div className="flex items-center gap-4">
                <div className="w-12 h-12 rounded-xl bg-green-50 flex items-center justify-center shrink-0">
                  <Banknote className="w-6 h-6 text-green-600" />
                </div>
                <div>
                  <p className="font-semibold text-gray-900">
                    {bankName || store.destination.sbp_bank_id} &bull; {store.destination.sbp_phone}
                  </p>
                  <p className="text-sm text-gray-500">{store.destination.recipient_name}</p>
                </div>
              </div>
              <button
                onClick={openModal}
                className="flex items-center gap-1 text-sm text-blue-600 hover:text-blue-700 font-medium shrink-0"
              >
                <Pencil className="w-4 h-4" />
                Изменить
              </button>
            </div>
          ) : (
            <div className="flex items-center justify-between gap-4">
              <div className="flex items-center gap-4">
                <div className="w-12 h-12 rounded-xl bg-yellow-50 flex items-center justify-center shrink-0">
                  <Banknote className="w-6 h-6 text-yellow-500" />
                </div>
                <div>
                  <p className="font-semibold text-gray-900">Реквизиты не настроены</p>
                  <p className="text-sm text-gray-500">Укажите банк СБП, чтобы получать выплаты</p>
                </div>
              </div>
              <Button onClick={openModal} className="flex items-center gap-1 shrink-0">
                Настроить <ChevronRight className="w-4 h-4" />
              </Button>
            </div>
          )}
        </section>

        {/* Payout history */}
        <section className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6">
          <h2 className="text-base font-semibold text-gray-900 mb-4">История выплат</h2>

          {store.loading && <p className="text-gray-400 text-sm">Загрузка...</p>}

          {!store.loading && store.items.length === 0 && (
            <p className="text-gray-400 text-sm">Выплат пока не было</p>
          )}

          {store.items.length > 0 && (
            <div className="overflow-x-auto -mx-6">
              <table className="w-full text-sm">
                <thead>
                  <tr className="text-left text-gray-400 border-b text-xs uppercase tracking-wide">
                    <th className="pb-3 pl-6 pr-4 font-medium">Период</th>
                    <th className="pb-3 pr-4 font-medium">Брутто</th>
                    <th className="pb-3 pr-4 font-medium">Комиссия</th>
                    <th className="pb-3 pr-4 font-medium">К выплате</th>
                    <th className="pb-3 pr-4 font-medium">Статус</th>
                    <th className="pb-3 pr-6 font-medium"></th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-50">
                  {store.items.map((p) => (
                    <tr key={p.id} className="hover:bg-gray-50 transition-colors">
                      <td className="py-3 pl-6 pr-4 whitespace-nowrap text-gray-700">
                        {formatDate(p.period_start)} – {formatDate(p.period_end)}
                      </td>
                      <td className="py-3 pr-4 text-gray-500">{formatAmount(p.gross_amount)}</td>
                      <td className="py-3 pr-4 text-gray-500">{formatAmount(p.commission_amount)}</td>
                      <td className="py-3 pr-4 font-semibold text-gray-900">{formatAmount(p.net_amount)}</td>
                      <td className="py-3 pr-4">
                        <StatusBadge status={p.status} label={STATUS_LABELS[p.status] || p.status} />
                      </td>
                      <td className="py-3 pr-6">
                        <button
                          onClick={() => handleShowDetail(p.id)}
                          className="text-blue-600 hover:text-blue-800 text-xs font-medium"
                        >
                          Детали
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>
      </div>

      {/* Requisites modal */}
      {modalOpen && (
        <div
          className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 p-4"
          onClick={() => setModalOpen(false)}
        >
          <div
            className="bg-white rounded-2xl shadow-xl w-full max-w-md"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex items-center justify-between p-6 border-b">
              <h3 className="font-semibold text-base">Реквизиты для выплат (СБП)</h3>
              <button onClick={() => setModalOpen(false)} className="text-gray-400 hover:text-gray-700">
                <X className="w-5 h-5" />
              </button>
            </div>

            <form onSubmit={handleSave} className="p-6 space-y-4">
              <div>
                <Label htmlFor="sbp_phone">Номер телефона СБП</Label>
                <PhoneInput
                  id="sbp_phone"
                  value={form.sbp_phone}
                  onChange={(e) => handleChange('sbp_phone', e.target.value)}
                  error={errors.sbp_phone}
                />
                {errors.sbp_phone && <p className="text-red-500 text-xs mt-1">{errors.sbp_phone}</p>}
              </div>

              <div>
                <Label htmlFor="sbp_bank_id">Банк получателя</Label>
                {store.banksLoading ? (
                  <p className="text-sm text-gray-400">Загрузка банков...</p>
                ) : (
                  <Select
                    id="sbp_bank_id"
                    value={form.sbp_bank_id}
                    onChange={(e) => handleChange('sbp_bank_id', e.target.value)}
                  >
                    <option value="">— выберите банк —</option>
                    {store.banks.map((b) => (
                      <option key={b.bank_id} value={b.bank_id}>
                        {b.name}
                      </option>
                    ))}
                  </Select>
                )}
                {errors.sbp_bank_id && <p className="text-red-500 text-xs mt-1">{errors.sbp_bank_id}</p>}
              </div>

              <div>
                <Label htmlFor="recipient_name">ФИО получателя</Label>
                <Input
                  id="recipient_name"
                  placeholder="Иванов Иван Иванович"
                  value={form.recipient_name}
                  onChange={(e) => handleChange('recipient_name', e.target.value)}
                  error={errors.recipient_name}
                />
                {errors.recipient_name && (
                  <p className="text-red-500 text-xs mt-1">{errors.recipient_name}</p>
                )}
              </div>

              <Button type="submit" disabled={store.submitting} className="w-full flex items-center justify-center gap-2">
                <Save className="w-4 h-4" />
                {store.submitting ? 'Сохраняем...' : 'Сохранить'}
              </Button>
            </form>
          </div>
        </div>
      )}

      {/* Payout detail modal */}
      {detailPayout && (
        <div
          className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 p-4"
          onClick={() => setDetailPayout(null)}
        >
          <div
            className="bg-white rounded-2xl shadow-xl w-full max-w-lg max-h-[85vh] overflow-y-auto"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex items-center justify-between p-6 border-b">
              <h3 className="font-semibold text-base">Детали выплаты</h3>
              <button onClick={() => setDetailPayout(null)} className="text-gray-400 hover:text-gray-700">
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="p-6 space-y-4">
              <dl className="space-y-3 text-sm">
                <div className="flex justify-between">
                  <dt className="text-gray-500">Период</dt>
                  <dd>{formatDate(detailPayout.period_start)} – {formatDate(detailPayout.period_end)}</dd>
                </div>
                <div className="flex justify-between">
                  <dt className="text-gray-500">Брутто</dt>
                  <dd>{formatAmount(detailPayout.gross_amount)}</dd>
                </div>
                <div className="flex justify-between">
                  <dt className="text-gray-500">
                    Комиссия ({(parseFloat(detailPayout.commission_rate_applied) * 100).toFixed(1)}%)
                  </dt>
                  <dd>{formatAmount(detailPayout.commission_amount)}</dd>
                </div>
                <div className="flex justify-between font-semibold border-t pt-3">
                  <dt>К выплате</dt>
                  <dd>{formatAmount(detailPayout.net_amount)}</dd>
                </div>
                <div className="flex justify-between items-center">
                  <dt className="text-gray-500">Статус</dt>
                  <dd>
                    <StatusBadge
                      status={detailPayout.status}
                      label={STATUS_LABELS[detailPayout.status] || detailPayout.status}
                    />
                  </dd>
                </div>
                {detailPayout.error_message && (
                  <div className="bg-red-50 rounded-lg p-3 text-sm text-red-700">
                    {friendlyError(detailPayout.error_message)}
                  </div>
                )}
              </dl>

              {detailPayout.orders && detailPayout.orders.length > 0 && (
                <div className="border-t pt-4">
                  <h4 className="font-medium text-sm mb-3">
                    Заказы в выплате ({detailPayout.orders.length})
                  </h4>
                  <div className="space-y-2">
                    {detailPayout.orders.map((o) => (
                      <div key={o.order_id} className="flex justify-between text-xs text-gray-600 py-1">
                        <span className="font-mono text-gray-400">{o.order_id.slice(0, 8)}…</span>
                        <span>{formatAmount(o.order_amount)}</span>
                        <span className="text-gray-400">−{formatAmount(o.commission_part)}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </PartnerLayout>
  )
})

export default PartnerPayoutsPage
