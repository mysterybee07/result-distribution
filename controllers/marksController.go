package controllers

import (
	"errors"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/middleware/validation"
	"github.com/mysterybee07/result-distribution-system/models"
	"gorm.io/gorm"
)

// Create a global validator
var validate = validator.New()

func AddMarks(c *fiber.Ctx) error {
	var batches []models.Batch
	if err := initializers.DB.Find(&batches).Error; err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error fetching batches")
		return err
	}

	var courses []models.Course
	if err := initializers.DB.Find(&courses).Error; err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error fetching courses")
		return err
	}

	var programs []models.Program
	if err := initializers.DB.Preload("Semesters").Find(&programs).Error; err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error fetching programs")
		return err
	}

	var semesters []models.Semester
	if err := initializers.DB.Find(&semesters).Error; err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error fetching semesters")
		return err
	}

	err := c.Render("dashboard/marks/add", fiber.Map{
		"Students":  []models.Student{}, // Empty initially
		"Courses":   courses,
		"Batches":   batches,
		"Programs":  programs,
		"Semesters": semesters,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
	}
	return nil
}

func CreateMarks(c *fiber.Ctx) error {
	// Parse incoming JSON request body into payload struct
	var payload models.MarksPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON",
		})
	}

	// Validate input using validator
	if err := validate.Struct(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate marks input
	if err := validation.ValidateMarksInput(&payload, false); err != nil {
		return c.Status(err.(*fiber.Error).Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Fetch course details
	var course models.Course
	if err := initializers.DB.First(&course, payload.CourseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Course not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Create marks for each student
	var marks []models.Mark
	for _, mark := range payload.Marks {
		// Check that obtained marks do not exceed total marks
		if mark.SemesterMarks > course.SemesterTotalMarks ||
			(course.PracticalTotalMarks != nil && mark.PracticalMarks > *course.PracticalTotalMarks) ||
			(course.AssistantTotalMarks != nil && mark.AssistantMarks > *course.AssistantTotalMarks) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Obtained marks cannot exceed total marks",
			})
		}

		newMark := models.Mark{
			BatchID:        payload.BatchID,
			ProgramID:      payload.ProgramID,
			SemesterID:     payload.SemesterID,
			CourseID:       payload.CourseID,
			StudentID:      mark.StudentID,
			SemesterMarks:  mark.SemesterMarks,
			AssistantMarks: mark.AssistantMarks,
			PracticalMarks: mark.PracticalMarks,
		}

		// Determine pass/fail status
		if newMark.SemesterMarks < course.SemesterPassMarks ||
			(course.PracticalPassMarks != nil && newMark.PracticalMarks < *course.PracticalPassMarks) ||
			(course.AssistantPassMarks != nil && newMark.AssistantMarks < *course.AssistantPassMarks) {
			newMark.Status = "failed"
		} else {
			newMark.Status = "pass"
		}

		marks = append(marks, newMark)
	}

	// Bulk insert the marks
	if err := initializers.DB.Create(&marks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create marks",
		})
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
	var input models.MarksPayload
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

func GetMarksBySymbolNumber(c *fiber.Ctx) error {
	// Get the symbol number from the request parameters
	symbolNumber := c.Params("symbolNumber")

	// Check if symbol number is provided
	if symbolNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Symbol number is required",
		})
	}

	// Find the student using the symbol number
	var student models.Student
	if err := initializers.DB.Where("symbol_number = ?", symbolNumber).First(&student).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Student not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve student",
		})
	}

	// Find the marks for the student
	var marks []models.Mark
	if err := initializers.DB.Where("student_id = ?", student.ID).
		Preload("Batch").
		Preload("Program").
		Preload("Semester").
		Preload("Course").
		Preload("Student").
		Find(&marks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve marks",
		})
	}

	totalMarks := 0
	status := "pass"
	for i := range marks {
		mark := &marks[i]

		if mark.SemesterMarks < mark.Course.SemesterPassMarks ||
			(mark.Course.AssistantPassMarks != nil && mark.AssistantMarks < *mark.Course.AssistantPassMarks) ||
			(mark.Course.PracticalPassMarks != nil && mark.PracticalMarks < *mark.Course.PracticalPassMarks) {
			status = "failed"
		}

		totalMarks += mark.TotalMarks
	}

	// Return the marks and overall status
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"student":    student,
		"marks":      marks,
		"totalMarks": totalMarks,
		"status":     status,
	})
}
