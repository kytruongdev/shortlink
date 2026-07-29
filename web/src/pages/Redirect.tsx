import { useParams } from 'react-router-dom'
import { Aurora } from '../components/ui/Aurora'
import { Spinner } from '../components/ui/Spinner'

// Placeholder; the real resolve + countdown interstitial arrives in the link-features commit.
export function Redirect() {
  const { code } = useParams()
  return (
    <>
      <Aurora />
      <main className="grid min-h-screen place-items-center text-muted">
        <div className="flex items-center gap-2 text-accent">
          <Spinner /> Resolving {code}…
        </div>
      </main>
    </>
  )
}
