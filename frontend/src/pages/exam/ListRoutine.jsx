import React from 'react';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import api from '../../api';
import { useQuery, useQueryClient } from '@tanstack/react-query';

const fetchExamRoutine = async () => {
    const { data } = await api.get("/exam/routines");
    return data;
};

const ListRoutine = () => {
    const queryClient = useQueryClient();

    // ✅ Toggle State with API Call
    const handleToggle = async (id, currentState) => {
        const newState = !currentState; // Toggle status
        console.log("🚀 ~ handleToggle ~ id:", id, "New State:", newState);

        try {
            const response = await api.post(`/exam/schedule/publish/${id}`, {
                status: newState, // Send the new state
            });

            console.log("✅ Toggle Successful:", response.data);

            // ✅ Refresh Data After Update
            queryClient.invalidateQueries(["examRoutine"]);
        } catch (error) {
            console.error("❌ Toggle Failed:", error);
        }
    };

    const { data, isLoading, error } = useQuery({
        queryKey: ["examRoutine"],
        queryFn: fetchExamRoutine,
    });

    console.log("🚀 ~ ListRoutine ~ data:", data);
    const examRoutine = data?.examRoutines;

    if (isLoading) return <p>Loading...</p>;
    if (error) return <p>Error fetching data: {error.message}</p>;

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
                        <TableHead>Status</TableHead>
                    </TableRow>
                </TableHeader>

                <TableBody>
                    {Array.isArray(examRoutine) &&
                        examRoutine.map((item, index) => (
                            <TableRow key={index} className="text-left">
                                <TableCell>{item.batch}</TableCell>
                                <TableCell>{item.program}</TableCell>
                                <TableCell>{item.semester}</TableCell>
                                <TableCell>{item.start_date.split("T")[0]}</TableCell>
                                <TableCell>{item.end_date.split("T")[0]}</TableCell>
                                <TableCell>
                                    <div className="flex items-center space-x-3">
                                        <span className="text-gray-700">
                                            {item.status ? "Active" : "Inactive"}
                                        </span>
                                        <button
                                            onClick={() => handleToggle(item.id, item.status)}
                                            className={`w-14 h-7 flex items-center rounded-full p-1 transition duration-300 ${
                                                item.status ? "bg-green-500" : "bg-gray-300"
                                            }`}
                                        >
                                            <div
                                                className={`w-6 h-6 bg-white rounded-full shadow-md transform transition-transform ${
                                                    item.status ? "translate-x-7" : "translate-x-0"
                                                }`}
                                            />
                                        </button>
                                    </div>
                                </TableCell>
                            </TableRow>
                        ))}
                </TableBody>
            </Table>
        </div>
    );
};

export default ListRoutine;
