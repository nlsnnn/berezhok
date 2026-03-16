import { Navigate, Outlet } from 'react-router-dom'
import { useAuth } from '@/context/AuthContext'

export default function RequireAuth() {
  const { isAuthenticated, user } = useAuth()

  if (!isAuthenticated) {
    return <Navigate to="/partner/login" replace />
  }

  // If must_change_password redirect to change-password page
  if (user?.must_change_password && window.location.pathname !== '/partner/change-password') {
    return <Navigate to="/partner/change-password" replace />
  }

  return <Outlet />
}
