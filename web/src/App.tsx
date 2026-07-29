import { Link as LinkIcon, Sparkles } from 'lucide-react'
import { Aurora } from './components/ui/Aurora'
import { Badge } from './components/ui/Badge'
import { Button } from './components/ui/Button'
import { Card } from './components/ui/Card'
import { Input } from './components/ui/Input'

function App() {
  return (
    <>
      <Aurora />
      <main className="mx-auto max-w-3xl px-6 py-24 text-center">
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
          Theme &amp; UI kit wired: Vite + React + Tailwind, aurora background, Fraunces + Inter.
        </p>

        <Card className="mx-auto mt-8 flex max-w-xl gap-2 p-2">
          <Input defaultValue="https://kenh14.vn/some/very/long/article/path" />
          <Button>
            <LinkIcon size={16} /> Shorten
          </Button>
        </Card>

        <div className="mt-6 flex justify-center gap-3">
          <Button variant="ghost">Ghost button</Button>
          <Button loading>Loading</Button>
        </div>
      </main>
    </>
  )
}

export default App
