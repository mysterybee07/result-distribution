import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import api from "../../api";
import { useData } from "../../context/DataContext";
import { useToast } from "@/hooks/use-toast";

export default function AssignCenters() {
    
    const { toast } = useToast();
    const { batches, programs } = useData(); // Get data from context
    const [formData, setFormData] = useState({
        batch_id: "",
        program_id: "",
    });

    // Handle form submission
    const mutation = useMutation({
        mutationFn: async (data) => {
            const response = await api.get(`/exam/assign-centers?batch_id=${data.batch_id}&program_id=${data.program_id}`);
            return response.data;
        },
        onSuccess: (data) => {
            alert("✅ Centers assigned successfully!");
            // toast.success("Centers assigned successfully!", {
            //     position: "top-right",
            //     autoClose: 3000,
            //     hideProgressBar: false,
            //     closeOnClick: true,
            //     pauseOnHover: true,
            //     draggable: true,
            // });
            console.log("🚀 ~ Success:", data);
        },
        onError: (error) => {
            alert("❌ Error assigning centers");
            // toast.error("Error assigning centers", {
            //     position: "top-right",
            //     autoClose: 3000,
            //     hideProgressBar: false,
            //     closeOnClick: true,
            //     pauseOnHover: true,
            //     draggable: true,
            // });
            console.error("❌ Error:", error);
        },
    });

    // Handle input change
    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: Number(e.target.value), // Ensure numeric values
        });
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-100">
            <div className="bg-white p-6 rounded-lg shadow-lg w-full max-w-md">
                <h2 className="text-2xl font-bold text-center mb-6">Assign Exam Centers</h2>

                <form
                    onSubmit={(e) => {
                        e.preventDefault();
                        mutation.mutate(formData);
                    }}
                >
                    {/* Batch Selection */}
                    <div className="mb-4">
                        <label className="block text-gray-700 font-medium mb-2">Select Batch</label>
                        <select
                            name="batch_id"
                            value={formData.batch_id}
                            onChange={handleChange}
                            className="w-full p-3 border rounded-lg"
                            required
                        >
                            <option value="">Choose a batch</option>
                            {batches?.map((batch) => (
                                <option key={batch.ID} value={batch.ID}>
                                    {batch.batch}
                                </option>
                            ))}
                        </select>
                    </div>

                    {/* Program Selection */}
                    <div className="mb-4">
                        <label className="block text-gray-700 font-medium mb-2">Select Program</label>
                        <select
                            name="program_id"
                            value={formData.program_id}
                            onChange={handleChange}
                            className="w-full p-3 border rounded-lg"
                            required
                        >
                            <option value="">Choose a program</option>
                            {programs?.map((program) => (
                                <option key={program.ID} value={program.ID}>
                                    {program.program_name}
                                </option>
                            ))}
                        </select>
                    </div>

                    {/* Submit Button */}
                    <button
                        type="submit"
                        disabled={mutation.isLoading}
                        className="w-full bg-blue-600 text-white p-3 rounded-lg font-semibold hover:bg-blue-700 transition"
                    >
                        {mutation.isLoading ? "Assigning Centers..." : "Assign Centers"}
                    </button>
                </form>
            </div>
        </div>
    );
}