export interface User {
  id: string
  email: string
}

export interface AuthResponse {
  access_token: string
  expires_in: number
  user?: User
}

export interface EncodeResponse {
  short_url: string
  code: string
}

export interface ResolveResponse {
  long_url: string
}

export interface LinkItem {
  code: string
  short_url: string
  long_url: string
  created_at: string
}

export interface ListLinksResponse {
  links: LinkItem[]
}

export interface ApiError {
  code: string
  message: string
}
