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

	// Fetch the existing ExamRoutine
	var examRoutine models.ExamRoutine
	if err := initializers.DB.First(&examRoutine, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Exam Routine with the given ID not found",
		})
	}

	// Parse the request body to get the status
	var requestBody struct {
		Status string `json:"status"`
	}

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update the status field
	examRoutine.Status = requestBody.Status

	// Save the updated ExamRoutine
	if err := initializers.DB.Save(&examRoutine).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to update status: %v", err),
		})
	}

	// Return a success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Exam Routine status updated successfully",
		"examRoutine": examRoutine,
	})
}

func ListExamsRoutine(c *fiber.Ctx) error {
	var examRoutines []models.ExamRoutine

	if err := initializers.DB.
		Preload("Batch").
		Preload("Program").
		Preload("Semester").
		Find(&examRoutines).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no exam routines available",
		})
	}

	// Transform response to include names instead of IDs
	var response []fiber.Map
	for _, routine := range examRoutines {
		response = append(response, fiber.Map{
			"start_date": routine.StartDate,
			"end_date":   routine.EndDate,
			"batch":      routine.Batch.Batch,           // Assuming Batch has a `Name` field
			"program":    routine.Program.ProgramName,   // Assuming Program has a `Name` field
			"semester":   routine.Semester.SemesterName, // Assuming Semester has a `Name` field
			"status":     routine.Status,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "routines retrieved successfully",
		"examRoutines": response,
	})
}

func ListExamSchedules(c *fiber.Ctx) error {
	var examSchedules []models.ExamSchedules

	// Preload related data to fetch names instead of IDs
	if err := initializers.DB.
		Preload("Course").
		Preload("ExamRoutine.Batch").
		Preload("ExamRoutine.Program").
		Preload("ExamRoutine.Semester").
		Find(&examSchedules).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no exam schedules available",
		})
	}

	// Transform response to include names instead of IDs
	var schedules []fiber.Map
	for _, schedule := range examSchedules {
		schedules = append(schedules, fiber.Map{
			"exam_date": schedule.ExamDate,
			"course":    schedule.Course.Name,
			"batch":     schedule.ExamRoutine.Batch.Batch,
			"program":   schedule.ExamRoutine.Program.ProgramName,
			"semester":  schedule.ExamRoutine.Semester.SemesterName,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "exam schedules retrieved successfully",
		"examSchedules": schedules,
	})
}

func GetFilteredExamSchedules(c *fiber.Ctx) error {
	var examSchedules []models.ExamSchedules

	// Extract filter parameters from the request
	batchID := c.Query("batch_id")
	programID := c.Query("program_id")
	semesterID := c.Query("semester_id")

	// Build the query with necessary filters
	query := initializers.DB.
		Preload("Course").
		Preload("ExamRoutine.Batch").
		Preload("ExamRoutine.Program").
		Preload("ExamRoutine.Semester").
		Joins("JOIN exam_routines ON exam_routines.id = exam_schedules.exam_routine_id")

	// Apply filters if provided
	if batchID != "" {
		query = query.Where("exam_routines.batch_id = ?", batchID)
	}
	if programID != "" {
		query = query.Where("exam_routines.program_id = ?", programID)
	}
	if semesterID != "" {
		query = query.Where("exam_routines.semester_id = ?", semesterID)
	}

	// Execute the query
	if err := query.Find(&examSchedules).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error retrieving exam schedules",
		})
	}

	// Check if no schedules are found
	if len(examSchedules) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no exam schedules found with the given filters",
		})
	}

	// Transform response to include names instead of IDs
	var schedulesCourse []fiber.Map
	for _, schedule := range examSchedules {
		schedulesCourse = append(schedulesCourse, fiber.Map{
			"courseCode": schedule.Course.CourseCode,
			"course":     schedule.Course.Name,
			"exam_date":  schedule.ExamDate.Format("2006-01-02"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "exam schedules retrieved successfully",
		"examSchedules": schedulesCourse,
	})
}
