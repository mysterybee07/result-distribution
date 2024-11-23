package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
	"github.com/mysterybee07/result-distribution-system/utils"
)

func CreateNotice(c *fiber.Ctx) error {
	// Upload the file and get the file path
	filePath, err := utils.UploadFile(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error uploading file: " + err.Error(),
		})
	}

	// Create a new notice instance
	var noticeInput models.NoticeInput
	if err := c.BodyParser(&noticeInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request payload",
		})
	}

	// Convert to the Notice model
	notice := models.Notice{
		Title:       noticeInput.Title,
		Description: noticeInput.Description,
		ProgramID:   noticeInput.ProgramID,
		BatchID:     noticeInput.BatchID,
		SemesterID:  noticeInput.SemesterID,
		FilePath:    filePath,
	}

	// Save the notice to the database
	if err := initializers.DB.Create(&notice).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error saving notice to database",
			"error":   err.Error(),
		})
	}

	// Return a success message
	return c.JSON(fiber.Map{
		"message": "Notice created successfully",
		"notice":  notice,
	})
}

func UpdateNotice(c *fiber.Ctx) error {
	id := c.Params("id")
	var notice models.Notice

	// Find the notice by ID
	if err := initializers.DB.First(&notice, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Notice not found",
		})
	}

	// Parse the new data
	var noticeInput models.NoticeInput
	if err := c.BodyParser(&noticeInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request payload",
		})
	}

	// Update the notice fields
	notice.Title = noticeInput.Title
	notice.Description = noticeInput.Description
	notice.ProgramID = noticeInput.ProgramID
	notice.BatchID = noticeInput.BatchID
	notice.SemesterID = noticeInput.SemesterID

	// Handle image upload using the UpdateImage function
	newFilePath, err := utils.UpdateFile(c, notice.FilePath)
	if err != nil && err.Error() != "nofile found in the form data" {
		// If there was an error other than no image being uploaded, return the error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update file: " + err.Error(),
		})
	}

	// If a new image was uploaded, update the ImageURL field
	if newFilePath != "" {
		notice.FilePath = newFilePath
	}

	// Save the updated notice
	if err := initializers.DB.Save(&notice).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating notice",
			"error":   err.Error(),
		})
	}

	// Return the updated notice
	return c.JSON(fiber.Map{
		"message": "Notice updated successfully",
		"notice":  notice,
	})
}

func GetAllNotices(c *fiber.Ctx) error {
	var allNotices []models.Notice

	// Retrieve all notices from the database with associated names
	if err := initializers.DB.Preload("Program").Preload("Batch").Preload("Semester").Find(&allNotices).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving notices",
			"error":   err.Error(),
		})
	}

	// Transform notices to include names instead of IDs
	var notices []map[string]interface{}
	for _, notice := range allNotices {
		notices = append(notices, map[string]interface{}{
			"ID":          notice.ID,
			"Title":       notice.Title,
			"Description": notice.Description,
			"Program":     notice.Program.ProgramName,
			"Batch":       notice.Batch.Batch,
			"Semester":    notice.Semester.SemesterName,
			"FilePath":    notice.FilePath,
			"Created_at":  notice.CreatedAt,
		})
	}

	// Return the transformed list of notices
	return c.JSON(fiber.Map{
		"notices": notices,
	})
}

func GetNoticeById(c *fiber.Ctx) error {
	id := c.Params("id")
	var notice models.Notice

	// Find the notice by ID
	if err := initializers.DB.First(&notice, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Notice not found",
		})
	}

	if err := initializers.DB.Where("id = ?", id).Find(&notice).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving notices",
			"error":   err.Error(),
		})
	}

	// Return a success message
	return c.JSON(fiber.Map{
		"notice":  notice,
		"message": "Notice deleted successfully",
	})
}

func DeleteNotice(c *fiber.Ctx) error {
	id := c.Params("id")
	var notice models.Notice

	// Find the notice by ID
	if err := initializers.DB.First(&notice, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Notice not found",
		})
	}

	// Delete the notice from the database
	if err := initializers.DB.Delete(&notice).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting notice",
			"error":   err.Error(),
		})
	}

	// Return a success message
	return c.JSON(fiber.Map{
		"message": "Notice deleted successfully",
	})
}

func GetNoticesByProgram(c *fiber.Ctx) error {
	var notices []models.Notice

	// Get the program ID from the query parameter
	programID := c.Query("program_id")

	// Check if the program_id is provided
	if programID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Program ID is required",
		})
	}

	// Filter notices by ProgramID
	if err := initializers.DB.Where("program_id = ?", programID).Find(&notices).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving notices",
			"error":   err.Error(),
		})
	}

	// Return the filtered notices
	return c.JSON(fiber.Map{
		"notices": notices,
	})
}

func GetNoticesByProgramAndBatch(c *fiber.Ctx) error {
	var notices []models.Notice

	// Get the program ID and batch ID from the query parameters
	programID := c.Query("program_id")
	batchID := c.Query("batch_id")

	// Check if both program_id and batch_id are provided
	if programID == "" || batchID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Both Program ID and Batch ID are required",
		})
	}

	// Filter notices by ProgramID and BatchID
	if err := initializers.DB.Where("program_id = ? AND batch_id = ?", programID, batchID).Find(&notices).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving notices",
			"error":   err.Error(),
		})
	}

	// Return the filtered notices
	return c.JSON(fiber.Map{
		"notices": notices,
	})
}
