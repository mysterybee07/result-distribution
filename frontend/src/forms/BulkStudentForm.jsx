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
});

const BulkStudentForm = ({ isDrawerOpen, setIsDrawerOpen }) => {
    const navigate = useNavigate();

    // fetching batch and program data
    const { programs, loadingPrograms, errorPrograms, batches, loadingBatches, errorBatches } = useData();

    if (loadingPrograms || loadingBatches) return <div>Loading...</div>;

    if (errorPrograms) return <div>Error loading programs: {errorPrograms.message}</div>;
    if (errorBatches) return <div>Error loading batches: {errorBatches.message}</div>;
    const [students, setStudents] = useState([]);

    const handleFileChange = (event) => {
        const file = event.target.files[0];

        Papa.parse(file, {
            header: true,
            complete: (results) => {
                setStudents(results.data); // Save parsed data
            },
        });
    };

    // define the form
    const form = useForm({
        resolver: zodResolver(formSchema),
        defaultValues: {
            batch_id: '',
            program_id: '',
        },
    });
    const { register, formState: { errors }, setValue, watch } = form;

    // mutation for creating student
    const { mutate: createStudent } = useMutation({
        mutationFn: async (student) => {
            const response = await api.post('/students/create', student);
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
        console.log("Submitting student data:", data);
        const student = {
            batch_id: data.batch_id,
            program_id: data.program_id,
            students: students.map(student => ({
                fullname: student.fullname,
                symbol_number: student.symbol_number,
                registration_number: student.registration_number,
            })),
        }
        createStudent(student);
    };
    return (
        // <div className='flex items-center justify-center'>
        <Card className="w-1/2 shadow-lg hover:shadow-2xl py-8 text-start">
            {/* <CardHeader>
                    <CardTitle>Add Student in Bulk</CardTitle>
                    <CardDescription>
                        <a href="/student_test.xlsx" download>
                            <Button variant="link">
                                Download Sample File
                            </Button>
                        </a>
                    </CardDescription>
                </CardHeader> */}
            <CardContent>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)}>
                        <Select
                            value={watch('batch_id')}
                            onValueChange={(value) => {
                                console.log("Selected batch ID:", Number(value)); // Debugging line
                                setValue('batch_id', Number(value));
                            }}                            >
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
        // </div>
    )
}

export default BulkStudentForm