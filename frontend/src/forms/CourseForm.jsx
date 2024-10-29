import { zodResolver } from '@hookform/resolvers/zod';
import React, { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form';
import { z } from 'zod'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '../components/ui/card';
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
import { ScrollArea } from "@/components/ui/scroll-area"

import { useData } from '../context/DataContext';
import { Input } from '../components/ui/input';
import { Button } from '../components/ui/button';
import { useMutation, useQuery } from '@tanstack/react-query';
import api from '../api';
import { useNavigate } from 'react-router-dom';
import { useToast } from "@/hooks/use-toast";

const formSchema = z.object({
    course_code: z.string().min(2, { message: 'Course code must be at least 2 characters long' }),
    name: z.string().min(2, { message: 'Course name must be at least 2 characters long' }),
    semester_pass_marks: z.string().min(0, { message: 'Semester pass marks must be at least 0' }),
    practical_pass_marks: z.string().min(0, { message: 'Practical pass marks must be at least 0' }),
    assistant_pass_marks: z.string().min(0, { message: 'Assistant pass marks must be at least 0' }),
    semester_total_marks: z.string().min(0, { message: 'Semester total marks must be at least 0' }).max(100, { message: 'Semester total marks must be at most 100' }),
    practical_total_marks: z.string().min(0, { message: 'Practical total marks must be at least 0' }),
    assistant_total_marks: z.string().min(0, { message: 'Assistant total marks must be at least 0' }),
    program_id: z.number(),
    semester_id: z.number(),
})

const CourseForm = () => {
    const navigate = useNavigate();
    const { toast } = useToast();
    // fetching data from context
    const { programs, loadingPrograms, errorPrograms,
        // semesters, loadingSemesters, errorSemesters,
    } = useData();
    const [courses, setCourses] = useState([]); // State to store the array of student data
    console.log("🚀 ~ CourseForm ~ courses:", courses)


    // define the form
    const form = useForm({
        resolver: zodResolver(formSchema),
        defaultValues: {
            course_code: '',
            name: '',
            semester_pass_marks: '',
            practical_pass_marks: '',
            assistant_pass_marks: '',
            semester_total_marks: '',
            practical_total_marks: '',
            assistant_total_marks: '',
            program_id: '',
            semester_id: '',
        }
    });
    const { register, handleSubmit, formState: { errors }, setValue, watch, reset } = form;
    // Watch for changes in 'program_id'
    const programId = watch('program_id');
    console.log("🚀 ~ CourseForm ~ programId:", programId)

    // UseQuery to fetch semesters based on selected program_id
    const { data: semesters, isLoading: loadingSemesters, error: errorSemesters } = useQuery({
        queryKey: ['semesters', programId],  // Use programId in query key for caching
        queryFn: async () => {
            const response = await api.get(`/semester/by-program/${programId}`);
            console.log("🚀 ~ response:", response);
            return response.data.semesters;
        },
        enabled: Boolean(programId), // Only fetch semesters if programId is available
        refetchOnWindowFocus: false, // Optional: avoid refetching on window focus
        refetchOnReconnect: false, // Optional: avoid refetching on reconnect
    });


    const formValues = watch();
    // Function to handle "Next" button
    const addOne = () => {
        console.log('Current Form Values:', formValues); // Display the current form data in console

        // Add the current form values to the students array
        setCourses((prevCourse) => [...prevCourse, formValues]);

        // Clear the form for the next input
        // reset();
    };
    // handle form submission
    const { mutate: createCourse } = useMutation({
        mutationFn: async (course) => {
            const response = await api.post('/courses/create', course);
            return response.data;
        },
        onSuccess: (data) => {
            navigate('/admin/courses');
            toast({
                // title: "Student Added",
                description: JSON.stringify(data.message),
                variant: "success",
            })
        },
        onError: (error) => {
            console.error(`Error: ---creating student error---:`, error);
        },
    }); 
    const onSubmit = async (data) => {
        console.log(data);
        createCourse(data);
    }

    return (
        <div className='flex items-center justify-center mt-16'>
            <Card className='w-full shadow-lg hover:shadow-2xl py-8 text-start'>
                <CardHeader>
                    <CardTitle className='self-center text-lg font-semibold'>
                        Course Add Form
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <Form {...form}>
                        <form onSubmit={form.handleSubmit(onSubmit)} className='flex gap-16'>
                            <div className='w-1/4'>
                                {/* Program and Semester */}
                                <p class="text-lg mb-2 font-semibold">Program and Semester:</p>
                                <div className='flex gap-2 flex-col'>
                                    <div className='w-full'>
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
                                    </div>

                                    <div className='w-full'>
                                        <Select
                                            value={watch('semester_id')}
                                            onValueChange={(value) => {
                                                console.log("Selected batch ID:", Number(value)); // Debugging line
                                                setValue('semester_id', Number(value));
                                            }}                            >
                                            <FormLabel>Select Semester: </FormLabel>
                                            <SelectTrigger>
                                                <SelectValue placeholder="Select Semester" />
                                            </SelectTrigger>
                                            <SelectContent>
                                                <SelectGroup>
                                                    <SelectLabel>Semesters</SelectLabel>
                                                    {Array.isArray(semesters) && semesters.map((semester, index) => (
                                                        <SelectItem key={index} value={semester.ID}>{semester.semester_name}</SelectItem>
                                                    ))}
                                                </SelectGroup>
                                            </SelectContent>
                                        </Select>
                                    </div>
                                </div>

                                {/* Course */}
                                <p class="text-lg mt-4 mb-2 font-semibold">Course:</p>
                                <div className='flex flex-col gap-2'>
                                    <FormField
                                        control={form.control}
                                        name="course_code"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>Course Code</FormLabel>
                                                <FormControl>
                                                    <Input placeholder="course_code" {...field} />
                                                </FormControl>
                                                <FormMessage>{errors.course_code?.message}</FormMessage>
                                            </FormItem>
                                        )}
                                    />
                                    <FormField
                                        control={form.control}
                                        name="name"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>Course name</FormLabel>
                                                <FormControl>
                                                    <Input placeholder="Course name" {...field} />
                                                </FormControl>
                                                <FormMessage>{errors.name?.message}</FormMessage>
                                            </FormItem>
                                        )}
                                    />
                                </div>
                            </div>

                            <div className='w-3/4'>
                                {/* Marks */}
                                <p class="text-lg mb-2 font-semibold">Marks:</p>
                                <div className='flex gap-16 mt-4'>
                                    <div className='w-1/2'>
                                        <FormField
                                            control={form.control}
                                            name="semester_total_marks"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormLabel>Semester Total Marks</FormLabel>
                                                    <FormControl>
                                                        <Input placeholder="semester_total_marks" {...field} type='number' min={0} max={100} />
                                                    </FormControl>
                                                    <FormMessage>{errors.semester_total_marks?.message}</FormMessage>
                                                </FormItem>
                                            )}
                                        />
                                    </div>
                                    <div className='w-1/2'>
                                        <FormField
                                            control={form.control}
                                            name="semester_pass_marks"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormLabel>Semester Pass Marks</FormLabel>
                                                    <FormControl>
                                                        <Input placeholder="semester_pass_marks" {...field} type='number' min={0} />
                                                    </FormControl>
                                                    <FormMessage>{errors.semester_pass_marks?.message}</FormMessage>
                                                </FormItem>
                                            )}
                                        />
                                    </div>
                                </div>

                                <div className='flex gap-16 mt-4'>
                                    <div className='w-1/2'>
                                        <FormField
                                            control={form.control}
                                            name="practical_total_marks"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormLabel>Pratical Total Marks</FormLabel>
                                                    <FormControl>
                                                        <Input placeholder="practical_total_marks" {...field} type='number' min={0} />
                                                    </FormControl>
                                                    <FormMessage>{errors.practical_total_marks?.message}</FormMessage>
                                                </FormItem>
                                            )}
                                        />
                                    </div>
                                    <div className='w-1/2'>
                                        <FormField
                                            control={form.control}
                                            name="practical_pass_marks"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormLabel>Pratical Pass Marks</FormLabel>
                                                    <FormControl>
                                                        <Input placeholder="practical_pass_marks" {...field} type='number' min={0} />
                                                    </FormControl>
                                                    <FormMessage>{errors.practical_pass_marks?.message}</FormMessage>
                                                </FormItem>
                                            )}
                                        />
                                    </div>
                                </div>

                                <div className='flex gap-16 mt-4'>
                                    <div className='w-1/2'>
                                        <FormField
                                            control={form.control}
                                            name="assistant_total_marks"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormLabel>Assistant Total Marks</FormLabel>
                                                    <FormControl>
                                                        <Input placeholder="assistant_total_marks" {...field} type='number' min={0} />
                                                    </FormControl>
                                                    <FormMessage>{errors.assistant_total_marks?.message}</FormMessage>
                                                </FormItem>
                                            )}
                                        />
                                    </div>
                                    <div className='w-1/2'>
                                        <FormField
                                            control={form.control}
                                            name="assistant_pass_marks"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormLabel>Assistant Pass Marks</FormLabel>
                                                    <FormControl>
                                                        <Input placeholder="assistant_pass_marks" {...field} type='number' min={0} />
                                                    </FormControl>
                                                    <FormMessage>{errors.assistant_pass_marks?.message}</FormMessage>
                                                </FormItem>
                                            )}
                                        />
                                    </div>
                                </div>


                                <div className='flex justify-end gap-4 mt-8'>
                                    <Button variant="secondary" onClick={addOne} type='button'>Next</Button>
                                    <Button type='submit'>Submit</Button>
                                </div>
                            </div>

                        </form>
                        {/* added courses */}
                        {courses.length > 0 && (
                            <div className='flex justify-end gap-4 mt-8'>
                                <ScrollArea className="h-[200px] w-full rounded-md border p-4">
                                    <ul className="w-full">
                                        {/* Header Row */}
                                        <li
                                            className="grid grid-cols-5 gap-4 font-semibold border-b pb-2 sticky top-0 bg-white z-10 text-center"
                                        >
                                            <div>
                                                <div className="grid grid-cols-2 gap-2">
                                                    <p>Course Code</p>
                                                    <p>Course Name</p>
                                                </div>
                                            </div>
                                            <div>
                                                <p className='text-center'>Semester</p>
                                                <div className="grid grid-cols-2 gap-2">
                                                    <p>Total Marks</p>
                                                    <p>Pass Marks</p>
                                                </div>
                                            </div>
                                            <div>
                                                <p className='text-center'>Practical</p>
                                                <div className="grid grid-cols-2 gap-2">
                                                    <p>Total Marks</p>
                                                    <p>Pass Marks</p>
                                                </div>
                                            </div>
                                            <div>
                                                <p className='text-center'>Assistant</p>
                                                <div className="grid grid-cols-2 gap-2">
                                                    <p>Total Marks</p>
                                                    <p>Pass Marks</p>
                                                </div>
                                            </div>
                                            <div>
                                                Action
                                            </div>
                                        </li>

                                        {/* Data Rows */}
                                        {courses.map((course, index) => (
                                            <li key={index} className="grid grid-cols-5 gap-4 py-2 border-b text-center">
                                                <div className="grid grid-cols-2 gap-2 text-left">
                                                    <p>{course.course_code}</p>
                                                    <p>{course.Name}</p>
                                                </div>
                                                <div className="grid grid-cols-2 gap-2">
                                                    <p>{course.semester_total_marks}</p>
                                                    <p>{course.semester_pass_marks}</p>
                                                </div>
                                                <div className="grid grid-cols-2 gap-2">
                                                    <p>{course.practical_total_marks}</p>
                                                    <p>{course.practical_pass_marks}</p>
                                                </div>
                                                <div className="grid grid-cols-2 gap-2">
                                                    <p>{course.assistant_total_marks}</p>
                                                    <p>{course.assistant_pass_marks}</p>
                                                </div>
                                                <div>
                                                    <Button variant="destructive" onClick={() => {
                                                        setCourses((prevCourse) => prevCourse.filter((_, i) => i !== index));
                                                    }}>Remove</Button>
                                                </div>
                                            </li>
                                        ))}
                                    </ul>
                                </ScrollArea>

                            </div>
                        )}
                    </Form>
                </CardContent>
            </Card>
        </div>
    )
}

export default CourseForm