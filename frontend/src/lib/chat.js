const ACTIVE_CHAT_STATUSES = new Set(['confirmed', 'picked_up'])

export function buildChatWebSocketUrl({ baseUrl, orderId, token, origin = window.location.origin }) {
  const resolved = new URL(`${baseUrl.replace(/\/$/, '')}/orders/${encodeURIComponent(orderId)}`, origin)

  if (resolved.protocol === 'http:') {
    resolved.protocol = 'ws:'
  } else if (resolved.protocol === 'https:') {
    resolved.protocol = 'wss:'
  }

  resolved.searchParams.set('token', token)
  return resolved.toString()
}

export function mergeMessages(existing = [], incoming = []) {
  const byID = new Map()

  for (const item of [...existing, ...incoming]) {
    if (item?.id) {
      byID.set(item.id, item)
    }
  }

  return [...byID.values()].sort((left, right) => {
    const leftTime = new Date(left.created_at || 0).getTime()
    const rightTime = new Date(right.created_at || 0).getTime()
    return leftTime - rightTime
  })
}

export function getLatestMessageId(messages = []) {
  return mergeMessages([], messages).at(-1)?.id || null
}

export function isOrderChatClosed(status) {
  return !ACTIVE_CHAT_STATUSES.has(status)
}
