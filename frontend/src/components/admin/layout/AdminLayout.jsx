import AdminSidebar from '@/components/admin/layout/AdminSidebar'

export default function AdminLayout({ title, subtitle, actions, children }) {
  return (
    <div className="min-h-screen bg-gradient-to-br from-cream-50 via-white to-brand-50/60 md:flex">
      <AdminSidebar />

      <div className="flex-1 min-w-0">
        <header className="border-b border-cream-200/80 bg-white/75 px-4 pb-5 pt-6 backdrop-blur-sm md:px-8 md:pb-6 md:pt-8">
          <div className="flex flex-wrap items-start justify-between gap-4">
            <div>
              <h1 className="text-2xl font-bold text-brand-900 md:text-3xl">{title}</h1>
              {subtitle && <p className="mt-1 text-sm text-brand-600">{subtitle}</p>}
            </div>
            {actions && <div className="flex flex-wrap items-center gap-2">{actions}</div>}
          </div>
        </header>

        <main className="px-4 py-6 pb-28 md:px-8 md:py-8 md:pb-8">{children}</main>
      </div>
    </div>
  )
}
