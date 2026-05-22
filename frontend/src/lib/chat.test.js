import assert from 'node:assert/strict'
import { test } from 'node:test'

import { buildChatWebSocketUrl, getLatestMessageId, isOrderChatClosed, mergeMessages } from './chat.js'

test('buildChatWebSocketUrl builds same-origin ws url with encoded token', () => {
  const url = buildChatWebSocketUrl({
    baseUrl: '/chat-ws',
    orderId: 'order-1',
    token: 'token with space',
    origin: 'http://localhost:5173',
  })

  assert.equal(url, 'ws://localhost:5173/chat-ws/orders/order-1?token=token+with+space')
})

test('buildChatWebSocketUrl keeps secure websocket protocol for https origins', () => {
  const url = buildChatWebSocketUrl({
    baseUrl: '/chat-ws',
    orderId: 'order-1',
    token: 'token',
    origin: 'https://partner.example',
  })

  assert.equal(url, 'wss://partner.example/chat-ws/orders/order-1?token=token')
})

test('mergeMessages deduplicates by id and sorts by created_at ascending', () => {
  const existing = [
    { id: '2', message: 'second', created_at: '2026-05-14T10:02:00Z' },
    { id: '1', message: 'first', created_at: '2026-05-14T10:01:00Z' },
  ]
  const incoming = [
    { id: '2', message: 'second duplicate', created_at: '2026-05-14T10:02:00Z' },
    { id: '3', message: 'third', created_at: '2026-05-14T10:03:00Z' },
  ]

  assert.deepEqual(
    mergeMessages(existing, incoming).map((item) => item.id),
    ['1', '2', '3']
  )
})

test('getLatestMessageId returns newest message id or null', () => {
  assert.equal(getLatestMessageId([]), null)
  assert.equal(
    getLatestMessageId([
      { id: 'old', created_at: '2026-05-14T10:01:00Z' },
      { id: 'new', created_at: '2026-05-14T10:02:00Z' },
    ]),
    'new'
  )
})

test('isOrderChatClosed treats non-active order statuses as closed', () => {
  assert.equal(isOrderChatClosed('confirmed'), false)
  assert.equal(isOrderChatClosed('picked_up'), false)
  assert.equal(isOrderChatClosed('completed'), true)
  assert.equal(isOrderChatClosed('cancelled'), true)
})
