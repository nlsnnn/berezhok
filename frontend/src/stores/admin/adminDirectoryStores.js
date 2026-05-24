import {
  createAdmin,
  deactivateAdmin,
  getAdminBox,
  getAdminCustomer,
  getAdminLocation,
  getAdminOrder,
  getAdminPartner,
  getAdminPayment,
  getAdminStats,
  listAdminBoxes,
  listAdminCustomers,
  listAdminLocations,
  listAdminOrders,
  listAdminPartners,
  listAdminPaymentEvents,
  listAdminPayments,
  listAdmins,
  listAuditEvents,
  rejectAdminPartnerLegalInfo,
  updateAdmin,
  updateAdminBoxStatus,
  updateAdminLocationStatus,
  updateAdminPartner,
  verifyAdminPartnerLegalInfo,
} from '@/api/admin'
import { AdminResourceStore } from '@/stores/admin/createAdminResourceStore'
import { makeAutoObservable, runInAction } from 'mobx'

class AdminPartnersStore extends AdminResourceStore {
  constructor() {
    super({ listFn: listAdminPartners, detailFn: getAdminPartner, defaultFilters: { search: '', status: 'all', legal_status: 'all' } })
  }

  update(id, payload) {
    return this.runAction(async () => {
      const data = await updateAdminPartner(id, payload)
      await this.load({ offset: this.pagination.offset })
      await this.loadDetail(id)
      return data
    })
  }

  updateStatusMany(ids, status) {
    return this.runAction(async () => {
      const results = []
      for (const id of ids) {
        results.push(await updateAdminPartner(id, { status }))
      }
      await this.load({ offset: this.pagination.offset })
      return results
    })
  }

  verifyLegalInfo(id) {
    return this.runAction(async () => {
      const data = await verifyAdminPartnerLegalInfo(id)
      await this.load({ offset: this.pagination.offset })
      await this.loadDetail(id)
      return data
    })
  }

  rejectLegalInfo(id, verificationComment) {
    return this.runAction(async () => {
      const data = await rejectAdminPartnerLegalInfo(id, verificationComment)
      await this.load({ offset: this.pagination.offset })
      await this.loadDetail(id)
      return data
    })
  }
}

class AdminLocationsStore extends AdminResourceStore {
  constructor() {
    super({ listFn: listAdminLocations, detailFn: getAdminLocation, defaultFilters: { search: '', status: 'all', partner_id: '' } })
  }

  updateStatus(id, status) {
    return this.runAction(async () => {
      const data = await updateAdminLocationStatus(id, status)
      await this.load({ offset: this.pagination.offset })
      await this.loadDetail(id)
      return data
    })
  }

  updateStatusMany(ids, status) {
    return this.runAction(async () => {
      const results = []
      for (const id of ids) {
        results.push(await updateAdminLocationStatus(id, status))
      }
      await this.load({ offset: this.pagination.offset })
      return results
    })
  }
}

class AdminBoxesStore extends AdminResourceStore {
  constructor() {
    super({ listFn: listAdminBoxes, detailFn: getAdminBox, defaultFilters: { search: '', status: 'all', location_id: '' } })
  }

  updateStatus(id, status) {
    return this.runAction(async () => {
      const data = await updateAdminBoxStatus(id, status)
      await this.load({ offset: this.pagination.offset })
      await this.loadDetail(id)
      return data
    })
  }

  updateStatusMany(ids, status) {
    return this.runAction(async () => {
      const results = []
      for (const id of ids) {
        results.push(await updateAdminBoxStatus(id, status))
      }
      await this.load({ offset: this.pagination.offset })
      return results
    })
  }
}

class AdminPaymentsStore extends AdminResourceStore {
  events = []

  constructor() {
    super({ listFn: listAdminPayments, detailFn: getAdminPayment, defaultFilters: { search: '', status: 'all' } })
  }

  async loadDetail(id) {
    const payment = await super.loadDetail(id)
    const eventsData = await listAdminPaymentEvents(id, { limit: 20, offset: 0 })
    this.events = Array.isArray(eventsData?.items) ? eventsData.items : []
    return payment
  }

  clearCurrent() {
    super.clearCurrent()
    this.events = []
  }
}

class AdminAdminsStore extends AdminResourceStore {
  constructor() {
    super({ listFn: listAdmins, detailFn: null, defaultFilters: { search: '' } })
  }

  create(payload) {
    return this.runAction(async () => {
      const data = await createAdmin(payload)
      await this.load({ offset: 0 })
      return data
    })
  }

  update(id, payload) {
    return this.runAction(async () => {
      const data = await updateAdmin(id, payload)
      await this.load({ offset: this.pagination.offset })
      return data
    })
  }

  deactivate(id) {
    return this.runAction(async () => {
      const data = await deactivateAdmin(id)
      await this.load({ offset: this.pagination.offset })
      return data
    })
  }

  deactivateMany(ids) {
    return this.runAction(async () => {
      const results = []
      for (const id of ids) {
        results.push(await deactivateAdmin(id))
      }
      await this.load({ offset: this.pagination.offset })
      return results
    })
  }
}

class AdminStatsStore {
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
      const data = await getAdminStats()
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

export const adminPartnersStore = new AdminPartnersStore()
export const adminLocationsStore = new AdminLocationsStore()
export const adminBoxesStore = new AdminBoxesStore()
export const adminCustomersStore = new AdminResourceStore({
  listFn: listAdminCustomers,
  detailFn: getAdminCustomer,
  defaultFilters: { search: '' },
})
export const adminOrdersStore = new AdminResourceStore({
  listFn: listAdminOrders,
  detailFn: getAdminOrder,
  defaultFilters: { search: '', status: 'all' },
})
export const adminPaymentsStore = new AdminPaymentsStore()
export const adminAuditStore = new AdminResourceStore({
  listFn: listAuditEvents,
  detailFn: null,
  defaultFilters: { action: '', entity_type: '' },
})
export const adminAdminsStore = new AdminAdminsStore()
export const adminStatsStore = new AdminStatsStore()
