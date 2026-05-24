import api from './client'

export const createApplication = (data) =>
  api.post('/applications', data).then((r) => r.data.data)
