import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { test } from 'node:test'

const source = readFileSync(new URL('./PartnerChatPage.jsx', import.meta.url), 'utf8')

test('message layout keeps long unbroken text inside the chat column', () => {
  assert.match(source, /<section className="[^"]*\bmin-w-0\b/)
  assert.match(source, /'[^']*\bmin-w-0\b[^']*\bmax-w-\[85%\][^']*'/)
  assert.match(source, /className="[^"]*\[overflow-wrap:anywhere\]/)
})

test('chat page scrolls messages inside the chat panel', () => {
  assert.match(source, /className="[^"]*\bh-\[calc\(100svh-230px\)\][^"]*\boverflow-hidden\b/)
  assert.match(source, /<section className="[^"]*\bmin-h-0\b[^"]*\boverflow-hidden\b/)
  assert.match(source, /className="[^"]*\bmin-h-0\b[^"]*\bflex-1\b[^"]*\boverflow-y-auto\b/)
})
