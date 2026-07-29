import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios'
import { session } from '../auth/session'
import type { ApiError, AuthResponse } from '../api/types'

const baseURL = import.meta.env.VITE_API_BASE_URL

export const api = axios.create({ baseURL, withCredentials: true })

// Attach the in-memory access token to every request.
api.interceptors.request.use((config) => {
  const token = session.getToken()
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

// --- refresh-on-401 with single-flight queue ---
let isRefreshing = false
let waiters: Array<(token: string | null) => void> = []

const notifyWaiters = (token: string | null) => {
  waiters.forEach((resolve) => resolve(token))
  waiters = []
}

// Refresh via a bare axios call so it never re-enters this interceptor.
async function refreshAccessToken(): Promise<string | null> {
  try {
    const res = await axios.post<AuthResponse>(`${baseURL}/auth/refresh`, null, {
      withCredentials: true,
    })
    session.setToken(res.data.access_token)
    return res.data.access_token
  } catch {
    session.setToken(null)
    return null
  }
}

type RetriableConfig = InternalAxiosRequestConfig & { _retried?: boolean }

api.interceptors.response.use(
  (res) => res,
  async (error: AxiosError) => {
    const original = error.config as RetriableConfig | undefined
    const url = original?.url ?? ''
    const isAuthCall =
      url.includes('/auth/login') || url.includes('/auth/register') || url.includes('/auth/refresh')

    if (error.response?.status !== 401 || !original || original._retried || isAuthCall) {
      return Promise.reject(error)
    }
    original._retried = true

    // A refresh is already in flight: queue this request until it resolves.
    if (isRefreshing) {
      return new Promise((resolve, reject) => {
        waiters.push((token) => {
          if (!token) return reject(error)
          original.headers.Authorization = `Bearer ${token}`
          resolve(api(original))
        })
      })
    }

    isRefreshing = true
    const token = await refreshAccessToken()
    isRefreshing = false
    notifyWaiters(token)

    if (!token) {
      session.fireUnauthorized()
      return Promise.reject(error)
    }
    original.headers.Authorization = `Bearer ${token}`
    return api(original)
  },
)

/** Extract the API's { code, message } error, with a safe fallback. */
export function toApiError(err: unknown): ApiError {
  if (axios.isAxiosError(err) && err.response?.data && typeof err.response.data === 'object') {
    const data = err.response.data as Partial<ApiError>
    if (data.code && data.message) return { code: data.code, message: data.message }
  }
  return { code: 'UNKNOWN', message: 'Something went wrong. Please try again.' }
}
