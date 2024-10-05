import React from 'react'
import { Button } from '../components/ui/button'
import { Link, useNavigate } from 'react-router-dom'

const Navigation = () => {
  const navigate = useNavigate()
  return (
    <nav className="fixed inset-x-0 top-0 z-50 bg-white shadow-sm dark:bg-gray-950/90">
      <div className="w-full max-w-7xl mx-auto px-4">
        <div className="flex justify-between h-14 items-center">
          <Link href="#" className="flex items-center" prefetch={false}>
            {/* <MountainIcon className="h-6 w-6" /> */}
            <span className="sr-only">Acme Inc</span>
          </Link>
          <nav className="hidden md:flex gap-4">
            <Link
              to="/"
              className="font-medium flex items-center text-sm transition-colors hover:underline"
              prefetch={false}
            >
              Home
            </Link>
            <Link
              to="/exam"
              className="font-medium flex items-center text-sm transition-colors hover:underline"
              prefetch={false}
            >
              Exam
            </Link>
            <Link
              to="/result"
              className="font-medium flex items-center text-sm transition-colors hover:underline"
              prefetch={false}
            >
              Result
            </Link>
          </nav>
          <div className="flex items-center gap-4">
            <Button variant="outline" size="sm" onClick={() => navigate('/login')}>
              Sign in
            </Button>
            <Button size="sm" onClick={() => navigate('/register')}>Sign up</Button>
          </div>
        </div>
      </div>  
    </nav>

  )
}

export default Navigation