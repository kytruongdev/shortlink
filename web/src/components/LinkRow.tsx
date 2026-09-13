import { Check, Pencil, Trash2, X } from 'lucide-react'
import { useState } from 'react'
import { deleteLink, updateLink } from '../api/links'
import type { LinkItem } from '../api/types'
import { toApiError } from '../lib/apiClient'
import { Button } from './ui/Button'
import { CopyButton } from './ui/CopyButton'
import { Input } from './ui/Input'
import { Spinner } from './ui/Spinner'

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

export function LinkRow({
  link,
  onChange,
  onRemove,
}: {
  link: LinkItem
  onChange: (updated: LinkItem) => void
  onRemove: () => void
}) {
  const [editing, setEditing] = useState(false)
  const [confirming, setConfirming] = useState(false)
  const [value, setValue] = useState(link.long_url)
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  async function save() {
    setError(null)
    setBusy(true)
    try {
      onChange(await updateLink(link.code, value.trim()))
      setEditing(false)
    } catch (err) {
      setError(toApiError(err).message)
    } finally {
      setBusy(false)
    }
  }

  async function remove() {
    setBusy(true)
    try {
      await deleteLink(link.code)
      onRemove()
    } catch (err) {
      setError(toApiError(err).message)
      setBusy(false)
      setConfirming(false)
    }
  }

  if (editing) {
    return (
      <div className="rounded-2xl border border-line bg-white/80 px-4 py-3 backdrop-blur-sm">
        <div className="flex items-center gap-2">
          <span className="min-w-[76px] shrink-0 font-mono font-semibold text-accent">{link.code}</span>
          <Input value={value} onChange={(e) => setValue(e.target.value)} invalid={!!error} />
          <Button onClick={save} loading={busy} className="shrink-0">
            <Check size={15} /> Save
          </Button>
          <button
            type="button"
            onClick={() => {
              setEditing(false)
              setValue(link.long_url)
              setError(null)
            }}
            className="shrink-0 rounded-xl border border-line bg-white p-2.5 text-muted transition hover:text-ink"
          >
            <X size={16} />
          </button>
        </div>
        {error && <p className="mt-2 text-sm text-danger">{error}</p>}
      </div>
    )
  }

  return (
    <div className="flex items-center gap-3 rounded-2xl border border-line bg-white/75 px-4 py-3.5 backdrop-blur-sm transition hover:-translate-y-0.5 hover:shadow-lg hover:shadow-accent/10">
      <span className="min-w-[76px] shrink-0 font-mono font-semibold text-accent">{link.code}</span>
      <span className="min-w-0 flex-1 truncate text-sm text-muted">{link.long_url}</span>
      <span className="hidden whitespace-nowrap text-sm font-medium text-ink sm:inline">
        {link.click_count} {link.click_count === 1 ? 'click' : 'clicks'}
      </span>
      <span className="hidden whitespace-nowrap text-sm text-muted md:inline">{formatDate(link.created_at)}</span>

      {confirming ? (
        <span className="flex shrink-0 items-center gap-1.5">
          <span className="hidden text-sm text-muted sm:inline">Delete?</span>
          <button
            type="button"
            onClick={() => setConfirming(false)}
            className="rounded-lg border border-line bg-white px-3 py-2 text-sm font-medium text-muted"
          >
            Cancel
          </button>
          <button
            type="button"
            onClick={remove}
            disabled={busy}
            className="inline-flex items-center gap-1.5 rounded-lg bg-danger px-3 py-2 text-sm font-semibold text-white transition active:scale-95 disabled:opacity-70"
          >
            {busy ? <Spinner /> : 'Delete'}
          </button>
        </span>
      ) : (
        <span className="flex shrink-0 items-center gap-1">
          <CopyButton value={link.short_url} />
          <button
            type="button"
            onClick={() => setEditing(true)}
            title="Edit destination"
            className="rounded-lg border border-line bg-white p-2 text-muted transition hover:text-ink"
          >
            <Pencil size={15} />
          </button>
          <button
            type="button"
            onClick={() => setConfirming(true)}
            title="Delete"
            className="rounded-lg border border-line bg-white p-2 text-muted transition hover:text-danger"
          >
            <Trash2 size={15} />
          </button>
        </span>
      )}
    </div>
  )
}
