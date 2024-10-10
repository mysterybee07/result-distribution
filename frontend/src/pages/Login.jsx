import React from 'react'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "../components/ui/card"
import { LoginForm } from '../forms/LoginForm'
import { Button } from '../components/ui/button'
import { useNavigate } from 'react-router-dom'

const Login = () => {
  const navigate = useNavigate();
  return (
    <div className='flex flex-col items-center justify-center mt-8'>
      <div className='mb-6'>
        <img src="../../public/logo.png" alt="nature" className="mt-8 w-16 h-16 object-cover mx-auto rounded-lg" />
        <p className='text-xl font-semibold'>Sign in to Result-e</p>
      </div>
      <Card className="w-1/3 shadow-lg hover:shadow-2xl py-6">
        {/* <CardHeader>
          <CardTitle>Login</CardTitle>
          <CardDescription>Login to our site using identifier and password.</CardDescription>
        </CardHeader> */}
        <CardContent>
          <LoginForm />
        </CardContent>
        <CardFooter>
          <p className='text-sm'>Don't have an account?
            <Button variant='link' onClick={() => navigate("/register")}>
              <p className='text-blue-600'>Sign up.</p>
            </Button>
          </p>
        </CardFooter>
      </Card>
    </div>
  )
}

export default Login