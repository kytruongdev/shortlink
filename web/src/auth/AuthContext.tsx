import { createContext, useCallback, useEffect, useRef, useState, type ReactNode } from 'react'
import * as authApi from '../api/auth'
import type { AuthResponse } from '../api/types'
import { decodeUser, type AuthUser } from './jwt'
import { session } from './session'

export type AuthStatus = 'loading' | 'authed' | 'anon'

export interface AuthContextValue {
  user: AuthUser | null
  status: AuthStatus
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string) => Promise<void>
  logout: () => Promise<void>
}

// eslint-disable-next-line react-refresh/only-export-components
export const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null)
  const [status, setStatus] = useState<AuthStatus>('loading')

  const applyAuth = useCallback((res: AuthResponse) => {
    session.setToken(res.access_token)
    setUser(decodeUser(res.access_token))
    setStatus('authed')
  }, [])

  const clearAuth = useCallback(() => {
    session.setToken(null)
    setUser(null)
    setStatus('anon')
  }, [])

  const login = useCallback(
    async (email: string, password: string) => applyAuth(await authApi.login(email, password)),
    [applyAuth],
  )

  const register = useCallback(
    async (email: string, password: string) => applyAuth(await authApi.register(email, password)),
    [applyAuth],
  )

  const logout = useCallback(async () => {
    try {
      await authApi.logout()
    } finally {
      clearAuth()
    }
  }, [clearAuth])

  // Restore the session from the refresh cookie once on load, and let the axios
  // interceptor sign us out if a later refresh fails.
  const didInit = useRef(false)
  useEffect(() => {
    session.setOnUnauthorized(clearAuth)
    if (!didInit.current) {
      didInit.current = true
      authApi.refresh().then(applyAuth).catch(clearAuth)
    }
    return () => session.setOnUnauthorized(null)
  }, [applyAuth, clearAuth])

  return (
    <AuthContext.Provider value={{ user, status, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  )
}
