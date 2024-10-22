import React from 'react'
import StudentTable from '../../components/StudentTable'
import { Button } from '@/components/ui/button'
import { useNavigate } from 'react-router-dom';

const Student = () => {
    const navigate = useNavigate();
    const handleClick = () => {
        navigate('/admin/students/create');
    };
    return (
        <div className='flex flex-col gap-8'>
            <Button onClick={handleClick} className="w-32 self-end">
                Add Student
            </Button>

            <StudentTable />
        </div>
    )
}

export default Student