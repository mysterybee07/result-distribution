import React, { useState } from 'react';
import Papa from 'papaparse';
import { Card, CardContent } from '../components/ui/card';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import {
    Form,
    FormLabel,
} from "@/components/ui/form"
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectLabel,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import {
    DrawerFooter,
    DrawerClose,
} from "@/components/ui/drawer"
import { useData } from '../context/DataContext';
import { Button } from '../components/ui/button';
import { useMutation } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import api from '../api';

const formSchema = z.object({
    batch_id: z.number(),
    program_id: z.number(),
    file: z.any(),
});

const BulkStudentForm = ({ isDrawerOpen, setIsDrawerOpen }) => {
    const navigate = useNavigate();

    // fetching batch and program data
    const { programs, loadingPrograms, errorPrograms, batches, loadingBatches, errorBatches } = useData();

    if (loadingPrograms || loadingBatches) return <div>Loading...</div>;

    if (errorPrograms) return <div>Error loading programs: {errorPrograms.message}</div>;
    if (errorBatches) return <div>Error loading batches: {errorBatches.message}</div>;
    const [students, setStudents] = useState([]);
    const [file, setFile] = useState(null);

    const handleFileChange = (event) => {
        const file = event.target.files[0];
        setFile(file);
        setValue('file', file); // Set the file in the form state

        Papa.parse(file, {
            header: true,
            complete: (results) => {
                setStudents(results.data);
            },
        });
    };

    // define the form
    const form = useForm({
        resolver: zodResolver(formSchema),
        defaultValues: {
            batch_id: '',
            program_id: '',
            file: null
        },
    });
    const { register, formState: { errors }, setValue, watch } = form;

    // mutation for creating student
    const { mutate: createStudent } = useMutation({
        mutationFn: async (formData) => {
            const response = await api.post('/students/create', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            });
            return response.data;
        },
        onSuccess: (data) => {
            console.log("Student added successfully:", data);
            setIsDrawerOpen(false);
            navigate('/admin/students');
        },
        onError: (error) => {
            console.error("Error adding student:", error);
        },
    });

    // define the submit handler
    const onSubmit = (data) => {
        console.log("Current file state:", file);
        console.log("Form data received:", data);

        const formData = new FormData();
        formData.append("batch_id", String(data.batch_id));
        formData.append("program_id", String(data.program_id));

        if (data.file && data.file instanceof File) {
            formData.append("file", data.file);
        }

        // Debugging FormData
        console.log("FormData entries:");
        for (const [key, value] of formData.entries()) {
            console.log(key, value);
        }

        createStudent(formData);
    };

    return (
        <Card className="w-1/2 shadow-lg hover:shadow-2xl py-8 text-start">
            <CardContent>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} encType='multipart/form-data' method='post'>
                        <Select
                            value={watch('batch_id')}
                            onValueChange={(value) => {
                                console.log("Selected batch ID:", Number(value)); // Debugging line
                                setValue('batch_id', Number(value));
                            }}
                        >
                            <FormLabel>Select batch: </FormLabel>
                            <SelectTrigger>
                                <SelectValue placeholder="Select batch" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectGroup>
                                    <SelectLabel>Batch</SelectLabel>
                                    {Array.isArray(batches) && batches.map((batch, index) => (
                                        <SelectItem key={index} value={batch.ID}>{batch.batch}</SelectItem>
                                    ))}
                                </SelectGroup>
                            </SelectContent>
                        </Select>

                        <Select
                            value={watch('program_id')}
                            onValueChange={(value) => setValue('program_id', Number(value))}
                        >
                            <FormLabel>Select program: </FormLabel>
                            <SelectTrigger>
                                <SelectValue placeholder="Select program" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectGroup>
                                    <SelectLabel>Program</SelectLabel>
                                    {Array.isArray(programs) && programs.map((program, index) => (
                                        <SelectItem key={index} value={program.ID}>{program.program_name}</SelectItem>
                                    ))}
                                </SelectGroup>
                            </SelectContent>
                        </Select>
                        <input
                            type="file"
                            accept=".csv"
                            onChange={handleFileChange}
                            className='mt-4'
                        />
                        {isDrawerOpen && (
                            <DrawerFooter>
                                <DrawerClose>
                                    <Button
                                        type="submit"
                                        className='w-full mt-6'
                                    >
                                        Submit
                                    </Button>
                                    <Button variant="outline" className='w-full mt-2'>Cancel</Button>
                                </DrawerClose>
                            </DrawerFooter>
                        )}

                    </form>
                </Form>
            </CardContent>
        </Card>
    )
}

export default BulkStudentForm;