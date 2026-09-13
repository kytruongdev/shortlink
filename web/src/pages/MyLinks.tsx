import { Link as LinkIcon } from 'lucide-react'
import { useCallback, useEffect, useState } from 'react'
import { listLinks } from '../api/links'
import type { LinkItem } from '../api/types'
import { EncodeForm } from '../components/EncodeForm'
import { LinkRow } from '../components/LinkRow'
import { Badge } from '../components/ui/Badge'
import { Spinner } from '../components/ui/Spinner'

function EmptyState() {
  return (
    <div className="mt-6 flex flex-col items-center gap-3 py-12 text-center text-muted">
      <LinkIcon size={34} className="text-accent/60" />
      <p>No links yet. Create your first short link and it'll show up here.</p>
    </div>
  )
}

export function MyLinks() {
  const [links, setLinks] = useState<LinkItem[] | null>(null)

  const load = useCallback(() => {
    listLinks()
      .then(setLinks)
      .catch(() => setLinks([]))
  }, [])

  useEffect(() => {
    load()
  }, [load])

  const onChange = (updated: LinkItem) =>
    setLinks((prev) => prev?.map((l) => (l.code === updated.code ? updated : l)) ?? null)

  const onRemove = (code: string) =>
    setLinks((prev) => prev?.filter((l) => l.code !== code) ?? null)

  return (
    <main className="mx-auto max-w-3xl px-6 py-14">
      <div className="text-center">
        <Badge>
          <LinkIcon size={13} /> New link
        </Badge>
        <h1 className="mt-4 font-serif text-3xl font-semibold">Create a short link</h1>
      </div>

      <div className="mt-6">
        <EncodeForm onCreated={load} />
      </div>

      <section className="mt-12">
        <h2 className="text-center text-sm uppercase tracking-wide text-muted">Your links</h2>
        {links === null ? (
          <div className="mt-6 grid place-items-center text-accent">
            <Spinner />
          </div>
        ) : links.length === 0 ? (
          <EmptyState />
        ) : (
          <div className="mt-4 space-y-2.5">
            {links.map((link) => (
              <LinkRow
                key={link.code}
                link={link}
                onChange={onChange}
                onRemove={() => onRemove(link.code)}
              />
            ))}
          </div>
        )}
      </section>
    </main>
  )
}
