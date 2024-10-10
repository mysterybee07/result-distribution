import React from 'react'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "../components/ui/card"
import { Button } from '../components/ui/button'
import Layout from '../layout/Layout'
import { useNavigate } from 'react-router-dom'
import { RegisterForm } from '../forms/RegisterForm'

const Register = () => {
  const navigate = useNavigate();
  return (
    <Layout>
      <div className='flex justify-center mt-16'>
        <Card className="flex flex-row w-4/5 shadow-lg hover:shadow-2xl">
          <div className='w-1/2 mt-8 p-4 flex flex-col items-center justify-center'>
            <img src="../../public/logo.png" alt="nature" className="w-48 h-48 object-cover mx-auto rounded-lg" />
            <p>Lorem ipsum dolor sit amet consectetur adipisicing elit. Vero quidem nulla suscipit tempore consequatur aliquam, perferendis obcaecati sit architecto in, quia eligendi tenetur excepturi eveniet alias veritatis culpa, ea qui.</p>
          </div>
          <div className='w-1/2 mt-8'>
            <CardHeader>
              <CardTitle>Register</CardTitle>
              <CardDescription>Register using your information.</CardDescription>
            </CardHeader>
            <CardContent>
              <RegisterForm />
            </CardContent>
            <CardFooter>
              <p className='text-sm'>Already have an account?
                <Button variant='link' onClick={() => navigate("/login")}>
                  <p className='text-blue-600'>Login.</p>
                </Button>
              </p>
            </CardFooter>
          </div>
        </Card>
      </div >
    </Layout>
  )
}

export default Register