import React, { useState, useEffect } from 'react';
import { ArrowDown, ArrowUp, BookOpen, Calendar, Download, Award, ArrowLeft } from 'lucide-react';

const ResultsPage = () => {
  const [loading, setLoading] = useState(true);
  const [student, setStudent] = useState(null);
  const [marks, setMarks] = useState([]);
  const [selectedSemester, setSelectedSemester] = useState(null);
  const [semesters, setSemesters] = useState([]);
  const [expandedCourse, setExpandedCourse] = useState(null);

  useEffect(() => {
    // Simulating API call to fetch student data
    setTimeout(() => {
      const mockStudent = {
        id: 1,
        name: 'John Doe',
        registrationNumber: 'S12345',
        program: 'Bachelor of Computer Science',
        batch: '2021-2025'
      };
      setStudent(mockStudent);
      
      // Simulating marks data
      const mockMarks = [
        {
          semesterId: 1,
          semesterName: 'Semester 1',
          courses: [
            {
              courseId: 101,
              courseName: 'Introduction to Programming',
              semesterMarks: 75,
              assistantMarks: 18,
              practicalMarks: 38,
              totalMarks: 131,
              status: 'Pass'
            },
            {
              courseId: 102,
              courseName: 'Computer Architecture',
              semesterMarks: 68,
              assistantMarks: 16,
              practicalMarks: 35,
              totalMarks: 119,
              status: 'Pass'
            }
          ]
        },
        {
          semesterId: 2,
          semesterName: 'Semester 2',
          courses: [
            {
              courseId: 201,
              courseName: 'Data Structures',
              semesterMarks: 72,
              assistantMarks: 17,
              practicalMarks: 36,
              totalMarks: 125,
              status: 'Pass'
            },
            {
              courseId: 202,
              courseName: 'Database Systems',
              semesterMarks: 70,
              assistantMarks: 16,
              practicalMarks: 37,
              totalMarks: 123,
              status: 'Pass'
            }
          ]
        }
      ];
      
      setMarks(mockMarks);
      setSemesters(mockMarks.map(mark => ({ 
        id: mark.semesterId, 
        name: mark.semesterName 
      })));
      setSelectedSemester(mockMarks[0].semesterId);
      setLoading(false);
    }, 1500);
  }, []);

  const handleToggleCourse = (courseId) => {
    if (expandedCourse === courseId) {
      setExpandedCourse(null);
    } else {
      setExpandedCourse(courseId);
    }
  };

  const calculateGrade = (totalMarks) => {
    if (totalMarks >= 140) return 'A';
    if (totalMarks >= 120) return 'B';
    if (totalMarks >= 100) return 'C';
    if (totalMarks >= 80) return 'D';
    return 'F';
  };

  const calculateCGPA = (semesterMarks) => {
    const totalCourses = semesterMarks.courses.length;
    const totalPoints = semesterMarks.courses.reduce((acc, course) => {
      const grade = calculateGrade(course.totalMarks);
      const points = grade === 'A' ? 4 : grade === 'B' ? 3 : grade === 'C' ? 2 : grade === 'D' ? 1 : 0;
      return acc + points;
    }, 0);
    
    return (totalPoints / totalCourses).toFixed(2);
  };

  const selectedSemesterData = marks.find(mark => mark.semesterId === selectedSemester);

  if (loading) {
    return (
      <div className="flex justify-center items-center h-screen">
        <div className="animate-spin rounded-full h-16 w-16 border-t-2 border-b-2 border-indigo-500"></div>
      </div>
    );
  }

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
            <h1 className="text-2xl font-bold text-gray-900">Result Portal</h1>
          </div>
        </div>
      </header>

      {/* Student Info Card - Mobile */}
      <div className="md:hidden bg-white shadow-md rounded-md p-4 mx-4 mt-4">
        <p className="font-medium">{student.name}</p>
        <p className="text-sm text-gray-600">{student.registrationNumber}</p>
        <p className="text-sm text-gray-600">{student.program}</p>
        <p className="text-sm text-gray-600">{student.batch}</p>
      </div>

      {/* Content */}
      <div className="container mx-auto px-4 py-6">
        {/* Semester Navigation */}
        <div className="flex overflow-x-auto pb-2 mb-6 scrollbar-hide">
          <div className="flex space-x-4">
            {semesters.map(semester => (
              <button
                key={semester.id}
                className={`px-4 py-2 rounded-md whitespace-nowrap transition-colors ${
                  selectedSemester === semester.id 
                    ? 'bg-indigo-100 text-indigo-700 font-medium'
                    : 'bg-white text-gray-600 hover:bg-gray-100'
                }`}
                onClick={() => setSelectedSemester(semester.id)}
              >
                {semester.name}
              </button>
            ))}
          </div>
        </div>

        {/* Summary Card */}
        {selectedSemesterData && (
          <div className="bg-white rounded-lg shadow-md p-6 mb-6">
            <h2 className="text-xl font-semibold mb-4 flex items-center">
              <Calendar className="h-5 w-5 mr-2 text-indigo-500" />
              {selectedSemesterData.semesterName} Summary
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="bg-green-50 p-4 rounded-md">
                <p className="text-sm text-green-600">CGPA</p>
                <p className="text-2xl font-bold text-green-700">
                  {calculateCGPA(selectedSemesterData)}
                </p>
              </div>
              <div className="bg-blue-50 p-4 rounded-md">
                <p className="text-sm text-blue-600">Courses</p>
                <p className="text-2xl font-bold text-blue-700">
                  {selectedSemesterData.courses.length}
                </p>
              </div>
              <div className="bg-purple-50 p-4 rounded-md">
                <p className="text-sm text-purple-600">Status</p>
                <p className="text-2xl font-bold text-purple-700">
                  {selectedSemesterData.courses.every(c => c.status === 'Pass') ? 'Passed' : 'Failed'}
                </p>
              </div>
            </div>
            <div className="mt-4 flex justify-end">
              <button className="flex items-center text-indigo-600 hover:text-indigo-800 text-sm font-medium">
                <Download className="h-4 w-4 mr-1" />
                Download Report
              </button>
            </div>
          </div>
        )}

        {/* Course Cards */}
        <div className="space-y-4">
          {selectedSemesterData?.courses.map(course => (
            <div key={course.courseId} className="bg-white rounded-lg shadow-md overflow-hidden transition-all duration-200">
              <div 
                className="p-4 cursor-pointer flex justify-between items-center"
                onClick={() => handleToggleCourse(course.courseId)}
              >
                <div className="flex items-center">
                  <BookOpen className="h-5 w-5 text-indigo-500 mr-3" />
                  <div>
                    <h3 className="font-medium text-gray-800">{course.courseName}</h3>
                    <p className="text-xs text-gray-500">Course ID: {course.courseId}</p>
                  </div>
                </div>
                <div className="flex items-center">
                  <div className="mr-4 text-right">
                    <p className="font-semibold text-gray-800">{course.totalMarks}/200</p>
                    <p className={`text-xs ${course.status === 'Pass' ? 'text-green-600' : 'text-red-600'}`}>
                      {course.status}
                    </p>
                  </div>
                  {expandedCourse === course.courseId ? 
                    <ArrowUp className="h-5 w-5 text-gray-400" /> : 
                    <ArrowDown className="h-5 w-5 text-gray-400" />
                  }
                </div>
              </div>
              
              {expandedCourse === course.courseId && (
                <div className="bg-gray-50 p-4 border-t border-gray-200">
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                    <div className="bg-white p-3 rounded shadow-sm">
                      <p className="text-xs text-gray-500">Semester Marks</p>
                      <p className="text-lg font-semibold">{course.semesterMarks}/100</p>
                    </div>
                    <div className="bg-white p-3 rounded shadow-sm">
                      <p className="text-xs text-gray-500">Assistant Marks</p>
                      <p className="text-lg font-semibold">{course.assistantMarks}/20</p>
                    </div>
                    <div className="bg-white p-3 rounded shadow-sm">
                      <p className="text-xs text-gray-500">Practical Marks</p>
                      <p className="text-lg font-semibold">{course.practicalMarks}/80</p>
                    </div>
                  </div>
                  
                  <div className="flex justify-between items-center bg-indigo-50 p-3 rounded">
                    <div className="flex items-center">
                      <Award className="h-5 w-5 text-indigo-600 mr-2" />
                      <div>
                        <p className="text-xs text-indigo-600">Final Grade</p>
                        <p className="text-xl font-bold text-indigo-700">{calculateGrade(course.totalMarks)}</p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="text-xs text-indigo-600">Total Score</p>
                      <p className="text-xl font-bold text-indigo-700">{course.totalMarks}/200</p>
                    </div>
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default ResultsPage;