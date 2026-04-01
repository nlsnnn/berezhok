import { makeAutoObservable, runInAction } from 'mobx'
import { createLocation, getPartnerProfile } from '@/api/partner'

class LocationsStore {
  profile = null
  loading = false
  submitting = false
  error = null

  constructor() {
    makeAutoObservable(this, {}, { autoBind: true })
  }

  get locations() {
    return this.profile?.locations || []
  }

  async loadProfile() {
    this.loading = true
    this.error = null
    try {
      const profile = await getPartnerProfile()
      runInAction(() => {
        this.profile = profile
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

  async create(payload) {
    this.submitting = true
    try {
      const created = await createLocation(payload)
      await this.loadProfile()
      return created
    } finally {
      runInAction(() => {
        this.submitting = false
      })
    }
  }
}

export const locationsStore = new LocationsStore()
