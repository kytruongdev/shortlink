import { jwtDecode } from 'jwt-decode'

export interface AuthUser {
  username: string
}

interface TokenPayload {
  sub: string
  username: string
  exp: number
}

/** Read the display user from an access token, or null if it can't be decoded. */
export function decodeUser(token: string): AuthUser | null {
  try {
    const payload = jwtDecode<TokenPayload>(token)
    return { username: payload.username }
  } catch {
    return null
  }
}
