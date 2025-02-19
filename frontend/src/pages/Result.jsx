import React, { useState, useEffect } from 'react';
import { ArrowDown, ArrowUp, BookOpen, Calendar, Download, Award, ArrowLeft } from 'lucide-react';
import { useAuth } from '../context/AuthContext';

const ResultsPage = () => {
  const [expandedCourse, setExpandedCourse] = useState(null);
  const [selectedSemester, setSelectedSemester] = useState(null);
  const [semesters, setSemesters] = useState([]);
  const [semesterMarks, setSemesterMarks] = useState({});

  const { profileData, profileLoading } = useAuth();

  useEffect(() => {
    if (profileData?.Marks) {
      // Get unique semesters from marks
      const uniqueSemesters = [...new Set(profileData.Marks.map(mark => mark.semester_id))];
      setSemesters(uniqueSemesters.map(id => ({
        id,
        name: `Semester ${id}`
      })));

      // Group marks by semester
      const marksBySemester = profileData.Marks.reduce((acc, mark) => {
        if (!acc[mark.semester_id]) {
          acc[mark.semester_id] = [];
        }
        acc[mark.semester_id].push(mark);
        return acc;
      }, {});

      setSemesterMarks(marksBySemester);
      setSelectedSemester(uniqueSemesters[0]);
    }
  }, [profileData]);

  const handleToggleCourse = (courseId) => {
    setExpandedCourse(expandedCourse === courseId ? null : courseId);
  };

  const calculateGrade = (totalMarks, totalPossibleMarks) => {
    const percentage = (totalMarks / totalPossibleMarks) * 100;
    if (percentage >= 80) return 'A';
    if (percentage >= 70) return 'B';
    if (percentage >= 60) return 'C';
    if (percentage >= 50) return 'D';
    return 'F';
  };

  const calculateCGPA = (marks) => {
    if (!marks || marks.length === 0) return "0.00";
    
    const totalPoints = marks.reduce((acc, course) => {
      const totalPossibleMarks = course.Course.semester_total_marks + 
                                course.Course.practical_total_marks + 
                                course.Course.assistant_total_marks;
      const grade = calculateGrade(course.total_marks, totalPossibleMarks);
      const points = grade === 'A' ? 4 : grade === 'B' ? 3 : grade === 'C' ? 2 : grade === 'D' ? 1 : 0;
      return acc + points;
    }, 0);

    return (totalPoints / marks.length).toFixed(2);
  };

  if (profileLoading) {
    return (
      <div className="flex justify-center items-center h-screen">
        <div className="animate-spin rounded-full h-16 w-16 border-t-2 border-b-2 border-indigo-500"></div>
      </div>
    );
  }

  const selectedSemesterMarks = semesterMarks[selectedSemester] || [];

  return (
    <div className="bg-gray-50 min-h-screen">
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

      <div className="md:hidden bg-white shadow-md rounded-md p-4 mx-4 mt-4">
        <p className="font-medium">{profileData?.Students?.fullname}</p>
        <p className="text-sm text-gray-600">Reg: {profileData?.Students?.registration_number}</p>
        <p className="text-sm text-gray-600">{profileData?.Students?.Program?.program_name}</p>
        <p className="text-sm text-gray-600">Batch: {profileData?.Students?.Batch?.batch}</p>
      </div>

      <div className="container mx-auto px-4 py-6">
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

        {selectedSemesterMarks.length > 0 && (
          <div className="bg-white rounded-lg shadow-md p-6 mb-6">
            <h2 className="text-xl font-semibold mb-4 flex items-center">
              <Calendar className="h-5 w-5 mr-2 text-indigo-500" />
              Semester {selectedSemester} Summary
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="bg-green-50 p-4 rounded-md">
                <p className="text-sm text-green-600">CGPA</p>
                <p className="text-2xl font-bold text-green-700">
                  {calculateCGPA(selectedSemesterMarks)}
                </p>
              </div>
              <div className="bg-blue-50 p-4 rounded-md">
                <p className="text-sm text-blue-600">Courses</p>
                <p className="text-2xl font-bold text-blue-700">
                  {selectedSemesterMarks.length}
                </p>
              </div>
              <div className="bg-purple-50 p-4 rounded-md">
                <p className="text-sm text-purple-600">Status</p>
                <p className="text-2xl font-bold text-purple-700">
                  {selectedSemesterMarks.every(m => m.status === 'pass') ? 'Passed' : 'Failed'}
                </p>
              </div>
            </div>
          </div>
        )}

        <div className="space-y-4">
          {selectedSemesterMarks.map(mark => (
            <div key={mark.ID} className="bg-white rounded-lg shadow-md overflow-hidden transition-all duration-200">
              <div
                className="p-4 cursor-pointer flex justify-between items-center"
                onClick={() => handleToggleCourse(mark.ID)}
              >
                <div className="flex items-center">
                  <BookOpen className="h-5 w-5 text-indigo-500 mr-3" />
                  <div>
                    <h3 className="font-medium text-gray-800">{mark.Course.name}</h3>
                    <p className="text-xs text-gray-500">Course Code: {mark.Course.course_code}</p>
                  </div>
                </div>
                <div className="flex items-center">
                  <div className="mr-4 text-right">
                    <p className="font-semibold text-gray-800">
                      {mark.total_marks}/{mark.Course.semester_total_marks + mark.Course.practical_total_marks + mark.Course.assistant_total_marks}
                    </p>
                    <p className={`text-xs ${mark.status === 'pass' ? 'text-green-600' : 'text-red-600'}`}>
                      {mark.status.toUpperCase()}
                    </p>
                  </div>
                  {expandedCourse === mark.ID ? 
                    <ArrowUp className="h-5 w-5 text-gray-400" /> : 
                    <ArrowDown className="h-5 w-5 text-gray-400" />
                  }
                </div>
              </div>

              {expandedCourse === mark.ID && (
                <div className="bg-gray-50 p-4 border-t border-gray-200">
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                    <div className="bg-white p-3 rounded shadow-sm">
                      <p className="text-xs text-gray-500">Semester Marks</p>
                      <p className="text-lg font-semibold">
                        {mark.semester_marks}/{mark.Course.semester_total_marks}
                      </p>
                    </div>
                    <div className="bg-white p-3 rounded shadow-sm">
                      <p className="text-xs text-gray-500">Assistant Marks</p>
                      <p className="text-lg font-semibold">
                        {mark.assistant_marks}/{mark.Course.assistant_total_marks}
                      </p>
                    </div>
                    <div className="bg-white p-3 rounded shadow-sm">
                      <p className="text-xs text-gray-500">Practical Marks</p>
                      <p className="text-lg font-semibold">
                        {mark.practical_marks}/{mark.Course.practical_total_marks}
                      </p>
                    </div>
                  </div>

                  <div className="flex justify-between items-center bg-indigo-50 p-3 rounded">
                    <div className="flex items-center">
                      <Award className="h-5 w-5 text-indigo-600 mr-2" />
                      <div>
                        <p className="text-xs text-indigo-600">Final Grade</p>
                        <p className="text-xl font-bold text-indigo-700">
                          {calculateGrade(
                            mark.total_marks,
                            mark.Course.semester_total_marks + 
                            mark.Course.practical_total_marks + 
                            mark.Course.assistant_total_marks
                          )}
                        </p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="text-xs text-indigo-600">Total Score</p>
                      <p className="text-xl font-bold text-indigo-700">
                        {mark.total_marks}/{mark.Course.semester_total_marks + mark.Course.practical_total_marks + mark.Course.assistant_total_marks}
                      </p>
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