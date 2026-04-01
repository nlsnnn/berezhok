import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import { Toaster } from 'sonner'
import { AuthProvider } from '@/context/AuthContext'
import RequireAuth from '@/components/RequireAuth'

import LandingPage from '@/pages/landing/LandingPage'
import AdminPage from '@/pages/admin/AdminPage'
import PartnerLoginPage from '@/pages/partner/PartnerLoginPage'
import PartnerDashboard from '@/pages/partner/PartnerDashboard'
import ChangePasswordPage from '@/pages/partner/ChangePasswordPage'
import CreateLocationPage from '@/pages/partner/CreateLocationPage'
import LocationsPage from '@/pages/partner/LocationsPage'
import BoxesPage from '@/pages/partner/BoxesPage'
import CreateBoxPage from '@/pages/partner/CreateBoxPage'
import EditBoxPage from '@/pages/partner/EditBoxPage'
import OrderPickupPage from '@/pages/partner/OrderPickupPage'

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Toaster position="top-right" richColors closeButton />
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route path="/admin" element={<AdminPage />} />
          <Route path="/partner/login" element={<PartnerLoginPage />} />

          <Route element={<RequireAuth />}>
            <Route path="/partner/dashboard" element={<PartnerDashboard />} />
            <Route path="/partner/change-password" element={<ChangePasswordPage />} />
            <Route path="/partner/locations" element={<LocationsPage />} />
            <Route path="/partner/locations/new" element={<CreateLocationPage />} />
            <Route path="/partner/boxes" element={<BoxesPage />} />
            <Route path="/partner/boxes/new" element={<CreateBoxPage />} />
            <Route path="/partner/boxes/:id/edit" element={<EditBoxPage />} />
            <Route path="/partner/orders/pickup" element={<OrderPickupPage />} />
          </Route>

          <Route path="/partner" element={<Navigate to="/partner/dashboard" replace />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  )
}
