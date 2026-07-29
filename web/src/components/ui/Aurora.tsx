/** Fixed animated aurora blobs painted behind all content. */
export function Aurora() {
  return (
    <div className="pointer-events-none fixed inset-0 -z-10 overflow-hidden">
      <span className="aurora-blob aurora-b1" />
      <span className="aurora-blob aurora-b2" />
      <span className="aurora-blob aurora-b3" />
    </div>
  )
}
