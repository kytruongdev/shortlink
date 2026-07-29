import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { resolve } from '../api/links'
import { Aurora } from '../components/ui/Aurora'
import { Button } from '../components/ui/Button'

function safeHost(url: string): string {
  try {
    return new URL(url).host
  } catch {
    return url
  }
}

export function Redirect() {
  const { code } = useParams()
  const [dest, setDest] = useState<string | null>(null)
  const [notFound, setNotFound] = useState(false)
  const [count, setCount] = useState(3)

  // Resolve the code once.
  useEffect(() => {
    if (!code) return
    let active = true
    resolve(code)
      .then((r) => {
        if (active) setDest(r.long_url)
      })
      .catch(() => {
        if (active) setNotFound(true)
      })
    return () => {
      active = false
    }
  }, [code])

  // Count down, then navigate.
  useEffect(() => {
    if (!dest) return
    if (count <= 0) {
      window.location.replace(dest)
      return
    }
    const timer = setTimeout(() => setCount((c) => c - 1), 1000)
    return () => clearTimeout(timer)
  }, [dest, count])

  return (
    <>
      <Aurora />
      <main className="grid min-h-screen place-items-center px-6">
        <div className="w-full max-w-md rounded-2xl border border-white/90 bg-white/80 p-9 text-center shadow-[0_24px_60px_-28px_rgba(80,50,160,.4)] backdrop-blur-md">
          {notFound ? (
            <>
              <h1 className="font-serif text-2xl font-semibold">Link not found</h1>
              <p className="mt-2 text-sm text-muted">
                This short link doesn't exist or has been removed.
              </p>
              <Button className="mt-5" onClick={() => (window.location.href = '/')}>
                Go to ShortLink
              </Button>
            </>
          ) : (
            <>
              <div className="mx-auto mb-5 h-14 w-14 animate-spin rounded-full border-[3px] border-accent/20 border-t-accent" />
              <h1 className="font-serif text-2xl font-semibold">Redirecting…</h1>
              <p className="mt-1 text-sm text-muted">You're being sent to:</p>
              <div className="my-4 break-all rounded-xl border border-[#e6dcff] bg-[#f6f3ff] px-4 py-2.5 text-sm font-medium text-accent">
                {dest ? safeHost(dest) : '…'}
              </div>
              <p className="text-sm text-muted">
                Going in <b>{Math.max(count, 0)}</b>s…{' '}
                {dest && (
                  <a href={dest} className="font-medium text-accent">
                    Go now →
                  </a>
                )}
              </p>
            </>
          )}
        </div>
      </main>
    </>
  )
}
