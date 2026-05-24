import assert from 'node:assert/strict'
import { test } from 'node:test'

import { getAllowedAdminLinks, getPrimaryAdminLinks } from './adminNav.js'

test('super admin sees admin management link', () => {
  assert.equal(getAllowedAdminLinks('super_admin').some((link) => link.to === '/admin/admins'), true)
})

test('admin and support do not see admin management link', () => {
  assert.equal(getAllowedAdminLinks('admin').some((link) => link.to === '/admin/admins'), false)
  assert.equal(getAllowedAdminLinks('support').some((link) => link.to === '/admin/admins'), false)
})

test('primary mobile links keep core operational sections first', () => {
  assert.deepEqual(
    getPrimaryAdminLinks('admin').map((link) => link.to),
    ['/admin/applications', '/admin/partners', '/admin/orders']
  )
})
