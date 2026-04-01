import { makeAutoObservable, runInAction } from 'mobx'
import { getOrderByPickupCode, listPendingConfirmationOrders, pickupOrder } from '@/api/partner'

class OrdersStore {
  pending = []
  current = null
  loading = false
  lookupLoading = false
  pickupLoading = false
  error = null

  constructor() {
    makeAutoObservable(this, {}, { autoBind: true })
  }

  async loadPending() {
    this.loading = true
    this.error = null
    try {
      const data = await listPendingConfirmationOrders()
      const items = Array.isArray(data?.items) ? data.items : data
      runInAction(() => {
        this.pending = items || []
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

  async lookupByCode(code) {
    this.lookupLoading = true
    this.error = null
    try {
      const data = await getOrderByPickupCode(code)
      runInAction(() => {
        this.current = data
      })
      return data
    } catch (error) {
      runInAction(() => {
        this.current = null
        this.error = error
      })
      throw error
    } finally {
      runInAction(() => {
        this.lookupLoading = false
      })
    }
  }

  async pickup(orderId) {
    this.pickupLoading = true
    try {
      const data = await pickupOrder(orderId)
      runInAction(() => {
        if (this.current) {
          this.current = {
            ...this.current,
            status: data?.status ?? 'picked_up',
          }
        }
      })
      return data
    } finally {
      runInAction(() => {
        this.pickupLoading = false
      })
    }
  }
}

export const ordersStore = new OrdersStore()
