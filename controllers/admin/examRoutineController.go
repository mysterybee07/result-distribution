package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/middleware/validation"
	"github.com/mysterybee07/result-distribution-system/models"
	"github.com/mysterybee07/result-distribution-system/utils"
)

func PublishExamRoutine(c *fiber.Ctx) error {

	// Parse request body
	var req models.ExamRoutineRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

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
