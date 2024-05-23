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

		// Check that obtained marks do not exceed total marks
		if markEntry.SemesterMarks > course.SemesterTotalMarks ||
			(course.PracticalTotalMarks != nil && markEntry.PracticalMarks > *course.PracticalTotalMarks) ||
			(course.AssistantTotalMarks != nil && markEntry.AssistantMarks > *course.AssistantTotalMarks) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Obtained marks cannot exceed total marks",
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
	if err := validation.ValidateMarksInput(&input, true); err != nil {
		return c.Status(err.(*fiber.Error).Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Fetch the course details to get the pass marks and total marks
	var course models.Course
	if err := initializers.DB.First(&course, input.CourseID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not find the course",
		})
	}

	// Create a slice to store the updated marks
	var updatedMarks []models.Mark
	for _, markEntry := range input.Marks {
		// Check that obtained marks do not exceed total marks
		if markEntry.SemesterMarks > course.SemesterTotalMarks ||
			(course.PracticalTotalMarks != nil && markEntry.PracticalMarks > *course.PracticalTotalMarks) ||
			(course.AssistantTotalMarks != nil && markEntry.AssistantMarks > *course.AssistantTotalMarks) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Obtained marks cannot exceed total marks",
			})
		}

		// Find the existing mark record
		var mark models.Mark
		if err := initializers.DB.Where("student_id = ? AND course_id = ?", markEntry.StudentID, input.CourseID).First(&mark).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not find the mark record",
			})
		}

		// Update the mark fields
		mark.SemesterMarks = markEntry.SemesterMarks
		mark.AssistantMarks = markEntry.AssistantMarks
		mark.PracticalMarks = markEntry.PracticalMarks

		// Check pass/fail status based on the course's pass marks
		if mark.SemesterMarks < course.SemesterPassMarks ||
			(course.PracticalPassMarks != nil && mark.PracticalMarks < *course.PracticalPassMarks) ||
			(course.AssistantPassMarks != nil && mark.AssistantMarks < *course.AssistantPassMarks) {
			mark.Status = "failed"
		} else {
			mark.Status = "pass"
		}

		// Save the updated mark
		if err := initializers.DB.Save(&mark).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not update the mark",
			})
		}

		updatedMarks = append(updatedMarks, mark)
	}

	// Preload associations for each updated mark
	for i := range updatedMarks {
		if err := initializers.DB.Preload("Batch").
			Preload("Program").
			Preload("Semester").
			Preload("Course").
			Preload("Student").
			First(&updatedMarks[i], updatedMarks[i].ID).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not preload associations",
			})
		}
	}

	// Return success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Marks updated successfully",
		"marks":   updatedMarks,
	})
}
