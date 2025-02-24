import React, { useEffect } from 'react'
import { z } from 'zod'
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
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
import { useForm } from 'react-hook-form';
import { zodResolver } from "@hookform/resolvers/zod";
import { Input } from '../components/ui/input';
import { Card, CardContent } from '../components/ui/card';
import { Button } from '../components/ui/button';
import { useData } from '../context/DataContext';
import { useMutation, useQuery } from '@tanstack/react-query';
import api from '../api';
import { useNavigate, useParams } from 'react-router-dom';
import { useToast } from "@/hooks/use-toast";


const formSchema = z.object({
    symbol_number: z.string().min(5, { message: "Symbol number must be at least 5 characters." }),
    registration_number: z.string().min(5, { message: "Registration number must be at least 5 characters." }),
    fullname: z.string().min(3, { message: "Full name must be at least 3 character." }),
    batch_id: z.number(),
    program_id: z.number(),
    current_semester: z.number(),
    college_id: z.number(),
});

const StudentForm = () => {
    const { toast } = useToast();

    const { id } = useParams();
    // Fetch student data based on the dynamic ID
    const { data: singleStudent, isLoading: loadingSingleStudent, error: errorSingleStudent } = useQuery({
        queryKey: ['singleStudent', id], // Add the ID to the query key for caching
        queryFn: async () => {
            const response = await api.get(`/students/${id}`);
            return response.data;
        },
        enabled: Boolean(id), // Only run the query if there's a valid ID (for create mode, it won't fetch)
    });
    const { student } = singleStudent || {};
    const initialData = student;
    console.log("🚀 ~ StudentForm ~ initialData:", initialData)

    // if (loadingSingleStudent) return <div>Loading...</div>;


    const isEditMode = !!initialData; // Check if the form is in edit mode
    const navigate = useNavigate();
    // fetching batch and program data
    const { programs, errorPrograms, batches, errorBatches, college } = useData();
    console.log("🚀 ~ StudentForm ~ college:", college)


    if (errorPrograms) return <div>Error loading programs: {errorPrograms.message}</div>;
    if (errorBatches) return <div>Error loading batches: {errorBatches.message}</div>;

    // define the form
    const form = useForm({
        resolver: zodResolver(formSchema),
        defaultValues: {
            symbol_number: '',
            registration_number: '',
            fullname: '',
            college_id: '',
            batch_id: '',
            program_id: '',
            current_semester: '',
        },
    });
    const { formState: { errors }, setValue, watch, reset } = form;

    useEffect(() => {
        if (initialData) {
            reset(initialData); // Populate form in edit mode
        }
    }, [initialData, reset]);

    // mutation for creating student
    const { mutate: createStudent } = useMutation({
        mutationFn: async (student) => {
            const endpoint = isEditMode ? `/students/update/${initialData.ID}` : '/students/create';
            const method = isEditMode ? 'put' : 'post';
            const response = await api[method](endpoint, student);
            return response.data;
        },
        onSuccess: (data) => {
            console.log(`Student ${isEditMode ? 'updated' : 'created'} successfully:`, data);
            navigate('/admin/students');
            toast({
                // title: "Student Added",
                description: JSON.stringify(data.message),
                variant: "success",
            })
        },
        onError: (error) => {
            console.error(`Error ${isEditMode ? 'updating' : 'creating'} student:`, error);
        },
    });

    // define the submit handler
    const onSubmit = (data) => {
        console.log("Submitting student data:", data);
        const student = {
            batch_id: data.batch_id,
            program_id: data.program_id,
            college_id: data.college_id,
            students: [
                {
                    fullname: data.fullname,
                    symbol_number: data.symbol_number,
                    registration_number: data.registration_number,
                    college_id: data.college_id,
                }
            ]
        }
        console.log("🚀 ~ onSubmit ~ student:", student)
        createStudent(student);
    };



    return (
        <div className='flex'>
            <div className="w-full py-8 text-start">
                <div className="space-y-2">
                    <h1 className="text-xl font-bold mb-4">
                        {isEditMode ? 'Edit Student' : 'Create Student'}
                    </h1>
                    <Form {...form}>
                        <form className="space-y-6 p-6 max-w-2xl mx-auto bg-white rounded-lg shadow-md" onSubmit={form.handleSubmit(onSubmit)}>
                            {/* Full Name */}
                            <FormField
                                control={form.control}
                                name="fullname"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel className="text-lg font-medium">Full Name</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder="Full Name"
                                                {...field}
                                                className="mt-1 p-2 border border-gray-300 rounded-md w-full"
                                            />
                                        </FormControl>
                                        <FormMessage className="text-red-500">{errors.fullname?.message}</FormMessage>
                                    </FormItem>
                                )}
                            />

                            {/* College Selection */}
                            <Select
                                value={watch('college_id')}
                                onValueChange={(value) => {
                                    console.log("Selected college ID:", Number(value));
                                    setValue('college_id', Number(value));
                                }}
                            >
                                <FormLabel className="text-lg font-medium">Select College</FormLabel>
                                <SelectTrigger className="mt-1 p-2 border border-gray-300 rounded-md w-full">
                                    <SelectValue placeholder="Select College" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectGroup>
                                        <SelectLabel>College</SelectLabel>
                                        {Array.isArray(college) && college.map((item, index) => (
                                            <SelectItem key={index} value={item.id}>
                                                {item.college_name}
                                            </SelectItem>
                                        ))}
                                    </SelectGroup>
                                </SelectContent>
                            </Select>

                            {/* Symbol Number */}
                            <FormField
                                control={form.control}
                                name="symbol_number"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel className="text-lg font-medium">Symbol Number</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder="Symbol Number"
                                                {...field}
                                                className="mt-1 p-2 border border-gray-300 rounded-md w-full"
                                            />
                                        </FormControl>
                                        <FormMessage className="text-red-500">{errors.symbol_number?.message}</FormMessage>
                                    </FormItem>
                                )}
                            />

                            {/* Registration Number */}
                            <FormField
                                control={form.control}
                                name="registration_number"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel className="text-lg font-medium">Register Number</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder="Register Number"
                                                {...field}
                                                className="mt-1 p-2 border border-gray-300 rounded-md w-full"
                                            />
                                        </FormControl>
                                        <FormMessage className="text-red-500">{errors.registration_number?.message}</FormMessage>
                                    </FormItem>
                                )}
                            />

                            {/* Batch Selection */}
                            <Select
                                value={watch('batch_id')}
                                onValueChange={(value) => {
                                    console.log("Selected batch ID:", Number(value));
                                    setValue('batch_id', Number(value));
                                }}
                            >
                                <FormLabel className="text-lg font-medium">Select Batch</FormLabel>
                                <SelectTrigger className="mt-1 p-2 border border-gray-300 rounded-md w-full">
                                    <SelectValue placeholder="Select Batch" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectGroup>
                                        <SelectLabel>Batch</SelectLabel>
                                        {Array.isArray(batches) && batches.map((batch, index) => (
                                            <SelectItem key={index} value={batch.ID}>
                                                {batch.batch}
                                            </SelectItem>
                                        ))}
                                    </SelectGroup>
                                </SelectContent>
                            </Select>

                            {/* Program Selection */}
                            <Select
                                value={watch('program_id')}
                                onValueChange={(value) => setValue('program_id', Number(value))}
                            >
                                <FormLabel className="text-lg font-medium">Select Program</FormLabel>
                                <SelectTrigger className="mt-1 p-2 border border-gray-300 rounded-md w-full">
                                    <SelectValue placeholder="Select Program" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectGroup>
                                        <SelectLabel>Program</SelectLabel>
                                        {Array.isArray(programs) && programs.map((program, index) => (
                                            <SelectItem key={index} value={program.ID}>
                                                {program.program_name}
                                            </SelectItem>
                                        ))}
                                    </SelectGroup>
                                </SelectContent>
                            </Select>

                            {/* Semester Selection */}
                            <Select
                                onValueChange={(value) => setValue('current_semester', Number(value))}
                            >
                                <FormLabel className="text-lg font-medium">Select Semester</FormLabel>
                                <SelectTrigger className="mt-1 p-2 border border-gray-300 rounded-md w-full">
                                    <SelectValue placeholder="Select Semester" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectGroup>
                                        <SelectLabel>Semester</SelectLabel>
                                        <SelectItem value="1">1st</SelectItem>
                                        <SelectItem value="2">2nd</SelectItem>
                                        <SelectItem value="3">3rd</SelectItem>
                                        <SelectItem value="4">4th</SelectItem>
                                    </SelectGroup>
                                </SelectContent>
                            </Select>

                            {/* Submit Button */}
                            <Button
                                type="submit"
                                className="mt-6 w-full py-3 px-6 bg-blue-600 text-white font-semibold rounded-md hover:bg-blue-700 transition-colors"
                            >
                                {isEditMode ? 'Update Student' : 'Create Student'}
                            </Button>
                        </form>
                    </Form>

                </div>
            </div>
        </div>
    )
}

export default StudentForm
