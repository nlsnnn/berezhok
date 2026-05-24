import assert from 'node:assert/strict'
import { test } from 'node:test'

import {
  ADMIN_ROLE_LABELS,
  buildAdminUser,
  canManageAdmins,
  canMutateOperations,
} from './admin.js'

test('buildAdminUser normalizes admin login response', () => {
  assert.deepEqual(
    buildAdminUser({
      user_id: 'admin-1',
      email: 'admin@berezhok.local',
      name: 'Главный администратор',
      role: 'super_admin',
    }),
    {
      id: 'admin-1',
      email: 'admin@berezhok.local',
      name: 'Главный администратор',
      role: 'super_admin',
    }
  )
})

test('buildAdminUser falls back to safe defaults for optional fields', () => {
  assert.deepEqual(buildAdminUser({ id: 'admin-2', email: 'support@example.com' }), {
    id: 'admin-2',
    email: 'support@example.com',
    name: 'Администратор',
    role: 'support',
  })
})

test('admin permissions match role boundaries', () => {
  assert.equal(canManageAdmins({ role: 'super_admin' }), true)
  assert.equal(canManageAdmins({ role: 'admin' }), false)
  assert.equal(canMutateOperations({ role: 'super_admin' }), true)
  assert.equal(canMutateOperations({ role: 'admin' }), true)
  assert.equal(canMutateOperations({ role: 'support' }), false)
})

test('admin role labels are Russian UI text', () => {
  assert.equal(ADMIN_ROLE_LABELS.super_admin, 'Суперадминистратор')
  assert.equal(ADMIN_ROLE_LABELS.admin, 'Администратор')
  assert.equal(ADMIN_ROLE_LABELS.support, 'Поддержка')
})
