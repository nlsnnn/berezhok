import axios from 'axios'
import { buildChatWebSocketUrl } from '@/lib/chat'

const chatApi = axios.create({
  baseURL: import.meta.env.VITE_CHAT_API_URL || '/chat-api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

chatApi.interceptors.request.use((config) => {
  const token = localStorage.getItem('partner_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

chatApi.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('partner_token')
      localStorage.removeItem('partner_user')
      if (window.location.pathname.startsWith('/partner') && window.location.pathname !== '/partner/login') {
        window.location.href = '/partner/login'
      }
    }
    return Promise.reject(error)
  }
)

export const listOrderMessages = (orderId, { limit = 50, before } = {}) =>
  chatApi
    .get(`/orders/${orderId}/messages`, {
      params: {
        limit,
        ...(before ? { before } : {}),
      },
    })
    .then((r) => r.data.data)

export const markOrderMessagesRead = (orderId, messageId) =>
  chatApi.post(`/orders/${orderId}/read`, { message_id: messageId }).then((r) => r.data?.data)

export function createOrderChatSocket(orderId) {
  const token = localStorage.getItem('partner_token')
  if (!token) {
    throw new Error('Требуется авторизация')
  }

  return new WebSocket(
    buildChatWebSocketUrl({
      baseUrl: import.meta.env.VITE_CHAT_WS_URL || '/chat-ws',
      orderId,
      token,
    })
  )
}
