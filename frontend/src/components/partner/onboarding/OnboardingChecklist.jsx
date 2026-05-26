import { useState } from 'react'
import { Link } from 'react-router-dom'
import {
  AlertTriangle,
  ArrowRight,
  CheckCircle2,
  Circle,
  Clock,
  MapPin,
  Package,
  ShoppingBag,
  Wallet,
  XCircle,
} from 'lucide-react'
import { useAuth } from '@/context/AuthContext'
import Button from '@/components/ui/actions/Button'

export default function OnboardingChecklist({ data }) {
  const { user } = useAuth()
  const status = data?.partner?.status
  const hasLegalInfo = data?.partner?.has_legal_info
  const legalInfoStatus = data?.partner?.legal_info_status
  const locations = data?.locations || []
  const hasActiveBoxes = locations.some((l) => l.active_boxes_count > 0)

  const dismissedKey = `partner_onboarding_v2_dismissed_${user?.partner_id}`
  const [dismissed, setDismissed] = useState(() => localStorage.getItem(dismissedKey) === 'true')

  const dismiss = () => {
    localStorage.setItem(dismissedKey, 'true')
    setDismissed(true)
  }

  if (status === 'pending_documents') {
    return <PendingDocsPanel hasLegalInfo={hasLegalInfo} legalInfoStatus={legalInfoStatus} />
  }

  if (status === 'active' && !dismissed) {
    const steps = [
      {
        icon: CheckCircle2,
        label: 'Юридические данные подтверждены',
        tip: 'Ваши данные проверены администратором — аккаунт полностью активен.',
        done: true,
        link: null,
      },
      {
        icon: MapPin,
        label: 'Создайте первую локацию',
        tip: 'Локация — это точка выдачи боксов. Укажите адрес, часы работы и фото.',
        done: locations.length > 0,
        link: '/partner/locations',
        actionLabel: 'Управление локациями',
      },
      {
        icon: Package,
        label: 'Добавьте первый бокс',
        tip: 'Бокс — набор продуктов по сниженной цене. Задайте состав, цену и количество.',
        done: hasActiveBoxes,
        link: '/partner/boxes/new',
        actionLabel: 'Создать бокс',
      },
      {
        icon: Wallet,
        label: 'Настройте реквизиты для выплат',
        tip: 'Укажите СБП-реквизиты или счёт, чтобы получать выплаты от продаж.',
        done: false,
        link: '/partner/payouts',
        actionLabel: 'Настроить',
      },
      {
        icon: ShoppingBag,
        label: 'Изучите управление заказами',
        tip: 'Подтверждайте заказы и выдавайте боксы клиентам по коду из приложения.',
        done: false,
        link: '/partner/orders',
        actionLabel: 'Перейти к заказам',
      },
    ]

    const completedCount = steps.filter((s) => s.done).length

    return (
      <section className="rounded-2xl border border-brand-200 bg-gradient-to-r from-brand-50 to-cream-50 p-5 md:p-6">
        <div className="flex items-start justify-between gap-4 mb-5">
          <div>
            <div className="inline-flex items-center gap-2 rounded-full border border-brand-200 bg-white/70 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-brand-700 mb-2">
              Быстрый старт
            </div>
            <h2 className="text-lg font-bold text-brand-900">Вы активированы! Вот что стоит сделать</h2>
            <p className="text-sm text-brand-600 mt-1">Выполните шаги ниже, чтобы начать принимать заказы</p>
          </div>
          <button
            onClick={dismiss}
            className="shrink-0 text-xs text-brand-500 hover:text-brand-700 underline underline-offset-2 transition-colors mt-1"
          >
            Скрыть
          </button>
        </div>

        <div className="mb-4">
          <div className="flex items-center justify-between text-xs font-medium text-brand-700 mb-1.5">
            <span>{completedCount} из {steps.length} выполнено</span>
            <span>{Math.round((completedCount / steps.length) * 100)}%</span>
          </div>
          <div className="h-2 overflow-hidden rounded-full bg-brand-200/60">
            <div
              className="h-full rounded-full bg-brand-500 transition-all duration-500"
              style={{ width: `${(completedCount / steps.length) * 100}%` }}
            />
          </div>
        </div>

        <div className="space-y-2">
          {steps.map((step, idx) => (
            <StepRow key={idx} step={step} />
          ))}
        </div>

        <div className="mt-4 pt-4 border-t border-brand-100 flex justify-end">
          <button
            onClick={dismiss}
            className="text-sm text-brand-500 hover:text-brand-700 transition-colors"
          >
            Понятно, не показывать подсказки
          </button>
        </div>
      </section>
    )
  }

  return null
}

function StepRow({ step }) {
  const { icon: Icon, label, tip, done, link, actionLabel } = step

  const content = (
    <div
      className={`flex items-start gap-3 rounded-xl border px-4 py-3 transition-colors ${
        done
          ? 'border-brand-200 bg-white/60'
          : 'border-cream-200 bg-white hover:border-brand-200'
      }`}
    >
      <div className={`mt-0.5 shrink-0 ${done ? 'text-brand-500' : 'text-cream-400'}`}>
        {done ? <CheckCircle2 size={18} /> : <Circle size={18} />}
      </div>
      <div className="flex-1 min-w-0">
        <p className={`text-sm font-medium ${done ? 'text-brand-700 line-through decoration-brand-300' : 'text-brand-900'}`}>
          {label}
        </p>
        {!done && <p className="text-xs text-brand-500 mt-0.5">{tip}</p>}
      </div>
      {!done && link && (
        <div className="shrink-0 flex items-center gap-1 text-xs font-medium text-brand-600">
          {actionLabel}
          <ArrowRight size={12} />
        </div>
      )}
    </div>
  )

  if (!done && link) {
    return <Link to={link}>{content}</Link>
  }

  return content
}

function PendingDocsPanel({ hasLegalInfo, legalInfoStatus }) {
  if (!hasLegalInfo) {
    return (
      <section className="rounded-2xl border border-amber-300 bg-gradient-to-r from-amber-50 to-orange-50 p-5 md:p-6">
        <div className="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
          <div className="space-y-2">
            <div className="inline-flex items-center gap-2 rounded-full border border-amber-300 bg-white/70 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-amber-800">
              <AlertTriangle size={14} />
              Онбординг партнёра
            </div>
            <h2 className="text-lg font-bold text-amber-950">Заполните юридические данные</h2>
            <p className="text-sm text-amber-900/90">
              Пока данные не заполнены, вы не можете принимать заказы.
              Активация недоступна до прохождения проверки.
            </p>
          </div>
          <Link to="/partner/legal-info" className="shrink-0">
            <Button className="w-full md:w-auto">Заполнить данные</Button>
          </Link>
        </div>
        <ProgressBar step={1} total={3} label="Шаг 1 из 3 — заполните данные" />
      </section>
    )
  }

  if (legalInfoStatus === 'failed') {
    return (
      <section className="rounded-2xl border border-red-300 bg-gradient-to-r from-red-50 to-orange-50 p-5 md:p-6">
        <div className="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
          <div className="space-y-2">
            <div className="inline-flex items-center gap-2 rounded-full border border-red-300 bg-white/70 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-red-700">
              <XCircle size={14} />
              Данные отклонены
            </div>
            <h2 className="text-lg font-bold text-red-950">Юридические данные не прошли проверку</h2>
            <p className="text-sm text-red-900/90">
              Администратор отклонил ваши данные. Проверьте корректность ИНН, ОГРН и юридического адреса и отправьте снова.
            </p>
          </div>
          <Link to="/partner/legal-info" className="shrink-0">
            <Button className="w-full md:w-auto bg-red-600 hover:bg-red-700">Исправить данные</Button>
          </Link>
        </div>
        <ProgressBar step={1} total={3} label="Шаг 1 из 3 — исправьте данные" color="red" />
      </section>
    )
  }

  return (
    <section className="rounded-2xl border border-blue-200 bg-gradient-to-r from-blue-50 to-indigo-50 p-5 md:p-6">
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div className="space-y-2">
          <div className="inline-flex items-center gap-2 rounded-full border border-blue-200 bg-white/70 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-blue-700">
            <Clock size={14} />
            На проверке
          </div>
          <h2 className="text-lg font-bold text-blue-950">Данные отправлены — ожидайте проверки</h2>
          <p className="text-sm text-blue-900/90">
            Администратор проверит ваши юридические данные и активирует аккаунт.
            Это обычно занимает до 1 рабочего дня.
          </p>
        </div>
        <div className="shrink-0 rounded-xl border border-blue-200 bg-white/70 px-4 py-3 text-center">
          <p className="text-xs text-blue-600 font-medium">Статус</p>
          <p className="text-sm font-bold text-blue-800 mt-0.5">На рассмотрении</p>
        </div>
      </div>
      <ProgressBar step={2} total={3} label="Шаг 2 из 3 — ожидание проверки" color="blue" />
    </section>
  )
}

function ProgressBar({ step, total, label, color = 'amber' }) {
  const colors = {
    amber: { bar: 'bg-amber-500', track: 'bg-amber-200/80', text: 'text-amber-900/80' },
    blue: { bar: 'bg-blue-500', track: 'bg-blue-200/80', text: 'text-blue-900/80' },
    red: { bar: 'bg-red-500', track: 'bg-red-200/80', text: 'text-red-900/80' },
  }
  const c = colors[color]

  return (
    <div className="mt-5">
      <div className={`mb-2 flex items-center justify-between text-xs font-medium ${c.text}`}>
        <span>{label}</span>
      </div>
      <div className={`h-2 overflow-hidden rounded-full ${c.track}`}>
        <div className={`h-full rounded-full ${c.bar}`} style={{ width: `${(step / total) * 100}%` }} />
      </div>
    </div>
  )
}
