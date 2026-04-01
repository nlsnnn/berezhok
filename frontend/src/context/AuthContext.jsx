import { createContext, useContext } from 'react'
import { observer } from 'mobx-react-lite'
import { useStores } from '@/context/StoresContext'

const AuthContext = createContext(null)

const AuthProviderBase = ({ children }) => {
  const { authStore } = useStores()

  return (
    <AuthContext.Provider
      value={{
        user: authStore.user,
        login: authStore.login,
        logout: authStore.logout,
        isAuthenticated: authStore.isAuthenticated,
        markPasswordChanged: authStore.markPasswordChanged,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export const AuthProvider = observer(AuthProviderBase)

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
