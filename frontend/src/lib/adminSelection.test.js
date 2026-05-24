import assert from 'node:assert/strict'
import { test } from 'node:test'

import { pruneSelectedIds, toggleAllPageIds, toggleSelectedId } from './adminSelection.js'

test('toggleSelectedId adds and removes ids', () => {
  assert.deepEqual(toggleSelectedId([], 'a'), ['a'])
  assert.deepEqual(toggleSelectedId(['a', 'b'], 'a'), ['b'])
})

test('toggleAllPageIds selects all visible ids or removes them when already selected', () => {
  assert.deepEqual(toggleAllPageIds(['x'], ['a', 'b']), ['x', 'a', 'b'])
  assert.deepEqual(toggleAllPageIds(['x', 'a', 'b'], ['a', 'b']), ['x'])
})

test('pruneSelectedIds keeps only ids from current page data', () => {
  assert.deepEqual(pruneSelectedIds(['a', 'b', 'c'], ['a', 'c']), ['a', 'c'])
})
