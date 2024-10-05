import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { Button } from "@/components/ui/button"
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

const formSchema = z.object({
    username: z.string()
        .email({ message: "Please enter a valid email address." })
        .min(5, { message: "Email must be at least 5 characters." }),
    password: z.string().min(6, {
        message: "Password must be at least 8 characters.",
    }),
})

export function ProfileForm() {
    // 1. Define your form.
    const form = useForm({
        resolver: zodResolver(formSchema),
        defaultValues: {
            username: "",
            password: "",
        },
    })

    // 2. Define a submit handler.
    function onSubmit(values) {
        // Do something with the form values.
        // ✅ This will be type-safe and validated.
        console.log(values)
    }

    return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
                <FormField
                    control={form.control}
                    name="username"
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
                                <Input placeholder="Password" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <div className="flex justify-between">
                    <label className="flex items-center">
                        <input
                            type="checkbox"
                            name="remember"
                            id="remember"
                            className="form-checkbox"
                        />
                        <span className="ml-2 text-sm">Remember me</span>
                    </label>
                    <Button type="submit" variant="link">Forgot password?</Button>
                </div>
                <Button type="submit" size="lg" className="w-full">Login</Button>
            </form>
        </Form>
    )
}