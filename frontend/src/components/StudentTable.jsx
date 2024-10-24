import {
    Table,
    TableBody,
    TableCaption,
    TableCell,
    TableFooter,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import React, { useState } from 'react';
import { useQuery } from "@tanstack/react-query";
import api from "../api";
import { FaEdit, FaTrash } from "react-icons/fa";
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
import BulkStudentForm from "../forms/BulkStudentForm";
// import BulkStudentForm from '../../forms/BulkStudentForm';


export default function StudentTable() {
    const navigate = useNavigate();
    const [isDrawerOpen, setIsDrawerOpen] = useState(false);
    console.log("🚀 ~ Student ~ isDrawerOpen:", isDrawerOpen)
    const [searchQuery, setSearchQuery] = useState('');  // State for search query
    const [selectedBatch, setSelectedBatch] = useState("");  // New state for batch filter
    const [selectedProgram, setSelectedProgram] = useState("");  // New state for program filter
    const [currentPage, setCurrentPage] = useState(1);   // State for pagination
    const [pageSize] = useState(10);  // Number of students per page (you can make this dynamic if needed)

    // Fetch students using useQuery
    const { data: students, isLoading, error } = useQuery({
        queryKey: ['students'],
        queryFn: async () => {
            const response = await api.get("/students");
            return response.data.students;
        },
    });

    if (isLoading) return <div>Loading...</div>;
    if (error) return <div>Error: {error.message}</div>;

    const uniqueBatches = [...new Set(students.map(student => student.Batch.batch))];  // Get unique batches
    const uniquePrograms = [...new Set(students.map(student => student.Program.program_name))];  // Get unique programs

    // Filter students by search query
    const filteredStudents = students.filter((student) => {
        const matchesSearch =
            student.fullname.toLowerCase().includes(searchQuery.toLowerCase()) ||
            student.symbol_number.toLowerCase().includes(searchQuery.toLowerCase()) ||
            student.registration_number.toLowerCase().includes(searchQuery.toLowerCase());

        const matchesBatch = selectedBatch ? student.Batch.batch.toString() === selectedBatch : true;
        const matchesProgram = selectedProgram ? student.Program.program_name === selectedProgram : true;

        return matchesSearch && matchesBatch && matchesProgram;
    });


    // Calculate total number of pages
    const totalPages = Math.ceil(filteredStudents.length / pageSize);

    // Get students for the current page
    const currentPageStudents = filteredStudents.slice((currentPage - 1) * pageSize, currentPage * pageSize);

    return (
        <>

            {/* Search Input */}
            <p className="font-bold text-2xl">Student Table</p>

            <div className='flex justify-between'>
                <div className="flex mb-4">
                    <input
                        type="text"
                        placeholder="Search..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="border p-2 mr-4 w-64"
                    />
                    <select
                        value={selectedBatch}
                        onChange={(e) => setSelectedBatch(e.target.value)}
                        className="border p-2 mr-4"
                    >
                        <option value="">All Batches</option>
                        {uniqueBatches.map((batch) => (
                            <option key={batch} value={batch}>{batch}</option>
                        ))}
                    </select>
                    <select
                        value={selectedProgram}
                        onChange={(e) => setSelectedProgram(e.target.value)}
                        className="border p-2"
                    >
                        <option value="">All Programs</option>
                        {uniquePrograms.map((program) => (
                            <option key={program} value={program}>{program}</option>
                        ))}
                    </select>
                </div>

                <div>
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
                        className="w-32 ml-4"
                    >
                        Add Student
                    </Button>
                </div>
            </div>


            {/* Student Table */}
            <Table>
                {/* <TableCaption className="font-bold text-xl">Student Table</TableCaption> */}
                <TableHeader className="text-left">
                    <TableRow>
                        <TableHead>Name</TableHead>
                        <TableHead>Symbol Number</TableHead>
                        <TableHead>Reg Number</TableHead>
                        <TableHead>Batch</TableHead>
                        <TableHead>Program</TableHead>
                        <TableHead>Semester</TableHead>
                        <TableHead>Status</TableHead>
                        <TableHead>Actions</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {currentPageStudents.map((student) => (
                        <TableRow key={student.ID} className="text-left">
                            <TableCell >{student.fullname}</TableCell>
                            <TableCell >{student.symbol_number}</TableCell>
                            <TableCell >{student.registration_number}</TableCell>
                            <TableCell >{student.Batch.batch}</TableCell>
                            <TableCell >{student.Program.program_name}</TableCell>
                            <TableCell >{student.current_semester}</TableCell>
                            <TableCell >{student.status}</TableCell>
                            <TableCell className="flex items-center gap-4">
                                <FaEdit />
                                <FaTrash className="text-red-600" />
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>

            {/* Pagination Controls */}
            <div className="flex justify-between mt-4">
                <button
                    disabled={currentPage === 1}
                    onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
                    className="px-4 py-2 bg-gray-300 rounded-md disabled:opacity-50 flex items-center"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                        strokeWidth={1.5}
                        stroke="currentColor"
                        className="w-5 h-5 mr-2"
                    >
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            d="M15.75 19.5L8.25 12l7.5-7.5"
                        />
                    </svg>
                    Prev
                </button>


                <span>Page {currentPage} of {totalPages}</span>

                <button
                    disabled={currentPage === totalPages}
                    onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
                    className="px-4 py-2 bg-gray-300 rounded-md disabled:opacity-50 flex items-center"
                >
                    Next
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                        strokeWidth={1.5}
                        stroke="currentColor"
                        className="w-5 h-5 ml-2"
                    >
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            d="M8.25 4.5l7.5 7.5-7.5 7.5"
                        />
                    </svg>                </button>
            </div>
        </>
    );
}

