import QRCode from 'qrcode'
import { useEffect, useState } from 'react'

/** Render a QR code for the given value, generated client-side. */
export function QRImage({ value, size = 96 }: { value: string; size?: number }) {
  const [src, setSrc] = useState<string | null>(null)

  useEffect(() => {
    let active = true
    QRCode.toDataURL(value, { margin: 0, width: size * 2 })
      .then((url) => {
        if (active) setSrc(url)
      })
      .catch(() => {
        /* ignore generation errors */
      })
    return () => {
      active = false
    }
  }, [value, size])

  return src ? (
    <img src={src} width={size} height={size} alt="QR code" className="rounded-lg" />
  ) : (
    <div style={{ width: size, height: size }} className="animate-pulse rounded-lg bg-line" />
  )
}
