import React from 'react'
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import StudentTable from '../../components/StudentTable'
import DashboardAsideCard from '../../components/DashboardAsideCard'

const Dashboard = () => {
    return (
        <div className='w-full flex flex-col gap-16'>
            <div className='h-32 bg-blue-700 mt-8 relative overflow-visible m-0 p-0'>
                <p className='text-white text-start p-8 text-xl'>Admin Dashboard</p>
                <div className='p-4 grid grid-cols-1 gap-16 md:grid-cols-4'>
                    <Card className="-mt-8 bg-red-200 border-none"> {/* Keep negative margin for overflow */}
                        <CardHeader>
                            <CardTitle>100+</CardTitle>
                            <CardDescription>Students</CardDescription>
                        </CardHeader>
                    </Card>
                    <Card className="-mt-8 bg-yellow-400 border-none">
                        <CardHeader>
                            <CardTitle>100+</CardTitle>
                            <CardDescription>Users</CardDescription>
                        </CardHeader>
                    </Card>
                    <Card className="-mt-8 bg-blue-200 border-none">
                        <CardHeader>
                            <CardTitle>5+</CardTitle>
                            <CardDescription>Courses</CardDescription>
                        </CardHeader>
                    </Card>
                    <Card className="-mt-8 bg-green-200 border-none">
                        <CardHeader>
                            <CardTitle>20+</CardTitle>
                            <CardDescription>Batches</CardDescription>
                        </CardHeader>
                    </Card>
                </div>
            </div>
            <div className='grid grid-rows-2 grid-flow-col gap-16'>
                <div className="row-span-2 col-span-5">
                    <StudentTable />
                </div>
                <div className='row-span-2 col-span-1'>
                    <DashboardAsideCard />
                </div>
            </div>
        </div>
    )
}

export default Dashboard