package controllers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

// UploadFile handles file uploads
// func UploadFile(c *fiber.Ctx) error {
// 	// Get the uploaded file from the request
// 	file, err := c.FormFile("file_path")
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "No file uploaded",
// 		})
// 	}

// 	// Create a unique filename to avoid overwriting existing files
// 	filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), file.Filename)
// 	uploadDir := "uploads"

// 	// Create the upload directory if it doesn't exist
// 	err = os.MkdirAll(uploadDir, os.ModePerm)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Error creating upload directory",
// 		})
// 	}

// 	// Open a file for writing in the upload directory
// 	out, err := os.Create(filepath.Join(uploadDir, filename))
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Error creating file",
// 		})
// 	}
// 	defer out.Close()

// 	// Copy the uploaded file to the destination
// 	fileContent, err := file.Open() // File must be opened for reading
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Error reading uploaded file",
// 		})
// 	}
// 	defer fileContent.Close()

// 	_, err = io.Copy(out, fileContent)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Error saving file",
// 		})
// 	}

// 	// Save the uploaded file metadata to the database
// 	uploadedNotice := models.Notice{
// 		FilePath: filepath.Join(uploadDir, filename), // Save the complete file path
// 	}

// 	if err := initializers.DB.Create(&uploadedNotice).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Error saving file metadata to database",
// 		})
// 	}

// 	// Return a success message
// 	return c.JSON(fiber.Map{
// 		"message":   "File uploaded successfully",
// 		"file_path": uploadedNotice.FilePath, // Include file path in the response
// 	})
// }

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
