import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Attach JWT token from localStorage to every request
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('partner_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Unwrap response envelope { success, data }
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('partner_token')
      localStorage.removeItem('partner_user')
      // redirect to login if inside partner routes
      if (window.location.pathname.startsWith('/partner') && window.location.pathname !== '/partner/login') {
        window.location.href = '/partner/login'
      }
    }
    return Promise.reject(error)
  }
)

export default api
