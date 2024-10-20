import React from 'react'
import StudentTable from '../../components/StudentTable'
import { Button } from '@/components/ui/button'

const Student = () => {
    const handleClick = () => {
        console.log("Button clicked");
    };
    return (
        <div className='flex flex-col gap-8'>
            <Button onClick={handleClick} className="w-32 self-end">Add Student</Button>

            <StudentTable />
        </div>
    )
}

export default Student