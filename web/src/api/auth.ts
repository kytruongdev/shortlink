import { api } from '../lib/apiClient'
import type { AuthResponse } from './types'

export function register(email: string, password: string): Promise<AuthResponse> {
  return api.post<AuthResponse>('/auth/register', { email, password }).then((r) => r.data)
}

export function login(email: string, password: string): Promise<AuthResponse> {
  return api.post<AuthResponse>('/auth/login', { email, password }).then((r) => r.data)
}

export function refresh(): Promise<AuthResponse> {
  return api.post<AuthResponse>('/auth/refresh').then((r) => r.data)
}

export function logout(): Promise<void> {
  return api.post('/auth/logout').then(() => undefined)
}
