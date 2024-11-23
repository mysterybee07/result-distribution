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

import { useEffect } from "react";

export default function CreateNoticeForm({ mode = 'create', notice = {} }) {
    console.log("🚀 ~ CreateNoticeForm ~ mode:", mode);
    console.log("🚀 ~ CreateNoticeForm ~ notice:", notice);
    const navigate = useNavigate();
    // Fetching data from context
    const { programs, batches } = useData();
    console.log("🚀 ~ CreateNoticeForm ~ batches:", batches);

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
        queryKey: ['semesters', programId], // Use programId in query key for caching
        queryFn: async () => {
            const response = await api.get(`/semester/by-program/${programId}`);
            return response.data.semesters;
        },
        enabled: Boolean(programId), // Only fetch semesters if programId is available
        refetchOnWindowFocus: false, // Optional: avoid refetching on window focus
        refetchOnReconnect: false, // Optional: avoid refetching on reconnect
    });

    const [isSubmitting, setIsSubmitting] = useState(false);

    const { mutate: handleSubmitAction } = useMutation({
        mutationFn: async (formData) => {
            const url =
                mode === "create"
                    ? "http://127.0.0.1:3000/notice/create"
                    : `http://127.0.0.1:3000/notice/update/${notice.ID}`;
            const method = mode === "create" ? "post" : "put";
            const response = await axios[method](url, formData);
            return response.data;
        },
        onSuccess: () => {
            navigate("/dashboard");
        },
        onError: (error) => {
            console.error(`Error ${mode === "create" ? "creating" : "editing"} notice:`, error);
        },
    });

    const onSubmit = (values) => {
        console.log("🚀 ~ onSubmit ~ values:", values);
        // Create a FormData object
        const formData = new FormData();
        formData.append("title", values.title);
        formData.append("description", values.description);
        formData.append("program_id", values.program_id.toString());
        formData.append("semester_id", values.semester_id.toString());
        formData.append("batch_id", values.batch_id.toString());
        formData.append("file_path", values.file_path);
        console.log("🚀 ~ onSubmit ~ formData:", formData);
        setIsSubmitting(true);
        handleSubmitAction(formData);
    };

    // Set form values after `notice` data is available
    useEffect(() => {
        if (mode === "edit" && notice) {
            setValue("title", notice.title || "");
            setValue("description", notice.description || "");
            setValue("program_id", notice.program_id || "");
            setValue("batch_id", notice.batch_id || "");
            setValue("semester_id", notice.semester_id || "");
        }
    }, [mode, notice, setValue, programId]);

    if (loadingSemesters) {
        return <p>Loading semesters...</p>;  // Show loading state for semesters
    }

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
                                {programs?.map((program) => (
                                    <SelectItem key={program.ID} value={program.ID}>
                                        {program.program_name}
                                    </SelectItem>
                                ))}
                            </SelectGroup>
                        </SelectContent>
                    </Select>

                    <Select
                        value={watch('semester_id')}
                        onValueChange={(value) => {
                            setValue('semester_id', Number(value));
                        }}
                    >
                        <FormLabel>Select Semester: </FormLabel>
                        <SelectTrigger>
                            <SelectValue placeholder="Select Semester" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectGroup>
                                <SelectLabel>Semesters</SelectLabel>
                                {semesters?.map((semester) => (
                                    <SelectItem key={semester.ID} value={semester.ID}>
                                        {semester.semester_name}
                                    </SelectItem>
                                ))}
                            </SelectGroup>
                        </SelectContent>
                    </Select>
                </div>
                <Select
                    value={watch('batch_id')}
                    onValueChange={(value) => {
                        setValue('batch_id', Number(value));
                    }}
                >
                    <FormLabel>Select Batch: </FormLabel>
                    <SelectTrigger>
                        <SelectValue placeholder="Select Batch" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectGroup>
                            <SelectLabel>Batches</SelectLabel>
                            {batches?.map((batch) => (
                                <SelectItem key={batch.ID} value={batch.ID}>
                                    {batch.batch}
                                </SelectItem>
                            ))}
                        </SelectGroup>
                    </SelectContent>
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
                    {isSubmitting
                        ? mode === "create"
                            ? "Creating..."
                            : "Updating..."
                        : mode === "create"
                            ? "Create Notice"
                            : "Update Notice"}
                </Button>
            </form>
        </Form>
    );
}

