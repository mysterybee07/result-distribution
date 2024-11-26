package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/utils"
)

func PublishExamRoutine(c *fiber.Ctx) error {
	// Request payload
	type Request struct {
		BatchID    uint      `json:"batch_id"`
		ProgramID  uint      `json:"program_id"`
		SemesterID uint      `json:"semester_id"`
		StartDate  time.Time `json:"start_date"`
		EndDate    time.Time `json:"end_date"`
	}

	// Parse request body
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Call the helper function
	fileName, err := utils.ExamRoutine(req.BatchID, req.ProgramID, req.SemesterID, req.StartDate, req.EndDate)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message":  "Exam routine published successfully",
		"fileName": fileName,
	})
}
