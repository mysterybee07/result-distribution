import { useParams, useNavigate } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import api from "../api";

const NoticeDetailPage = () => {
    const { id } = useParams();
    const navigate = useNavigate();

    const fetchNoticeById = async () => {
        const response = await api.get(`/notice/by-id/${id}`);
        return response.data.notice;
    };

    const { data: notice = {}, isLoading, isError, error } = useQuery({
        queryKey: ["notice", id],
        queryFn: fetchNoticeById,
        enabled: !!id,
    });
    console.log("🚀 ~ NoticeDetailPage ~ notice:", notice)

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
        <div className="min-h-screen bg-gray-100 py-10 px-4 sm:px-6 lg:px-8">
            <div className="max-w-4xl mx-auto bg-white shadow-lg rounded-lg p-6">
                <div className="flex items-center mb-6">
                    <button onClick={() => navigate(-1)} className="text-gray-500 hover:text-blue-500">
                        ← Back
                    </button>
                </div>

                <h2 className="text-3xl font-bold text-gray-800 mb-4">{notice?.title || "No Title"}</h2>


                <div className="grid grid-cols-1 md:grid-cols-1 gap-6">
                    <p className="text-gray-700">
                        <span className="font-semibold">Published At:</span>{" "}
                        {notice?.Created_at ? new Date(notice.CreatedAt).toLocaleString() : "N/A"}
                    </p>
                    <div className="flex justify-center flex-wrap items-center gap-2">
                        <span className="px-3 py-1 bg-gray-200 text-gray-700 rounded-full text-sm font-semibold">
                            Batch: {notice?.batch_id || "N/A"}
                        </span>
                        <span className="px-3 py-1 bg-gray-200 text-gray-700 rounded-full text-sm font-semibold">
                            Program: {notice?.program_id || "N/A"}
                        </span>
                        <span className="px-3 py-1 bg-gray-200 text-gray-700 rounded-full text-sm font-semibold">
                            Semester: {notice?.semester_id || "N/A"}
                        </span>
                    </div>

                    <div className="space-y-4">
                        <p className="text-gray-700 font-semibold">Description:</p>
                        <p className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                            {notice?.description || "No Description Available"}
                        </p>

                        {notice?.file_path && (
                            <div className="mt-6">
                                <p className="text-gray-700 font-semibold mb-2">Attached File:</p>

                                {notice.file_path.endsWith(".pdf") ? (
                                    <div className="border border-gray-300 p-4 rounded-lg shadow-sm bg-gray-50">
                                        <iframe
                                            src={`http://127.0.0.1:3000/${notice.file_path}`}
                                            className="w-full h-[500px]"
                                            title="PDF Attachment"
                                        />
                                        <div className="mt-2 text-center">
                                            <a
                                                href={`http://127.0.0.1:3000/${notice.file_path}`}
                                                target="_blank"
                                                rel="noopener noreferrer"
                                                className="text-blue-500 hover:underline font-semibold"
                                            >
                                                View PDF in new tab
                                            </a>
                                        </div>
                                    </div>
                                ) : (
                                    <img
                                        src={`http://127.0.0.1:3000/${notice.file_path}`}
                                        alt="Notice Attachment"
                                        className="w-full h-auto rounded-lg border border-gray-300 shadow-sm"
                                    />
                                )}
                            </div>
                        )}

                    </div>
                </div>

                <div className="mt-8">
                    <button
                        onClick={() => navigate(-1)}
                        className="px-6 py-2 bg-blue-500 text-white font-semibold rounded-lg shadow-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-400"
                    >
                        Back to Notices
                    </button>
                </div>
            </div>
        </div>
    );
};

export default NoticeDetailPage;
