import {
  PARTNER_OFFER_EFFECTIVE_DATE,
  PARTNER_OFFER_TITLE,
  PARTNER_OFFER_VERSION,
  partnerOfferSections,
} from '@/lib/partnerOffer'

export default function PartnerOfferPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-cream-50 via-white to-brand-50/60 px-4 py-8 md:px-8">
      <div className="mx-auto max-w-4xl">
        <div className="card space-y-6">
          <header className="border-b border-cream-200 pb-5">
            <p className="text-sm font-medium uppercase tracking-wide text-brand-600">Публичный документ</p>
            <h1 className="mt-2 text-3xl font-bold text-brand-950">{PARTNER_OFFER_TITLE}</h1>
            <p className="mt-2 text-sm text-brand-700">
              Редакция {PARTNER_OFFER_VERSION} от {PARTNER_OFFER_EFFECTIVE_DATE}
            </p>
          </header>

          <section className="space-y-4 text-sm leading-7 text-brand-900">
            <p>
              Настоящий документ регулирует использование IT-платформы Berezhok партнёрами, размещающими боксы для продажи клиентам.
              Платформа действует как агрегатор, а договор купли-продажи товара заключается напрямую между клиентом и партнёром.
            </p>

            {partnerOfferSections.map((section) => (
              <article key={section.title} className="space-y-2">
                <h2 className="text-lg font-semibold text-brand-950">{section.title}</h2>
                {section.body.map((paragraph) => (
                  <p key={paragraph}>{paragraph}</p>
                ))}
              </article>
            ))}

            <article className="space-y-2">
              <h2 className="text-lg font-semibold text-brand-950">8. Электронный акцепт</h2>
              <p>
                Проставление отметки о согласии в интерфейсе личного кабинета и нажатие кнопки подтверждения рассматриваются как акцепт оферты в электронной форме.
              </p>
              <p>
                Для юридически значимой фиксации принятия оферты платформа вправе хранить сведения о версии оферты, времени акцепта, идентификаторе партнёра и технических параметрах сессии.
              </p>
            </article>
          </section>
        </div>
      </div>
    </div>
  )
}
