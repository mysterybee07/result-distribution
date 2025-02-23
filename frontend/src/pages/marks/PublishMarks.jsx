import { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "../../api";
import { useData } from "../../context/DataContext";

export default function PublishMarks() {
    const { batches, programs } = useData(); // Get data from context
    const [formData, setFormData] = useState({
        batch_id: "",
        program_id: "",
        semester_id: "",
    });
    // UseQuery to fetch semesters based on selected program_id
    const { data: semesters, isLoading: loadingSemesters, error: errorSemesters } = useQuery({
        queryKey: ['semesters', formData.program_id],
        queryFn: async () => {
            const response = await api.get(`/semester/by-program/${formData.program_id}`);
            return response.data.semesters;
        },
        enabled: Boolean(formData.program_id),
        refetchOnWindowFocus: false,
        refetchOnReconnect: false,
    });
    console.log("🚀 ~ PublishMarks ~ semesters:", semesters)
    // Handle form submission
    const mutation = useMutation({
        mutationFn: async (data) => {
            const response = await api.post("/result/publish", data);
            return response.data;
        },
        onSuccess: (data) => {
            alert("✅ Result published successfully!");
            console.log("🚀 ~ Success:", data);
        },
        onError: (error) => {
            alert("❌ Error publishing result.", errorMessage);
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
                <h2 className="text-2xl font-bold text-center mb-6">Publish Results</h2>

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

                    {/* Semester Selection */}
                    <div className="mb-4">
                        <label className="block text-gray-700 font-medium mb-2">Select Semester</label>
                        <select
                            name="semester_id"
                            value={formData.semester_id}
                            onChange={handleChange}
                            className="w-full p-3 border rounded-lg"
                            required
                        >
                            <option value="">Choose a semester</option>
                            {semesters?.map((semester) => (
                                <option key={semester.ID} value={semester.ID}>
                                    {semester.semester_name}
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
                        {mutation.isLoading ? "Publishing..." : "Publish Result"}
                    </button>
                </form>
            </div>
        </div>
    );
}
