import React, { useState } from 'react'
import StudentTable from '../../components/StudentTable'
import { useNavigate } from 'react-router-dom';

const Student = () => {
    const navigate = useNavigate();
    const [isDrawerOpen, setIsDrawerOpen] = useState(false);
    console.log("🚀 ~ Student ~ isDrawerOpen:", isDrawerOpen)

    return (
        <div className='flex flex-col gap-8'>
            <StudentTable />
        </div>
    )
}

export default Student