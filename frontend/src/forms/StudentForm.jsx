import React from 'react'
import { z } from 'zod'
import {
    Form,
    FormControl,
    FormDescription,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form"

const formSchema = z.object({
    symbol_number: z.string(),
    registration_number: z.string(),
    fullname: z.string(),
    batch_id: z.number(),
    program_id: z.number(),
    current_semester: z.number(),
    status: z.string(),
});
const StudentForm = () => {
    // define the form
    const form = useForm({
        resolver: zodResolver(formSchema),
        defaultValues: {
            symbol_number: '',
            registration_number: '',
            fullname: '',
            batch_id: 0,
            program_id: 0,
            current_semester: 0,
            status: 'inactive',
        },
    });
    const { register, handleSubmit, formState: { errors } } = form;

    // define the submit handler
    const onSubmit = async (data) => {
        console.log(data);
        // try {
        //     await api.post('/students', data);
        // } catch (error) {
        //     console.error(error);
        // }
    };
    return (
        <Form {...form}>
            <form onSubmit={handleSubmit(onSubmit)}>
                <FormItem>
                    <FormLabel htmlFor="symbol_number">Symbol Number</FormLabel>
                    <FormControl
                        type="text"
                        id="symbol_number"
                        {...register('symbol_number')}
                    />
                    <FormDescription>Enter the symbol number of the student</FormDescription>
                    <FormMessage>{errors.symbol_number?.message}</FormMessage>
                </FormItem>

                <FormItem>
                    <FormLabel htmlFor="registration_number">Registration Number</FormLabel>
                    <FormControl
                        type="text"
                        id="registration_number"
                        {...register('registration_number')}
                    />
                    <FormDescription>Enter the registration number of the student</FormDescription>
                    <FormMessage>{errors.registration_number?.message}</FormMessage>
                </FormItem>

                <FormItem>
                    <FormLabel htmlFor="fullname">Full Name</FormLabel>
                    <FormControl
                        type="text"
                        id="fullname"
                        {...register('fullname')}
                    />
                    <FormDescription>Enter the full name of the student</FormDescription>
                    <FormMessage>{errors.fullname?.message}</FormMessage>
                </FormItem>

                <FormItem>
                    <FormLabel htmlFor="batch_id">Batch</FormLabel>
                    <FormControl
                        type="number"
                        id="batch_id"
                        {...register('batch_id')}
                    />
                    <FormDescription>Enter the batch of the student</FormDescription>
                    <FormMessage>{errors.batch_id?.message}</FormMessage>
                </FormItem>

                <FormItem>
                    <FormLabel htmlFor="program_id">Program</FormLabel>
                    <FormControl
                        type="number"
                        id="program_id"
                        {...register('program_id')}
                    />
                    <FormDescription>Enter the program of the student</FormDescription>
                    <FormMessage>{errors.program_id?.message}</FormMessage>
                </FormItem>

                <FormItem>
                    <FormLabel htmlFor="current_semester">Current Semester</FormLabel>
                    <FormControl
                        type="number"
                        id="current_semester"
                        {...register('current_semester')}
                    />
                    <FormDescription>Enter the current semester of the student</FormDescription>
                    <FormMessage>{errors.current_semester?.message}</FormMessage>
                </FormItem>

                <FormItem>
                    <FormLabel htmlFor="status">Status</FormLabel>
                    <FormControl
                        type="text"
                        id="status"
                        {...register('status')}
                    />
                    <FormDescription>Enter the status of the student</FormDescription>
                    <FormMessage>{errors.status?.message}</FormMessage>
                </FormItem>

                <FormItem>
                    <FormControl type="submit" />
                </FormItem>

            </form>
        </Form>
    )
}

export default StudentForm