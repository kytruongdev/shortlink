import { Download } from 'lucide-react'
import QRCode from 'qrcode'
import { useEffect, useState } from 'react'

/** Render a QR code for the given value, generated client-side. */
export function QRImage({
  value,
  size = 96,
  downloadName,
}: {
  value: string
  size?: number
  downloadName?: string
}) {
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

  async function download() {
    try {
      const hiRes = await QRCode.toDataURL(value, { margin: 2, width: 512 })
      const a = document.createElement('a')
      a.href = hiRes
      a.download = downloadName ?? 'qr.png'
      a.click()
    } catch {
      /* ignore */
    }
  }

  return (
    <div className="flex flex-col items-center gap-1.5">
      {src ? (
        <img src={src} width={size} height={size} alt="QR code" className="rounded-lg" />
      ) : (
        <div style={{ width: size, height: size }} className="animate-pulse rounded-lg bg-line" />
      )}
      {downloadName && (
        <button
          type="button"
          onClick={download}
          className="inline-flex items-center gap-1 text-xs font-medium text-accent hover:underline"
        >
          <Download size={12} /> Download
        </button>
      )}
    </div>
  )
}
