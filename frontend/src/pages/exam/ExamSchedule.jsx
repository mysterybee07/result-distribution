import { useMutation } from "@tanstack/react-query";
import React, { useState } from "react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import api from "../../api";

const createExamRoutine = async (examData) => {
    try {
        // Make a POST request with axios
        const response = await api.post("/exam/schedule/create", examData);

        // Check for a successful response
        return response.data;
    } catch (error) {
        // Handle any errors
        throw new Error(error.response?.data?.message || "Failed to create exam routine");
    }
};

const useCreateExamRoutine = () => {
    return useMutation({ mutationFn: createExamRoutine });
};

const ExamSchedule = () => {
    const { mutate, data, isLoading, error } = useCreateExamRoutine();

    const [examData, setExamData] = useState({
        batchId: "",
        programId: "",
        semesterId: "",
        startDate: "",
        endDate: "",
    });

    const handleChange = (e) => {
        setExamData({ ...examData, [e.target.name]: e.target.value });
    };

    const handleSubmit = (e) => {
        e.preventDefault();

        // Format startDate and endDate properly
        const formatDateTime = (dateString) => {
            const date = new Date(dateString);
            return date.toISOString().slice(0, 19).replace("T", " "); // Format: YYYY-MM-DD HH:mm:ss
        };

        mutate({
            ...examData,
            batchId: Number(examData.batchId),
            programId: Number(examData.programId),
            semesterId: Number(examData.semesterId),
            // startDate: Date(examData.startDate),
            // endDate: Date(examData.endDate),
            startDate: examData.startDate,
            endDate: examData.endDate,
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
                            <Input
                                id="batchId"
                                type="number"
                                name="batchId"
                                value={examData.batchId}
                                onChange={handleChange}
                                required
                            />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="programId">Program ID</Label>
                            <Input
                                id="programId"
                                type="number"
                                name="programId"
                                value={examData.programId}
                                onChange={handleChange}
                                required
                            />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="semesterId">Semester ID</Label>
                            <Input
                                id="semesterId"
                                type="number"
                                name="semesterId"
                                value={examData.semesterId}
                                onChange={handleChange}
                                required
                            />
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
