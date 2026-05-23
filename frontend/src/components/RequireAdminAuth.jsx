import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import { useStores } from '@/context/StoresContext'

function RequireAdminAuthBase() {
  const { adminAuthStore } = useStores()
  const location = useLocation()

  if (!adminAuthStore.isAuthenticated) {
    return <Navigate to="/admin/login" replace state={{ from: location }} />
  }

  return <Outlet />
}

export default observer(RequireAdminAuthBase)
