import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import { Toaster } from 'sonner'
import { AuthProvider } from '@/context/AuthContext'
import RequireAuth from '@/components/RequireAuth'
import RequireAdminAuth from '@/components/RequireAdminAuth'

import LandingPage from '@/pages/landing/LandingPage'
import AdminApplicationsPage from '@/pages/admin/AdminApplicationsPage'
import AdminLoginPage from '@/pages/admin/AdminLoginPage'
import AdminAdminsPage from '@/pages/admin/AdminAdminsPage'
import AdminAuditPage from '@/pages/admin/AdminAuditPage'
import {
  AdminBoxesPage,
  AdminCustomersPage,
  AdminLocationsPage,
  AdminOrdersPage,
  AdminPartnersPage,
  AdminPaymentsPage,
} from '@/pages/admin/AdminDirectoryPages'
import AdminStatsPage from '@/pages/admin/AdminStatsPage'
import PartnerLoginPage from '@/pages/partner/PartnerLoginPage'
import PartnerDashboard from '@/pages/partner/PartnerDashboard'
import ChangePasswordPage from '@/pages/partner/ChangePasswordPage'
import CreateLocationPage from '@/pages/partner/CreateLocationPage'
import LocationsPage from '@/pages/partner/LocationsPage'
import BoxesPage from '@/pages/partner/BoxesPage'
import CreateBoxPage from '@/pages/partner/CreateBoxPage'
import EditBoxPage from '@/pages/partner/EditBoxPage'
import OrderPickupPage from '@/pages/partner/OrderPickupPage'
import OrdersPage from '@/pages/partner/OrdersPage'
import EmployeesPage from '@/pages/partner/EmployeesPage'
import StatsPage from '@/pages/partner/StatsPage'
import ProfilePage from '@/pages/partner/ProfilePage'
import PartnerLegalInfoPage from '@/pages/partner/PartnerLegalInfoPage'
import PartnerOfferPage from '@/pages/partner/PartnerOfferPage'
import PartnerChatPage from '@/pages/partner/PartnerChatPage'
import PromoPage from '@/pages/partner/PromoPage'

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Toaster position="top-right" richColors closeButton />
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route path="/admin/login" element={<AdminLoginPage />} />
          <Route path="/partner/login" element={<PartnerLoginPage />} />
          <Route path="/partner/offer" element={<PartnerOfferPage />} />

            {/* Admin — protected */}
            <Route element={<RequireAdminAuth />}>
              <Route path="/admin" element={<Navigate to="/admin/applications" replace />} />
              <Route path="/admin/applications" element={<AdminApplicationsPage />} />
              <Route path="/admin/partners" element={<AdminPartnersPage />} />
              <Route path="/admin/locations" element={<AdminLocationsPage />} />
              <Route path="/admin/boxes" element={<AdminBoxesPage />} />
              <Route path="/admin/customers" element={<AdminCustomersPage />} />
              <Route path="/admin/orders" element={<AdminOrdersPage />} />
              <Route path="/admin/payments" element={<AdminPaymentsPage />} />
              <Route path="/admin/stats" element={<AdminStatsPage />} />
              <Route path="/admin/audit" element={<AdminAuditPage />} />
              <Route path="/admin/admins" element={<AdminAdminsPage />} />
            </Route>

            {/* Partner — protected */}
            <Route element={<RequireAuth />}>
              <Route path="/partner/dashboard" element={<PartnerDashboard />} />
              <Route path="/partner/change-password" element={<ChangePasswordPage />} />

              {/* Locations */}
              <Route path="/partner/locations" element={<LocationsPage />} />
              <Route path="/partner/locations/new" element={<CreateLocationPage />} />

              {/* Boxes */}
              <Route path="/partner/boxes" element={<BoxesPage />} />
              <Route path="/partner/boxes/new" element={<CreateBoxPage />} />
              <Route path="/partner/boxes/:id/edit" element={<EditBoxPage />} />

              {/* Orders */}
              <Route path="/partner/orders" element={<OrdersPage />} />
              <Route path="/partner/orders/pickup" element={<OrderPickupPage />} />

              {/* Chat */}
              <Route path="/partner/chat" element={<PartnerChatPage />} />

              {/* Promo */}
              <Route path="/partner/promo" element={<PromoPage />} />

              {/* Employees */}
              <Route path="/partner/employees" element={<EmployeesPage />} />

              {/* Stats */}
              <Route path="/partner/stats" element={<StatsPage />} />

              {/* Profile */}
              <Route path="/partner/profile" element={<ProfilePage />} />
              <Route path="/partner/legal-info" element={<PartnerLegalInfoPage />} />
            </Route>

          <Route path="/partner" element={<Navigate to="/partner/dashboard" replace />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  )
}
