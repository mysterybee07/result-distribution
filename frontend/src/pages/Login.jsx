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
    <div className='flex flex-col justify-center items-center mt-8 bg-gray-100'>
      <div className='flex flex-col flex-1 mb-6 gap-2'>
        <img
          src="../../public/logo.webp"
          alt="Result-e logo"
          className="mt-8 w-32 h-32 object-cover mx-auto rounded-lg"
        />
        <p className='text-xl font-semibold text-center'>Sign in to Result-e</p>
      </div>
      <div className="flex-1 w-full h-screen flex justify-center items-center">
        <div className="w-full max-w-md mx-auto p-8 bg-white">
          <LoginForm />
        </div>
      </div>
    </div>

  )
}

export default Login