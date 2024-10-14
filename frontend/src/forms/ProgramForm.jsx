import { Button } from "@/components/ui/button"
import {
    Dialog,
    DialogContent,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useState } from "react";
import api from "../api";
import { useMutation } from "@tanstack/react-query";
import { useToast } from "@/hooks/use-toast";

export default function ProgramForm() {
    const { toast } = useToast();
    const [program_name, setProgram_name] = useState("");
    const [isDialogOpen, setIsDialogOpen] = useState(false); 

    // Mutation to handle form submission
    const { mutate, isLoading, isError, isSuccess, error } = useMutation({
        mutationFn: async (newProgram) => {
            // Replace with your API endpoint
            const response = await api.post("/program/create", newProgram);
            return response.data;
        },
        onSuccess: (data) => {
            console.log("Program added successfully:", data);
            setIsDialogOpen(false); 
            toast({
                title: "Program Added",
                description: JSON.stringify(data.message),
                variant: "success",
              })
              setProgram_name("");
        },
        onError: (error) => {
            console.error("Error adding program:", error);
        },
    });

    const handleSubmit = (e) => {
        e.preventDefault();
        // Call the mutation function with the form data
        mutate({ program_name });
    };

    return (
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
            <DialogTrigger asChild>
                <Button variant="" size="lg" className="w-full" onClick={() => setIsDialogOpen(true)}>
                    Add New
                </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>Add new program</DialogTitle>
                </DialogHeader>
                <form onSubmit={handleSubmit}>
                    <div className="grid gap-4 py-4">
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="program_name" className="text-start">
                                Program Name:
                            </Label>
                            <Input
                                id="program_name"
                                value={program_name}
                                onChange={(e) => setProgram_name(e.target.value)}
                                className="col-span-3"
                            />
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
    );
}