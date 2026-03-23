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

// Boxes
export const createBox = (data) =>
  api.post('/partner/boxes', data).then((r) => r.data.data)

export const getBoxById = (id) =>
  api.get(`/partner/boxes/${id}`).then((r) => r.data.data)

export const updateBox = (id, data) =>
  api.put(`/partner/boxes/${id}`, data).then((r) => r.data.data)

export const deleteBox = (id) =>
  api.delete(`/partner/boxes/${id}`).then((r) => r.data.data)

export const listBoxes = () =>
  api.get('/partner/boxes').then((r) => r.data.data)

export const listLocationBoxes = (locationId) =>
  api.get(`/locations/${locationId}/boxes`).then((r) => r.data.data)

// Media
export const uploadMedia = (file) => {
  const formData = new FormData()
  formData.append('file', file)
  return api.post('/media/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  }).then((r) => r.data.data)
}
