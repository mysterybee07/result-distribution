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
    <div>
      <form onSubmit={handleSubmit} encType="multipart/form-data" method="post">
        <div
          {...getRootProps()}
          style={{
            border: "2px dashed #ccc",
            borderRadius: "8px",
            padding: "20px",
            textAlign: "center",
            backgroundColor: isDragActive ? "#e0f7fa" : "#f9f9f9",
            cursor: "pointer",
            transition: "background-color 0.2s",
          }}
        >
          <input {...getInputProps()} />
          {isDragActive ? (
            <p style={{ color: "#00796b" }}>Drop the files here...</p>
          ) : (
            <p style={{ color: "#666" }}>Drag and drop files here, or click to select files</p>
          )}
        </div>

        {uploadedFiles.length > 0 && (
          <div style={{ marginTop: "20px" }}>
            <h3>Uploaded Files:</h3>
            <ul>
              {uploadedFiles.map((file, index) => (
                <li key={index}>
                  <strong>Name:</strong> {file.name} <br />
                  <strong>Size:</strong> {(file.size / 1024).toFixed(2)} KB <br />
                  <strong>Type:</strong> {file.type} <br />
                  <strong>Last Modified:</strong> {file.lastModified}
                </li>
              ))}
            </ul>
          </div>
        )}

        <Button type="submit" disabled={isLoading}>
          {isLoading ? 'Uploading...' : 'Submit'}
        </Button>
      </form>
    </div>
  );
};

export default CreateCollege;
