import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import { useAuth } from '@/context/AuthContext'

function RequireAuthBase() {
  const { isAuthenticated, user } = useAuth()
  const location = useLocation()

  const isEmployeeAllowedPath =
    location.pathname === '/partner/change-password' ||
    location.pathname.startsWith('/partner/orders')

  if (!isAuthenticated) {
    return <Navigate to="/partner/login" replace />
  }

  if (user?.must_change_password && location.pathname !== '/partner/change-password') {
    return <Navigate to="/partner/change-password" replace />
  }

  if (user?.role === 'employee' && !isEmployeeAllowedPath) {
    return <Navigate to="/partner/orders/pickup" replace />
  }

  return <Outlet />
}

export default observer(RequireAuthBase)
