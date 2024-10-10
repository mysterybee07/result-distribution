package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// To upload file
func UploadFile(c *fiber.Ctx) (string, error) {
	// Get the uploaded file from the request
	file, err := c.FormFile("file_path")
	if err != nil {
		return "", err // Return the error directly
	}

	// Create a unique filename to avoid overwriting existing files
	filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), file.Filename)
	uploadDir := "uploads"

	// Create the upload directory if it doesn't exist
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return "", err // Return the error directly
	}

	// Open a file for writing in the upload directory
	out, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		return "", err // Return the error directly
	}
	defer out.Close()

	// Copy the uploaded file to the destination
	fileContent, err := file.Open() // File must be opened for reading
	if err != nil {
		return "", err // Return the error directly
	}
	defer fileContent.Close()

	_, err = io.Copy(out, fileContent)
	if err != nil {
		return "", err // Return the error directly
	}

	// Return the file path
	return filepath.Join(uploadDir, filename), nil
}

// Function to upload image
func UploadImage(c *fiber.Ctx) (string, error) {
	// Retrieve the multipart form data
	form, err := c.MultipartForm()
	if err != nil {
		log.Println("Failed to get multipart form data:", err)
		return "", fmt.Errorf("failed to get multipart form data: %v", err)
	}

	files := form.File["image_url"]
	if len(files) == 0 {
		log.Println("No image file found in the form data")
		return "", fmt.Errorf("no image file found in the form data")
	}

	file := files[0] // assuming only one image file is uploaded

	// Validate that the file is an image by checking its MIME type
	fileType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(fileType, "image/") {
		log.Println("Uploaded file is not an image:", fileType)
		return "", fmt.Errorf("uploaded file is not an image: %s", fileType)
	}

	// Generate a random file name and sanitize it
	fileName := RandLetter(5) + "-" + SanitizeFileName(file.Filename)
	filePath := filepath.Join("./uploads", fileName)

	// Ensure the upload directory exists
	uploadDir := filepath.Dir(filePath)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Println("Failed to create upload directory:", err)
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Save the image file
	if err := c.SaveFile(file, filePath); err != nil {
		log.Println("Failed to save image file:", err)
		return "", fmt.Errorf("failed to save image file: %v", err)
	}

	// Return the file path of the saved image
	return filePath, nil
}
