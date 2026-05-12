import { makeAutoObservable, runInAction } from 'mobx'
import { createEmployee, deleteEmployee, getPartnerProfile, listEmployees } from '@/api/partner'

class EmployeesStore {
  items = []
  locations = []
  loading = false
  submitting = false
  deletingId = null
  error = null

  constructor() {
    makeAutoObservable(this, {}, { autoBind: true })
  }

  async load() {
    this.loading = true
    this.error = null

    try {
      const [employees, profile] = await Promise.all([
        listEmployees(),
        getPartnerProfile(),
      ])

      runInAction(() => {
        this.items = Array.isArray(employees) ? employees : []
        this.locations = Array.isArray(profile?.locations) ? profile.locations : []
      })
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

  async create(payload) {
    this.submitting = true
    try {
      const created = await createEmployee(payload)
      await this.load()
      return created
    } finally {
      runInAction(() => {
        this.submitting = false
      })
    }
  }

  async remove(employeeId) {
    this.deletingId = employeeId
    try {
      const result = await deleteEmployee(employeeId)
      runInAction(() => {
        this.items = this.items.filter((item) => item.id !== employeeId)
      })
      return result
    } finally {
      runInAction(() => {
        this.deletingId = null
      })
    }
  }
}

export const employeesStore = new EmployeesStore()
