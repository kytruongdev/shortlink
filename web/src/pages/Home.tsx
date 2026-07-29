import { Sparkles } from 'lucide-react'
import { EncodeForm } from '../components/EncodeForm'
import { Badge } from '../components/ui/Badge'

export function Home() {
  return (
    <main className="mx-auto max-w-3xl px-6 py-16 text-center">
      <Badge>
        <Sparkles size={13} /> Fast, private URL shortener
      </Badge>
      <h1 className="mt-6 font-serif text-5xl font-semibold tracking-tight">
        Shorten links. Keep it{' '}
        <span className="bg-gradient-to-br from-accent to-accent-2 bg-clip-text text-transparent">
          clean.
        </span>
      </h1>
      <p className="mx-auto mt-4 max-w-md text-muted">
        Turn long URLs into short, shareable links in a single click.
      </p>

      <div className="mt-8">
        <EncodeForm />
      </div>
    </main>
  )
}
