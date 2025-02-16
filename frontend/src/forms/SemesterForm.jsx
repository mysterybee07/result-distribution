import { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
    DialogFooter,
} from "@/components/ui/dialog";
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import api from "../api";
import { useToast } from "@/hooks/use-toast";
import { useData } from "../context/DataContext";

export default function SemesterForm() {
    const { toast } = useToast();
    const [semesterName, setSemesterName] = useState("");
    const [programId, setProgramId] = useState("");
    const [isDialogOpen, setIsDialogOpen] = useState(false);

    // Fetch programs for selection
    const { programs, loadingPrograms } = useData();

    // Mutation for adding semester
    const { mutate, isLoading } = useMutation({
        mutationFn: async (newSemester) => {
            const response = await api.post("/semester/create", newSemester);
            return response.data;
        },
        onSuccess: (data) => {
            setIsDialogOpen(false);
            toast({
                title: "Semester Added",
                description: JSON.stringify(data.message),
                variant: "success",
            });
        },
        onError: (error) => {
            toast({
                title: "Error",
                description: error.response?.data?.error || "An unexpected error occurred.",
                variant: "destructive",
            });
        },
    });

    const handleSubmit = (e) => {
        e.preventDefault();
        mutate({ semester_name: Number(semesterName), program_id: Number(programId) });
    };

    return (
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
            <DialogTrigger asChild>
                <Button size="lg" className="w-full">Add new semester</Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>Add Semester</DialogTitle>
                </DialogHeader>
                <form onSubmit={handleSubmit}>
                    <div className="grid gap-4 py-4">
                        {/* Program Selection */}
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label className="text-right">Program</Label>
                            <Select onValueChange={setProgramId} disabled={loadingPrograms}>
                                <SelectTrigger className="col-span-3">
                                    <SelectValue placeholder="Select a program" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectGroup>
                                        {programs?.map((program) => (
                                            <SelectItem key={program.ID} value={program.ID.toString()}>
                                                {program.program_name}
                                            </SelectItem>
                                        ))}
                                    </SelectGroup>
                                </SelectContent>
                            </Select>
                        </div>
                        {/* Semester Name Input */}
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="semester" className="text-right">
                                Semester Name
                            </Label>
                            <Input
                                id="semester"
                                value={semesterName}
                                onChange={(e) => setSemesterName(e.target.value)}
                                className="col-span-3"
                                placeholder="Enter semester name"
                            />
                        </div>


                    </div>
                    <DialogFooter>
                        <Button type="submit" disabled={isLoading}>
                            {isLoading ? "Saving..." : "Save changes"}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    );
}
