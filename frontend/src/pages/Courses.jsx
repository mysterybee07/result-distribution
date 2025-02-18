import React, { useState } from 'react';
import { BookOpen, ArrowLeft, Search, Filter, ChevronLeft, ChevronRight, BookMarked, User, Clock } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import { useData } from '../context/DataContext';
import api from '../api';

// Fetch courses function 
const fetchAllCourse = async () => {
  const response = await api.get(`/courses`);
  return response.data.courses;
};

const fetchCourse = async ({ queryKey }) => {
  const [, program_id, semester_id] = queryKey;
  if (!program_id || !semester_id) return [];

  const response = await api.get(`/courses/filter?program_id=${program_id}&semester_id=${semester_id}`);
  return response.data.courses;
};

const CoursesPage = () => {
  // State variables
  const [selectedProgram, setSelectedProgram] = useState("all");
  const [selectedSemester, setSelectedSemester] = useState("all");
  const [searchTerm, setSearchTerm] = useState('');
  const [search, setSearch] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const coursesPerPage = 6;

  // Use data hook to get programs
  const { programs } = useData();
  console.log("🚀 ~ CoursesPage ~ programs:", programs)

  // Fetch semesters based on selected program
  const {
    data: semesters = [],
    isLoading: loadingSemesters
  } = useQuery({
    queryKey: ["semesters", selectedProgram],
    queryFn: async () => {
      const response = await api.get(`/semester/by-program/${selectedProgram}`);
      return response.data.semesters;
    },
    enabled: selectedProgram !== "all",
  });
  console.log("🚀 ~ CoursesPage ~ semesters:", semesters)


  // Fetch all courses initially
  const {
    data: allCourses = [],
    isLoading: loadingAllCourses
  } = useQuery({
    queryKey: ["allCourses"],
    queryFn: fetchAllCourse,
  });

  // Fetch courses based on selected program & semester when search is triggered
  const {
    data: filteredApiCourses = [],
    isLoading: loadingFilteredCourses
  } = useQuery({
    queryKey: ["courses", selectedProgram, selectedSemester],
    queryFn: fetchCourse,
    enabled: search,
  });

  // Determine which courses to display
  const displayCourses = search ? filteredApiCourses : allCourses;

  // Client-side filtering for search term
  const searchFilteredCourses = displayCourses.filter(course =>
    searchTerm ? (
      course.code.toLowerCase().includes(searchTerm.toLowerCase()) ||
      course.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      course.instructor.toLowerCase().includes(searchTerm.toLowerCase())
    ) : true
  );

  // Handle search and filter submission
  const handleApplyFilters = () => {
    setSearch(true);
    setCurrentPage(1);
  };

  // Reset filters
  const handleResetFilters = () => {
    setSelectedProgram("all");
    setSelectedSemester("all");
    setSearchTerm('');
    setSearch(false);
    setCurrentPage(1);
  };

  // Pagination logic
  const indexOfLastCourse = currentPage * coursesPerPage;
  const indexOfFirstCourse = indexOfLastCourse - coursesPerPage;
  const currentCourses = searchFilteredCourses.slice(indexOfFirstCourse, indexOfLastCourse);
  const totalPages = Math.ceil(searchFilteredCourses.length / coursesPerPage);

  // Loading state
  const isLoading = loadingAllCourses || (search && loadingFilteredCourses);

  return (
    <div className="bg-gray-50 min-h-screen">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex items-center">
            <button
              onClick={() => window.history.back()}
              className="mr-4 p-1 rounded-full text-gray-400 hover:text-gray-500 hover:bg-gray-100"
            >
              <ArrowLeft size={20} />
            </button>
            <h1 className="text-2xl font-bold text-gray-900">All Courses</h1>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Search and Filter Bar */}
        <div className="mb-6 flex flex-col sm:flex-row space-y-4 sm:space-y-0 sm:space-x-4">
          <div className="relative flex-grow">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <Search size={18} className="text-gray-400" />
            </div>
            <input
              type="text"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              placeholder="Search courses by code, name, or instructor..."
              className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          {/* Program Filter */}
          <div className="sm:w-56">
            <select
              value={selectedProgram}
              onChange={(e) => {
                setSelectedProgram(e.target.value);
                setSelectedSemester("all"); // Reset semester when program changes
              }}
              className="block w-full rounded-md border border-gray-300 shadow-sm py-2 pl-3 pr-10 text-base focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="all">All Programs</option>
              {programs && programs.map(program => (
                <option key={program.ID} value={program.ID}>
                  {program.program_name}
                </option>
              ))}
            </select>
          </div>

          {/* Semester Filter - Enabled only when a program is selected */}
          <div className="sm:w-56">
            <select
              value={selectedSemester}
              onChange={(e) => setSelectedSemester(e.target.value)}
              disabled={selectedProgram === "all" || loadingSemesters}
              className="block w-full rounded-md border border-gray-300 shadow-sm py-2 pl-3 pr-10 text-base focus:outline-none focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 disabled:text-gray-500"
            >
              <option value="all">All Semesters</option>
              {semesters.map(semester => (
                <option key={semester.ID} value={semester.ID}>
                  Semester {semester.semester_name}
                </option>
              ))}
            </select>
          </div>

          {/* Filter Actions */}
          <div className="flex space-x-2">
            <button
              onClick={handleApplyFilters}
              className="px-4 py-2 bg-blue-600 text-white rounded-md shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              Apply Filters
            </button>
            <button
              onClick={handleResetFilters}
              className="px-4 py-2 border border-gray-300 text-gray-700 rounded-md shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              Reset
            </button>
          </div>
        </div>

        {/* Loading State */}
        {isLoading && (
          <div className="flex justify-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
          </div>
        )}

        {/* Courses Grid */}
        {!isLoading && (
          <div className="text-left grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {currentCourses.length > 0 ? (
              currentCourses.map((course) => (
                <div key={course.id} className="bg-white shadow rounded-lg overflow-hidden hover:shadow-md transition-shadow duration-300">
                  <div className="p-6">
                    <div className="flex items-start justify-between">
                      <div>
                        <p className="inline-block px-2 py-1 rounded text-xs font-semibold bg-blue-100 text-blue-800 mb-2">
                          {course.course_code}
                        </p>
                        <h3 className="text-lg font-semibold text-gray-900 mb-2">{course.name}</h3>
                      </div>
                      <p className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                        {course.credits} Credits
                      </p>
                    </div>

                    <p className="text-sm text-gray-600 mt-2 line-clamp-3">{course.description ?? "This course is very good"}</p>

                    <div className="mt-4 space-y-2">

                      <div className="flex items-center text-sm text-gray-500">
                        <BookMarked size={16} className="mr-2 text-gray-400" />
                        Semester {course.semester_id}
                      </div>
                      <div className="flex items-center text-sm text-gray-500">
                        <User size={16} className="mr-2 text-gray-400" />
                        {course.is_compulsory ? "Compulsory" : "Optional"}
                      </div>
                    </div>

                    {/* <div className="mt-4 pt-3 border-t border-gray-200">
                      <button className="w-full text-blue-600 hover:text-blue-800 text-sm font-medium">
                        View Course Details
                      </button>
                    </div> */}
                  </div>
                </div>
              ))
            ) : (
              <div className="col-span-full text-center py-10 text-gray-500">
                No courses found matching your criteria
              </div>
            )}
          </div>
        )}

        {/* Pagination */}
        {!isLoading && searchFilteredCourses.length > coursesPerPage && (
          <div className="flex items-center justify-between border-t border-gray-200 bg-white px-4 py-3 sm:px-6 mt-6 rounded-md shadow">
            <div className="flex flex-1 justify-between sm:hidden">
              <button
                onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
                disabled={currentPage === 1}
                className={`relative inline-flex items-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 ${currentPage === 1 ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-50'}`}
              >
                Previous
              </button>
              <button
                onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
                disabled={currentPage === totalPages}
                className={`relative ml-3 inline-flex items-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 ${currentPage === totalPages ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-50'}`}
              >
                Next
              </button>
            </div>
            <div className="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
              <div>
                <p className="text-sm text-gray-700">
                  Showing <span className="font-medium">{indexOfFirstCourse + 1}</span> to <span className="font-medium">{Math.min(indexOfLastCourse, searchFilteredCourses.length)}</span> of{' '}
                  <span className="font-medium">{searchFilteredCourses.length}</span> courses
                </p>
              </div>
              <div>
                <nav className="isolate inline-flex -space-x-px rounded-md shadow-sm" aria-label="Pagination">
                  <button
                    onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
                    disabled={currentPage === 1}
                    className={`relative inline-flex items-center rounded-l-md border border-gray-300 bg-white px-2 py-2 text-sm font-medium text-gray-500 ${currentPage === 1 ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-50'}`}
                  >
                    <span className="sr-only">Previous</span>
                    <ChevronLeft className="h-5 w-5" aria-hidden="true" />
                  </button>
                  {[...Array(totalPages)].map((_, i) => (
                    <button
                      key={i}
                      onClick={() => setCurrentPage(i + 1)}
                      className={`relative inline-flex items-center border ${currentPage === i + 1 ? 'bg-blue-50 border-blue-500 text-blue-600 z-10' : 'border-gray-300 bg-white text-gray-500 hover:bg-gray-50'} px-4 py-2 text-sm font-medium`}
                    >
                      {i + 1}
                    </button>
                  ))}
                  <button
                    onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
                    disabled={currentPage === totalPages}
                    className={`relative inline-flex items-center rounded-r-md border border-gray-300 bg-white px-2 py-2 text-sm font-medium text-gray-500 ${currentPage === totalPages ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-50'}`}
                  >
                    <span className="sr-only">Next</span>
                    <ChevronRight className="h-5 w-5" aria-hidden="true" />
                  </button>
                </nav>
              </div>
            </div>
          </div>
        )}
      </main>

      {/* Footer */}
      <footer className="bg-white border-t border-gray-200 p-4 mt-12">
        <div className="max-w-7xl mx-auto text-center text-sm text-gray-500">
          <p>© {new Date().getFullYear()} University Course Portal. All rights reserved.</p>
        </div>
      </footer>
    </div>
  );
};

export default CoursesPage;