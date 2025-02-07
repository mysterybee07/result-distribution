import React from 'react'
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import api from '../../api';
import { useQuery } from '@tanstack/react-query';

const fetchExamSchedules = async () => {
    const { data } = await api.get("/exam/schedules");
    return data;
};

const ListExams = () => {
    const { data, isLoading, error } = useQuery({
        queryKey: ["examSchedules"],
        queryFn: fetchExamSchedules,
    });
    console.log("🚀 ~ ListExams ~ data:", data)
    const exam_schedules = data?.examSchedules;
    console.log("🚀 ~ ListExams ~ exam_schedules:", exam_schedules)

    if (isLoading) return <p>Loading...</p>;
    if (error) return <p>Error fetching data: {error.message}</p>;
    return (
        <div>
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead>Course</TableHead>
                        <TableHead>Batch</TableHead>
                        <TableHead>Program</TableHead>
                        <TableHead>Semester</TableHead>
                        <TableHead>Exam Date</TableHead>
                    </TableRow>
                </TableHeader>

                <TableBody>
                    {Array.isArray(exam_schedules) &&
                        exam_schedules.map((item, index) => (
                            <TableRow key={index} className="text-left">
                                <TableCell>{item.course}</TableCell>
                                <TableCell>{item.batch}</TableCell>
                                <TableCell>{item.program}</TableCell>
                                <TableCell>{item.semester}</TableCell>
                                <TableCell>{item.exam_date.split("T")[0]}</TableCell>
                            </TableRow>
                        ))}
                </TableBody>
            </Table>
        </div>
    )
}

export default ListExams