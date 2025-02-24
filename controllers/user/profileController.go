package controllers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

func GetUserProfile(c *fiber.Ctx) error {
	// Get userID from context
	userID, ok := c.Locals("userID").(string)
	if !ok {
		log.Println("userID not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Convert userID to int
	id, err := strconv.Atoi(userID)
	if err != nil {
		log.Printf("Failed to convert userID to integer: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error"})
	}

	// Fetch user from database
	var user models.User
	if err := initializers.DB.First(&user, id).Error; err != nil {
		log.Printf("User not found: %v\n", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Fetch student information
	var student models.Student
	if err := initializers.DB.Where("symbol_number = ?", user.SymbolNumber).
		Preload("Batch").
		Preload("Program").
		First(&student).Error; err != nil {
		log.Printf("Student not found for user with email %s: %v\n", user.Email, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Student not found"})
	}

	// Determine the previous semester
	previousSemester := student.CurrentSemester - 1

	var marks []models.Mark
	var passStatus string = "pass"
	totalMarksObtained := 0
	totalFullMarks := 0

	// Loop through previous semesters until marks are found
	for previousSemester > 0 {
		// Fetch marks for the previous semester from the results table
		if err := initializers.DB.
			Joins("JOIN results ON results.batch_id = ? AND results.program_id = ? AND results.semester_id = ? AND results.status = 'Published'", student.BatchID, student.ProgramID, previousSemester).
			Where("marks.student_id = ?", student.ID).
			Preload("Course").
			Find(&marks).Error; err != nil {
			log.Printf("Failed to find marks for student with ID %d in semester %d: %v\n", student.ID, previousSemester, err)
		}

		// If marks are found, process them
		if len(marks) > 0 {
			// Determine pass status
			for _, mark := range marks {
				if mark.Status != "pass" {
					passStatus = "fail"
					break
				}
			}

			// Calculate total obtained marks
			for _, mark := range marks {
				totalMarksObtained += mark.TotalMarks
			}

			// Calculate total full marks
			for _, mark := range marks {
				if mark.Course.SemesterTotalMarks != 0 {
					totalFullMarks += mark.Course.SemesterTotalMarks
				}
				if mark.Course.AssistantTotalMarks != nil {
					totalFullMarks += *mark.Course.AssistantTotalMarks
				}
				if mark.Course.PracticalTotalMarks != nil {
					totalFullMarks += *mark.Course.PracticalTotalMarks
				}
			}
			break
		}

		// Move to an earlier semester
		previousSemester--
	}

	// If no marks were found, return a response without marks data
	if len(marks) == 0 {
		return c.JSON(fiber.Map{
			"Users":    user,
			"Students": student,
			"message":  "No marks found for previous semesters",
		})
	}

	// Return user profile with marks data
	return c.JSON(fiber.Map{
		"Users":              user,
		"Students":           student,
		"Marks":              marks,
		"PassStatus":         passStatus,
		"TotalMarksObtained": totalMarksObtained,
		"TotalFullMarks":     totalFullMarks,
	})
}
