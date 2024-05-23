package controllers

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/middleware/validation"
	"github.com/mysterybee07/result-distribution-system/models"
)

// Create a global validator
var validate = validator.New()

func CreateMarks(c *fiber.Ctx) error {
	var input validation.CreateMarkInput
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

	// Validate marks input
	if err := validation.ValidateMarksInput(&input, false); err != nil {
		return c.Status(err.(*fiber.Error).Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	//fetching the pass marks
	var course models.Course
	if err := initializers.DB.First(&course, input.CourseID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not find the course",
		})
	}

	// Create marks for each student
	var marks []models.Mark
	for _, markEntry := range input.Marks {
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
		if mark.SemesterMarks < course.SemesterPassMarks ||
			(course.PracticalPassMarks != nil && mark.PracticalMarks < *course.PracticalPassMarks) ||
			(course.AssistantPassMarks != nil && mark.AssistantMarks < *course.AssistantPassMarks) {
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

func UpdateMarks(c *fiber.Ctx) error {
	// Parse request body
	var input validation.CreateMarkInput
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

	// Find marks by CourseID
	var marks []models.Mark
	if err := initializers.DB.Where("course_id = ?", input.CourseID).Find(&marks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not find marks",
		})
	}

	// Update marks
	for _, mark := range marks {
		for _, updatedMark := range input.Marks {
			if mark.StudentID == updatedMark.StudentID {
				mark.SemesterMarks = updatedMark.SemesterMarks
				mark.AssistantMarks = updatedMark.AssistantMarks
				mark.PracticalMarks = updatedMark.PracticalMarks

				if mark.SemesterMarks < 24 && mark.AssistantMarks < 8 && mark.PracticalMarks < 8 {
					mark.Status = "failed"
				} else {
					mark.Status = "pass"
				}

				if err := initializers.DB.Save(&mark).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Could not update marks",
					})
				}
				break
			}
		}
	}

	// Return success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Marks updated successfully",
		"marks":   marks,
	})
}
