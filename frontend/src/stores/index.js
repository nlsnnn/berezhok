import { authStore } from '@/stores/partner/authStore'
import { dashboardStore } from '@/stores/partner/dashboardStore'
import { statsStore } from '@/stores/partner/statsStore'
import { boxesStore } from '@/stores/partner/boxesStore'
import { locationsStore } from '@/stores/partner/locationsStore'
import { ordersStore } from '@/stores/partner/ordersStore'
import { employeesStore } from '@/stores/partner/employeesStore'
import { chatStore } from '@/stores/partner/chatStore'
import { payoutsStore } from '@/stores/partner/payoutsStore'
import { pinsStore } from '@/stores/partner/pinsStore'
import { adminAuthStore } from '@/stores/admin/adminAuthStore'
import { adminApplicationsStore } from '@/stores/admin/adminApplicationsStore'
import {
  adminAdminsStore,
  adminAuditStore,
  adminBoxesStore,
  adminCustomersStore,
  adminLocationsStore,
  adminOrdersStore,
  adminPartnersStore,
  adminPaymentsStore,
  adminStatsStore,
} from '@/stores/admin/adminDirectoryStores'
import { applicationStore } from '@/stores/landing/applicationStore'

export const stores = {
  authStore,
  dashboardStore,
  statsStore,
  boxesStore,
  locationsStore,
  ordersStore,
  employeesStore,
  chatStore,
  payoutsStore,
  pinsStore,
  adminAuthStore,
  adminApplicationsStore,
  adminPartnersStore,
  adminLocationsStore,
  adminBoxesStore,
  adminCustomersStore,
  adminOrdersStore,
  adminPaymentsStore,
  adminAuditStore,
  adminAdminsStore,
  adminStatsStore,
  applicationStore,
}
