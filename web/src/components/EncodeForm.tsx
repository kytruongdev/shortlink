import { Link as LinkIcon } from 'lucide-react'
import { useState, type FormEvent } from 'react'
import { Link as RouterLink } from 'react-router-dom'
import { encode } from '../api/links'
import type { EncodeResponse } from '../api/types'
import { useAuth } from '../auth/useAuth'
import { toApiError } from '../lib/apiClient'
import { QRImage } from './QRImage'
import { Button } from './ui/Button'
import { Card } from './ui/Card'
import { CopyButton } from './ui/CopyButton'
import { Input } from './ui/Input'

export function EncodeForm({ onCreated }: { onCreated?: () => void }) {
  const { status } = useAuth()
  const [url, setUrl] = useState('')
  const [result, setResult] = useState<EncodeResponse | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function onSubmit(e: FormEvent) {
    e.preventDefault()
    setError(null)

    const value = url.trim()
    if (!/^https?:\/\/.+/i.test(value)) {
      setError('Please enter a valid http(s) URL.')
      setResult(null)
      return
    }

    setLoading(true)
    try {
      const res = await encode(value)
      setResult(res)
      onCreated?.()
    } catch (err) {
      const apiErr = toApiError(err)
      setError(
        apiErr.code === 'QUOTA_EXCEEDED'
          ? "You've reached the daily limit of 10 links. Log in for unlimited."
          : apiErr.message,
      )
      setResult(null)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="mx-auto max-w-xl text-left">
      <Card className="p-2">
        <form onSubmit={onSubmit} className="flex gap-2">
          <Input
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            invalid={!!error}
            placeholder="https://example.com/very/long/link"
            aria-label="URL to shorten"
          />
          <Button type="submit" loading={loading}>
            <LinkIcon size={16} /> Shorten
          </Button>
        </form>
      </Card>

      {error && <p className="mt-2.5 text-sm text-danger">{error}</p>}

      {result && (
        <div className="mt-3.5">
          <div className="flex items-center justify-between gap-3 rounded-xl border border-line bg-white px-4 py-3">
            <a
              href={result.short_url}
              target="_blank"
              rel="noreferrer"
              className="truncate font-semibold text-accent"
            >
              {result.short_url}
            </a>
            <CopyButton value={result.short_url} />
          </div>
          <div className="mt-2.5 flex items-center gap-3.5 rounded-xl border border-line bg-white px-3.5 py-3">
            <QRImage value={result.short_url} />
            <div className="text-sm text-muted">
              <span className="font-semibold text-ink">Scan to open</span>
              <br />
              Point a camera at the code to visit the short link.
            </div>
          </div>
        </div>
      )}

      <p className="mt-4 text-center text-sm text-muted">
        {status === 'anon' ? (
          <>
            Anonymous — up to 10 links/day.{' '}
            <RouterLink to="/register" className="font-medium text-accent">
              Sign up
            </RouterLink>{' '}
            for unlimited.
          </>
        ) : status === 'authed' ? (
          'Signed in — unlimited links.'
        ) : null}
      </p>
    </div>
  )
}
