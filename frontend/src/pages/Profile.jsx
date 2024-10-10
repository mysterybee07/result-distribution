import React from 'react'
import Layout from '../layout/Layout'
import { Button } from '../components/ui/button'

const Profile = () => {
    return (
        <Layout>
            <div className='p-4 bg-blue-500 text-white'>Profile</div>
            <Button variant="secondary">Click me</Button>
        </Layout>
    )
}

export default Profile