import api from './client'

export const createApplication = (data) =>
  api.post('/applications', data).then((r) => r.data.data)

export const listApplications = () =>
  api.get('/applications').then((r) => r.data.data)

export const getApplication = (id) =>
  api.get(`/applications/${id}`).then((r) => r.data.data)

export const approveApplication = (id) =>
  api.post(`/applications/${id}/approve`).then((r) => r.data.data)

export const rejectApplication = (id, rejection_reason) =>
  api.post(`/applications/${id}/reject`, { rejection_reason }).then((r) => r.data.data)

export const deleteApplication = (id) =>
  api.delete(`/applications/${id}`).then((r) => r.data.data)
