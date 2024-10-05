import React from 'react'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "../components/ui/card"
import { ProfileForm } from '../forms/LoginForm'
import { Button } from '../components/ui/button'
import { useNavigate } from 'react-router-dom'

const Login = () => {
  const navigate = useNavigate();
  return (
    // <div className="flex items-center justify-center h-screen bg-gray-200">
    //   <div className="w-1/3 bg-white rounded shadow-lg p-8 m-4">
    //     <h1 className="block w-full text-center text-2xl font-bold mb-6">Login</h1>
    //     <form className="mb-4 md:flex md:flex-wrap md:justify-between">
    //       <div className="flex flex-col mb-4 md:w-full">
    //         <label className="mb-2 font-bold text-lg text-grey-darkest" htmlFor="email">
    //           Email
    //         </label>
    //         <input className="border py-2 px-3 text-grey-darkest" type="email" name="email" id="email" />
    //       </div>
    //       <div className="flex flex-col mb-6 md:w-full">
    //         <label className="mb-2 font-bold text-lg text-grey-darkest" htmlFor="password">
    //           Password
    //         </label>
    //         <input className="border py-2 px-3 text-grey-darkest" type="password" name="password" id="password" />
    //       </div>
    //       <div className="flex items-center justify-between">
    //         <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded" type="button">
    //           Login
    //         </button>
    //         <a className="inline-block align-baseline font-bold text-sm text-blue-500 hover:text-blue-800" href="#">
    //           Forgot Password?
    //         </a>
    //       </div>
    //     </form>
    //     <div className="text-center">
    //       <p className="text-grey-dark text-sm">
    //         Don't have an account?{" "}
    //         <a href="#" className="no-underline text-blue font-bold">
    //           Sign up
    //         </a>
    //         .
    //       </p>
    //     </div>
    //   </div>
    // </div>
    <div className='flex items-center justify-center'>
      <Card className="w-1/3">
        <CardHeader>
          <CardTitle>Login</CardTitle>
          <CardDescription>Login to our site using identifier and password.</CardDescription>
        </CardHeader>
        <CardContent>
          <ProfileForm />
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