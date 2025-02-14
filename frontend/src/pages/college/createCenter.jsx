import React, { useState } from "react";
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectLabel,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";

import { useData } from "../../context/DataContext";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Card, CardContent } from "../../components/ui/card";
import { Button } from "../../components/ui/button";
import { useMutation } from "@tanstack/react-query";
import { useLocation, useNavigate } from "react-router-dom";
import api from "../../api";

const formSchema = z.object({
    batch_id: z.number(),
    program_id: z.number(),
});

const CreateCenter = () => {
    const location = useLocation();
    const navigate = useNavigate();
    const queryParams = new URLSearchParams(location.search);
    const colleges = JSON.parse(decodeURIComponent(queryParams.get("colleges")));
    console.log("🚀 ~ CreateCenter ~ colleges:", colleges);

    const [records, setRecords] = useState(
        colleges.map((college) => ({
            collegeName: college,
            isCenter: false,
            capacity: 0,
        }))
    );
    const [message, setMessage] = useState("");

    const { programs, loadingPrograms, errorPrograms, batches, loadingBatches, errorBatches } = useData();

    if (loadingPrograms || loadingBatches) return <div>Loading...</div>;
    if (errorPrograms) return <div>Error loading programs: {errorPrograms.message}</div>;
    if (errorBatches) return <div>Error loading batches: {errorBatches.message}</div>;

    const form = useForm({
        resolver: zodResolver(formSchema),
        defaultValues: {
            batch_id: "",
            program_id: "",
        },
    });
    const { setValue, watch } = form;

    const handleRecordChange = (index, field, value) => {
        const updatedRecords = records.map((record, i) =>
            i === index ? { ...record, [field]: field === "capacity" ? parseInt(value, 10) : value } : record
        );
        setRecords(updatedRecords);
    };

    const onSubmit = async (data) => {
        const requestData = {
            batch_id: parseInt(data.batch_id, 10),
            program_id: parseInt(data.program_id, 10),
            records: records.map((r) => ({
                college_name: r.collegeName,
                is_center: true,
                capacity: parseInt(r.capacity, 10),
            })),
        };

        try {
            const response = await api.post("/exam/update-center-and-capacity", requestData);
            if (response.status === 200) {
                const message = response.data.message || "Success!";
                setMessage(message); // Display success message
                console.log("🚀 ~ onSubmit ~ response:", response);
                return navigate('/admin/college')
            } else {
                setMessage("An error occurred while updating capacities.");
            }
        } catch (err) {
            setMessage("Failed to send request.");
            console.error(err);
        }
    };


    return (
        <div className="flex items-center justify-center mt-16">
            <Card className="w-1/2 shadow-lg hover:shadow-2xl py-8 text-start">
                <CardContent>
                    <h1 className="text-xl font-bold mb-4">Create Center</h1>
                    <Form {...form}>
                        <form onSubmit={form.handleSubmit(onSubmit)}>
                            <div className="flex flex-col gap-4">
                                <Select
                                    value={watch("batch_id")}
                                    onValueChange={(value) => setValue("batch_id", Number(value))}
                                >
                                    <FormLabel>Select batch: </FormLabel>
                                    <SelectTrigger>
                                        <SelectValue placeholder="Select batch" />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectGroup>
                                            <SelectLabel>Batch</SelectLabel>
                                            {Array.isArray(batches) &&
                                                batches.map((batch, index) => (
                                                    <SelectItem key={index} value={batch.ID}>
                                                        {batch.batch}
                                                    </SelectItem>
                                                ))}
                                        </SelectGroup>
                                    </SelectContent>
                                </Select>

                                <Select
                                    value={watch("program_id")}
                                    onValueChange={(value) => setValue("program_id", Number(value))}
                                >
                                    <FormLabel>Select program: </FormLabel>
                                    <SelectTrigger>
                                        <SelectValue placeholder="Select  program" />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectGroup>
                                            <SelectLabel>Program</SelectLabel>
                                            {Array.isArray(programs) &&
                                                programs.map((program, index) => (
                                                    <SelectItem key={index} value={program.ID}>
                                                        {program.program_name}
                                                    </SelectItem>
                                                ))}
                                        </SelectGroup>
                                    </SelectContent>
                                </Select>

                                <div className="overflow-x-auto">
                                    <table className="table-auto border-collapse border border-gray-300 w-full">
                                        <thead>
                                            <tr className="bg-gray-200">
                                                <th className="border border-gray-300 px-4 py-2">College Name</th>
                                                <th className="border border-gray-300 px-4 py-2">Capacity</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {records.map((record, index) => (
                                                <tr key={index}>
                                                    <td className="border border-gray-300 px-4 py-2">
                                                        {record.collegeName}
                                                    </td>
                                                    <td className="border border-gray-300 px-4 py-2">
                                                        <input
                                                            type="number"
                                                            min="0"
                                                            value={record.capacity}
                                                            onChange={(e) =>
                                                                handleRecordChange(index, "capacity", e.target.value)
                                                            }
                                                            className="w-full border border-gray-300 rounded-md p-2"
                                                            placeholder="Enter capacity"
                                                        />
                                                    </td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                </div>

                                <Button type="submit" className="self-end mt-4">
                                    Create Center
                                </Button>

                            </div>
                        </form>
                    </Form>
                </CardContent>
            </Card>
        </div>
    );
};

export default CreateCenter;
