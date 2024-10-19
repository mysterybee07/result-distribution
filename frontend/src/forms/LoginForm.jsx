import { useMutation } from '@tanstack/react-query';
import api from '../api';
import { useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";
import { Input } from "../components/ui/input";
import { FaRegEye, FaRegEyeSlash } from "react-icons/fa";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

// Define your form schema
const formSchema = z.object({
    identifier: z.string()
        .min(5, { message: "Email must be at least 5 characters." }),
    password: z.string().min(6, {
        message: "Password must be at least 6 characters.",
    }),
});

export function LoginForm() {
    const { login, userRole } = useAuth();
    const navigate = useNavigate();

    const form = useForm({
        resolver: zodResolver(formSchema),
        defaultValues: {
            identifier: "",
            password: "",
        },
    });

    // Mutation for logging in
    const { mutate: loginMutation } = useMutation({
        mutationFn: async (values) => {
            console.log("Submitting login values:", values); // Debugging
            const { identifier, password } = values;
            const response = await api.post("/user/login", { identifier, password });
            return response.data; // Ensure this returns the correct data structure
        },
        onSuccess: (data) => {
            console.log("Login successful:", data); // Debugging

            const { user, token } = data; // Destructure user and token from response
            login(user); // Call the login function from context

            // Navigate based on user role
            if (user.role === "admin") {
                navigate("/admin");
            } else {
                navigate("/"); // Navigate to home or another route for non-admins
            }
        },
        onError: (error) => {
            console.error("Login failed:", error); // Log error
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
                console.error("Unexpected error:", error.message);
            }
        },
    });

    const onSubmit = (values) => {
        console.log("Form submitted:", values); // Debugging
        loginMutation(values);
    };

    const [showPassword, setShowPassword] = useState(false);
    const togglePasswordVisibility = () => setShowPassword(!showPassword);

    return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
                <FormField
                    control={form.control}
                    name="identifier"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Identifier</FormLabel>
                            <FormControl>
                                <Input placeholder="Email or Symbol number" {...field} />
                            </FormControl>
                            <FormMessage error={form.formState.errors.identifier} />
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="password"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Password</FormLabel>
                            <FormControl>
                                <Input
                                    placeholder="Password"
                                    {...field}
                                    type={showPassword ? 'text' : 'password'}
                                    endAdornment={
                                        <span onClick={togglePasswordVisibility} className="cursor-pointer">
                                            {showPassword ? <FaRegEye /> : <FaRegEyeSlash />}
                                        </span>
                                    }
                                />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <div className="flex justify-between items-center">
                    <label
                        htmlFor="remember_me"
                        className="text-sm font-medium leading-none"
                    >
                        <Checkbox id="remember_me" className="mr-2" />
                        Remember me
                    </label>
                    <Button type="submit" variant="link">Forgot password?</Button>
                </div>
                <Button type="submit" size="lg" className="w-full">Login</Button>
            </form>
        </Form>
    );
}
