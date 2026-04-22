import { makeAutoObservable, runInAction } from 'mobx'
import { getPartnerStats } from '@/api/partner'

const DEFAULT_FILTERS = {
  period: 'last_7_days',
  date_from: '',
  date_to: '',
  location_id: '',
  status: '',
  top_locations_sort: 'revenue_desc',
  top_boxes_sort: 'revenue_desc',
  orders_sort: 'created_at_desc',
  limit: 20,
  offset: 0,
}

class StatsStore {
  filters = { ...DEFAULT_FILTERS }
  data = null
  loading = false
  error = null

  constructor() {
    makeAutoObservable(this, {}, { autoBind: true })
  }

  setPeriod(period) {
    this.filters.period = period
    if (period !== 'custom') {
      this.filters.date_from = ''
      this.filters.date_to = ''
    }
    this.filters.offset = 0
  }

  setCustomRange(dateFrom, dateTo) {
    this.filters.period = 'custom'
    this.filters.date_from = dateFrom
    this.filters.date_to = dateTo
    this.filters.offset = 0
  }

  setDraftDateFrom(value) {
    this.filters.date_from = value
  }

  setDraftDateTo(value) {
    this.filters.date_to = value
  }

  setLocation(locationId) {
    this.filters.location_id = locationId || ''
    this.filters.offset = 0
  }

  setStatus(status) {
    this.filters.status = status || ''
    this.filters.offset = 0
  }

  setTopLocationsSort(value) {
    this.filters.top_locations_sort = value
  }

  setTopBoxesSort(value) {
    this.filters.top_boxes_sort = value
  }

  setOrdersSort(value) {
    this.filters.orders_sort = value
    this.filters.offset = 0
  }

  setPage(offset) {
    this.filters.offset = Math.max(0, offset)
  }

  get page() {
    return Math.floor(this.filters.offset / this.filters.limit) + 1
  }

  get totalPages() {
    const total = this.data?.meta?.pagination?.total ?? 0
    return Math.max(1, Math.ceil(total / this.filters.limit))
  }

  async load() {
    this.loading = true
    this.error = null
    try {
      const data = await getPartnerStats(this.filters)
      runInAction(() => {
        this.data = data
      })
      return data
    } catch (error) {
      runInAction(() => {
        this.error = error
      })
      throw error
    } finally {
      runInAction(() => {
        this.loading = false
      })
    }
  }
}

export const statsStore = new StatsStore()
