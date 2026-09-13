import { Link } from 'react-router-dom'

export function NotFound() {
  return (
    <main className="mx-auto max-w-md px-6 py-24 text-center">
      <h1 className="font-serif text-5xl font-semibold">404</h1>
      <p className="mt-2 text-muted">This page doesn't exist.</p>
      <Link to="/" className="mt-5 inline-block font-medium text-accent">
        ← Back home
      </Link>
    </main>
  )
}
