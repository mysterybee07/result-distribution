package controllers

import (
	"github.com/gofiber/fiber/v2"
	fileController "github.com/mysterybee07/result-distribution-system/controllers/file"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

func CreateNotice(c *fiber.Ctx) error {
	// Upload the file and get the file path
	filePath, err := fileController.UploadFile(c)
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
