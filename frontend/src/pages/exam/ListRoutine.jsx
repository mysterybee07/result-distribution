import React, { useState } from 'react';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog"

import api from '../../api';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { Button } from "@/components/ui/button"
import { useNavigate } from "react-router-dom"


const fetchExamRoutine = async () => {
    const { data } = await api.get("/exam/routines");
    return data;
};

const fetchExamSchedules = async ({ queryKey }) => {
    const [, { batchID, programID, semesterID }] = queryKey;
    const { data } = await api.get("/exam/schedules/by-batch-program", {
        params: { batch_id: batchID, program_id: programID, semester_id: semesterID },
    });
    return data.examSchedules; // Assuming response contains { examSchedules: [...] }
};

const ListRoutine = () => {
    const queryClient = useQueryClient();
    const [selectedRoutine, setSelectedRoutine] = useState(null); // Track selected routine

    // Fetch exam schedules when dialog opens
    const { data: examSchedules, isLoading, error } = useQuery({
        queryKey: ["examSchedules", selectedRoutine],
        queryFn: fetchExamSchedules,
        enabled: !!selectedRoutine, // Fetch only when selectedRoutine is set
    });

    const handleToggle = async (id, currentState) => {
        try {
            const newState = !currentState;
            await api.post(`/exam/schedule/publish/${id}`, { status: newState });
            queryClient.invalidateQueries(["examRoutine"]); // Refresh data
        } catch (error) {
            console.error("Toggle Failed:", error);
        }
    };

    const { data, isLoading: isRoutineLoading, error: routineError } = useQuery({
        queryKey: ["examRoutine"],
        queryFn: fetchExamRoutine,
    });

    if (isRoutineLoading) return <p>Loading...</p>;
    if (routineError) return <p>Error fetching data: {routineError.message}</p>;

    const examRoutine = data?.examRoutines;

    return (
        <div>
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead>Status</TableHead>
                        <TableHead>Batch</TableHead>
                        <TableHead>Program</TableHead>
                        <TableHead>Semester</TableHead>
                        <TableHead>Start Date</TableHead>
                        <TableHead>End Date</TableHead>
                        <TableHead>Action</TableHead>
                    </TableRow>
                </TableHeader>

                <TableBody>
                    {examRoutine?.map((item) => (
                        <TableRow key={item.id}>
                            <TableCell>
                                <div className="flex items-center space-x-3">
                                    <span className="text-gray-700">
                                        {item.status ? "Active" : "Inactive"}
                                    </span>
                                    <button
                                        onClick={() => handleToggle(item.id, item.status)}
                                        className={`w-14 h-7 flex items-center rounded-full p-1 transition duration-300 ${item.status ? "bg-green-500" : "bg-gray-300"
                                            }`}
                                    >
                                        <div
                                            className={`w-6 h-6 bg-white rounded-full shadow-md transform transition-transform ${item.status ? "translate-x-7" : "translate-x-0"
                                                }`}
                                        />
                                    </button>
                                </div>
                            </TableCell>
                            <TableCell>{item.batch}</TableCell>
                            <TableCell>{item.program}</TableCell>
                            <TableCell>{item.semester}</TableCell>
                            <TableCell>{item.start_date.split("T")[0]}</TableCell>
                            <TableCell>{item.end_date.split("T")[0]}</TableCell>

                            <TableCell>
                                <Dialog onOpenChange={(isOpen) => isOpen && setSelectedRoutine(item)}>
                                    <DialogTrigger><Button>View Schedule</Button></DialogTrigger>
                                    <DialogContent>
                                        <DialogHeader>
                                            <DialogTitle>Exam Schedule</DialogTitle>
                                        </DialogHeader>
                                        {isLoading ? (
                                            <p>Loading...</p>
                                        ) : error ? (
                                            <p>Error: {error.message}</p>
                                        ) : (
                                            <div className="overflow-x-auto">
                                                <table className="w-full border-collapse border border-gray-300 rounded-lg shadow-md">
                                                    <thead>
                                                        <tr className="bg-gray-200 text-gray-700">
                                                            <th className="border border-gray-300 px-4 py-2">Course</th>
                                                            <th className="border border-gray-300 px-4 py-2">Exam Date</th>
                                                            <th className="border border-gray-300 px-4 py-2">Start Time</th>
                                                            <th className="border border-gray-300 px-4 py-2">End Time</th>
                                                            <th className="border border-gray-300 px-4 py-2">Venue</th>
                                                        </tr>
                                                    </thead>
                                                    <tbody>
                                                        {examSchedules?.map((schedule, index) => (
                                                            <tr key={index} className="border border-gray-300">
                                                                <td className="border border-gray-300 px-4 py-2">{schedule.course}</td>
                                                                <td className="border border-gray-300 px-4 py-2">{schedule.exam_date}</td>
                                                                <td className="border border-gray-300 px-4 py-2">12:00 PM</td>
                                                                <td className="border border-gray-300 px-4 py-2">03:00 PM</td>
                                                                <td className="border border-gray-300 px-4 py-2">{schedule.venue}</td>
                                                            </tr>
                                                        ))}
                                                    </tbody>
                                                </table>
                                            </div>
                                        )}
                                    </DialogContent>
                                </Dialog>
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </div>
    );
};


export default ListRoutine;
