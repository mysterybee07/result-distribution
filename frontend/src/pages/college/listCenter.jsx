import React from 'react'
import api from '../../api';
import { useQuery } from '@tanstack/react-query';
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
import { FaEdit, FaTrash } from 'react-icons/fa';
const ListCenter = () => {
    const fetchCenter = async () => {
        const response = await api.get("/college/all-centers");
        console.log("🚀 ~ fetchCenter ~ response:", response.data);
        return response.data.centers;
    };

    const { data: centers = [], isLoading, error } = useQuery({
        queryKey: ["centers"],
        queryFn: fetchCenter,
    });
    console.log("🚀 ~ ListCenter ~ centers:", centers)
    if (isLoading) return <div>Loading...</div>;
    if (error) return <div>Error: {error.message}</div>;
    return (
        <div>
            <Table>
                {/* <TableCaption className="font-bold text-xl">Student Table</TableCaption> */}
                <TableHeader className="text-left">
                    <TableRow>
                        <TableHead>S.N</TableHead>
                        <TableHead>Collge Name</TableHead>
                        <TableHead>Adddress</TableHead>
                        <TableHead>Capacity</TableHead>
                        <TableHead>Student Count</TableHead>
                        {/* <TableHead>Actions</TableHead> */}
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {centers.map((center, index) => (
                        <TableRow key={center.ID} className="text-left">
                            <TableCell >{index+1}</TableCell>
                            <TableCell >{center.college_name}</TableCell>
                            <TableCell >{center.address}</TableCell>
                            <TableCell >{center.capacity}</TableCell>
                            <TableCell >{center.students_count}</TableCell>
                            {/* <TableCell className="flex items-center gap-4">
                                <FaEdit
                                    className="text-blue-600 cursor-pointer"
                                    onClick={() => navigate(`/admin/students/edit/${student.ID}`)}
                                />
                                <FaTrash
                                    onClick={() => navigate(`/admin/students/${student.ID}`)}
                                    className="text-red-600"
                                />
                            </TableCell> */}
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </div>
    )
}

export default ListCenter