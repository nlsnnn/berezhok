import { makeAutoObservable, runInAction } from 'mobx'

const DEFAULT_PAGINATION = {
  total: 0,
  limit: 20,
  offset: 0,
  has_more: false,
}

function cleanParams(params) {
  return Object.fromEntries(
    Object.entries(params).filter(([, value]) => value !== '' && value !== 'all' && value !== null && value !== undefined)
  )
}

export class AdminResourceStore {
  items = []
  pagination = DEFAULT_PAGINATION
  filters = {}
  current = null
  loading = false
  detailLoading = false
  actionLoading = false
  error = null

  constructor({ listFn, detailFn, defaultFilters = {}, limit = 20 }) {
    this.listFn = listFn
    this.detailFn = detailFn
    this.defaultFilters = defaultFilters
    this.pagination = { ...DEFAULT_PAGINATION, limit }
    this.filters = { ...defaultFilters }
    makeAutoObservable(this, { listFn: false, detailFn: false, defaultFilters: false }, { autoBind: true })
  }

  setFilter(name, value) {
    this.filters = {
      ...this.filters,
      [name]: value,
    }
  }

  resetFilters() {
    this.filters = { ...this.defaultFilters }
  }

  buildParams(offset = 0) {
    return cleanParams({
      ...this.filters,
      limit: this.pagination.limit,
      offset,
    })
  }

  async load({ offset = 0 } = {}) {
    this.loading = true
    this.error = null
    try {
      const data = await this.listFn(this.buildParams(offset))
      const items = Array.isArray(data?.items) ? data.items : []

      runInAction(() => {
        this.items = items
        this.pagination = {
          total: data?.pagination?.total ?? items.length,
          limit: data?.pagination?.limit ?? this.pagination.limit,
          offset: data?.pagination?.offset ?? offset,
          has_more: Boolean(data?.pagination?.has_more),
        }
      })

      return data
    } catch (error) {
      runInAction(() => {
        this.items = []
        this.pagination = { ...DEFAULT_PAGINATION, limit: this.pagination.limit }
        this.error = error
      })
      throw error
    } finally {
      runInAction(() => {
        this.loading = false
      })
    }
  }

  async loadDetail(id) {
    if (!this.detailFn) return null

    this.detailLoading = true
    try {
      const data = await this.detailFn(id)
      runInAction(() => {
        this.current = data
      })
      return data
    } finally {
      runInAction(() => {
        this.detailLoading = false
      })
    }
  }

  clearCurrent() {
    this.current = null
  }

  nextPage() {
    if (!this.pagination.has_more || this.loading) return null
    return this.load({ offset: this.pagination.offset + this.pagination.limit })
  }

  prevPage() {
    if (this.pagination.offset <= 0 || this.loading) return null
    return this.load({ offset: Math.max(0, this.pagination.offset - this.pagination.limit) })
  }

  async runAction(action) {
    this.actionLoading = true
    try {
      return await action()
    } finally {
      runInAction(() => {
        this.actionLoading = false
      })
    }
  }
}
