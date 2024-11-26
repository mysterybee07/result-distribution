import React, { useState } from 'react'
import StudentTable from '../../components/StudentTable'
import { Button } from '@/components/ui/button'
import { useNavigate } from 'react-router-dom';
import {
    Drawer,
    DrawerContent,
    DrawerDescription,
    DrawerHeader,
    DrawerTitle,
    DrawerTrigger,
} from "@/components/ui/drawer"
import BulkStudentForm from '../../forms/BulkStudentForm';


const Student = () => {
    const navigate = useNavigate();
    const [isDrawerOpen, setIsDrawerOpen] = useState(false);
    console.log("🚀 ~ Student ~ isDrawerOpen:", isDrawerOpen)

    return (
        <div className='flex flex-col gap-8'>
            
            {/* <div className='flex justify-end gap-4'>
                <Drawer>
                    <DrawerTrigger  >
                        <Button
                            onClick={() => setIsDrawerOpen(true)}
                            className="w-32"
                            variant="secondary"
                        >
                            Bulk Add
                        </Button>
                    </DrawerTrigger>
                    <DrawerContent className="flex items-center justify-center">
                        <DrawerHeader>
                            <DrawerTitle>Add Students in Bulk</DrawerTitle>
                            <DrawerDescription>
                                <a href="/student_test.xlsx" download>
                                    <Button variant="link">
                                        Download Sample File
                                    </Button>
                                </a>
                            </DrawerDescription>
                        </DrawerHeader>
                        <BulkStudentForm isDrawerOpen={isDrawerOpen} setIsDrawerOpen={setIsDrawerOpen} />
                    </DrawerContent>
                </Drawer>


                <Button
                    onClick={() => navigate('/admin/students/create')}
                    className="w-32"
                >
                    Add Student
                </Button>
            </div> */}

            <StudentTable />
        </div>
    )
}

export default Student