import { Navigate, Outlet } from 'react-router-dom'
import { Spinner } from '../components/ui/Spinner'
import { useAuth } from './useAuth'

export function ProtectedRoute() {
  const { status } = useAuth()

  if (status === 'loading') {
    return (
      <div className="grid min-h-[50vh] place-items-center text-accent">
        <Spinner />
      </div>
    )
  }
  if (status === 'anon') return <Navigate to="/login" replace />
  return <Outlet />
}
