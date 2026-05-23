import { makeAutoObservable, runInAction } from 'mobx'
import { adminLogin, getAdminMe } from '@/api/admin'
import { buildAdminUser } from '@/lib/admin'

class AdminAuthStore {
  user = null
  loading = false

  constructor() {
    makeAutoObservable(this, {}, { autoBind: true })
    this.restore()
  }

  get isAuthenticated() {
    return Boolean(this.user)
  }

  restore() {
    try {
      const stored = localStorage.getItem('admin_user')
      this.user = stored ? buildAdminUser(JSON.parse(stored)) : null
    } catch {
      this.user = null
    }
  }

  async login(email, password) {
    this.loading = true
    try {
      const data = await adminLogin(email, password)
      const user = buildAdminUser(data)
      localStorage.setItem('admin_token', data.token)
      localStorage.setItem('admin_user', JSON.stringify(user))

      runInAction(() => {
        this.user = user
      })

      return data
    } finally {
      runInAction(() => {
        this.loading = false
      })
    }
  }

  async refresh() {
    if (!localStorage.getItem('admin_token')) return null

    this.loading = true
    try {
      const data = await getAdminMe()
      const user = buildAdminUser(data)
      localStorage.setItem('admin_user', JSON.stringify(user))

      runInAction(() => {
        this.user = user
      })

      return user
    } finally {
      runInAction(() => {
        this.loading = false
      })
    }
  }

  logout() {
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_user')
    this.user = null
  }
}

export const adminAuthStore = new AdminAuthStore()
