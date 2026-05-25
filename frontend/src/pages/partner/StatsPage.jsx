import { useEffect, useCallback } from 'react'
import { observer } from 'mobx-react-lite'
import {
  BarChart3, CalendarRange, Download, MapPin, Package, RefreshCcw,
  ShoppingBag, Star, Store, Wallet,
} from 'lucide-react'
import {
  Area, Bar, CartesianGrid, Cell, ComposedChart, Legend,
  Pie, PieChart, ResponsiveContainer, Tooltip, XAxis, YAxis,
} from 'recharts'
import * as XLSX from 'xlsx'
import { toast } from 'sonner'
import PartnerLayout from '@/components/partner/layout/PartnerLayout'
import Button from '@/components/ui/actions/Button'
import Select from '@/components/ui/form/Select'
import Spinner from '@/components/ui/feedback/Spinner'
import { useStores } from '@/context/StoresContext'
import { getOrderStatusMeta, PARTNER_ORDER_FILTERS } from '@/lib/orderStatus'
import { formatDate, formatDateTime, getErrorMessage } from '@/lib/utils'

const PERIOD_OPTIONS = [
  { value: 'today', label: 'Сегодня' },
  { value: 'last_7_days', label: 'Последние 7 дней' },
  { value: 'last_30_days', label: 'Последние 30 дней' },
  { value: 'this_month', label: 'Этот месяц' },
  { value: 'last_month', label: 'Прошлый месяц' },
  { value: 'custom', label: 'Свой период' },
]

const TOP_LOCATION_SORTS = [
  { value: 'revenue_desc', label: 'По выручке' },
  { value: 'orders_desc', label: 'По заказам' },
  { value: 'rating_desc', label: 'По рейтингу' },
  { value: 'name_asc', label: 'По названию' },
]

const TOP_BOX_SORTS = [
  { value: 'revenue_desc', label: 'По выручке' },
  { value: 'orders_desc', label: 'По заказам' },
  { value: 'name_asc', label: 'По названию' },
]

const ORDER_SORTS = [
  { value: 'created_at_desc', label: 'Сначала новые' },
  { value: 'created_at_asc', label: 'Сначала старые' },
  { value: 'amount_desc', label: 'Сумма по убыванию' },
  { value: 'amount_asc', label: 'Сумма по возрастанию' },
  { value: 'pickup_time_desc', label: 'Поздняя выдача' },
  { value: 'pickup_time_asc', label: 'Ранняя выдача' },
]

const STATUS_PIE_COLORS = ['#6d7bd1', '#a3e635', '#f97316', '#ef4444', '#8b5cf6', '#06b6d4']

const CustomTooltip = ({ active, payload }) => {
  if (!active || !payload?.length) return null
  const d = payload[0].payload
  return (
    <div className="rounded-xl border border-cream-200 bg-white px-3 py-2 text-sm shadow-md">
      <p className="font-semibold text-brand-900">{d.name}</p>
      <p className="text-brand-600">{d.value} заказов · {d.share}%</p>
    </div>
  )
}

function exportToExcel(data, filters) {
  const wb = XLSX.utils.book_new()

  const summary = data?.summary ?? {}
  const summaryRows = [
    ['Показатель', 'Значение'],
    ['Всего заказов', summary.orders_total ?? 0],
    ['Завершено', summary.orders_completed ?? 0],
    ['Отменено', summary.orders_cancelled ?? 0],
    ['Валовая выручка (₽)', summary.gross_revenue ?? 0],
    ['Чистая выручка (₽)', summary.net_revenue ?? 0],
    ['Средний чек (₽)', Math.round(summary.avg_order_value ?? 0)],
    ['Средний рейтинг', (summary.avg_rating ?? 0).toFixed(2)],
    ['Количество отзывов', summary.reviews_count ?? 0],
  ]
  XLSX.utils.book_append_sheet(wb, XLSX.utils.aoa_to_sheet(summaryRows), 'Сводка')

  if (data?.timeline?.length) {
    const timelineRows = [['Дата', 'Заказов', 'Выручка (₽)', 'Чистыми (₽)', 'Рейтинг']]
    data.timeline.forEach((item) => {
      timelineRows.push([
        item.date,
        item.orders_total ?? 0,
        item.gross_revenue ?? 0,
        item.net_revenue ?? 0,
        item.avg_rating ?? 0,
      ])
    })
    XLSX.utils.book_append_sheet(wb, XLSX.utils.aoa_to_sheet(timelineRows), 'Динамика')
  }

  if (data?.top_locations?.length) {
    const locRows = [['Локация', 'Адрес', 'Заказов', 'Выручка (₽)', 'Чистыми (₽)', 'Рейтинг']]
    data.top_locations.forEach((item) => {
      locRows.push([item.name, item.address, item.orders_total, item.gross_revenue, item.net_revenue, item.avg_rating?.toFixed(2) ?? '—'])
    })
    XLSX.utils.book_append_sheet(wb, XLSX.utils.aoa_to_sheet(locRows), 'Топ локаций')
  }

  if (data?.top_boxes?.length) {
    const boxRows = [['Бокс', 'Локация', 'Заказов', 'Выполнено', 'Выручка (₽)', 'Чистыми (₽)']]
    data.top_boxes.forEach((item) => {
      boxRows.push([item.name, item.location_name, item.orders_total, item.orders_completed, item.gross_revenue, item.net_revenue])
    })
    XLSX.utils.book_append_sheet(wb, XLSX.utils.aoa_to_sheet(boxRows), 'Топ боксов')
  }

  if (data?.orders?.length) {
    const orderRows = [['Код', 'Статус', 'Бокс', 'Локация', 'Клиент', 'Сумма (₽)', 'Выдача с', 'Выдача по', 'Создан']]
    data.orders.forEach((o) => {
      orderRows.push([
        o.pickup_code,
        getOrderStatusMeta(o.status).label,
        o.box_name ?? '—',
        o.location_name ?? '—',
        o.customer_name ?? o.customer_phone ?? '—',
        Math.round(o.amount ?? 0),
        o.pickup_time_start ? formatDateTime(o.pickup_time_start) : '—',
        o.pickup_time_end ? formatDateTime(o.pickup_time_end) : '—',
        o.created_at ? formatDateTime(o.created_at) : '—',
      ])
    })
    XLSX.utils.book_append_sheet(wb, XLSX.utils.aoa_to_sheet(orderRows), 'Заказы')
  }

  const period = filters.period === 'custom'
    ? `${filters.date_from}_${filters.date_to}`
    : filters.period
  XLSX.writeFile(wb, `berezhok_stats_${period}.xlsx`)
}

function StatsPageBase() {
  const { statsStore, locationsStore } = useStores()

  useEffect(() => {
    statsStore.load().catch((error) => {
      toast.error(getErrorMessage(error))
    })
    if (!locationsStore.profile) {
      locationsStore.loadProfile().catch(() => {})
    }
  }, [statsStore, locationsStore])

  const data = statsStore.data
  const summary = data?.summary

  const handleReload = useCallback(async () => {
    try {
      await statsStore.load()
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }, [statsStore])

  const handlePeriodChange = async (event) => {
    statsStore.setPeriod(event.target.value)
    if (event.target.value !== 'custom') await handleReload()
  }

  const handleCustomRangeSubmit = async (event) => {
    event.preventDefault()
    statsStore.setCustomRange(statsStore.filters.date_from, statsStore.filters.date_to)
    await handleReload()
  }

  const handleLocationChange = async (event) => {
    statsStore.setLocation(event.target.value)
    await handleReload()
  }

  const handleStatusChange = async (event) => {
    statsStore.setStatus(event.target.value)
    await handleReload()
  }

  const handleTopLocationsSortChange = async (event) => {
    statsStore.setTopLocationsSort(event.target.value)
    await handleReload()
  }

  const handleTopBoxesSortChange = async (event) => {
    statsStore.setTopBoxesSort(event.target.value)
    await handleReload()
  }

  const handleOrdersSortChange = async (event) => {
    statsStore.setOrdersSort(event.target.value)
    await handleReload()
  }

  const handlePageChange = async (nextOffset) => {
    statsStore.setPage(nextOffset)
    await handleReload()
  }

  const handleExport = () => {
    if (!data) return
    try {
      exportToExcel(data, statsStore.filters)
      toast.success('Файл Excel сохранён')
    } catch {
      toast.error('Не удалось создать файл')
    }
  }

  return (
    <PartnerLayout
      title="Статистика"
      subtitle="Аналитика по заказам, выручке, качеству сервиса и эффективности точек"
      actions={
        <div className="flex gap-2">
          {data && (
            <Button variant="secondary" className="gap-2" onClick={handleExport}>
              <Download size={16} />
              Excel
            </Button>
          )}
          <Button variant="secondary" className="gap-2" onClick={handleReload} disabled={statsStore.loading}>
            <RefreshCcw size={16} />
            Обновить
          </Button>
        </div>
      }
    >
      <div className="space-y-6">
        {/* Фильтры */}
        <section className="card space-y-4">
          <div className="grid gap-4 xl:grid-cols-[1.2fr_1fr_1fr]">
            <div>
              <p className="mb-2 text-sm font-medium text-brand-700">Период</p>
              <Select value={statsStore.filters.period} onChange={handlePeriodChange}>
                {PERIOD_OPTIONS.map((o) => (
                  <option key={o.value} value={o.value}>{o.label}</option>
                ))}
              </Select>
            </div>
            <div>
              <p className="mb-2 text-sm font-medium text-brand-700">Локация</p>
              <Select value={statsStore.filters.location_id} onChange={handleLocationChange}>
                <option value="">Все локации</option>
                {locationsStore.locations.map((l) => (
                  <option key={l.id} value={l.id}>{l.name}</option>
                ))}
              </Select>
            </div>
            <div>
              <p className="mb-2 text-sm font-medium text-brand-700">Статус заказа</p>
              <Select value={statsStore.filters.status} onChange={handleStatusChange}>
                <option value="">Все статусы</option>
                {PARTNER_ORDER_FILTERS.filter((i) => i.value !== 'all').map((i) => (
                  <option key={i.value} value={i.value}>{i.label}</option>
                ))}
              </Select>
            </div>
          </div>

          {statsStore.filters.period === 'custom' && (
            <form className="grid gap-4 md:grid-cols-[1fr_1fr_auto]" onSubmit={handleCustomRangeSubmit}>
              <label className="space-y-2">
                <span className="text-sm font-medium text-brand-700">Дата с</span>
                <input
                  type="date"
                  className="input-base"
                  value={statsStore.filters.date_from}
                  onChange={(e) => statsStore.setDraftDateFrom(e.target.value)}
                />
              </label>
              <label className="space-y-2">
                <span className="text-sm font-medium text-brand-700">Дата по</span>
                <input
                  type="date"
                  className="input-base"
                  value={statsStore.filters.date_to}
                  onChange={(e) => statsStore.setDraftDateTo(e.target.value)}
                />
              </label>
              <div className="flex items-end">
                <Button type="submit" className="w-full md:w-auto">Применить</Button>
              </div>
            </form>
          )}

          {data?.meta && (
            <div className="rounded-xl border border-cream-200 bg-cream-50 px-4 py-3 text-sm text-brand-700">
              Период отчёта: {formatDate(data.meta.date_from)} — {formatDate(data.meta.date_to)}
            </div>
          )}
        </section>

        {statsStore.loading && !data && (
          <div className="card flex justify-center py-20">
            <Spinner size={34} />
          </div>
        )}

        {statsStore.error && !data && (
          <div className="card flex flex-col items-center gap-4 py-16 text-center">
            <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-red-50">
              <RefreshCcw size={24} className="text-red-500" />
            </div>
            <div>
              <h3 className="text-lg font-semibold text-brand-900">Не удалось загрузить статистику</h3>
              <p className="mt-2 max-w-md text-sm text-brand-600">{getErrorMessage(statsStore.error)}</p>
            </div>
            <Button onClick={handleReload}>Повторить</Button>
          </div>
        )}

        {data && (
          <>
            {/* Метрики */}
            <section className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
              <MetricCard icon={ShoppingBag} label="Всего заказов" value={summary?.orders_total ?? 0} hint={`Завершено: ${summary?.orders_completed ?? 0}`} />
              <MetricCard icon={Wallet} label="Выручка" value={`${summary?.gross_revenue ?? 0} ₽`} hint={`Чистыми: ${summary?.net_revenue ?? 0} ₽`} />
              <MetricCard icon={CalendarRange} label="Средний чек" value={`${Math.round(summary?.avg_order_value ?? 0)} ₽`} hint={`Отменено: ${summary?.orders_cancelled ?? 0}`} />
              <MetricCard icon={Star} label="Рейтинг" value={(summary?.avg_rating ?? 0).toFixed(1)} hint={`Отзывов: ${summary?.reviews_count ?? 0}`} />
            </section>

            {/* Графики: динамика + статусы */}
            <section className="grid gap-5 xl:grid-cols-[1.4fr_1fr]">
              <article className="card">
                <div className="mb-4 flex items-center justify-between gap-3">
                  <div>
                    <h2 className="text-lg font-semibold text-brand-900">Динамика по дням</h2>
                    <p className="text-sm text-brand-600">Заказы и выручка по выбранному периоду</p>
                  </div>
                  <div className="rounded-xl bg-brand-50 px-3 py-2 text-sm text-brand-700">
                    Всего дней: {data.timeline?.length ?? 0}
                  </div>
                </div>
                <TimelineChart timeline={data.timeline || []} />
              </article>

              <article className="card">
                <div className="mb-4">
                  <h2 className="text-lg font-semibold text-brand-900">Срез по статусам</h2>
                  <p className="text-sm text-brand-600">Распределение заказов в периоде</p>
                </div>
                <StatusPieChart breakdown={data.status_breakdown || []} />
              </article>
            </section>

            {/* Топы */}
            <section className="grid gap-5 xl:grid-cols-2">
              <article className="card">
                <div className="mb-4 flex items-start justify-between gap-4">
                  <div>
                    <h2 className="text-lg font-semibold text-brand-900">Топ локаций</h2>
                    <p className="text-sm text-brand-600">Сравнение точек по ключевым показателям</p>
                  </div>
                  <div className="w-48">
                    <Select value={statsStore.filters.top_locations_sort} onChange={handleTopLocationsSortChange}>
                      {TOP_LOCATION_SORTS.map((o) => (
                        <option key={o.value} value={o.value}>{o.label}</option>
                      ))}
                    </Select>
                  </div>
                </div>
                <div className="space-y-3">
                  {(data.top_locations || []).map((item) => (
                    <LocationCard key={item.location_id} item={item} />
                  ))}
                  {(!data.top_locations || data.top_locations.length === 0) && <EmptyBlock text="Нет данных по локациям." />}
                </div>
              </article>

              <article className="card">
                <div className="mb-4 flex items-start justify-between gap-4">
                  <div>
                    <h2 className="text-lg font-semibold text-brand-900">Топ боксов</h2>
                    <p className="text-sm text-brand-600">Самые сильные позиции по продажам</p>
                  </div>
                  <div className="w-48">
                    <Select value={statsStore.filters.top_boxes_sort} onChange={handleTopBoxesSortChange}>
                      {TOP_BOX_SORTS.map((o) => (
                        <option key={o.value} value={o.value}>{o.label}</option>
                      ))}
                    </Select>
                  </div>
                </div>
                <div className="space-y-3">
                  {(data.top_boxes || []).map((item) => (
                    <BoxCard key={item.box_id} item={item} />
                  ))}
                  {(!data.top_boxes || data.top_boxes.length === 0) && <EmptyBlock text="Нет данных по боксам." />}
                </div>
              </article>
            </section>

            {/* Список заказов */}
            <section className="card">
              <div className="mb-4 flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
                <div>
                  <h2 className="text-lg font-semibold text-brand-900">Заказы в периоде</h2>
                  <p className="text-sm text-brand-600">Детальный аналитический список по текущим фильтрам</p>
                </div>
                <div className="w-full max-w-xs">
                  <Select value={statsStore.filters.orders_sort} onChange={handleOrdersSortChange}>
                    {ORDER_SORTS.map((o) => (
                      <option key={o.value} value={o.value}>{o.label}</option>
                    ))}
                  </Select>
                </div>
              </div>

              <div className="space-y-4">
                {(data.orders || []).map((order) => (
                  <OrderCard key={order.id} order={order} />
                ))}
                {(!data.orders || data.orders.length === 0) && <EmptyBlock text="Заказов по текущим фильтрам нет." />}
              </div>

              <div className="mt-5 flex flex-col gap-3 border-t border-cream-200 pt-4 sm:flex-row sm:items-center sm:justify-between">
                <p className="text-sm text-brand-600">
                  Показано {(data.orders || []).length} из {data.meta?.pagination?.total ?? 0}
                </p>
                <div className="flex gap-2">
                  <Button
                    variant="secondary"
                    disabled={statsStore.loading || (data.meta?.pagination?.offset ?? 0) <= 0}
                    onClick={() => handlePageChange(Math.max(0, (data.meta?.pagination?.offset ?? 0) - (data.meta?.pagination?.limit ?? 20)))}
                  >
                    Назад
                  </Button>
                  <Button
                    variant="secondary"
                    disabled={statsStore.loading || !data.meta?.pagination?.has_more}
                    onClick={() => handlePageChange((data.meta?.pagination?.offset ?? 0) + (data.meta?.pagination?.limit ?? 20))}
                  >
                    Вперёд
                  </Button>
                </div>
              </div>
            </section>
          </>
        )}
      </div>
    </PartnerLayout>
  )
}

function MetricCard({ icon: Icon, label, value, hint }) {
  return (
    <article className="rounded-2xl border border-cream-200 bg-white p-4 shadow-sm">
      <div className="mb-3 flex h-10 w-10 items-center justify-center rounded-xl bg-brand-100">
        <Icon size={18} className="text-brand-600" />
      </div>
      <p className="text-sm text-brand-600">{label}</p>
      <p className="mt-1 text-2xl font-bold text-brand-950">{value}</p>
      <p className="mt-2 text-xs text-brand-500">{hint}</p>
    </article>
  )
}

function TimelineChart({ timeline }) {
  if (!timeline.length) {
    return <EmptyBlock text="Нет данных для построения графика." />
  }

  const chartData = timeline.map((item) => ({
    date: formatDate(item.date),
    orders: item.orders_total ?? 0,
    revenue: item.gross_revenue ?? 0,
    net: item.net_revenue ?? 0,
  }))

  return (
    <ResponsiveContainer width="100%" height={240}>
      <ComposedChart data={chartData} margin={{ top: 4, right: 8, left: 0, bottom: 0 }}>
        <CartesianGrid strokeDasharray="3 3" stroke="#f0ece8" />
        <XAxis
          dataKey="date"
          tick={{ fontSize: 11, fill: '#7c6f62' }}
          tickLine={false}
          axisLine={false}
          interval="preserveStartEnd"
        />
        <YAxis
          yAxisId="left"
          tick={{ fontSize: 11, fill: '#7c6f62' }}
          tickLine={false}
          axisLine={false}
          width={40}
        />
        <YAxis
          yAxisId="right"
          orientation="right"
          tick={{ fontSize: 11, fill: '#7c6f62' }}
          tickLine={false}
          axisLine={false}
          width={56}
          tickFormatter={(v) => `${v} ₽`}
        />
        <Tooltip
          contentStyle={{ borderRadius: 12, border: '1px solid #e8e2dc', fontSize: 13 }}
          formatter={(value, name) => {
            if (name === 'Выручка' || name === 'Чистыми') return [`${value} ₽`, name]
            return [value, name]
          }}
        />
        <Legend wrapperStyle={{ fontSize: 12, paddingTop: 8 }} />
        <Bar yAxisId="left" dataKey="orders" name="Заказы" fill="#c7d2fe" radius={[4, 4, 0, 0]} maxBarSize={32} />
        <Area yAxisId="right" type="monotone" dataKey="revenue" name="Выручка" stroke="#6366f1" fill="#eef2ff" strokeWidth={2} dot={false} />
        <Area yAxisId="right" type="monotone" dataKey="net" name="Чистыми" stroke="#22c55e" fill="none" strokeWidth={1.5} strokeDasharray="4 2" dot={false} />
      </ComposedChart>
    </ResponsiveContainer>
  )
}

function StatusPieChart({ breakdown }) {
  if (!breakdown.length) {
    return <EmptyBlock text="За выбранный период нет заказов по этому фильтру." />
  }

  const chartData = breakdown.map((item) => ({
    name: getOrderStatusMeta(item.status).label,
    value: item.count,
    share: Math.round((item.share || 0) * 100),
  }))

  return (
    <div className="flex flex-col items-center gap-4">
      <ResponsiveContainer width="100%" height={200}>
        <PieChart>
          <Pie
            data={chartData}
            cx="50%"
            cy="50%"
            innerRadius={55}
            outerRadius={85}
            paddingAngle={2}
            dataKey="value"
          >
            {chartData.map((_, i) => (
              <Cell key={i} fill={STATUS_PIE_COLORS[i % STATUS_PIE_COLORS.length]} />
            ))}
          </Pie>
          <Tooltip content={<CustomTooltip />} />
        </PieChart>
      </ResponsiveContainer>
      <div className="flex w-full flex-wrap justify-center gap-x-4 gap-y-2">
        {chartData.map((d, i) => (
          <div key={d.name} className="flex items-center gap-1.5 text-xs text-brand-700">
            <span className="h-2.5 w-2.5 rounded-full" style={{ background: STATUS_PIE_COLORS[i % STATUS_PIE_COLORS.length] }} />
            {d.name}: {d.share}%
          </div>
        ))}
      </div>
    </div>
  )
}

function LocationCard({ item }) {
  return (
    <div className="rounded-xl border border-cream-200 bg-white p-4">
      <div className="flex items-start justify-between gap-3">
        <div>
          <p className="font-semibold text-brand-900">{item.name}</p>
          <p className="mt-1 flex items-center gap-1 text-sm text-brand-600">
            <MapPin size={13} />
            {item.address}
          </p>
        </div>
        <span className="rounded-full bg-brand-100 px-3 py-1 text-xs font-semibold text-brand-700">{item.orders_total} заказов</span>
      </div>
      <div className="mt-3 grid gap-2 sm:grid-cols-3">
        <MiniInfo icon={Wallet} label="Выручка" value={`${item.gross_revenue} ₽`} />
        <MiniInfo icon={ShoppingBag} label="Выдано" value={item.orders_completed} />
        <MiniInfo icon={Star} label="Рейтинг" value={item.avg_rating?.toFixed?.(1) ?? '0.0'} />
      </div>
    </div>
  )
}

function BoxCard({ item }) {
  return (
    <div className="rounded-xl border border-cream-200 bg-white p-4">
      <div className="flex items-start justify-between gap-3">
        <div>
          <p className="font-semibold text-brand-900">{item.name}</p>
          <p className="mt-1 flex items-center gap-1 text-sm text-brand-600">
            <Store size={13} />
            {item.location_name}
          </p>
        </div>
        <span className="rounded-full bg-brand-100 px-3 py-1 text-xs font-semibold text-brand-700">{item.orders_total} заказов</span>
      </div>
      <div className="mt-3 grid gap-2 sm:grid-cols-3">
        <MiniInfo icon={Wallet} label="Выручка" value={`${item.gross_revenue} ₽`} />
        <MiniInfo icon={Package} label="Выдано" value={item.orders_completed} />
        <MiniInfo icon={BarChart3} label="Чистыми" value={`${item.net_revenue} ₽`} />
      </div>
    </div>
  )
}

function OrderCard({ order }) {
  const statusMeta = getOrderStatusMeta(order.status)
  return (
    <article className="rounded-2xl border border-cream-200 bg-white p-4">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <p className="text-xs uppercase tracking-[0.22em] text-cream-500">Код получения</p>
          <p className="mt-1 text-lg font-semibold tracking-[0.12em] text-brand-950">{order.pickup_code}</p>
        </div>
        <div className="text-right">
          <span className={`badge ${statusMeta.className}`}>{statusMeta.label}</span>
          <p className="mt-2 text-lg font-bold text-brand-950">{Math.round(order.amount || 0)} ₽</p>
        </div>
      </div>
      <div className="mt-4 grid gap-3 md:grid-cols-2 xl:grid-cols-4">
        <MiniInfo icon={Package} label="Бокс" value={order.box_name || '—'} />
        <MiniInfo icon={Store} label="Локация" value={order.location_name || '—'} />
        <MiniInfo icon={ShoppingBag} label="Клиент" value={order.customer_name || order.customer_phone || '—'} />
        <MiniInfo icon={CalendarRange} label="Выдача" value={`${formatDateTime(order.pickup_time_start)} — ${formatDateTime(order.pickup_time_end)}`} />
      </div>
      <p className="mt-3 text-xs text-brand-500">Создан {formatDateTime(order.created_at)}</p>
    </article>
  )
}

function MiniInfo({ icon: Icon, label, value }) {
  return (
    <div className="rounded-xl border border-cream-200 bg-cream-50 px-3 py-2.5">
      <p className="mb-1 text-xs uppercase tracking-wider text-cream-500">{label}</p>
      <p className="flex items-center gap-2 text-sm font-medium text-brand-800">
        <Icon size={14} className="text-brand-500" />
        {value}
      </p>
    </div>
  )
}

function EmptyBlock({ text }) {
  return (
    <div className="rounded-xl border border-dashed border-cream-300 bg-cream-50 px-4 py-8 text-center text-sm text-brand-600">
      {text}
    </div>
  )
}

export default observer(StatsPageBase)
