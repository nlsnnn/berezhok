import assert from 'node:assert/strict'
import { test } from 'node:test'

import { AdminResourceStore } from './createAdminResourceStore.js'

test('AdminResourceStore can be subclassed by resource-specific stores', () => {
  class TestResourceStore extends AdminResourceStore {
    constructor() {
      super({
        listFn: async () => ({ items: [], pagination: { total: 0, limit: 20, offset: 0, has_more: false } }),
        detailFn: async (id) => ({ id }),
        defaultFilters: { search: '', status: 'all' },
      })
    }

    customAction() {
      return 'ok'
    }
  }

  const store = new TestResourceStore()

  assert.equal(store.customAction(), 'ok')
  assert.deepEqual(store.buildParams(), { limit: 20, offset: 0 })
})
