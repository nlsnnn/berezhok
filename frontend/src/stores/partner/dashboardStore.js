import { makeAutoObservable, runInAction } from 'mobx'
import { getPartnerDashboard } from '@/api/partner'

class DashboardStore {
  data = null
  loading = false
  error = null

  constructor() {
    makeAutoObservable(this, {}, { autoBind: true })
  }

  async load() {
    this.loading = true
    this.error = null
    try {
      const response = await getPartnerDashboard()
      runInAction(() => {
        this.data = response
      })
    } catch (error) {
      runInAction(() => {
        this.error = error
      })
    } finally {
      runInAction(() => {
        this.loading = false
      })
    }
  }
}

export const dashboardStore = new DashboardStore()
