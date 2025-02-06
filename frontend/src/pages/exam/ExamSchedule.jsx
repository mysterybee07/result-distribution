import { useMutation, useQuery } from "@tanstack/react-query";
import React, { useState } from "react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import api from "../../api";
import { useData } from "../../context/DataContext";

const createExamRoutine = async (examData) => {
    try {
        // Make a POST request with axios
        const response = await api.post("/exam/schedule/create", examData);

        // Check for a successful response
        return response.data;
    } catch (error) {
        console.log("🚀 ~ createExamRoutine ~ error:", error)
        // Handle any errors
        throw new Error(error.response.data.error || "Failed to create exam routine");
    }
};

const useCreateExamRoutine = () => {
    return useMutation({ mutationFn: createExamRoutine });
};

const ExamSchedule = () => {
    const { mutate, data, isLoading, error } = useCreateExamRoutine();
    const { programs, batches } = useData();
    console.log("🚀 ~ ExamSchedule ~ programs:", programs)

    const [examData, setExamData] = useState({
        batchId: "",
        programId: "",
        semesterId: "",
        startDate: "",
        endDate: "",
    });
    console.log("🚀 ~ ExamSchedule ~ examData:", examData)

    const {
        data: semesters,
        isLoading: loadingSemesters,
        error: errorSemesters,
    } = useQuery({
        queryKey: ["semesters", examData.programId], // Add `selectedProgram` to the query key
        queryFn: async () => {
            const response = await api.get(`/semester/by-program/${examData.programId}`);
            return response.data.semesters;
        },
        enabled: !!examData.programId,
    });
        console.log("🚀 ~ ExamSchedule ~ semesters:", semesters)
    const handleChange = (e) => {
        setExamData({ ...examData, [e.target.name]: e.target.value });
    };

    const handleSubmit = (e) => {
        e.preventDefault();

        // Format startDate and endDate properly
        const startDate = `${examData.startDate}T00:00:00Z`
        const endDate = `${examData.endDate}T00:00:00Z`

        mutate({
            batch_id: Number(examData.batchId),
            program_id: Number(examData.programId),
            semester_id: Number(examData.semesterId),
            start_date: startDate,
            end_date: endDate,
        });
    };


    return (
        <div className="flex justify-center items-center min-h-screen">
            <Card className="w-full max-w-lg shadow-lg">
                <CardHeader>
                    <CardTitle className="text-xl font-bold">Create Exam Routine</CardTitle>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleSubmit} className="space-y-4">


                        <div className="space-y-2">
                            <Label htmlFor="batchId">Batch ID</Label>
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
                            <Label htmlFor="programId">Program ID</Label>
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
                            <Label htmlFor="semesterId">Semester ID</Label>
                            <select
                                id="semesterId"
                                name="semesterId"
                                value={examData.semesterId}
                                onChange={handleChange}
                                required
                                className="border rounded p-2 w-full"
                            >
                                <option value="" disabled>Select Semester</option>
                                {semesters.map((semester) => (
                                    <option key={semester.ID} value={semester.ID}>{semester.semester_name}</option>
                                ))}
                            </select>
                        </div>


                        {/* Start Date & Time */}
                        <div className="space-y-2">
                            <Label htmlFor="startDate">Start Date</Label>
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
                            <Label htmlFor="endDate">End Date</Label>
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
                            <p className="text-sm text-gray-600">File Name: {data.fileName}</p>
                            <pre className="mt-2 p-2 bg-gray-100 rounded text-xs">{JSON.stringify(data.examSchedules, null, 2)}</pre>
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    );
};

export default ExamSchedule;
