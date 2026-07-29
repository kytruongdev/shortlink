import { api } from '../lib/apiClient'
import type { EncodeResponse, LinkItem, ListLinksResponse, ResolveResponse } from './types'

export function encode(longUrl: string): Promise<EncodeResponse> {
  return api.post<EncodeResponse>('/encode', { long_url: longUrl }).then((r) => r.data)
}

export function resolve(code: string): Promise<ResolveResponse> {
  return api.get<ResolveResponse>(`/${code}`).then((r) => r.data)
}

export function listLinks(): Promise<LinkItem[]> {
  return api.get<ListLinksResponse>('/links').then((r) => r.data.links)
}
