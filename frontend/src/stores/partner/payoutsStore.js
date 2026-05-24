import { makeAutoObservable, runInAction } from 'mobx'
import {
  getPayoutById,
  getPayoutDestination,
  getSBPBanks,
  listPayouts,
  savePayoutDestination,
} from '@/api/partner'

class PayoutsStore {
  destination = null
  banks = []
  items = []
  current = null
  total = 0
  loading = false
  banksLoading = false
  submitting = false
  error = null

  constructor() {
    makeAutoObservable(this, {}, { autoBind: true })
  }

  async loadDestination() {
    this.loading = true
    this.error = null
    try {
      const data = await getPayoutDestination()
      runInAction(() => {
        this.destination = data
      })
    } catch (error) {
      runInAction(() => {
        if (error?.response?.status === 404) {
          this.destination = null
        } else {
          this.error = error
        }
      })
    } finally {
      runInAction(() => {
        this.loading = false
      })
    }
  }

  async loadBanks() {
    if (this.banks.length > 0) return
    this.banksLoading = true
    try {
      const data = await getSBPBanks()
      runInAction(() => {
        this.banks = data || []
      })
    } catch {
      // non-critical — user can still type bank ID manually if needed
    } finally {
      runInAction(() => {
        this.banksLoading = false
      })
    }
  }

  async saveDestination(payload) {
    this.submitting = true
    this.error = null
    try {
      const data = await savePayoutDestination(payload)
      runInAction(() => {
        this.destination = data
      })
      return data
    } catch (error) {
      runInAction(() => {
        this.error = error
      })
      throw error
    } finally {
      runInAction(() => {
        this.submitting = false
      })
    }
  }

  async loadHistory(params = {}) {
    this.loading = true
    this.error = null
    try {
      const data = await listPayouts(params)
      runInAction(() => {
        this.items = data?.payouts || []
        this.total = data?.total || 0
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

  async loadById(id) {
    this.loading = true
    this.error = null
    try {
      const data = await getPayoutById(id)
      runInAction(() => {
        this.current = data
      })
      return data
    } catch (error) {
      runInAction(() => {
        this.error = error
        this.current = null
      })
      throw error
    } finally {
      runInAction(() => {
        this.loading = false
      })
    }
  }
}

export const payoutsStore = new PayoutsStore()
