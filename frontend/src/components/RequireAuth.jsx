import { Navigate, Outlet } from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import { useAuth } from '@/context/AuthContext'

function RequireAuthBase() {
  const { isAuthenticated, user } = useAuth()

  if (!isAuthenticated) {
    return <Navigate to="/partner/login" replace />
  }

  if (user?.must_change_password && window.location.pathname !== '/partner/change-password') {
    return <Navigate to="/partner/change-password" replace />
  }

  return <Outlet />
}

export default observer(RequireAuthBase)
