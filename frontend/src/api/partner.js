import api from './client'

export const partnerLogin = (email, password) =>
  api.post('/partner/auth/login', { email, password }).then((r) => r.data.data)

export const getPartnerProfile = () =>
  api.get('/partner/profile').then((r) => r.data.data)

export const changePassword = (current_password, new_password) =>
  api.post('/partner/change-password', { current_password, new_password }).then((r) => r.data.data)

export const createLocation = (data) =>
  api.post('/partner/locations', data).then((r) => r.data.data)

export const listLocations = () =>
  api.get('/partner/locations').then((r) => r.data.data)
