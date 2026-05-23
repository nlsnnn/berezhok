import {
  approveAdminApplication,
  deleteAdminApplication,
  getAdminApplication,
  listAdminApplications,
  rejectAdminApplication,
} from '@/api/admin'
import { AdminResourceStore } from '@/stores/admin/createAdminResourceStore'

class AdminApplicationsStore extends AdminResourceStore {
  constructor() {
    super({
      listFn: listAdminApplications,
      detailFn: getAdminApplication,
      defaultFilters: {
        search: '',
        status: 'all',
      },
    })
  }

  async approve(id) {
    return this.runAction(async () => {
      const data = await approveAdminApplication(id)
      await this.load({ offset: this.pagination.offset })
      return data
    })
  }

  async reject(id, reason) {
    return this.runAction(async () => {
      const data = await rejectAdminApplication(id, reason)
      await this.load({ offset: this.pagination.offset })
      return data
    })
  }

  async delete(id) {
    return this.runAction(async () => {
      const data = await deleteAdminApplication(id)
      await this.load({ offset: this.pagination.offset })
      return data
    })
  }
}

export const adminApplicationsStore = new AdminApplicationsStore()
