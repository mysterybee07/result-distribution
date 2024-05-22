package controllers

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

// Create a global validator
var validate = validator.New()

type CreateMarkInput struct {
	BatchID    uint `json:"batch_id" validate:"required"`
	ProgramID  uint `json:"program_id" validate:"required"`
	SemesterID uint `json:"semester_id" validate:"required"`
	CourseID   uint `json:"course_id" validate:"required"`
	Marks      []struct {
		StudentID      uint `json:"student_id" validate:"required"`
		SemesterMarks  int  `json:"semester_marks" validate:"required"`
		AssistantMarks int  `json:"assistant_marks" validate:"required"`
		PracticalMarks int  `json:"practical_marks" validate:"required"`
	} `json:"marks" validate:"required,dive"`
}

func CreateMarks(c *fiber.Ctx) error {
	// db := initializers.DB // Use the global database connection

	var input CreateMarkInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Validate input using validator
	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check if the program exists
	var program models.Program
	if err := initializers.DB.First(&program, input.ProgramID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Program not found",
		})
	}

	// Check if the semester exists
	var semester models.Semester
	if err := initializers.DB.First(&semester, input.SemesterID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Semester not found",
		})
	}

	// Check if the Course exists for the given program and semester
	var existingCourse models.Course
	if err := initializers.DB.Where("id = ? AND program_id = ? AND semester_id = ?", input.CourseID, input.ProgramID, input.SemesterID).First(&existingCourse).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Course not found for the given program and semester",
		})
	}

	// Create marks for each student
	var marks []models.Mark
	for _, markEntry := range input.Marks {
		var existingMark models.Mark
		err := initializers.DB.Where("batch_id = ? AND program_id = ? AND semester_id = ? AND course_id = ? AND student_id = ?",
			input.BatchID, input.ProgramID, input.SemesterID, input.CourseID, markEntry.StudentID).First(&existingMark).Error
		if err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Mark entry already exists for the student",
			})
		}
		mark := models.Mark{
			BatchID:        input.BatchID,
			ProgramID:      input.ProgramID,
			SemesterID:     input.SemesterID,
			CourseID:       input.CourseID,
			StudentID:      markEntry.StudentID,
			SemesterMarks:  markEntry.SemesterMarks,
			AssistantMarks: markEntry.AssistantMarks,
			PracticalMarks: markEntry.PracticalMarks,
		}
		if mark.SemesterMarks < 24 && mark.AssistantMarks < 8 && mark.PracticalMarks < 8 {
			mark.Status = "failed"
		} else {
			mark.Status = "pass"
		}
		marks = append(marks, mark)
	}

	// Bulk insert the marks
	if err := initializers.DB.Create(&marks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create marks"})
	}

	// Preload associations for each mark
	for i := range marks {
		if err := initializers.DB.Preload("Batch").
			Preload("Program").
			Preload("Semester").
			Preload("Course").
			Preload("Student").
			First(&marks[i], marks[i].ID).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not preload associations",
			})
		}
	}

	// Return success message
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Marks created successfully",
		"marks":   marks,
	})
}
