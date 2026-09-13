import { Link as LinkIcon } from 'lucide-react'
import { Link, NavLink, useNavigate } from 'react-router-dom'
import { useAuth } from '../auth/useAuth'
import { cn } from '../lib/cn'
import { Button } from './ui/Button'

export function Navbar() {
  const { user, status, logout } = useAuth()
  const navigate = useNavigate()

  return (
    <header className="mx-auto flex max-w-4xl items-center justify-between px-6 py-4">
      <Link to="/" className="flex items-center gap-2 text-lg font-bold tracking-tight">
        <span className="grid h-7 w-7 place-items-center rounded-lg bg-gradient-to-br from-accent to-accent-2 text-white shadow-lg shadow-accent/40">
          <LinkIcon size={15} />
        </span>
        short
        <span className="bg-gradient-to-br from-accent to-accent-2 bg-clip-text text-transparent">
          link
        </span>
      </Link>

      <div className="flex items-center gap-2.5 text-sm">
        {status === 'authed' && user ? (
          <>
            <NavLink
              to="/links"
              className={({ isActive }) =>
                cn(
                  'rounded-lg px-3 py-2 font-medium transition',
                  isActive ? 'text-accent' : 'text-muted hover:text-ink',
                )
              }
            >
              My links
            </NavLink>
            <Link
              to="/links"
              className="flex items-center gap-2 font-semibold text-ink"
              title="My links"
            >
              <span className="grid h-7 w-7 place-items-center rounded-full bg-gradient-to-br from-accent to-accent-2 text-xs font-bold text-white">
                {user.username.charAt(0).toUpperCase()}
              </span>
              <span className="hidden sm:inline">{user.username}</span>
            </Link>
            <Button variant="ghost" onClick={() => logout()}>
              Log out
            </Button>
          </>
        ) : status === 'anon' ? (
          <>
            <Button variant="ghost" onClick={() => navigate('/login')}>
              Log in
            </Button>
            <Button onClick={() => navigate('/register')}>Sign up</Button>
          </>
        ) : null}
      </div>
    </header>
  )
}
