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
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableFooter,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Button } from '../../components/ui/button';
import { FaEdit } from 'react-icons/fa';
import { useNavigate } from 'react-router-dom';

// Fetch courses function (moved outside for better reusability)
const fetchCourse = async ({ queryKey }) => {
  const [, program_id, semester_id] = queryKey;
  if (!program_id || !semester_id) return [];

  const response = await api.get(`/courses/filter?program_id=${program_id}&semester_id=${semester_id}`);
  return response.data.courses;
};

const ListCourse = () => {
  const navigate = useNavigate();
  const { programs } = useData();

  // State variables
  const [selectedProgram, setSelectedProgram] = useState("");
  const [selectedSemester, setSelectedSemester] = useState("");
  const [search, setSearch] = useState(false);

  // Fetch semesters based on selected program
  const { data: semesters = [], isLoading: loadingSemesters, error: errorSemesters } = useQuery({
    queryKey: ["semesters", selectedProgram],
    queryFn: async () => {
      const response = await api.get(`/semester/by-program/${selectedProgram}`);
      return response.data.semesters;
    },
    enabled: !!selectedProgram,
  });

  // Fetch courses based on selected program & semester
  const { data: courses = [], isLoading, error } = useQuery({
    queryKey: ["courses", selectedProgram, selectedSemester],
    queryFn: fetchCourse,
    enabled: search,
  });

  // Show loading or error messages
  if (isLoading) return <div>Loading courses...</div>;
  if (error) return <div>Error: {error?.message || "Something went wrong"}</div>;
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
            setSearch(true);
            // No need to call fetchCourse directly; react-query handles it
            // The `enabled` option in useQuery ensures the query runs when values are valid
          }}
        >
          Search
        </Button>

      </div>
      <div className='mt-4'>
        {courses.length === 0 ?
          <div>
            No data found.
            <p>Try selecting a program and semester.</p>
          </div> : (
            <Table>
              <TableHeader>
                {/* First Row for Main Headers */}
                <TableRow>
                  <TableHead rowSpan={2}>Course Code</TableHead>
                  <TableHead rowSpan={2}>Course Name</TableHead>
                  <TableHead rowSpan={2}>Is Compulsory</TableHead>
                  <TableHead colSpan={2} className="text-center">Semester</TableHead>
                  <TableHead colSpan={2} className="text-center">Practical</TableHead>
                  <TableHead colSpan={2} className="text-center">Assistant</TableHead>
                  <TableHead rowSpan={2}>Action</TableHead>
                </TableRow>
                {/* Second Row for Sub-Headers */}
                <TableRow className="text-center">
                  <TableHead>Total Marks</TableHead>
                  <TableHead>Pass Marks</TableHead>
                  <TableHead>Total Marks</TableHead>
                  <TableHead>Pass Marks</TableHead>
                  <TableHead>Total Marks</TableHead>
                  <TableHead>Pass Marks</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {/* Map through courses */}
                {Array.isArray(courses) &&
                  courses.map((course, index) => (
                    <TableRow key={index}>
                      <TableCell>{course.course_code}</TableCell>
                      <TableCell>{course.name}</TableCell>
                      <TableCell>{course.is_compulsory ? 'True' : 'False'}</TableCell>
                      <TableCell>{course.semester_total_marks}</TableCell>
                      <TableCell>{course.semester_pass_marks}</TableCell>
                      <TableCell>{course.practical_total_marks}</TableCell>
                      <TableCell>{course.practical_pass_marks}</TableCell>
                      <TableCell>{course.assistant_total_marks}</TableCell>
                      <TableCell>{course.assistant_pass_marks}</TableCell>
                      <TableCell>
                        <FaEdit
                          className="text-blue-600 cursor-pointer"
                          onClick={() => navigate(`/admin/courses/edit/${course.ID}`)}
                        />
                      </TableCell>

                    </TableRow>
                  ))}
              </TableBody>
            </Table>
          )}

      </div>
    </>
  )
}

export default ListCourse

