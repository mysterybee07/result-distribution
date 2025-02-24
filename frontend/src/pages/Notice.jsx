import React, { useState, useEffect } from 'react';
import { Bell, ArrowLeft, Search, Filter, ChevronLeft, ChevronRight } from 'lucide-react';
import { Link } from 'react-router-dom';
import { useNotices } from '../hooks/useNotice';

const NoticesPage = () => {
  const [notices, setNotices] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const noticesPerPage = 5;

  let selectedProgram = null;
  let selectedBatch = null;

  const { 
    data, 
    isLoading, 
    isError, 
    error 
  } = useNotices({
    programId: selectedProgram,
    batchId: selectedBatch
  });

  useEffect(() => {
    if (data) {
      setNotices(data);
    }
  }, [data]);

  // Helper function to format dates
  const formatDate = (dateString) => {
    const options = { year: 'numeric', month: 'long', day: 'numeric' };
    return new Date(dateString).toLocaleDateString(undefined, options);
  };

  // Filter and search notices (no priority filter)
  const filteredNotices = notices?.filter(notice => {
    const matchesSearch = notice.Title.toLowerCase().includes(searchTerm.toLowerCase()) ||
      notice.Description.toLowerCase().includes(searchTerm.toLowerCase());
    return matchesSearch;
  });

  // Pagination logic
  const indexOfLastNotice = currentPage * noticesPerPage;
  const indexOfFirstNotice = indexOfLastNotice - noticesPerPage;
  const currentNotices = filteredNotices?.slice(indexOfFirstNotice, indexOfLastNotice);
  const totalPages = Math.ceil(filteredNotices?.length / noticesPerPage);

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
            <h1 className="text-2xl font-bold text-gray-900">All Notices</h1>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Search Bar */}
        <div className="mb-6">
          <div className="relative flex-grow">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <Search size={18} className="text-gray-400" />
            </div>
            <input
              type="text"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              placeholder="Search notices..."
              className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
        </div>

        {/* Notices List */}
        <div className="bg-white shadow overflow-hidden sm:rounded-md">
          <ul className="divide-y divide-gray-200">
            {currentNotices?.length > 0 ? (
              currentNotices?.map((notice) => (
                <li key={notice.ID} className="cursor-pointer hover:bg-gray-100 transition py-2 ">
                  <Link to={`/notice/${notice.ID}`} className="block px-6 py-4">
                    <div className="flex justify-between items-start">
                      <div>
                        <div className='flex gap-4'>
                          <h3 className="text-lg text-left font-medium text-gray-900">{notice.Title}</h3>
                          <div className="flex justify-center flex-wrap items-center gap-2">
                            <span className="px-3 py-1 bg-gray-200 text-gray-700 rounded-full text-sm font-semibold">
                              {notice?.Batch || "N/A"}
                            </span>
                            <span className="px-3 py-1 bg-gray-200 text-gray-700 rounded-full text-sm font-semibold">
                              {notice?.Program || "N/A"}
                            </span>
                            <span className="px-3 py-1 bg-gray-200 text-gray-700 rounded-full text-sm font-semibold">
                              Semester: {notice?.Semester || "N/A"}
                            </span>
                          </div>
                        </div>
                        <p className="mt-2 text-left text-sm text-gray-600">
                          {notice?.Description?.length > 300
                            ? `${notice.Description.slice(0, 300)}...`
                            : notice.Description}
                        </p>

                      </div>
                      <div className="flex flex-col text-right">
                        <span className="text-sm text-gray-500">{formatDate(notice.Created_at)}</span>
                      </div>
                    </div>
                  </Link>
                </li>
              ))
            ) : (
              <li className="px-6 py-4 text-center text-gray-500">
                No notices found matching your criteria
              </li>
            )}
          </ul>
        </div>

        {/* Pagination */}
        {filteredNotices?.length > noticesPerPage && (
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
                  Showing <span className="font-medium">{indexOfFirstNotice + 1}</span> to <span className="font-medium">{Math.min(indexOfLastNotice, filteredNotices?.length)}</span> of{' '}
                  <span className="font-medium">{filteredNotices?.length}</span> results
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


export default NoticesPage;