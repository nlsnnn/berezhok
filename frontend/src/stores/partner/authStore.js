import { makeAutoObservable, runInAction } from 'mobx'
import { partnerLogin } from '@/api/partner'
import {
  hasAcceptedCurrentPartnerOffer,
  PARTNER_OFFER_VERSION,
  writeAcceptedPartnerOfferVersion,
} from '@/lib/partnerOffer'

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
      this.user = stored ? this.buildUser(JSON.parse(stored)) : null
    } catch {
      this.user = null
    }
  }

  buildUser(rawUser = {}) {
    const partnerId = rawUser?.partner_id || null

    return {
      id: rawUser?.id || rawUser?.employee_id || rawUser?.user_id || null,
      employee_id: rawUser?.employee_id || rawUser?.user_id || null,
      email: rawUser?.email || null,
      role: rawUser?.role || null,
      partner_id: partnerId,
      location_id: rawUser?.location_id || null,
      must_change_password: Boolean(rawUser?.must_change_password),
      accepted_offer_version: hasAcceptedCurrentPartnerOffer(partnerId) ? PARTNER_OFFER_VERSION : null,
    }
  }

  async login(email, password) {
    this.loading = true
    try {
      const data = await partnerLogin(email, password)
      const payloadUser = this.buildUser({
        id: data?.employee_id || data?.user?.id || data?.user_id || null,
        employee_id: data?.employee_id || data?.user_id || null,
        email: data?.user?.email || email,
        role: data?.role || data?.user?.role || null,
        partner_id: data?.partner_id || data?.user?.partner_id || null,
        location_id: data?.location_id || data?.user?.location_id || null,
        must_change_password: Boolean(data?.must_change_password),
      })

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
    this.user = {
      ...this.user,
      must_change_password: false,
    }
    localStorage.setItem('partner_user', JSON.stringify(this.user))
  }

  markOfferAccepted() {
    if (!this.user?.partner_id) return

    writeAcceptedPartnerOfferVersion(this.user.partner_id)
    this.user = {
      ...this.user,
      accepted_offer_version: PARTNER_OFFER_VERSION,
    }
    localStorage.setItem('partner_user', JSON.stringify(this.user))
  }
}

export const authStore = new AuthStore()
