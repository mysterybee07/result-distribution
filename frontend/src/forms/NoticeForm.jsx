import { useMutation, useQuery } from "@tanstack/react-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { Button } from "@/components/ui/button";
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
} from "@/components/ui/select"
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import api from "../api";
import { useData } from "../context/DataContext";
import axios from "axios";

// Define form schema
const formSchema = z.object({
    title: z.string().nonempty({ message: "Title is required" }),
    description: z.string().nonempty({ message: "Description is required" }),
    program_id: z.number().min(1, { message: "Program is required." }),
    semester_id: z.number().min(1, { message: "Semester is required." }),
    batch_id: z.number().min(1, { message: "Batch is required." }),
    file_path: z.any().refine((file) => file instanceof File, {
        message: "A valid file must be uploaded",
    }),
});

export default function CreateNoticeForm() {
    const navigate = useNavigate();
    // fetching data from context
    const { programs, batches } = useData();
    console.log("🚀 ~ CreateNoticeForm ~ batches:", batches)


    const form = useForm({
        resolver: zodResolver(formSchema),
        defaultValues: {
            title: "",
            description: "",
            programID: "",
            batchID: "",
            semesterID: "",
            file_path: null,
        },
    });
    const { watch, setValue } = form;
    // Watch for changes in 'program_id'
    const programId = watch('program_id');

    // UseQuery to fetch semesters based on selected program_id
    const { data: semesters, isLoading: loadingSemesters, error: errorSemesters } = useQuery({
        queryKey: ['semesters', programId],  // Use programId in query key for caching
        queryFn: async () => {
            const response = await api.get(`/semester/by-program/${programId}`);
            // console.log("🚀 ~ response:", response);
            return response.data.semesters;
        },
        enabled: Boolean(programId), // Only fetch semesters if programId is available
        refetchOnWindowFocus: false, // Optional: avoid refetching on window focus
        refetchOnReconnect: false, // Optional: avoid refetching on reconnect
    });

    const [isSubmitting, setIsSubmitting] = useState(false);

    const { mutate: createNotice } = useMutation({
        mutationFn: async (formData) => {
            const data = new FormData();
            // for (const key in formData) {
            //     data.append(key, formData[key]);
            // }
            // const response = await api.post("/notice/create", formData);
            const response = await axios.post('http://127.0.0.1:3000/notice/create', formData);
            console.log("🚀 ~ response:", response);
            return response.data;
        },
        onSuccess: () => {
            navigate("/admin/notice");
        },
        onError: (error) => {
            console.error("Error creating notice:", error);
        },
    });

    const onSubmit = (values) => {
        console.log("🚀 ~ onSubmit ~ values:", values)
        // Create a FormData object
        const formData = new FormData();
        formData.append("title", values.title);
        formData.append("description", values.description);
        formData.append("program_id", values.program_id.toString());
        formData.append("semester_id", values.semester_id.toString());
        formData.append("batch_id", values.batch_id.toString());
        formData.append("file_path", values.file_path);
        console.log("🚀 ~ onSubmit ~ formData:", formData)
        setIsSubmitting(true);
        createNotice(formData);
    };

    return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8" encType="multipart/form-data">
                <FormField
                    control={form.control}
                    name="title"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Title</FormLabel>
                            <FormControl>
                                <Input placeholder="Notice Title" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="description"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Description</FormLabel>
                            <FormControl>
                                <Textarea placeholder="Notice Description" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                {/* Program and Semester */}
                <p class="text-lg mb-2 font-semibold">Program and Semester:</p>
                <div className='flex gap-2 flex-col'>
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
                        {/* <FormMessage>{errors.program_id?.message}</FormMessage> */}
                    </Select>

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
                        {/* <FormMessage>{errors.semester_id?.message}</FormMessage> */}
                    </Select>
                </div>
                <Select
                    value={watch('batch_id')}
                    onValueChange={(value) => {
                        console.log("Selected batch ID:", Number(value)); // Debugging line
                        setValue('batch_id', Number(value));
                    }}                            >
                    <FormLabel>Select Batch: </FormLabel>
                    <SelectTrigger>
                        <SelectValue placeholder="Select Batch" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectGroup>
                            <SelectLabel>Batchs</SelectLabel>
                            {Array.isArray(batches) && batches.map((batch, index) => (
                                <SelectItem key={index} value={batch.ID}>{batch.batch}</SelectItem>
                            ))}
                        </SelectGroup>
                    </SelectContent>
                    {/* <FormMessage>{errors.semester_id?.message}</FormMessage> */}
                </Select>
                <FormField
                    control={form.control}
                    name="file_path"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Upload File</FormLabel>
                            <FormControl>
                                <Input
                                    type="file"
                                    onChange={(e) => field.onChange(e.target.files[0])}
                                />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <Button type="submit" size="lg" className="w-full" disabled={isSubmitting}>
                    {isSubmitting ? "Creating..." : "Create Notice"}
                </Button>
            </form>
        </Form>
    );
}
