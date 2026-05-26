import { makeAutoObservable, runInAction } from 'mobx'
import { getLocationPins, listAvailablePins, updateLocationPins } from '@/api/partner'

class PinsStore {
  availablePins = []
  locationPins = {} // locationId -> pin[]
  loadingAvailable = false
  savingFor = null // locationId being saved

  constructor() {
    makeAutoObservable(this, {}, { autoBind: true })
  }

  async loadAvailable() {
    if (this.availablePins.length > 0) return
    this.loadingAvailable = true
    try {
      const pins = await listAvailablePins()
      runInAction(() => {
        this.availablePins = pins || []
      })
    } finally {
      runInAction(() => {
        this.loadingAvailable = false
      })
    }
  }

  seedFromLocations(locations) {
    runInAction(() => {
      for (const loc of locations) {
        if (loc.pins !== undefined) {
          this.locationPins[loc.id] = loc.pins || []
        }
      }
    })
  }

  async loadForLocation(locationId) {
    const pins = await getLocationPins(locationId)
    runInAction(() => {
      this.locationPins[locationId] = pins || []
    })
  }

  async save(locationId, pinCodes) {
    this.savingFor = locationId
    try {
      await updateLocationPins(locationId, pinCodes)
      runInAction(() => {
        this.locationPins[locationId] = this.availablePins.filter((p) =>
          pinCodes.includes(p.code)
        )
      })
    } finally {
      runInAction(() => {
        this.savingFor = null
      })
    }
  }
}

export const pinsStore = new PinsStore()
