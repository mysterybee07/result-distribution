import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import {
    Form,
    FormControl,
    FormDescription,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form"
import { Input } from "../components/ui/input"
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
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
} from "@/components/ui/tooltip"
import { FaRegEye, FaRegEyeSlash } from "react-icons/fa";


import { useState } from "react"
import { useData } from "../context/DataContext"

const formSchema = z
    .object({
        symbol_number: z.string().min(5, { message: "Symbol number must be at least 5 characters." }),
        registration_number: z.string().min(5, { message: "Registration number must be at least 5 characters." }),
        batch_id: z.number(),
        programIDStr: z.number(),
        email: z
            .string()
            .email({ message: "Please enter a valid email address." })
            .min(5, { message: "Email must be at least 5 characters." }),
        password: z
            .string()
            .min(4, { message: "Password must be at least 4 characters." }),
        confirm_password: z.string(),
        // image: z.string(),
        image: z
            .instanceof(File) // Check that the input is an instance of File
            .refine(file => file.size > 0, { message: "Image file is required." })
    })
    .refine((data) => data.password === data.confirm_password, {
        path: ["confirm_password"],
        message: "Passwords do not match",
    });


export function RegisterForm() {
    // 1. Define your form.
    const form = useForm({
        resolver: zodResolver(formSchema),
        defaultValues: {
            symbol_number: "",
            registration_number: "",
            batch_id: "",
            programIDStr: "",
            email: "",
            password: "",
            confirm_password: "",
            image: ""
        },
    })
    const { control, setValue, setError, formState: { errors }, register, watch } = form;

    // 2. Define a submit handler.
    const onSubmit = async (values) => {
        console.log("🚀 ~ onSubmit ~ values:", values);

        const { email, password, symbol_number, registration_number, batch_id, programIDStr, image } = values;

        const formData = new FormData();
        formData.append('symbol_number', symbol_number);
        formData.append('registration_number', registration_number);
        formData.append('batch_id', batch_id);
        formData.append('program_id', programIDStr);
        formData.append('email', email);
        formData.append('password', password);

        if (image instanceof File) {
            formData.append('image_url', image); // Properly append file
        } else {
            console.error('No valid image file selected.');
        }

        try {
            const response = await fetch('http://127.0.0.1:3000/user/register', {
                method: 'POST',
                body: formData,
            });

            const responseData = await response.json();

            if (!response.ok) {
                console.log("🚀 ~ onSubmit ~ errorData:", responseData);
                if (responseData.message.includes("Symbol Number")) {
                    setError("symbol_number", { message: responseData.message });
                } else if (responseData.message.includes("Registration Number")) {
                    setError("registration_number", { message: responseData.message });
                }
            } else {
                alert("Student registered successfully!");
                form.reset(); // Reset form after successful submission
            }
        } catch (error) {
            console.error("Error registering student:", error);
        }
    };


    // 3. password visibility
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);

    const togglePasswordVisibility = () => {
        setShowPassword(!showPassword);
    };
    const toggleConfirmPasswordVisibility = () => {
        setShowConfirmPassword(!showConfirmPassword);
    }
    // fetch program and batch data
    const { programs, errorPrograms, batches, errorBatches, college } = useData();


    return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8" encType="multipart/form-data">
                {/* <input type="file" name="image" accept="images/*" /> */}
                {errors.symbol_number && <p role="alert" className="text-red-500">{errors.symbol_number.message}</p>}
                {errors.registration_number && <p role="alert" className="text-red-500">{errors.registration_number.message}</p>}
                <FormField
                    control={control}
                    name="image"
                    render={({ field }) => (
                        <FormItem className="text-start">
                            <FormLabel className="text-start">Image</FormLabel>
                            <FormControl>
                                <Input
                                    type="file"
                                    accept="image/*"
                                    onChange={(e) => {
                                        const file = e.target.files[0];
                                        if (file) {
                                            setValue("image", file); // Store file in React Hook Form
                                        }
                                    }}
                                />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />


                <FormField
                    control={control}
                    name="email"
                    render={({ field }) => (
                        <FormItem className="text-start">
                            <FormLabel className="text-start">email</FormLabel>
                            <FormControl>
                                <Input placeholder="Email or Symbol number" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <FormField
                    control={control}
                    name="password"
                    render={({ field }) => (
                        <FormItem className="text-start">
                            <FormLabel>Password</FormLabel>
                            <FormControl>
                                <Input
                                    type={showPassword ? 'text' : 'password'}
                                    placeholder="Password" {...field}
                                    endAdornment={
                                        <TooltipProvider>
                                            <Tooltip>
                                                <TooltipTrigger asChild>
                                                    {showPassword ?
                                                        <FaRegEye className="h-4 w-4 cursor-pointer" onClick={togglePasswordVisibility} /> :
                                                        <FaRegEyeSlash className="h-4 w-4 cursor-pointer" onClick={togglePasswordVisibility} />}
                                                </TooltipTrigger>
                                            </Tooltip>
                                        </TooltipProvider>
                                    }
                                />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <FormField
                    control={control}
                    name="confirm_password"
                    rules={{
                        required: 'Please confirm your password',
                        validate: (value) => value === form.getValues('password') || 'Passwords do not match',
                    }}
                    render={({ field }) => (
                        <FormItem className="text-start">
                            <FormLabel>Confirm Password</FormLabel>
                            <FormControl>
                                <Input
                                    type={showConfirmPassword ? 'text' : 'password'}
                                    placeholder="Password" {...field}
                                    endAdornment={
                                        <TooltipProvider>
                                            <Tooltip>
                                                <TooltipTrigger asChild>
                                                    {showConfirmPassword ?
                                                        <FaRegEye className="h-4 w-4 cursor-pointer" onClick={toggleConfirmPasswordVisibility} /> :
                                                        <FaRegEyeSlash className="h-4 w-4 cursor-pointer" onClick={toggleConfirmPasswordVisibility} />}
                                                </TooltipTrigger>
                                            </Tooltip>
                                        </TooltipProvider>
                                    }
                                />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <div className="flex justify-between gap-4">
                    <FormField
                        control={control}
                        name="symbol_number"
                        render={({ field }) => (
                            <FormItem className="text-start">
                                <FormLabel className="text-start">Symbol number</FormLabel>
                                <FormControl>
                                    <Input placeholder="Enter your symbol number" {...field} />
                                </FormControl>
                                {errors.symbol_number && <p role="alert" className="text-red-500">{errors.symbol_number.message}</p>}
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                    <FormField
                        control={control}
                        name="registration_number"
                        render={({ field }) => (
                            <FormItem className="text-start">
                                <FormLabel className="text-start">Registration number</FormLabel>
                                <FormControl>
                                    <Input placeholder="Enter your registration number" {...field} />
                                </FormControl>
                                {errors.symbol_number && <p role="alert" className="text-red-500">{errors.symbol_number.message}</p>}                                <FormMessage />
                            </FormItem>
                        )}
                    />
                </div>
                <div className="flex flex-row items-center gap-4">
                    <Select
                        value={watch('batch_id')}
                        onValueChange={(value) => setValue('batch_id', Number(value))} // Update form value when batch changes
                    >
                        <FormLabel className="text-start">Select your batch: </FormLabel>
                        <SelectTrigger className="w-[180px]">
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
                </div>

                <div className="flex flex-row items-center gap-4">
                    <Select
                        value={watch('programIDStr')}
                        onValueChange={(value) => setValue('programIDStr', Number(value))} // Update form value when program changes
                    >
                        <FormLabel className="text-start">Select your program: </FormLabel>
                        <SelectTrigger className="w-[180px]">
                            <SelectValue placeholder="Select your program" />
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

                <Button type="submit" size="lg" className="w-full">Register</Button>
            </form>
        </Form>
    )
}