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
        batchIDStr: z.string(),
        programIDStr: z.string(),
        identifier: z
            .string()
            .email({ message: "Please enter a valid email address." })
            .min(5, { message: "Email must be at least 5 characters." }),
        password: z
            .string()
            .min(6, { message: "Password must be at least 6 characters." }),
        confirm_password: z.string(),
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
            batchIDStr: "",
            programIDStr: "",
            identifier: "",
            password: "",
        },
    })
    const { setValue } = form;

    // 2. Define a submit handler.
    const onSubmit = async (values) => {
        // Do something with the form values.
        // ✅ This will be type-safe and validated.
        const { identifier, password, symbol_number, registration_number, batchIDStr, programIDStr } = values;
        // const response = await fetch('http://127.0.0.1:3000/login', {
        //     method: 'POST',
        //     headers: {
        //         'Content-Type': 'application/json',
        //     },
        //     body: JSON.stringify({ identifier, password }),
        // });

        // const data = await response.json();
        // console.log(data);
        console.log(values)
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
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
                <FormField
                    control={form.control}
                    name="identifier"
                    render={({ field }) => (
                        <FormItem className="text-start">
                            <FormLabel className="text-start">Identifier</FormLabel>
                            <FormControl>
                                <Input placeholder="Email or Symbol number" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
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
                    control={form.control}
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
                        control={form.control}
                        name="symbol_number"
                        render={({ field }) => (
                            <FormItem className="text-start">
                                <FormLabel className="text-start">Symbol number</FormLabel>
                                <FormControl>
                                    <Input placeholder="Enter your symbol number" {...field} />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                    <FormField
                        control={form.control}
                        name="registration_number"
                        render={({ field }) => (
                            <FormItem className="text-start">
                                <FormLabel className="text-start">Registration number</FormLabel>
                                <FormControl>
                                    <Input placeholder="Enter your registration number" {...field} />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                </div>
                <div className="flex flex-row items-center gap-4">
                    <Select
                        onValueChange={(value) => setValue('batchIDStr', value)} // Update form value when batch changes
                    >
                        <FormLabel className="text-start">Select your batch: </FormLabel>
                        <SelectTrigger className="w-[180px]">
                            <SelectValue placeholder="Select your batch" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectGroup>
                                <SelectLabel>Batch</SelectLabel>
                                <SelectItem value="2024">2024</SelectItem>
                                <SelectItem value="2023">2023</SelectItem>
                                <SelectItem value="2022">2022</SelectItem>
                                <SelectItem value="2021">2021</SelectItem>
                                <SelectItem value="2020">2020</SelectItem>
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
                                <SelectItem value="CSIT">CSIT</SelectItem>
                                <SelectItem value="BBA">BBA</SelectItem>
                                <SelectItem value="BIM">BIM</SelectItem>
                                <SelectItem value="BBS">BBS</SelectItem>
                            </SelectGroup>
                        </SelectContent>
                    </Select>
                </div>

                <Button type="submit" size="lg" className="w-full">Register</Button>
            </form>
        </Form>
    )
}