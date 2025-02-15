import React, { useState } from 'react';
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "../../api";
import { Select, SelectLabel, SelectGroup, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useData } from '../../context/DataContext';


const createMarks = async () => {
    const { data } = await api.get("/marks/create");
    return data;
};

const fetchCourse = async ({ queryKey }) => {
    const [, program_id, semester_id] = queryKey; // Destructure from queryKey
    if (!program_id || !semester_id) {
        return []; // Return an empty array if required params are missing
    }
    const response = await api.get(`/courses/filter?program_id=${program_id}&semester_id=${semester_id}`);
    console.log("🚀 ~ fetchCourse ~ response:", response.data.courses);
    return response.data.courses;
};

const fetchFilteredStudents = async ({ queryKey }) => {
    const [, batch_id, program_id, semester_id] = queryKey;
    const response = await api.get(`/students/filter?batch_id=${batch_id}&program_id=${program_id}`);
    console.log("🚀 ~ fetchFilteredStudents ~ response:", response)
    return response.data.students;
};

const AddMarks = () => {
    const { batches, programs } = useData();

    // console.log("🚀 ~ AddMarks ~ batches:", batches)
    // console.log("🚀 ~ AddMarks ~ programs:", programs)
    const [formData, setFormData] = useState({
        batchID: '',
        programID: '',
        semesterID: '',
        courseID: '',
        studentMarks: [
            {
                studentID: '',
                semesterMarks: '',
                assistantMarks: '',
                practicalMarks: ''
            }
        ]
    });

    // query to fetch the semester
    const {
        data: semesters,
    } = useQuery({
        queryKey: ["semesters", formData.programID], // Add `selectedProgram` to the query key
        queryFn: async () => {
            const response = await api.get(`/semester/by-program/${formData.programID}`);
            return response.data.semesters;
        },
        enabled: !!formData.programID, // Run the query only if `selectedProgram` is truthy
    });

    // query to fetch releated course
    const { data: courses = [] } = useQuery({
        queryKey: ["courses", formData.programID, formData.semesterID], // Include program and semester in the queryKey
        queryFn: fetchCourse,
        // enabled: search // Run only when both values are truthy
    });
    // console.log("🚀 ~ AddMarks ~ courses:", courses)

    // query to fetch students
    const { data: students = [] } = useQuery({
        queryKey: ["students", formData.batchID, formData.programID, formData.semesterID],
        queryFn: fetchFilteredStudents,
        enabled: !!formData.batchID && !!formData.programID && !!formData.semesterID
    });
    console.log("🚀 ~ AddMarks ~ students:", students)

    const mutation = useMutation({
        mutationFn: createMarks,
        onSuccess: (data) => {
            console.log("Success:", data);
        },
        onError: (error) => {
            console.error("Error:", error.response?.data || error.message);
        },
    });

    const handleInputChange = (e, index = null) => {
        const { name, value } = e.target;

        if (index !== null) {
            // Handle student marks array
            setFormData(prev => ({
                ...prev,
                studentMarks: prev.studentMarks.map((mark, i) =>
                    i === index ? { ...mark, [name]: value } : mark
                )
            }));
        } else {
            // Handle other fields
            setFormData(prev => ({
                ...prev,
                [name]: value
            }));
        }
    };

    const addStudent = () => {
        setFormData(prev => ({
            ...prev,
            studentMarks: [
                ...prev.studentMarks,
                {
                    studentID: '',
                    semesterMarks: '',
                    assistantMarks: '',
                    practicalMarks: ''
                }
            ]
        }));
    };

    const removeStudent = (index) => {
        setFormData(prev => ({
            ...prev,
            studentMarks: prev.studentMarks.filter((_, i) => i !== index)
        }));
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        const marksData = {
            batchID: parseInt(formData.batchID),
            programID: parseInt(formData.programID),
            semesterID: parseInt(formData.semesterID),
            courseID: parseInt(formData.courseID),
            marks: formData.studentMarks.map(mark => ({
                studentID: parseInt(mark.studentID),
                semesterMarks: parseInt(mark.semesterMarks),
                assistantMarks: parseInt(mark.assistantMarks),
                practicalMarks: parseInt(mark.practicalMarks)
            }))
        };
        mutation.mutate(marksData);
    };

    const batchOptions = ["2023", "2024", "2025"];
    const programOptions = ["Computer Science", "Business", "Engineering"];
    const semesterOptions = ["Semester 1", "Semester 2", "Semester 3"];
    const courseOptions = ["Mathematics", "Physics", "Chemistry"];

    return (
        <div className="min-h-screen bg-gray-50 py-8 px-4 sm:px-6 lg:px-8">
            <div className="w-full mx-auto">
                <div className="">
                    <div className="px-4 py-5 sm:p-6">
                        <h2 className="text-lg font-medium text-gray-900 mb-6">Add Student Marks</h2>

                        <form onSubmit={handleSubmit} className="space-y-6">
                            {/* Course Information */}
                            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
                                <div>
                                    <label className="text-left block text-sm font-semibold text-gray-700">
                                        Batch
                                    </label>
                                    <Select value={formData.batchID} onValueChange={(value) => setFormData({ ...formData, batchID: Number(value) })}>
                                        <SelectTrigger>
                                            <SelectValue placeholder="Select Batch" />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectGroup>
                                                <SelectLabel>Programs</SelectLabel>
                                                {Array.isArray(batches) && batches.map((batch, index) => (
                                                    <SelectItem key={index} value={batch.ID}>{batch.batch}</SelectItem>
                                                ))}
                                            </SelectGroup>
                                        </SelectContent>
                                    </Select>
                                </div>

                                <div>
                                    <label className="text-left block text-sm font-semibold text-gray-700">
                                        Program
                                    </label>
                                    <Select value={formData.programID} onValueChange={(value) => setFormData({ ...formData, programID: Number(value) })}>
                                        <SelectTrigger>
                                            <SelectValue placeholder="Select Program" />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectGroup>
                                                <SelectLabel>Programs</SelectLabel>
                                                {Array.isArray(programs) && programs.map((program, index) => (
                                                    <SelectItem key={index} value={program.ID}>{program.program_name}</SelectItem>
                                                ))}
                                            </SelectGroup>
                                        </SelectContent>
                                    </Select>
                                </div>

                                <div>
                                    <label className="text-left block text-sm font-semibold text-gray-700">
                                        Semester
                                    </label>
                                    <Select value={formData.semesterID} onValueChange={(value) => setFormData({ ...formData, semesterID: Number(value) })}>
                                        <SelectTrigger>
                                            <SelectValue placeholder="Select Semester" />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectGroup>
                                                <SelectLabel>Semesters</SelectLabel>
                                                {Array.isArray(semesters) && semesters.map((semester, index) => (
                                                    <SelectItem key={index} value={semester.ID}>Semester {semester.semester_name}</SelectItem>
                                                ))}
                                            </SelectGroup>
                                        </SelectContent>
                                    </Select>
                                </div>

                                <div>
                                    <label className="text-left block text-sm font-semibold text-gray-700">
                                        Course
                                    </label>
                                    <Select value={formData.courseID} onValueChange={(value) => setFormData({ ...formData, courseID: Number(value) })}>
                                        <SelectTrigger>
                                            <SelectValue placeholder="Select Course" />
                                        </SelectTrigger>
                                        <SelectContent>
                                            {Array.isArray(courses) && courses.map((course, index) => (
                                                <SelectItem key={index} value={course.ID}>{course.name}</SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                </div>
                            </div>

                            {/* Student Marks */}
                            <div className="space-y-4">
                                <div className="flex justify-between items-center">
                                    <h3 className="text-md font-medium text-gray-900">Student Marks</h3>
                                    <button
                                        type="button"
                                        onClick={addStudent}
                                        className="inline-flex items-center px-3 py-1.5 border border-transparent text-xs font-medium rounded-md text-indigo-700 bg-indigo-100 hover:bg-indigo-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                                    >
                                        Add Student
                                    </button>
                                </div>

                                {formData.studentMarks.map((student, index) => (
                                    <div key={index} className="border rounded-lg p-4 bg-gray-50">
                                        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-5">
                                            <div className="lg:col-span-2">
                                                <label className="text-left block text-sm font-semibold text-gray-700">
                                                    Student ID
                                                </label>
                                                <input
                                                    type="text"
                                                    name="studentID"
                                                    value={student.studentID}
                                                    onChange={(e) => handleInputChange(e, index)}
                                                    className="p-2 mt-1 block w-full rounded-md border border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                                                    required
                                                />
                                            </div>

                                            <div>
                                                <label className="text-left block text-sm font-semibold text-gray-700">
                                                    Semester Marks
                                                </label>
                                                <input
                                                    type="number"
                                                    name="semesterMarks"
                                                    value={student.semesterMarks}
                                                    onChange={(e) => handleInputChange(e, index)}
                                                    className="p-2 mt-1 block w-full rounded-md border border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                                                    required
                                                />
                                            </div>

                                            <div>
                                                <label className="text-left block text-sm font-semibold text-gray-700">
                                                    Assistant Marks
                                                </label>
                                                <input
                                                    type="number"
                                                    name="assistantMarks"
                                                    value={student.assistantMarks}
                                                    onChange={(e) => handleInputChange(e, index)}
                                                    className="p-2 mt-1 block w-full rounded-md border border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                                                    required
                                                />
                                            </div>

                                            <div>
                                                <label className="text-left block text-sm font-semibold text-gray-700">
                                                    Practical Marks
                                                </label>
                                                <div className="flex gap-2">
                                                    <input
                                                        type="number"
                                                        name="practicalMarks"
                                                        value={student.practicalMarks}
                                                        onChange={(e) => handleInputChange(e, index)}
                                                        className="p-2 mt-1 block w-full rounded-md border border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                                                        required
                                                    />
                                                    {formData.studentMarks.length > 1 && (
                                                        <button
                                                            type="button"
                                                            onClick={() => removeStudent(index)}
                                                            className="mt-1 inline-flex items-center p-1.5 border border-transparent rounded-md text-red-700 bg-red-100 hover:bg-red-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                                                        >
                                                            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                                                <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
                                                            </svg>
                                                        </button>
                                                    )}
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </div>

                            {/* Submit Button */}
                            <div className="flex justify-end">
                                <button
                                    type="submit"
                                    disabled={mutation.isLoading}
                                    className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
                                >
                                    {mutation.isLoading ? "Submitting..." : "Submit Marks"}
                                </button>
                            </div>

                            {/* Status Messages */}
                            {mutation.isError && (
                                <div className="rounded-md bg-red-50 p-4">
                                    <div className="flex">
                                        <div className="text-sm text-red-700">
                                            Error: {mutation.error.message}
                                        </div>
                                    </div>
                                </div>
                            )}

                            {mutation.isSuccess && (
                                <div className="rounded-md bg-green-50 p-4">
                                    <div className="flex">
                                        <div className="text-sm text-green-700">
                                            Marks Created Successfully!
                                        </div>
                                    </div>
                                </div>
                            )}
                        </form>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default AddMarks;