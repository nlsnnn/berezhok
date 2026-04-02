import PartnerSidebar from '@/components/partner/layout/PartnerSidebar'

export default function PartnerLayout({ title, subtitle, actions, children }) {
  return (
    <div className="min-h-screen bg-gradient-to-br from-cream-50 via-white to-brand-50/60 md:flex">
      <PartnerSidebar />

      <div className="flex-1 min-w-0">
        <header className="px-4 pt-16 pb-5 md:px-8 md:pt-8 md:pb-6 border-b border-cream-200/80 bg-white/70 backdrop-blur-sm">
          <div className="flex items-start justify-between gap-4 flex-wrap">
            <div>
              <h1 className="text-2xl md:text-3xl font-bold text-brand-900">{title}</h1>
              {subtitle && <p className="text-sm text-brand-600 mt-1">{subtitle}</p>}
            </div>
            {actions && <div className="flex items-center gap-2">{actions}</div>}
          </div>
        </header>

        <main className="px-4 py-6 md:px-8 md:py-8">{children}</main>
      </div>
    </div>
  )
}
