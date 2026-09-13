import type { InputHTMLAttributes } from 'react'
import { cn } from '../../lib/cn'

type Props = InputHTMLAttributes<HTMLInputElement> & { invalid?: boolean }

export function Input({ invalid = false, className, ...rest }: Props) {
  return (
    <input
      className={cn(
        'w-full rounded-xl border bg-white px-4 py-3 text-[15px] outline-none transition',
        invalid
          ? 'border-danger focus:border-danger focus:ring-4 focus:ring-danger/10'
          : 'border-line focus:border-accent focus:ring-4 focus:ring-accent/15',
        className,
      )}
      {...rest}
    />
  )
}
