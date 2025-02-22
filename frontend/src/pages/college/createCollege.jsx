import { useMutation } from "@tanstack/react-query";
import React, { useCallback, useState } from "react";
import { useDropzone } from "react-dropzone";
import { useNavigate } from 'react-router-dom';
import { Button } from "../../components/ui/button";
import imageApi from "../../imageApi";

const CreateCollege = () => {
  const navigate = useNavigate();
  const [uploadedFiles, setUploadedFiles] = useState([]);

  const onDrop = useCallback((acceptedFiles) => {
    setUploadedFiles(acceptedFiles.map((file) => ({
      file, // Keep the actual file object
      name: file.name,
      size: file.size,
      type: file.type,
      lastModified: new Date(file.lastModified).toLocaleString(),
    })));
  }, []);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: 'text/csv', // Allow only CSV files
    maxSize: 5 * 1024 * 1024, // Limit file size to 5 MB
    onDropRejected: () => alert("Only CSV files under 5 MB are allowed."),
  });

  const { mutate: createCollege, isLoading } = useMutation({
    mutationFn: async (college) => {
      const res = await imageApi.post('/college/upload-college', college);
      return res.data;
    },
    onSuccess: () => navigate('/admin/college'),
    onError: (err) => {
      console.error('Error uploading file:', err.response?.data || err.message);
      alert(err.response?.data?.message || 'Failed to upload file. Please try again.');
    },

  });

  const handleSubmit = async (event) => {
    event.preventDefault(); // Prevent form reload
    if (uploadedFiles.length === 0) {
      alert("Please upload at least one file before submitting.");
      return;
    }
    console.log("submiting form")

    const formData = new FormData();
    uploadedFiles.forEach((file) => {
      console.log("🚀 ~ uploadedFiles.forEach ~ file:", file)
      formData.append("file", file.file);
    });
    await createCollege(formData);
  };

  return (
    <div className="w-full h-screen flex justify-center">
      <form onSubmit={handleSubmit} encType="multipart/form-data" method="post" className="w-full max-w-4xl p-4">
        <div
          {...getRootProps()}
          className={`border-2 border-dashed rounded-lg p-12 text-center cursor-pointer transition-colors w-full h-128 ${isDragActive ? "bg-teal-100" : "bg-gray-100"
            }`}
        >
          <input {...getInputProps()} />
          {isDragActive ? (
            <p className="text-teal-700">Drop the files here...</p>
          ) : (
            <p className="text-gray-600">Drag and drop files here, or click to select files</p>
          )}
        </div>

        {uploadedFiles.length > 0 && (
          <div className="mt-6">
            <h3 className="text-lg font-semibold">Uploaded Files:</h3>
            <ul className="list-disc list-inside">
              {uploadedFiles.map((file, index) => (
                <li key={index} className="mt-2">
                  <strong>Name:</strong> {file.name} <br />
                  <strong>Size:</strong> {(file.size / 1024).toFixed(2)} KB <br />
                  <strong>Type:</strong> {file.type} <br />
                  <strong>Last Modified:</strong> {new Date(file.lastModified).toLocaleDateString()}
                </li>
              ))}
            </ul>
          </div>
        )}

        <button
          type="submit"
          disabled={isLoading}
          className={`mt-4 w-full py-2 px-4 rounded-md text-white font-semibold ${isLoading ? "bg-gray-400 cursor-not-allowed" : "bg-blue-600 hover:bg-blue-700"
            } transition-colors`}
        >
          {isLoading ? "Uploading..." : "Submit"}
        </button>
      </form>
    </div>
  );
};

export default CreateCollege;
