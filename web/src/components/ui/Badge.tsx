import type { ReactNode } from 'react'

export function Badge({ children }: { children: ReactNode }) {
  return (
    <span className="inline-flex items-center gap-1.5 rounded-full border border-[#e6dcff] bg-[#f3edff] px-3 py-1.5 text-xs font-semibold text-accent">
      {children}
    </span>
  )
}
