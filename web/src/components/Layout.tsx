import { Outlet } from 'react-router-dom'
import { Aurora } from './ui/Aurora'
import { Navbar } from './Navbar'

export function Layout() {
  return (
    <>
      <Aurora />
      <Navbar />
      <Outlet />
    </>
  )
}
