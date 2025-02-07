import React from 'react'
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"


const ListExams = () => {
    return (
        <div>
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead>Batch</TableHead>
                        <TableHead>Program</TableHead>
                        <TableHead>Semester</TableHead>
                        <TableHead>Start Date</TableHead>
                        <TableHead>End Date</TableHead>
                        <TableHead>Published</TableHead>
                    </TableRow>
                </TableHeader>

                {/* <TableBody>
                    {Array.isArray(courses) &&
                        courses.map((course, index) => (
                            <TableRow key={index}>
                                <TableCell>{course.course_code}</TableCell>
                                <TableCell>{course.name}</TableCell>
                                <TableCell>{course.is_compulsory ? 'True' : 'False'}</TableCell>
                                <TableCell>{course.semester_total_marks}</TableCell>
                                <TableCell>{course.semester_pass_marks}</TableCell>
                                <TableCell>{course.practical_total_marks}</TableCell>
                                <TableCell>{course.practical_pass_marks}</TableCell>
                                <TableCell>{course.assistant_total_marks}</TableCell>
                                <TableCell>{course.assistant_pass_marks}</TableCell>
                                <TableCell>
                                    <FaEdit
                                        className="text-blue-600 cursor-pointer"
                                        onClick={() => navigate(`/admin/courses/edit/${course.ID}`)}
                                    />
                                </TableCell>

                            </TableRow>
                        ))}
                </TableBody> */}
            </Table>
        </div>
    )
}

export default ListExams