import { makeAutoObservable, runInAction } from 'mobx'
import { approveApplication, listApplications, rejectApplication } from '@/api/applications'

class AdminApplicationsStore {
  items = []
  loading = false
  actionLoading = false
  error = null

  constructor() {
    makeAutoObservable(this, {}, { autoBind: true })
  }

  async load() {
    this.loading = true
    this.error = null
    try {
      const data = await listApplications()
      const items = Array.isArray(data?.items) ? data.items : data
      runInAction(() => {
        this.items = items || []
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

  async approve(id) {
    this.actionLoading = true
    try {
      const data = await approveApplication(id)
      await this.load()
      return data
    } finally {
      runInAction(() => {
        this.actionLoading = false
      })
    }
  }

  async reject(id, reason) {
    this.actionLoading = true
    try {
      const data = await rejectApplication(id, reason)
      await this.load()
      return data
    } finally {
      runInAction(() => {
        this.actionLoading = false
      })
    }
  }
}

export const adminApplicationsStore = new AdminApplicationsStore()
