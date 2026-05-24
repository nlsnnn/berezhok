import adminApi from './adminClient'

const dataOf = (request) => request.then((response) => response.data.data)

export const adminLogin = (email, password) =>
  dataOf(adminApi.post('/auth/admin/login', { email, password }))

export const getAdminMe = () =>
  dataOf(adminApi.get('/admin/me'))

export const listAdminApplications = (params = {}) =>
  dataOf(adminApi.get('/admin/applications', { params }))

export const getAdminApplication = (id) =>
  dataOf(adminApi.get(`/admin/applications/${id}`))

export const approveAdminApplication = (id) =>
  dataOf(adminApi.post(`/admin/applications/${id}/approve`))

export const rejectAdminApplication = (id, rejection_reason) =>
  dataOf(adminApi.post(`/admin/applications/${id}/reject`, { rejection_reason }))

export const deleteAdminApplication = (id) =>
  dataOf(adminApi.delete(`/admin/applications/${id}`))

export const listAdmins = (params = {}) =>
  dataOf(adminApi.get('/admin/admins', { params }))

export const createAdmin = (payload) =>
  dataOf(adminApi.post('/admin/admins', payload))

export const updateAdmin = (id, payload) =>
  dataOf(adminApi.patch(`/admin/admins/${id}`, payload))

export const deactivateAdmin = (id) =>
  dataOf(adminApi.post(`/admin/admins/${id}/deactivate`))

export const listAuditEvents = (params = {}) =>
  dataOf(adminApi.get('/admin/audit', { params }))

export const listAdminPartners = (params = {}) =>
  dataOf(adminApi.get('/admin/partners', { params }))

export const getAdminPartner = (id) =>
  dataOf(adminApi.get(`/admin/partners/${id}`))

export const updateAdminPartner = (id, payload) =>
  dataOf(adminApi.patch(`/admin/partners/${id}`, payload))

export const listAdminLocations = (params = {}) =>
  dataOf(adminApi.get('/admin/locations', { params }))

export const getAdminLocation = (id) =>
  dataOf(adminApi.get(`/admin/locations/${id}`))

export const updateAdminLocationStatus = (id, status) =>
  dataOf(adminApi.patch(`/admin/locations/${id}/status`, { status }))

export const listAdminBoxes = (params = {}) =>
  dataOf(adminApi.get('/admin/boxes', { params }))

export const getAdminBox = (id) =>
  dataOf(adminApi.get(`/admin/boxes/${id}`))

export const updateAdminBoxStatus = (id, status) =>
  dataOf(adminApi.patch(`/admin/boxes/${id}/status`, { status }))

export const listAdminCustomers = (params = {}) =>
  dataOf(adminApi.get('/admin/customers', { params }))

export const getAdminCustomer = (id) =>
  dataOf(adminApi.get(`/admin/customers/${id}`))

export const listAdminOrders = (params = {}) =>
  dataOf(adminApi.get('/admin/orders', { params }))

export const getAdminOrder = (id) =>
  dataOf(adminApi.get(`/admin/orders/${id}`))

export const listAdminPayments = (params = {}) =>
  dataOf(adminApi.get('/admin/payments', { params }))

export const getAdminPayment = (id) =>
  dataOf(adminApi.get(`/admin/payments/${id}`))

export const listAdminPaymentEvents = (id, params = {}) =>
  dataOf(adminApi.get(`/admin/payments/${id}/events`, { params }))

export const getAdminStats = () =>
  dataOf(adminApi.get('/admin/stats'))
