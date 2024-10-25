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
    const { student } = singleStudent || {}; // Destructure the student data
    const initialData = student; // Set initial data if student data is available
    console.log("🚀 ~ StudentForm ~ initialData:", initialData)
    const isEditMode = !!initialData; // Check if the form is in edit mode
    const navigate = useNavigate();
    // fetching batch and program data
    const { programs, loadingPrograms, errorPrograms, batches, loadingBatches, errorBatches } = useData();

    if (loadingPrograms || loadingBatches) return <div>Loading...</div>;

    if (errorPrograms) return <div>Error loading programs: {errorPrograms.message}</div>;
    if (errorBatches) return <div>Error loading batches: {errorBatches.message}</div>;

    // define the form
    const form = useForm({
        resolver: zodResolver(formSchema),
        defaultValues: {
            symbol_number: '',
            registration_number: '',
            fullname: '',
            batch_id: '',
            program_id: '',
            current_semester: '',
        },
    });
    const { register, handleSubmit, formState: { errors }, setValue, watch, reset } = form;

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
            students: [
                {
                    fullname: data.fullname,
                    symbol_number: data.symbol_number,
                    registration_number: data.registration_number,
                }
            ]
        }
        console.log("🚀 ~ onSubmit ~ student:", student)
        createStudent(student);
    };

    return (
        <div className='flex items-center justify-center mt-16'>
            <Card className="w-1/2 shadow-lg hover:shadow-2xl py-8 text-start">
                <CardContent>
                    <h1 className="text-xl font-bold mb-4">
                        {isEditMode ? 'Edit Student' : 'Create Student'}
                    </h1>
                    <Form {...form}>
                        <form onSubmit={form.handleSubmit(onSubmit)}>
                            <FormField
                                control={form.control}
                                name="fullname"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Full Name</FormLabel>
                                        <FormControl>
                                            <Input placeholder="Full Name" {...field} />
                                        </FormControl>
                                        <FormMessage>{errors.fullname?.message}</FormMessage>
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="symbol_number"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Symbol Number</FormLabel>
                                        <FormControl>
                                            <Input placeholder="Symbol number" {...field} />
                                        </FormControl>
                                        <FormMessage>{errors.symbol_number?.message}</FormMessage>
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="registration_number"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Register Number</FormLabel>
                                        <FormControl>
                                            <Input placeholder="Register number" {...field} />
                                        </FormControl>
                                        <FormMessage>{errors.registration_number?.message}</FormMessage>
                                    </FormItem>
                                )}
                            />

                            <Select
                                value={watch('batch_id')}
                                onValueChange={(value) => {
                                    console.log("Selected batch ID:", Number(value)); // Debugging line
                                    setValue('batch_id', Number(value));
                                }}                            >
                                <FormLabel>Select your batch: </FormLabel>
                                <SelectTrigger>
                                    <SelectValue placeholder="Select your batch" />
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
                                <FormLabel>Select your program: </FormLabel>
                                <SelectTrigger>
                                    <SelectValue placeholder="Select your program" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectGroup>
                                        <SelectLabel>Program</SelectLabel>
                                        {Array.isArray(programs) && programs.map((program, index) => (
                                            <SelectItem key={index} value={program.ID}>{program.program_name}</SelectItem>
                                        ))}
                                        {/* <SelectItem value="1">CSIT</SelectItem>
                                        <SelectItem value="2">BBA</SelectItem>
                                        <SelectItem value="3">BIM</SelectItem>
                                        <SelectItem value="4">BBS</SelectItem> */}
                                    </SelectGroup>
                                </SelectContent>
                            </Select>

                            <Select
                                // value={watch('current_semester')}
                                onValueChange={(value) => setValue('current_semester', Number(value))}
                            >
                                <FormLabel>Select your semester: </FormLabel>
                                <SelectTrigger>
                                    <SelectValue placeholder="Select your semester" />
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

                            <Button type="submit" className="self-end mt-4">
                                {isEditMode ? 'Update Student' : 'Create Student'}
                            </Button>
                        </form>
                    </Form>
                </CardContent>
            </Card>
        </div>
    )
}

export default StudentForm
