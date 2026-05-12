import { useMemo, useState } from 'react'
import { FileText, ExternalLink, ShieldCheck } from 'lucide-react'
import { toast } from 'sonner'
import { Link } from 'react-router-dom'
import Modal from '@/components/ui/overlay/Modal'
import Button from '@/components/ui/actions/Button'
import { useAuth } from '@/context/AuthContext'
import {
  PARTNER_OFFER_EFFECTIVE_DATE,
  PARTNER_OFFER_TITLE,
  PARTNER_OFFER_VERSION,
  partnerOfferSections,
} from '@/lib/partnerOffer'

export default function PartnerOfferModal() {
  const { user, markOfferAccepted } = useAuth()
  const [confirmed, setConfirmed] = useState(false)
  const [submitting, setSubmitting] = useState(false)

  const shouldShow = useMemo(() => {
    if (!user) return false
    if (user.must_change_password) return false
    if (user.role === 'employee') return false
    return user.accepted_offer_version !== PARTNER_OFFER_VERSION
  }, [user])

  if (!shouldShow) return null

  const handleAccept = async () => {
    if (!confirmed) {
      toast.error('Подтвердите согласие с условиями оферты')
      return
    }

    setSubmitting(true)
    try {
      markOfferAccepted()
      toast.success('Оферта принята. Доступ к кабинету открыт.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <Modal
      open={shouldShow}
      onClose={() => {}}
      title="Подтвердите акцепт оферты"
      className="max-w-4xl"
      closeOnBackdrop={false}
      closeOnEscape={false}
      showCloseButton={false}
    >
      <div className="space-y-6">
        <section className="rounded-2xl border border-brand-200 bg-gradient-to-r from-brand-50 to-white p-5">
          <div className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
            <div className="space-y-2">
              <div className="inline-flex items-center gap-2 rounded-full bg-white px-3 py-1 text-xs font-semibold uppercase tracking-wide text-brand-700 shadow-sm">
                <ShieldCheck size={14} />
                Обязательное подтверждение
              </div>
              <div>
                <h2 className="text-xl font-bold text-brand-950">{PARTNER_OFFER_TITLE}</h2>
                <p className="mt-1 text-sm text-brand-700">
                  Редакция от {PARTNER_OFFER_EFFECTIVE_DATE}. Акцепт обязателен для продолжения работы в кабинете партнёра.
                </p>
              </div>
            </div>

            <Link
              to="/partner/offer"
              target="_blank"
              rel="noreferrer"
              className="inline-flex items-center gap-2 text-sm font-medium text-brand-700 hover:text-brand-900"
            >
              <ExternalLink size={16} />
              Открыть полную редакцию
            </Link>
          </div>
        </section>

        <section className="max-h-[45vh] space-y-4 overflow-y-auto pr-2">
          {partnerOfferSections.map((section) => (
            <article key={section.title} className="rounded-2xl border border-cream-200 bg-cream-50/70 p-4">
              <h3 className="flex items-center gap-2 text-base font-semibold text-brand-900">
                <FileText size={16} className="text-brand-500" />
                {section.title}
              </h3>
              <div className="mt-3 space-y-2 text-sm leading-6 text-brand-800">
                {section.body.map((paragraph) => (
                  <p key={paragraph}>{paragraph}</p>
                ))}
              </div>
            </article>
          ))}
        </section>

        <label className="flex items-start gap-3 rounded-2xl border border-cream-200 bg-white p-4">
          <input
            type="checkbox"
            checked={confirmed}
            onChange={(e) => setConfirmed(e.target.checked)}
            className="mt-1 h-4 w-4 rounded border-cream-300 text-brand-600 focus:ring-brand-500"
          />
          <span className="text-sm leading-6 text-brand-900">
            Я принимаю условия соглашения и подтверждаю акцепт публичной оферты для партнёров Berezhok, редакция {PARTNER_OFFER_VERSION}.
          </span>
        </label>

        <div className="flex flex-col-reverse gap-3 border-t border-cream-200 pt-4 md:flex-row md:items-center md:justify-between">
          <p className="text-xs leading-5 text-brand-600">
            До акцепта оферты доступ к функциям кабинета ограничен. Текущая реализация сохраняет подтверждение в браузере для данного партнёра.
          </p>
          <Button onClick={handleAccept} disabled={!confirmed || submitting} className="w-full md:w-auto">
            {submitting ? 'Подтверждаем...' : 'Принять условия'}
          </Button>
        </div>
      </div>
    </Modal>
  )
}
