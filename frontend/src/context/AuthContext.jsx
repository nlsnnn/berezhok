import { createContext, useContext, useState, useCallback } from 'react'
import { partnerLogin } from '@/api/partner'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [user, setUser] = useState(() => {
    try {
      const stored = localStorage.getItem('partner_user')
      return stored ? JSON.parse(stored) : null
    } catch {
      return null
    }
  })

  const login = useCallback(async (email, password) => {
    const data = await partnerLogin(email, password)
    const userData = {
      user_id: data.user_id,
      must_change_password: data.must_change_password,
    }
    localStorage.setItem('partner_token', data.token)
    localStorage.setItem('partner_user', JSON.stringify(userData))
    setUser(userData)
    return data
  }, [])

  const logout = useCallback(() => {
    localStorage.removeItem('partner_token')
    localStorage.removeItem('partner_user')
    setUser(null)
  }, [])

  return (
    <AuthContext.Provider value={{ user, login, logout, isAuthenticated: !!user }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
