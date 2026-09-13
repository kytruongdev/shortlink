import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../auth/useAuth'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { toApiError } from '../lib/apiClient'

const EMAIL_RE = /^\S+@\S+\.\S+$/

export function Register() {
  const { register } = useAuth()
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function onSubmit(e: FormEvent) {
    e.preventDefault()
    setError(null)

    const trimmed = email.trim()
    if (!EMAIL_RE.test(trimmed)) {
      setError('Enter a valid email address.')
      return
    }
    if (password.length < 8 || password.length > 24) {
      setError('Password must be 8–24 characters.')
      return
    }

    setLoading(true)
    try {
      await register(trimmed, password)
      navigate('/links')
    } catch (err) {
      setError(toApiError(err).message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="mx-auto max-w-sm px-6 py-14">
      <div className="rounded-2xl border border-white/90 bg-white/75 p-8 shadow-[0_24px_60px_-28px_rgba(80,50,160,.4)] backdrop-blur-md">
        <h1 className="text-center font-serif text-2xl font-semibold">Create account</h1>
        <p className="mt-1 text-center text-sm text-muted">Free — takes ten seconds.</p>

        {error && (
          <div className="mt-5 rounded-lg border border-danger/30 bg-danger/5 px-4 py-2.5 text-sm text-danger">
            {error}
          </div>
        )}

        <form onSubmit={onSubmit} className="mt-5 space-y-3.5">
          <label className="block">
            <span className="mb-1.5 block text-sm text-muted">Email</span>
            <Input
              type="email"
              autoComplete="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </label>
          <label className="block">
            <span className="mb-1.5 block text-sm text-muted">Password</span>
            <Input
              type="password"
              autoComplete="new-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </label>
          <Button type="submit" loading={loading} className="w-full">
            Sign up
          </Button>
        </form>

        <p className="mt-5 text-center text-sm text-muted">
          Have an account?{' '}
          <Link to="/login" className="font-medium text-accent">
            Log in
          </Link>
        </p>
      </div>
    </div>
  )
}
