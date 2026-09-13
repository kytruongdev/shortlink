/**
 * Module-level bridge for the in-memory access token, shared between the axios
 * client (reads/refreshes it) and the AuthContext (owns the React state).
 * Kept outside React to avoid a circular import between them.
 */
let accessToken: string | null = null
let onUnauthorized: (() => void) | null = null

export const session = {
  getToken: () => accessToken,
  setToken: (token: string | null) => {
    accessToken = token
  },
  /** AuthContext registers this to react when a refresh finally fails. */
  setOnUnauthorized: (cb: (() => void) | null) => {
    onUnauthorized = cb
  },
  fireUnauthorized: () => onUnauthorized?.(),
}
