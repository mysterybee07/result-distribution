import React, { useState, useEffect } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { User, Bell, BookOpen, Calendar, FileText } from 'lucide-react';
import { useCourses } from '../hooks/useCourses';
import { useNotices } from '../hooks/useNotice';
import { useExamSchedules } from '../hooks/useExamSchedules';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import api from '../api';

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
};

const StudentDashboard = () => {
  const [student, setStudent] = useState(studentData);// Debugging log

  const { profileData, isLoading: profileLoading } = useAuth();
  console.log("🚀 ~ StudentDashboard ~ profileData:", profileData)

  useEffect(() => {
    if (profileData) {
      setStudent(profileData.Students);
    }
  }, [profileData]);

  const navigate = useNavigate();

  // storing data in local storage
  localStorage.setItem("selectedFilters", JSON.stringify({
    batchId: 1,
    programId: 1,
    semesterId: 1
  }));

  const storedFilters = JSON.parse(localStorage.getItem("selectedFilters")) || {
    batchId: 1,
    programId: 1,
    semesterId: 1
  };

  const {
    filteredCourses,
    loadingFilteredCourses,
  } = useCourses({
    programId: storedFilters.programId,
    semesterId: storedFilters.semesterId,
    enableFiltering: true // Always enabled
  });

  const {
    data: notices,
    isLoading,
    isError,
    error
  } = useNotices({
    programId: storedFilters.programId,
    batchId: storedFilters.semesterId
  });
  const { data: examSchedule, isLoading: examScheduleLoading, error: examScheduleError } = useExamSchedules(storedFilters);
  // console.log("🚀 ~ StudentDashboard ~ examSchedule:", examSchedule)

  // Helper function to format dates
  const formatDate = (dateString) => {
    const options = { year: 'numeric', month: 'short', day: 'numeric' };
    return new Date(dateString).toLocaleDateString(undefined, options);
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
              {/* <img
                src={student.profile.avatarUrl}
                alt={student.profile.name}
                className="w-20 h-20 rounded-full border-4 border-blue-100"
              /> */}
              <div className="">
                <h2 className="text-left text-xl font-semibold text-gray-900">{student.fullname}</h2>
                <p className="text-sm text-gray-500">{student?.Users?.email ?? "user@gmail.com"}</p>
              </div>
            </div>

            <div className="mt-6 space-y-4">
              <div className="flex justify-between border-b pb-2">
                <span className="text-sm font-medium text-gray-500">Symbol Number</span>
                <span className="text-sm text-gray-900">{student?.symbol_number}</span>
              </div>
              <div className="flex justify-between border-b pb-2">
                <span className="text-sm font-medium text-gray-500">Registration Number</span>
                <span className="text-sm text-gray-900">{student?.registration_number}</span>
              </div>
              <div className="flex justify-between border-b pb-2">
                <span className="text-sm font-medium text-gray-500">Program</span>
                <span className="text-sm text-gray-900">{student?.Program?.program_name}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-sm font-medium text-gray-500">Current Semester</span>
                <span className="text-sm text-gray-900">{student?.Semester?.semester_name}</span>
              </div>
            </div>
          </div>

          {/* Recent Notices */}
          <div className="bg-white shadow rounded-lg p-6 col-span-2">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-900">Recent Notices</h2>
              <button
                className="text-sm text-blue-600 hover:text-blue-800"
                onClick={() => navigate("/notices")}
              >
                View All
              </button>
            </div>
            <div className="space-y-4">
              {notices?.slice(0, 4).map((notice) => (
                <div key={notice.ID} className="p-3 rounded-md border">
                  <div className="flex justify-between items-start">
                    <h3 className="text-sm font-medium">{notice.Title}</h3>
                    <span className="text-xs">{formatDate(notice.Created_at)}</span>
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
              {filteredCourses.map((course) => (
                <div key={course.course_code} className="p-3 border border-gray-200 rounded-md hover:bg-gray-50">
                  <div className="flex justify-between">
                    <span className="text-sm font-medium text-gray-900">{course.course_code}</span>
                    <span className="text-xs bg-blue-100 text-blue-800 px-2 py-0.5 rounded-full">4 Credits</span>
                  </div>
                  <p className="text-left text-sm text-gray-600 mt-1">{course.name}</p>
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
                  {examSchedule?.map((exam) => {
                    const courseInfo = filteredCourses.find(c => c.code === exam.course);
                    return (
                      <tr key={exam.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-left text-sm font-medium text-gray-900">{exam.course}</div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm text-gray-900">{formatDate(exam.exam_date)}</div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm text-gray-900">{exam.time ?? "12:00 - 03:00"}</div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm text-gray-900">{exam.venue ?? "CAB"}</div>
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