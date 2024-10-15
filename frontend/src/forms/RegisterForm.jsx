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


import { useEffect, useState } from "react"

const formSchema = z
    .object({
        symbol_number: z.string().min(5, { message: "Symbol number must be at least 5 characters." }),
        registration_number: z.string().min(5, { message: "Registration number must be at least 5 characters." }),
        batch_id: z.string(),
        programIDStr: z.string(),
        email: z
            .string()
            .email({ message: "Please enter a valid email address." })
            .min(5, { message: "Email must be at least 5 characters." }),
        password: z
            .string()
            .min(8, { message: "Password must be at least 6 characters." }),
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
    const { control, setValue, setError, formState: { errors }, register } = form;

    // 2. Define a submit handler.
    const onSubmit = async (values) => {
        // Do something with the form values.
        const { email, password, symbol_number, registration_number, batch_id, programIDStr, image } = values;

        const formData = new FormData();
        formData.append('symbol_number', symbol_number);
        formData.append('registration_number', registration_number);
        formData.append('batch_id', batch_id);
        formData.append('program_id', programIDStr);
        formData.append('email', email);
        formData.append('password', password);
        const fileInput = document.querySelector('input[type="file"]');
        if (fileInput && fileInput.files.length > 0) {
            formData.append('image_url', fileInput.files[0]); // use 'image' as the key
        } else {
            console.error('No image file selected.');
        }

        console.log("🚀 ~ onSubmit ~ formData:", formData)

        // ✅ This will be type-safe and validated.
        try {
            const response = await fetch('http://127.0.0.1:3000/user/register', {
                method: 'POST',
                // headers: {
                //     'Content-Type': 'application/json',
                // },
                body: formData,
            });

            if (!response.ok) {
                const errorData = await response.json();
                console.log("🚀 ~ onSubmit ~ errorData:", errorData.message)
                if (errorData.message.includes("Symbol Number")) {
                    setError("symbol_number", { message: errorData.message });
                } else if (errorData.message.includes("Registration Number")) {
                    setError("registration_number", { message: errorData.message });
                } // Set error from backend
            } else {
                alert("Student registered successfully!");
            }
        } catch (error) {
            console.error("Error registering student:", err);
        }

    }

    // 3. password visibility
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);

    const togglePasswordVisibility = () => {
        setShowPassword(!showPassword);
    };
    const toggleConfirmPasswordVisibility = () => {
        setShowConfirmPassword(!showConfirmPassword);
    }

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
                                        // Capture the file and set it in the form
                                        field.onChange(e.target.files[0]);
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
                                    <Input placeholder="Enter your symbol number" {...field} {...register('symbol_number')} />
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
                        onValueChange={(value) => setValue('batch_id', value)} // Update form value when batch changes
                    >
                        <FormLabel className="text-start">Select your batch: </FormLabel>
                        <SelectTrigger className="w-[180px]">
                            <SelectValue placeholder="Select your batch" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectGroup>
                                <SelectLabel>Batch</SelectLabel>
                                <SelectItem value="1">2024</SelectItem>
                                <SelectItem value="2">2023</SelectItem>
                                <SelectItem value="3">2022</SelectItem>
                                <SelectItem value="4">2021</SelectItem>
                                <SelectItem value="5">2020</SelectItem>
                            </SelectGroup>
                        </SelectContent>
                    </Select>
                </div>

                <div className="flex flex-row items-center gap-4">
                    <Select
                        onValueChange={(value) => setValue('programIDStr', value)} // Update form value when program changes
                    >
                        <FormLabel className="text-start">Select your program: </FormLabel>
                        <SelectTrigger className="w-[180px]">
                            <SelectValue placeholder="Select your program" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectGroup>
                                <SelectLabel>Program</SelectLabel>
                                <SelectItem value="1">CSIT</SelectItem>
                                <SelectItem value="2">BBA</SelectItem>
                                <SelectItem value="3">BIM</SelectItem>
                                <SelectItem value="4">BBS</SelectItem>
                            </SelectGroup>
                        </SelectContent>
                    </Select>
                </div>

                <Button type="submit" size="lg" className="w-full">Register</Button>
            </form>
        </Form>
    )
}