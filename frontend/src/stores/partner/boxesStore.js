import { makeAutoObservable, runInAction } from 'mobx'
import { createBox, deleteBox, getBoxById, listBoxes, updateBox } from '@/api/partner'

class BoxesStore {
  items = []
  current = null
  loading = false
  submitting = false
  error = null

  constructor() {
    makeAutoObservable(this, {}, { autoBind: true })
  }

  async load() {
    this.loading = true
    this.error = null
    try {
      const data = await listBoxes()
      const items = Array.isArray(data?.items) ? data.items : data
      runInAction(() => {
        this.items = items || []
      })
    } catch (error) {
      runInAction(() => {
        this.error = error
      })
    } finally {
      runInAction(() => {
        this.loading = false
      })
    }
  }

  async loadById(id) {
    this.loading = true
    this.error = null
    try {
      const box = await getBoxById(id)
      runInAction(() => {
        this.current = box
      })
      return box
    } catch (error) {
      runInAction(() => {
        this.error = error
        this.current = null
      })
      throw error
    } finally {
      runInAction(() => {
        this.loading = false
      })
    }
  }

  async create(payload) {
    this.submitting = true
    try {
      const created = await createBox(payload)
      await this.load()
      return created
    } finally {
      runInAction(() => {
        this.submitting = false
      })
    }
  }

  async update(id, payload) {
    this.submitting = true
    try {
      const updated = await updateBox(id, payload)
      await this.load()
      return updated
    } finally {
      runInAction(() => {
        this.submitting = false
      })
    }
  }

  async remove(id) {
    this.submitting = true
    try {
      await deleteBox(id)
      runInAction(() => {
        this.items = this.items.filter((item) => item.id !== id)
      })
    } finally {
      runInAction(() => {
        this.submitting = false
      })
    }
  }
}

export const boxesStore = new BoxesStore()
