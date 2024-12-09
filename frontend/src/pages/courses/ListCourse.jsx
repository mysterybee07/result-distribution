import { useQuery } from '@tanstack/react-query';
import React, { useState } from 'react'
import api from '../../api';
import { useData } from '../../context/DataContext';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Button } from '../../components/ui/button';
const ListCourse = () => {
  const { programs } = useData();
  // console.log("🚀 ~ ListCourse ~ semester:", semester)
  const [selectedProgram, setSelectedProgram] = useState("");
  const [selectedSemester, setSelectedSemester] = useState("");
  console.log("🚀 ~ ListCourse ~ selectedProgram:", selectedProgram)
  const {
    data: semesters,
    isLoading: loadingSemesters,
    error: errorSemesters,
  } = useQuery({
    queryKey: ["semesters", selectedProgram], // Add `selectedProgram` to the query key
    queryFn: async () => {
      const response = await api.get(`/semester/by-program/${selectedProgram}`);
      return response.data.semesters;
    },
    enabled: !!selectedProgram, // Run the query only if `selectedProgram` is truthy
  });
  console.log("🚀 ~ ListCourse ~ semesters:", semesters)
  const fetchCourse = async ({ queryKey }) => {
    const [, program_id, semester_id] = queryKey; // Destructure from queryKey
    if (!program_id || !semester_id) {
      return []; // Return an empty array if required params are missing
    }
    const response = await api.get(`/courses/filter?program_id=${program_id}&semester_id=${semester_id}`);
    console.log("🚀 ~ fetchCourse ~ response:", response.data.courses);
    return response.data.courses;
  };

  // Using useQuery with proper dependencies
  const { data: courses = [], isLoading, error } = useQuery({
    queryKey: ["courses", selectedProgram, selectedSemester], // Include program and semester in the queryKey
    queryFn: fetchCourse,
    enabled: !!selectedProgram && !!selectedSemester, // Run only when both values are truthy
  });

  console.log(courses);

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;
  return (
    <>
      <div className='flex justify-start space-x-4'>
        <Select
          value={selectedProgram}
          onValueChange={(value) => setSelectedProgram(value)}
        >
          <SelectTrigger>
            <SelectValue placeholder="Select program" />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              <SelectLabel>Program</SelectLabel>
              {Array.isArray(programs) && programs.map((program, index) => (
                <SelectItem key={index} value={program.ID}>{program.program_name}</SelectItem>
              ))}
            </SelectGroup>
          </SelectContent>
        </Select>
        <Select
          value={selectedSemester}
          onValueChange={(value) => setSelectedSemester(value)}
        >
          <SelectTrigger>
            <SelectValue placeholder="Select Semester" />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              <SelectLabel>Semester</SelectLabel>
              {Array.isArray(semesters) && semesters.map((semester, index) => (
                <SelectItem key={index} value={semester.ID}>{semester.semester_name}</SelectItem>
              ))}
            </SelectGroup>
          </SelectContent>
        </Select>
        <Button disabled={!selectedProgram || !selectedSemester}
          onClick={() => {
            // Ensure selectedProgram and selectedSemester are updated first
            if (!selectedProgram || !selectedSemester) {
              alert("Please select both a program and a semester.");
              return;
            }
            // No need to call fetchCourse directly; react-query handles it
            // The `enabled` option in useQuery ensures the query runs when values are valid
          }}
        >
          Search
        </Button>

      </div>
      <div>
        {Array.isArray(courses) && courses.map((course, index) => (
          <div key={index}>
            <h1>{course.name}</h1>
            <p>{course.course_code}</p>
            <p>{course.is_compulsory ? 'true' : 'false'}</p>
            <p>{course.semester_total_marks}</p>
            <p>{course.practical_total_marks}</p>
            <p>{course.assistant_total_marks}</p>
          </div>
        ))}
      </div>
    </>
  )
}

export default ListCourse