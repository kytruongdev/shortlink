import type { ButtonHTMLAttributes, ReactNode } from 'react'
import { cn } from '../../lib/cn'
import { Spinner } from './Spinner'

type Variant = 'primary' | 'ghost'

type Props = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: Variant
  loading?: boolean
  children?: ReactNode
}

export function Button({
  variant = 'primary',
  loading = false,
  className,
  children,
  disabled,
  ...rest
}: Props) {
  return (
    <button
      className={cn(
        'inline-flex items-center justify-center gap-2 rounded-xl px-4 py-2.5 text-sm font-semibold',
        'transition duration-150 active:translate-y-px active:scale-[.98] disabled:opacity-70',
        variant === 'primary' &&
          'bg-gradient-to-br from-accent to-indigo-500 text-white shadow-lg shadow-accent/30 hover:brightness-105',
        variant === 'ghost' &&
          'border border-line bg-white text-ink hover:shadow-md hover:shadow-black/5',
        className,
      )}
      disabled={disabled || loading}
      {...rest}
    >
      {loading ? <Spinner /> : children}
    </button>
  )
}
