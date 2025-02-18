import React, { useState, useEffect } from 'react';
import { BookOpen, ArrowLeft, Search, Filter, ChevronLeft, ChevronRight, BookMarked, User, Clock } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import api from '../api';

// Mock data - in a real app, you'd fetch this from your API
const allCourses = [
  { id: 1, code: "CS401", name: "Data Structures and Algorithms", credits: 4, instructor: "Dr. Emily Chen", semester: 4, description: "Advanced data structures and algorithm design techniques including analysis of algorithms, recursion, divide-and-conquer, dynamic programming, and graph algorithms." },
  { id: 2, code: "CS402", name: "Database Management Systems", credits: 3, instructor: "Prof. Michael Brown", semester: 4, description: "Fundamentals of database design, development, and management. Topics include ER modeling, relational model, SQL, indexing, and transaction processing." },
  { id: 3, code: "CS403", name: "Computer Networks", credits: 3, instructor: "Dr. Sarah Johnson", semester: 4, description: "Principles of computer networks, protocols, architectures, and the Internet. Topics include TCP/IP, routing, and network applications." },
  { id: 4, code: "CS404", name: "Software Engineering", credits: 4, instructor: "Prof. David Wilson", semester: 4, description: "Software development lifecycle, project management, requirements analysis, design, implementation, testing, and maintenance." },
  { id: 5, code: "MATH301", name: "Discrete Mathematics", credits: 3, instructor: "Dr. Robert Taylor", semester: 3, description: "Mathematical foundations including logic, sets, relations, functions, induction, recursion, counting, and graph theory with applications to computer science." },
  { id: 6, code: "CS301", name: "Operating Systems", credits: 4, instructor: "Dr. Lisa Wang", semester: 3, description: "Concepts and principles of operating systems including process management, memory management, file systems, and resource allocation." },
  { id: 7, code: "CS302", name: "Computer Architecture", credits: 3, instructor: "Prof. James Miller", semester: 3, description: "Organization and design of computer systems including processor architecture, memory hierarchy, I/O systems, and performance evaluation." },
  { id: 8, code: "CS201", name: "Object-Oriented Programming", credits: 4, instructor: "Dr. Karen Lee", semester: 2, description: "Advanced concepts of object-oriented programming including inheritance, polymorphism, interfaces, and design patterns." },
  { id: 9, code: "CS202", name: "Data Communications", credits: 3, instructor: "Prof. Thomas White", semester: 2, description: "Fundamentals of data communications, network architectures, protocols, and technologies with emphasis on Internet protocols." },
  { id: 10, code: "MATH201", name: "Calculus II", credits: 4, instructor: "Dr. Jennifer Adams", semester: 2, description: "Techniques of integration, applications of the definite integral, infinite series, polar coordinates, and parametric equations." },
  { id: 11, code: "CS101", name: "Introduction to Programming", credits: 4, instructor: "Prof. Richard Green", semester: 1, description: "Fundamentals of programming concepts, problem-solving, algorithm development, and implementation using a high-level programming language." },
  { id: 12, code: "CS102", name: "Discrete Structures", credits: 3, instructor: "Dr. Michelle Scott", semester: 1, description: "Introduction to mathematical structures for computer science including logic, sets, relations, functions, induction, and combinatorics." }
];


const CoursesPage = () => {
  const [courses, setCourses] = useState(allCourses);
  const [searchTerm, setSearchTerm] = useState('');
  const [currentFilter, setCurrentFilter] = useState('all');
  const [currentPage, setCurrentPage] = useState(1);
  const coursesPerPage = 6;

  // In a real app, you'd fetch data like this:
  // useEffect(() => {
  //   const fetchCourses = async () => {
  //     const response = await fetch('/api/courses');
  //     const data = await response.json();
  //     setCourses(data);
  //   };
  //   fetchCourses();
  // }, []);

  // Filter and search courses
  const filteredCourses = courses.filter(course => {
    const matchesSearch =
      course.code.toLowerCase().includes(searchTerm.toLowerCase()) ||
      course.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      course.instructor.toLowerCase().includes(searchTerm.toLowerCase());

    const matchesFilter = currentFilter === 'all' || currentFilter === course.semester.toString();

    return matchesSearch && matchesFilter;
  });

  // Pagination logic
  const indexOfLastCourse = currentPage * coursesPerPage;
  const indexOfFirstCourse = indexOfLastCourse - coursesPerPage;
  const currentCourses = filteredCourses.slice(indexOfFirstCourse, indexOfLastCourse);
  const totalPages = Math.ceil(filteredCourses.length / coursesPerPage);

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
          <div className="flex space-x-2">
            <div className="relative inline-block text-left">
              <div className="group">
                <button className="inline-flex justify-center w-full rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-100 focus:ring-blue-500">
                  <Filter size={18} className="mr-2" />
                  <span>Filter by Semester</span>
                </button>
                <div className="absolute right-0 mt-2 w-56 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 hidden group-hover:block z-10">
                  <div className="py-1" role="menu" aria-orientation="vertical">
                    <button
                      onClick={() => setCurrentFilter('all')}
                      className={`block px-4 py-2 text-sm w-full text-left ${currentFilter === 'all' ? 'bg-gray-100' : ''}`}
                    >
                      All Semesters
                    </button>
                    {[1, 2, 3, 4].map(semester => (
                      <button
                        key={semester}
                        onClick={() => setCurrentFilter(semester.toString())}
                        className={`block px-4 py-2 text-sm w-full text-left ${currentFilter === semester.toString() ? 'bg-gray-100' : ''}`}
                      >
                        Semester {semester}
                      </button>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Courses Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {currentCourses.length > 0 ? (
            currentCourses.map((course) => (
              <div key={course.id} className="bg-white shadow rounded-lg overflow-hidden hover:shadow-md transition-shadow duration-300">
                <div className="p-6">
                  <div className="flex items-start justify-between">
                    <div>
                      <p className="inline-block px-2 py-1 rounded text-xs font-semibold bg-blue-100 text-blue-800 mb-2">
                        {course.code}
                      </p>
                      <h3 className="text-lg font-semibold text-gray-900 mb-2">{course.name}</h3>
                    </div>
                    <p className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                      {course.credits} Credits
                    </p>
                  </div>

                  <p className="text-sm text-gray-600 mt-2 line-clamp-3">{course.description}</p>

                  <div className="mt-4 space-y-2">
                    <div className="flex items-center text-sm text-gray-500">
                      <User size={16} className="mr-2 text-gray-400" />
                      {course.instructor}
                    </div>
                    <div className="flex items-center text-sm text-gray-500">
                      <BookMarked size={16} className="mr-2 text-gray-400" />
                      Semester {course.semester}
                    </div>
                  </div>

                  <div className="mt-4 pt-3 border-t border-gray-200">
                    <button className="w-full text-blue-600 hover:text-blue-800 text-sm font-medium">
                      View Course Details
                    </button>
                  </div>
                </div>
              </div>
            ))
          ) : (
            <div className="col-span-full text-center py-10 text-gray-500">
              No courses found matching your criteria
            </div>
          )}
        </div>

        {/* Pagination */}
        {filteredCourses.length > coursesPerPage && (
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
                  Showing <span className="font-medium">{indexOfFirstCourse + 1}</span> to <span className="font-medium">{Math.min(indexOfLastCourse, filteredCourses.length)}</span> of{' '}
                  <span className="font-medium">{filteredCourses.length}</span> courses
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
      <footer className="bg-white border-t border-gray-200 mt-8">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <p className="text-sm text-center text-gray-500">© 2025 Result-E Student Portal. All rights reserved.</p>
        </div>
      </footer>
    </div>
  );
};

export default CoursesPage;