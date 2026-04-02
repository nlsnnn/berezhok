import { makeAutoObservable, runInAction } from 'mobx'
import { createApplication } from '@/api/applications'

class ApplicationStore {
  submitting = false

  constructor() {
    makeAutoObservable(this, {}, { autoBind: true })
  }

  async create(payload) {
    this.submitting = true
    try {
      return await createApplication(payload)
    } finally {
      runInAction(() => {
        this.submitting = false
      })
    }
  }
}

export const applicationStore = new ApplicationStore()
