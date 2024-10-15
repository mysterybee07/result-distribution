import { useMutation } from '@tanstack/react-query';
import api from '../api';
import { useState } from "react"
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
import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
} from "@/components/ui/tooltip"
import { Input } from "../components/ui/input"
import { FaRegEye, FaRegEyeSlash } from "react-icons/fa";
import { useNavigate } from "react-router-dom"
import { useAuth } from "../context/AuthContext"

const formSchema = z.object({
    identifier: z.string()
        .email({ message: "Please enter a valid email address." })
        .min(5, { message: "Email must be at least 5 characters." }),
    password: z.string().min(6, {
        message: "Password must be at least 8 characters.",
    }),
})

export function LoginForm() {
    const { login, userRole } = useAuth();
    const navigate = useNavigate();
    // 1. Define your form.
    const form = useForm({
        resolver: zodResolver(formSchema),
        defaultValues: {
            identifier: "",
            password: "",
        },
    })

    // 2. Define the mutation to handle login
    const { mutate: loginMutation } = useMutation({
        mutationFn: async (values) => {
            const { identifier, password } = values;
            const { data } = await api.post("/user/login", { identifier, password });
            return data; // The result of the login request
        },
        onSuccess: (data) => {
            const { user } = data;
            console.log("🚀 ~ Login Successful ~ data:", data);

            if (user && user.role === "admin") {
                navigate("/admin");
            } else {
                navigate("/");
            }

            login();
            userRole(user.role);
        },
        onError: (error) => {
            const errorData = error.response?.data;

            if (errorData?.errors) {
                // Set field errors in the form
                for (const [key, message] of Object.entries(errorData.errors)) {
                    form.setError(key, {
                        type: "manual",
                        message,
                    });
                }
            } else {
                console.error("Login failed:", error.message);
            }
        },
    });

    // 3. Define a submit handler.
    const onSubmit = async (values) => {
        loginMutation(values);
    }
    const [showPassword, setShowPassword] = useState(false);

    const togglePasswordVisibility = () => {
        setShowPassword(!showPassword);
    };

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
                            {/* <FormMessage /> */}
                            <FormMessage error={form.formState.errors.identifier} /> {/* Display error */}
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="password"
                    type="password"
                    render={({ field }) => (
                        <FormItem className="text-start">
                            <FormLabel>Password</FormLabel>
                            <FormControl>
                                <Input
                                    placeholder="Password" {...field}
                                    type={showPassword ? 'text' : 'password'}
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
                <div className="flex justify-between items-center">
                    {/* <label className="flex items-center">
                        <input
                            type="checkbox"
                            name="remember"
                            id="remember"
                            className="form-checkbox"
                        />
                        <span className="ml-2 text-sm">Remember me</span>
                    </label> */}
                    <label
                        htmlFor="remember_me"
                        className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                    >
                        <Checkbox id="remember_me" className="mr-2" />
                        Remember me
                    </label>
                    <Button type="submit" variant="link">Forgot password?</Button>
                </div>
                <Button type="submit" size="lg" className="w-full">Login</Button>
            </form>
        </Form>
    )
}