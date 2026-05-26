import { makeAutoObservable, runInAction } from 'mobx'
import {
  createLocation,
  getLocationById,
  getPartnerProfile,
  listLocations,
  updateLocation,
} from '@/api/partner'
import { pinsStore } from './pinsStore'

class LocationsStore {
  profile = null
  items = []
  current = null
  loading = false
  submitting = false
  error = null

  constructor() {
    makeAutoObservable(this, {}, { autoBind: true })
  }

  get locations() {
    return this.items
  }

  async load() {
    this.loading = true
    this.error = null
    try {
      const locations = await listLocations()
      runInAction(() => {
        this.items = locations
      })
      pinsStore.seedFromLocations(locations)
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

  async loadById(id) {
    this.loading = true
    this.error = null
    try {
      const location = await getLocationById(id)
      runInAction(() => {
        this.current = location
      })
      return location
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

  async create(payload) {
    this.submitting = true
    try {
      const created = await createLocation(payload)
      await this.load()
      return created
    } finally {
      runInAction(() => {
        this.submitting = false
      })
    }
  }

  async update(id, payload) {
    this.submitting = true
    try {
      const updated = await updateLocation(id, payload)
      runInAction(() => {
        this.current = updated
      })
      return updated
    } finally {
      runInAction(() => {
        this.submitting = false
      })
    }
  }
}

export const locationsStore = new LocationsStore()
