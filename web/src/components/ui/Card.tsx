import type { ReactNode } from 'react'
import { cn } from '../../lib/cn'

export function Card({ className, children }: { className?: string; children: ReactNode }) {
  return (
    <div
      className={cn(
        'rounded-2xl border border-white/90 bg-white/70 shadow-[0_20px_50px_-24px_rgba(80,50,160,.35)] backdrop-blur-md',
        className,
      )}
    >
      {children}
    </div>
  )
}
