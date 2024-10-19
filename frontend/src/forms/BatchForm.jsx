import { Button } from "@/components/ui/button"
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog"
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectLabel,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useMutation } from "@tanstack/react-query";
import { useState } from "react";
import api from "../api";
import { useToast } from "@/hooks/use-toast";
import { FaEdit } from "react-icons/fa";

export default function BatchForm({ batch }) {
    const { toast } = useToast();
    const [batch_name, setBatch_name] = useState(batch?.batch || "");
    console.log("🚀 ~ BatchForm ~ batch_name:", batch_name)
    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const currentYear = new Date().getFullYear();
    const years = Array.from({ length: 15 }, (_, i) => currentYear - i);

    const editbutton = () => {
        return (
            <span onClick={() => setIsDialogOpen(true)} className="cursor-pointer flex justify-center items-center">
                <FaEdit /> Edit
            </span>
        )
    }
    const addbutton = () => {
        return (
            <Button size="lg" className="w-full" onClick={() => setIsDialogOpen(true)}>
                Add new batch
            </Button>
        )
    }
    const button = batch ? editbutton : addbutton;

    // Mutation to handle form submission
    const { mutate, isLoading, isError, isSuccess, error } = useMutation({
        mutationFn: async (newBatch) => {
            // Replace with your API endpoint
            const response = await api.post("/batch/create", newBatch);
            return response.data;
        },
        onSuccess: (data) => {
            console.log("Program added successfully:", data);
            setIsDialogOpen(false);
            toast({
                title: "Batch Added",
                description: JSON.stringify(data.message),
                variant: "success",
            })
        },
        onError: (error) => {
            console.error("Error adding program:", error);
        },
    });

    const { mutate: updateBatch, isError: updateIsError, error: updateError } = useMutation({
        mutationFn: async (updatedBatch) => {
            // Replace with your API endpoint
            const response = await api.put(`/batch/update/${batch.ID}`, updatedBatch);
            return response.data;
        },
        onSuccess: (data) => {
            console.log("Program updated successfully:", data);
            setIsDialogOpen(false);
            toast({
                title: "Batch Updated",
                description: JSON.stringify(data.message),
                variant: "success",
            })
        },
        onError: (updateError) => {
            // Check if error response exists and display the message from the server
            if (updateError.response && updateError.response.data && updateError.response.data.error) {
                // Display the specific error message from the server
                const serverErrorMessage = updateError.response.data.error;
                toast({
                    title: "Error",
                    description: serverErrorMessage,
                    variant: "destructive",
                });
            } else {
                // Fallback for a generic error message
                console.error("Error updating program:", updateError);
                toast({
                    title: "Error",
                    description: "An unexpected error occurred.",
                    variant: "destructive",
                });
            }
        },
    });

    const handleSubmit = (e) => {
        e.preventDefault();
        // Call the mutation function with the form data
        const intBatch = parseInt(batch_name, 10); // '10' ensures decimal parsing
        console.log("🚀 ~ handleSubmit ~ intBatch:", intBatch);
        batch ? updateBatch({ batch: intBatch }) : mutate({ batch: intBatch });
    };

    return (
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
            <DialogTrigger asChild>
                {button()}
            </DialogTrigger>
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>Add batch</DialogTitle>
                </DialogHeader>
                <form onSubmit={handleSubmit}>
                    <div className="grid gap-4 py-4">
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="batch" className="text-right">
                                Batch Year
                            </Label>
                            {/* <Input
                                id="batch"
                                defaultValue="Pedro Duarte"
                                className="col-span-3"
                                onClick={(e) => setBatch(e.target.value)}
                            /> */}
                            <Select onValueChange={(value) => setBatch_name(value)}>
                                <SelectTrigger className="w-[180px]">
                                    <SelectValue placeholder={batch?.batch || "Select a year"} />
                                </SelectTrigger>
                                <SelectContent>
                                    {years.map((year) => (
                                        <SelectItem key={year} value={year.toString()}>
                                            {year}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>

                        </div>
                    </div>
                    <DialogFooter>
                        <Button type="submit" disabled={isLoading}>
                            {isLoading ? "Saving..." : "Save changes"}
                        </Button>
                    </DialogFooter>
                </form>
                {isError && <p className="text-red-500 mt-2">Error: {error.message}</p>}
                {/* {isSuccess && <p className="text-green-500 mt-2">Program added successfully!</p>} */}
            </DialogContent>
        </Dialog>
    )
}
