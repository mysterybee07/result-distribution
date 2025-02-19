import React, { useState, useEffect } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { User, Bell, BookOpen, Calendar, FileText } from 'lucide-react';

// Mock data - in a real app, you'd fetch this from your API
const studentData = {
  profile: {
    name: "John Doe",
    symbolNumber: "SYM12345",
    registrationNumber: "REG67890",
    currentSemester: 4,
    program: "Bachelor of Computer Science",
    email: "john.doe@example.com",
    avatarUrl: "/api/placeholder/150/150"
  },
  notices: [
    { id: 1, title: "Final Exam Schedule Published", date: "2025-02-15", priority: "high" },
    { id: 2, title: "Registration for Next Semester", date: "2025-02-12", priority: "medium" },
    { id: 3, title: "Holiday Notice: Founder's Day", date: "2025-02-10", priority: "low" },
    { id: 4, title: "Scholarship Applications Open", date: "2025-02-05", priority: "medium" }
  ],
  semesterResults: [
    { semester: "Semester 1", gpa: 3.5 },
    { semester: "Semester 2", gpa: 3.7 },
    { semester: "Semester 3", gpa: 3.6 },
    { semester: "Semester 4", gpa: 3.8 }
  ],
  currentCourses: [
    { code: "CS401", name: "Data Structures and Algorithms", credits: 4, instructor: "Dr. Emily Chen" },
    { code: "CS402", name: "Database Management Systems", credits: 3, instructor: "Prof. Michael Brown" },
    { code: "CS403", name: "Computer Networks", credits: 3, instructor: "Dr. Sarah Johnson" },
    { code: "CS404", name: "Software Engineering", credits: 4, instructor: "Prof. David Wilson" },
    { code: "MATH301", name: "Discrete Mathematics", credits: 3, instructor: "Dr. Robert Taylor" }
  ],
  upcomingExams: [
    { id: 1, course: "CS401", date: "2025-03-12", time: "10:00 AM", venue: "Hall A" },
    { id: 2, course: "CS402", date: "2025-03-15", time: "2:00 PM", venue: "Hall B" },
    { id: 3, course: "CS403", date: "2025-03-18", time: "10:00 AM", venue: "Hall C" },
    { id: 4, course: "CS404", date: "2025-03-21", time: "2:00 PM", venue: "Hall A" },
    { id: 5, course: "MATH301", date: "2025-03-24", time: "10:00 AM", venue: "Hall B" }
  ]
};

const StudentDashboard = () => {
  const [student, setStudent] = useState(studentData);
  
  // In a real app, you'd fetch data like this:
  // useEffect(() => {
  //   const fetchData = async () => {
  //     const response = await fetch('/api/student-data');
  //     const data = await response.json();
  //     setStudent(data);
  //   };
  //   fetchData();
  // }, []);

  // Helper function to format dates
  const formatDate = (dateString) => {
    const options = { year: 'numeric', month: 'short', day: 'numeric' };
    return new Date(dateString).toLocaleDateString(undefined, options);
  };

  // Helper function to determine notice priority classes
  const getPriorityClass = (priority) => {
    switch (priority) {
      case 'high': return 'bg-red-100 text-red-800 border-red-200';
      case 'medium': return 'bg-yellow-100 text-yellow-800 border-yellow-200';
      case 'low': return 'bg-green-100 text-green-800 border-green-200';
      default: return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  return (
    <div className="bg-gray-50 min-h-screen">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex justify-between items-center">
            <h1 className="text-2xl font-bold text-gray-900">Student Dashboard</h1>
            <div className="flex items-center space-x-4">
              <span className="text-sm text-gray-500">Last updated: {new Date().toLocaleDateString()}</span>
              <button className="p-1 rounded-full text-gray-400 hover:text-gray-500">
                <Bell size={20} />
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          
          {/* Student Profile */}
          <div className="bg-white shadow rounded-lg p-6 col-span-1">
            <div className="flex items-start">
              <img
                src={student.profile.avatarUrl}
                alt={student.profile.name}
                className="w-20 h-20 rounded-full border-4 border-blue-100"
              />
              <div className="ml-4">
                <h2 className="text-xl font-semibold text-gray-900">{student.profile.name}</h2>
                <p className="text-sm text-gray-500">{student.profile.email}</p>
              </div>
            </div>
            
            <div className="mt-6 space-y-4">
              <div className="flex justify-between border-b pb-2">
                <span className="text-sm font-medium text-gray-500">Symbol Number</span>
                <span className="text-sm text-gray-900">{student.profile.symbolNumber}</span>
              </div>
              <div className="flex justify-between border-b pb-2">
                <span className="text-sm font-medium text-gray-500">Registration Number</span>
                <span className="text-sm text-gray-900">{student.profile.registrationNumber}</span>
              </div>
              <div className="flex justify-between border-b pb-2">
                <span className="text-sm font-medium text-gray-500">Program</span>
                <span className="text-sm text-gray-900">{student.profile.program}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-sm font-medium text-gray-500">Current Semester</span>
                <span className="text-sm text-gray-900">{student.profile.currentSemester}</span>
              </div>
            </div>
          </div>
          
          {/* Recent Notices */}
          <div className="bg-white shadow rounded-lg p-6 col-span-2">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-900">Recent Notices</h2>
              <button className="text-sm text-blue-600 hover:text-blue-800">View All</button>
            </div>
            <div className="space-y-4">
              {student.notices.map((notice) => (
                <div key={notice.id} className={`p-3 rounded-md border ${getPriorityClass(notice.priority)}`}>
                  <div className="flex justify-between items-start">
                    <h3 className="text-sm font-medium">{notice.title}</h3>
                    <span className="text-xs">{formatDate(notice.date)}</span>
                  </div>
                  <div className="mt-1">
                    <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${
                      notice.priority === 'high' ? 'bg-red-100 text-red-800' : 
                      notice.priority === 'medium' ? 'bg-yellow-100 text-yellow-800' : 
                      'bg-green-100 text-green-800'
                    }`}>
                      {notice.priority.charAt(0).toUpperCase() + notice.priority.slice(1)} Priority
                    </span>
                  </div>
                </div>
              ))}
            </div>
          </div>
          
          {/* Progress Graph */}
          <div className="bg-white shadow rounded-lg p-6 col-span-1 lg:col-span-2">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Academic Progress</h2>
            <div className="h-64">
              <ResponsiveContainer width="100%" height="100%">
                <LineChart
                  data={student.semesterResults}
                  margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
                >
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="semester" />
                  <YAxis domain={[0, 4]} />
                  <Tooltip />
                  <Legend />
                  <Line
                    type="monotone"
                    dataKey="gpa"
                    stroke="#4F46E5"
                    activeDot={{ r: 8 }}
                    strokeWidth={2}
                  />
                </LineChart>
              </ResponsiveContainer>
            </div>
            <div className="mt-4 text-center">
              <p className="text-sm text-gray-500">Current CGPA: <span className="font-semibold text-blue-600">3.65</span></p>
            </div>
          </div>
          
          {/* Current Courses */}
          <div className="bg-white shadow rounded-lg p-6 col-span-1">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-900">Current Courses</h2>
              <BookOpen size={18} className="text-gray-400" />
            </div>
            <div className="space-y-3">
              {student.currentCourses.map((course) => (
                <div key={course.code} className="p-3 border border-gray-200 rounded-md hover:bg-gray-50">
                  <div className="flex justify-between">
                    <span className="text-sm font-medium text-gray-900">{course.code}</span>
                    <span className="text-xs bg-blue-100 text-blue-800 px-2 py-0.5 rounded-full">{course.credits} Credits</span>
                  </div>
                  <p className="text-sm text-gray-600 mt-1">{course.name}</p>
                  <p className="text-xs text-gray-500 mt-1">Instructor: {course.instructor}</p>
                </div>
              ))}
            </div>
          </div>
          
          {/* Upcoming Exams */}
          <div className="bg-white shadow rounded-lg p-6 col-span-1 lg:col-span-2">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-900">Upcoming Examinations</h2>
              <Calendar size={18} className="text-gray-400" />
            </div>
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Course</th>
                    <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                    <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Time</th>
                    <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Venue</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {student.upcomingExams.map((exam) => {
                    const courseInfo = student.currentCourses.find(c => c.code === exam.course);
                    return (
                      <tr key={exam.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm font-medium text-gray-900">{exam.course}</div>
                          <div className="text-sm text-gray-500">{courseInfo ? courseInfo.name : ''}</div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm text-gray-900">{formatDate(exam.date)}</div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm text-gray-900">{exam.time}</div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm text-gray-900">{exam.venue}</div>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </div>
        </div>
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

export default StudentDashboard;