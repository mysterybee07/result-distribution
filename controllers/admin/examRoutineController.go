package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/middleware/validation"
	"github.com/mysterybee07/result-distribution-system/models"
	"github.com/mysterybee07/result-distribution-system/utils"
)

func CreateExamRoutine(c *fiber.Ctx) error {

	// Parse request body
	var req models.ExamRoutineRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	fmt.Println(req)

	// Call the validation function
	if err := validation.ValidateExamScheduleRequest(&req); err != nil {
		return c.Status(err.(*fiber.Error).Code).JSON(fiber.Map{"error": err.(*fiber.Error).Message})
	}

	// Call the function to generate the exam routine
	fileName, examSchedules, err := utils.ExamRoutine(req.BatchID, req.ProgramID, req.SemesterID, req.StartDate, req.EndDate)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return both the CSV file name and the generated schedules
	return c.JSON(fiber.Map{
		"message":       "Exams routine published successfully",
		"fileName":      fileName,
		"examSchedules": examSchedules,
	})
}

func PublishExamRoutine(c *fiber.Ctx) error {

	id := c.Params("id")

	var examRoutine models.ExamRoutine

	if err := initializers.DB.First(&examRoutine, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"errror": "Exam Routine with id not found",
		})
	}
	examRoutine.Status = "Published"

	// Save the updated ExamRoutine
	if err := initializers.DB.Save(&examRoutine).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to update status: %v", err),
		})
	}

	// Return a success response
	return c.JSON(fiber.Map{
		"message":     "ExamRoutine status updated to Published",
		"examRoutine": examRoutine,
	})
}
