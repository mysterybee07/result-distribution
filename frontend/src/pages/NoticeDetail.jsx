import { useParams, useNavigate } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { ArrowLeft } from "lucide-react"; // Ensure you have lucide-react installed
import api from "../api";

const NoticeDetailPage = () => {
  const { id } = useParams(); // Get ID from the URL
  const navigate = useNavigate(); // Navigation for the back button

  // Function to fetch notice details
  const fetchNoticeById = async () => {
    const response = await api.get(`/notice/by-id/${id}`);
    return response.data;
  };

  const { data: notice, isLoading, isError, error } = useQuery({
    queryKey: ["notice", id],
    queryFn: fetchNoticeById,
    enabled: !!id, // Ensure it only runs if ID exists
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50">
        <p className="text-xl font-semibold text-gray-600 animate-pulse">Loading notice...</p>
      </div>
    );
  }

  if (isError) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50">
        <p className="text-xl font-semibold text-red-600">Error: {error.message}</p>
      </div>
    );
  }

  return (
    <div className="bg-gray-50 min-h-screen p-4 sm:p-8">
      {/* Header */}
      <header className="max-w-4xl mx-auto mb-8">
        <button
          onClick={() => navigate(-1)}
          className="flex items-center text-blue-600 hover:underline mb-6"
        >
          <ArrowLeft className="w-5 h-5 mr-2" />
          Back to Notices
        </button>
        <h1 className="text-4xl font-bold text-gray-900">{notice.Title}</h1>
      </header>

      {/* Notice Content */}
      <main className="max-w-4xl mx-auto bg-white rounded-2xl shadow-lg overflow-hidden">
        <div className="p-8">
          {/* Notice Details */}
          <div className="mb-6">
            <p className="text-gray-700 mb-4">{notice.Description}</p>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-6">
              <div>
                <p className="text-sm text-gray-500 uppercase">Batch</p>
                <p className="text-lg font-medium text-gray-800">{notice.Batch}</p>
              </div>

              <div>
                <p className="text-sm text-gray-500 uppercase">Program</p>
                <p className="text-lg font-medium text-gray-800">{notice.Program}</p>
              </div>

              <div>
                <p className="text-sm text-gray-500 uppercase">Semester</p>
                <p className="text-lg font-medium text-gray-800">{notice.Semester}</p>
              </div>

              <div>
                <p className="text-sm text-gray-500 uppercase">Created On</p>
                <p className="text-lg font-medium text-gray-800">
                  {new Date(notice.Created_at).toLocaleDateString(undefined, {
                    year: "numeric",
                    month: "long",
                    day: "numeric",
                    hour: "2-digit",
                    minute: "2-digit",
                  })}
                </p>
              </div>
            </div>
          </div>

          {/* File Section */}
          {notice.FilePath && (
            <div className="mt-8">
              <p className="text-sm text-gray-500 uppercase mb-2">Attached File</p>
              <div className="flex items-center space-x-4">
                <img
                  src={`/${notice.FilePath}`}
                  alt="Attached file"
                  className="max-w-[200px] rounded-lg shadow-md"
                />
                <a
                  href={`/${notice.FilePath}`}
                  download
                  className="text-blue-600 hover:underline"
                >
                  Download File
                </a>
              </div>
            </div>
          )}
        </div>
      </main>

      {/* Footer */}
      <footer className="mt-12 text-center text-gray-600">
        <p>© 2025 Result-E Student Portal. All rights reserved.</p>
      </footer>
    </div>
  );
};

export default NoticeDetailPage;
