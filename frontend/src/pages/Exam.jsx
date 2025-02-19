import React, { useState } from 'react';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import { Calendar, Clock, MapPin, Download, ArrowLeft, CalendarPlus } from 'lucide-react';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';
import { useQuery } from '@tanstack/react-query';
import api from '../api';
import { useExamSchedules } from '../hooks/useExamSchedules';

// Fix for the Leaflet icon issue in React
delete L.Icon.Default.prototype._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png',
  iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png',
  shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png',
});


const ExamPage = () => {
  // Define default values for batchID, programID, and semesterID
  const defaultBatchID = 1;
  const defaultProgramID = 1;
  const defaultSemesterID = 1;
  // For a random location in Kathmandu
  const examCenter = { lat: 27.7172, lng: 85.3240, name: "Kathmandu University Examination Center" };

  const [filters, setFilters] = useState({
    batchID: 1,
    programID: 1,
    semesterID: 1,
  });
  
  const { data: examSchedule, isLoading, error } = useExamSchedules(filters);
  console.log("🚀 ~ ExamPage ~ examSchedule:", examSchedule)

  // Function to add an exam to Google Calendar
  const addToGoogleCalendar = (exam) => {
    const startDate = new Date(exam.date);
    const endDate = new Date(startDate.getTime() + exam.duration * 60000);

    const startTimeFormatted = startDate.toISOString().replace(/-|:|\.\d+/g, '');
    const endTimeFormatted = endDate.toISOString().replace(/-|:|\.\d+/g, '');

    const googleCalendarUrl = `https://www.google.com/calendar/render?action=TEMPLATE&text=${encodeURIComponent(`${exam.courseCode} Exam: ${exam.courseName}`)}&dates=${startTimeFormatted}/${endTimeFormatted}&details=${encodeURIComponent(`${exam.courseCode} Final Examination at ${examCenter.name}`)}&location=${encodeURIComponent(`${examCenter.name} (${examCenter.lat}, ${examCenter.lng})`)}&sf=true&output=xml`;

    window.open(googleCalendarUrl, '_blank');
  };

  // Function to download the entire exam schedule as an ICS file
  const downloadExamScheduleICS = () => {
    let icsContent = [
      'BEGIN:VCALENDAR',
      'VERSION:2.0',
      'PRODID:-//University//Exam Schedule//EN'
    ];

    examSchedule.forEach(exam => {
      const startDate = new Date(exam.date);
      const endDate = new Date(startDate.getTime() + exam.duration * 60000);

      const startTimeFormatted = startDate.toISOString().replace(/[-:]/g, '').split('.')[0] + 'Z';
      const endTimeFormatted = endDate.toISOString().replace(/[-:]/g, '').split('.')[0] + 'Z';

      icsContent = [
        ...icsContent,
        'BEGIN:VEVENT',
        `DTSTART:${startTimeFormatted}`,
        `DTEND:${endTimeFormatted}`,
        `SUMMARY:${exam.courseCode} Exam: ${exam.courseName}`,
        `DESCRIPTION:${exam.courseCode} Final Examination at ${examCenter.name}`,
        `LOCATION:${examCenter.name} (${examCenter.lat}, ${examCenter.lng})`,
        'END:VEVENT'
      ];
    });

    icsContent.push('END:VCALENDAR');

    const blob = new Blob([icsContent.join('\r\n')], { type: 'text/calendar' });
    const url = URL.createObjectURL(blob);

    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', 'exam-schedule.ics');
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  // Format date for display
  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: 'numeric',
      minute: 'numeric',
      hour12: true
    }).format(date);
  };

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
            <h1 className="text-2xl font-bold text-gray-900">Examination Schedule</h1>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Schedule Section */}
        <section className="mb-12">
          <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6">
            <h2 className="text-xl font-semibold text-gray-800 mb-4 md:mb-0">Your Upcoming Exams</h2>
            <div className="flex space-x-4">
              <button
                onClick={downloadExamScheduleICS}
                className="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                <Download className="mr-2 h-4 w-4" />
                Download Schedule
              </button>
            </div>
          </div>

          {isLoading ? (
            <div className="flex justify-center py-12">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
            </div>
          ) : (
            <div className="bg-white shadow overflow-hidden sm:rounded-md">
              {examSchedule.length > 0 ? (
                <ul role="list" className=" text-left divide-y divide-gray-200">
                  {examSchedule.map((exam) => (
                    <li key={exam.id}>
                      <div className="px-4 py-4 sm:px-6">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center">
                            <div className="flex-shrink-0 h-10 w-10 bg-blue-100 rounded-full flex items-center justify-center">
                              <Calendar className="h-5 w-5 text-blue-600" />
                            </div>
                            <div className="ml-4">
                              <p className="text-sm font-medium text-gray-900">{exam.courseCode}</p>
                              <p className="text-sm text-gray-500">{exam.course}</p>
                            </div>
                          </div>
                          <button
                            onClick={() => addToGoogleCalendar(exam)}
                            className="inline-flex items-center px-3 py-1 border border-transparent text-xs font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                          >
                            <CalendarPlus className="h-4 w-4 mr-1" />
                            Add to Calendar
                          </button>
                        </div>
                        <div className="mt-2 sm:flex sm:justify-between">
                          <div className="sm:flex">
                            <p className="flex items-center text-sm text-gray-500">
                              <Clock className="flex-shrink-0 mr-1.5 h-4 w-4 text-gray-400" />
                              {formatDate(exam.exam_date)}
                            </p>
                          </div>
                          <div className="mt-2 flex items-center text-sm text-gray-500 sm:mt-0">
                            <p>Duration: 3 hours</p>
                          </div>
                        </div>
                      </div>
                    </li>
                  ))}
                </ul>
              ) : (
                <div className="text-center py-12">
                  <p className="text-gray-500">No upcoming exams scheduled.</p>
                </div>
              )}
            </div>
          )}
        </section>

        {/* Exam Center Map Section */}
        <section className="mb-12">
          <h2 className="text-xl font-semibold text-gray-800 mb-6">Exam Center Location</h2>

          <div className="bg-white shadow-md rounded-lg overflow-hidden">
            <div className="p-4 border-b">
              <div className="flex items-start">
                <MapPin className="h-5 w-5 text-red-500 mt-1 flex-shrink-0" />
                <div className="ml-3">
                  <h3 className="text-lg font-medium text-gray-900">{examCenter.name}</h3>
                  <p className="mt-1 text-sm text-gray-500">
                    Coordinates: {examCenter.lat.toFixed(4)}, {examCenter.lng.toFixed(4)}
                  </p>
                </div>
              </div>
            </div>

            <div className="h-96 z-0">
              <MapContainer
                center={[examCenter.lat, examCenter.lng]}
                zoom={15}
                scrollWheelZoom={false}
                style={{ height: '100%', width: '100%' }}
              >
                <TileLayer
                  attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                  url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                />
                <Marker position={[examCenter.lat, examCenter.lng]}>
                  <Popup>
                    <div className="text-center">
                      <h3 className="font-semibold">{examCenter.name}</h3>
                      <p className="text-sm text-gray-600">This is your examination center</p>
                    </div>
                  </Popup>
                </Marker>
              </MapContainer>
            </div>

            <div className="p-4 bg-gray-50">
              <h4 className="text-md font-medium text-gray-900 mb-2">Getting Here</h4>
              <p className="text-sm text-gray-600">
                The examination center is located in central Kathmandu. We recommend arriving at least 30 minutes before your scheduled exam time. Public transportation is available, with bus routes 101 and 205 stopping nearby.
              </p>
            </div>
          </div>
        </section>

        {/* Additional Information */}
        <section className="mb-12">
          <h2 className="text-xl font-semibold text-gray-800 mb-6">Examination Guidelines</h2>

          <div className="bg-white shadow-md rounded-lg overflow-hidden">
            <div className="p-6">
              <div className="space-y-4">
                <div>
                  <h3 className="text-lg font-medium text-gray-900">What to Bring</h3>
                  <ul className="mt-2 text-sm text-gray-600 list-disc pl-5 space-y-1">
                    <li>Your university ID card</li>
                    <li>Admission card / hall ticket</li>
                    <li>Blue or black pens</li>
                    <li>Calculator (if allowed for your exam)</li>
                    <li>Water bottle (clear, no labels)</li>
                  </ul>
                </div>

                <div>
                  <h3 className="text-lg font-medium text-gray-900">Rules During Examination</h3>
                  <ul className="mt-2 text-sm text-gray-600 list-disc pl-5 space-y-1">
                    <li>No electronic devices allowed (except approved calculators)</li>
                    <li>Students arriving more than 30 minutes late will not be admitted</li>
                    <li>No leaving the examination hall during the first hour and last 30 minutes</li>
                    <li>Maintain silence at all times</li>
                    <li>All materials provided must be returned at the end of the exam</li>
                  </ul>
                </div>

                <div className="bg-yellow-50 p-4 rounded-md">
                  <h3 className="text-lg font-medium text-yellow-800">Important Notice</h3>
                  <p className="mt-2 text-sm text-yellow-700">
                    COVID-19 safety protocols will be in effect. Please wear a mask and maintain social distancing. Temperature checks may be conducted at entry points.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </section>
      </main>

      {/* Footer */}
      <footer className="bg-white border-t border-gray-200 p-4 mt-12">
        <div className="max-w-7xl mx-auto text-center text-sm text-gray-500">
          <p>© {new Date().getFullYear()} University Examination Portal. All rights reserved.</p>
        </div>
      </footer>
    </div>
  );
};

export default ExamPage;