import { useMutation, useQuery } from "@tanstack/react-query";
import React, { useState } from "react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import api from "../../api";
import { useData } from "../../context/DataContext";
import { useToast } from "@/hooks/use-toast";

const ExamSchedule = () => {
    const { toast } = useToast(); // Now we can use the hook here

    const createExamRoutine = async (examData) => {
        try {
            // Make a POST request with axios
            const response = await api.post("/exam/schedule/create", examData);
            return response.data; // Return the response data for useMutation
        } catch (error) {
            console.log("🚀 ~ createExamRoutine ~ error:", error);
            throw new Error(error.response?.data?.error || "Failed to create exam routine");
        }
    };
    const useCreateExamRoutine = () => {
        return useMutation({
            mutationFn: createExamRoutine,
            onSuccess: (data) => {
                // Show success toast when submission is successful
                toast({
                    title: "Center Updated",
                    description: data.message,
                    status: "success",
                    position: "top-right",
                    duration: 5000,
                });
            },
            onError: (error) => {
                // Show error toast if there's an issue
                toast({
                    title: "Error",
                    description: error.message,
                    status: "error",
                    position: "top-right",
                    duration: 5000,
                });
            },
        });
    };
    const { mutate, data, isLoading, error } = useCreateExamRoutine();
    const { programs, batches } = useData();
    
    const [examData, setExamData] = useState({
        batchId: "",
        programId: "",
        semesterId: "",
        startDate: "",
        endDate: "",
    });

    const { data: semesters } = useQuery({
        queryKey: ["semesters", examData.programId],
        queryFn: async () => {
            const response = await api.get(`/semester/by-program/${examData.programId}`);
            return response?.data?.semesters;
        },
        enabled: !!examData.programId,
    });

    const handleChange = (e) => {
        setExamData({ ...examData, [e.target.name]: e.target.value });
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        const startDate = `${examData.startDate}T00:00:00Z`;
        const endDate = `${examData.endDate}T00:00:00Z`;

        mutate({
            batch_id: Number(examData.batchId),
            program_id: Number(examData.programId),
            semester_id: Number(examData.semesterId),
            start_date: startDate,
            end_date: endDate,
        });
    };

    return (
        <div className="flex justify-center text-left">
            <Card className="w-3/4">
                <CardHeader>
                    <CardTitle className="text-xl font-bold text-center">Create Exam Routine</CardTitle>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleSubmit} className="space-y-8">
                        <div className="space-y-2">
                            <Label htmlFor="batchId" className="font-bold">Batch ID</Label>
                            <select
                                id="batchId"
                                name="batchId"
                                value={examData.batchId}
                                onChange={handleChange}
                                required
                                className="border rounded p-2 w-full"
                            >
                                <option value="" disabled>Select Batch</option>
                                {batches.map((batch) => (
                                    <option key={batch.ID} value={batch.ID}>{batch.batch}</option>
                                ))}
                            </select>
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="programId" className="font-bold">Program ID</Label>
                            <select
                                id="programId"
                                name="programId"
                                value={examData.programId}
                                onChange={handleChange}
                                required
                                className="border rounded p-2 w-full"
                            >
                                <option value="" disabled>Select Program</option>
                                {programs.map((program) => (
                                    <option key={program.ID} value={program.ID}>{program.program_name}</option>
                                ))}
                            </select>
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="semesterId" className="font-bold">Semester ID</Label>
                            <select
                                id="semesterId"
                                name="semesterId"
                                value={examData.semesterId}
                                onChange={handleChange}
                                required
                                className="border rounded p-2 w-full"
                            >
                                <option value="" disabled>Select Semester</option>
                                {semesters?.map((semester) => (
                                    <option key={semester.ID} value={semester.ID}>{semester.semester_name}</option>
                                ))}
                            </select>
                        </div>

                        {/* Start Date & Time */}
                        <div className="space-y-2">
                            <Label htmlFor="startDate" className="font-bold">Start Date</Label>
                            <Input
                                id="startDate"
                                type="date"
                                name="startDate"
                                value={examData.startDate}
                                onChange={handleChange}
                                required
                            />
                        </div>

                        {/* End Date & Time */}
                        <div className="space-y-2">
                            <Label htmlFor="endDate" className="font-bold">End Date</Label>
                            <Input
                                id="endDate"
                                type="date"
                                name="endDate"
                                value={examData.endDate}
                                onChange={handleChange}
                                required
                            />
                        </div>

                        <Button type="submit" disabled={isLoading} className="w-full">
                            {isLoading ? "Creating..." : "Create Routine"}
                        </Button>
                    </form>

                    {error && <p className="text-red-500 mt-2">Error: {error.message}</p>}

                    {data && (
                        <div className="mt-4 p-4 bg-gray-50 rounded-md">
                            <p className="font-semibold">{data.message}</p>
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    );
};


export default ExamSchedule;
