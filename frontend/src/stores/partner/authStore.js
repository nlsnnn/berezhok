import { makeAutoObservable, runInAction } from 'mobx'
import { partnerLogin } from '@/api/partner'

class AuthStore {
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
      const stored = localStorage.getItem('partner_user')
      this.user = stored ? JSON.parse(stored) : null
    } catch {
      this.user = null
    }
  }

  async login(email, password) {
    this.loading = true
    try {
      const data = await partnerLogin(email, password)
      const payloadUser = {
        id: data?.user?.id || data?.user_id || null,
        email: data?.user?.email || email,
        role: data?.user?.role || null,
        partner_id: data?.user?.partner_id || null,
        location_id: data?.user?.location_id || null,
        must_change_password: Boolean(data?.must_change_password),
      }

      localStorage.setItem('partner_token', data.token)
      localStorage.setItem('partner_user', JSON.stringify(payloadUser))

      runInAction(() => {
        this.user = payloadUser
      })

      return data
    } finally {
      runInAction(() => {
        this.loading = false
      })
    }
  }

  logout() {
    localStorage.removeItem('partner_token')
    localStorage.removeItem('partner_user')
    this.user = null
  }

  markPasswordChanged() {
    if (!this.user) return
    this.user.must_change_password = false
    localStorage.setItem('partner_user', JSON.stringify(this.user))
  }
}

export const authStore = new AuthStore()
