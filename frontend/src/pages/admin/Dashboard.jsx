import React, { useState } from 'react'
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
import { useAuth } from '../../context/AuthContext'
import { Navigate } from 'react-router-dom'
import ListNotice from '../notice/ListNotice'
import { useQuery } from '@tanstack/react-query'
import api from '../../api'
import { useData } from '../../context/DataContext'

const Dashboard = () => {
    const { isAuthenticated } = useAuth();
    const { programs, batches } = useData();

    const [selectedProgram, setSelectedProgram] = useState();
    const [selectedBatch, setSelectedBatch] = useState();
    // console.log("🚀 ~ Dashboard ~ filterParams:", filterParams)

    // // Modify this function to set the program ID when a program is selected
    // const filterByProgram = (id) => {
    //     console.log(id);
    //     setProgramId(id); // Update the programId
    // };

    // Fetch notices based on the selected programId
    const fetchNotices = async () => {
        const apiUrl =
            selectedProgram && selectedBatch ? `/notice/by-program-and-batch?program_id=${selectedProgram}&batch_id=${selectedBatch}`
                : selectedProgram ? `/notice/by-program?program_id=${selectedProgram}`
                    : '/notice'; // Conditional URL based on programId
        const response = await api.get(apiUrl);
        console.log(response.data.notices);
        return Array.isArray(response.data.notices) ? response.data.notices : []; // Ensure it's always an array
    };

    // Using useQuery to fetch notices based on programId
    const { data, isLoading, isError, error } = useQuery({
        queryKey: ['notices', selectedProgram],  // Include programId in the query key
        queryFn: fetchNotices,
        enabled: true,  // This ensures the query runs immediately
    });

    if (isLoading) {
        return <p>Loading...</p>;
    }

    if (isError) {
        return <p>Error: {error.message}</p>;
    }

    return isAuthenticated ? (
        <div className='w-full flex flex-col gap-16'>
            <div className='h-32 bg-blue-700 relative overflow-visible m-0 p-0'>
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
                    {data.length > 0 ?
                        <ListNotice
                            data={data}
                            programs={programs}
                            batches={batches}
                            selectedProgram={selectedProgram}
                            setSelectedProgram={setSelectedProgram}
                            setSelectedBatch={setSelectedBatch}
                        />
                        :
                        <p>No notices found</p>
                    }
                </div>
                <div className='row-span-2 col-span-1'>
                    <DashboardAsideCard />
                </div>
            </div>
        </div>
    ) : <Navigate to='/login' />
}

export default Dashboard