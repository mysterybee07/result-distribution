package controllers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

func Profile(c *fiber.Ctx) error {
	err := c.Render("users/profile", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}
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
	if err := initializers.DB.Where("symbol_number = ?", user.SymbolNumber).Preload("Batch").Preload("Program").First(&student).Error; err != nil {
		log.Printf("Student not found for user with email %s: %v\n", user.Email, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Student not found"})
	}

	// Determine the previous semester
	previousSemester := student.CurrentSemester - 1

	// Fetch marks for the student that belong to the previous semester in the results table
	var marks []models.Mark
	if err := initializers.DB.Joins("JOIN results ON results.batch_id = ? AND results.program_id = ? AND results.semester_id = ? AND results.status = 'Published'", student.BatchID, student.ProgramID, previousSemester).
		Where("marks.student_id = ?", student.ID).
		Preload("Course").
		Find(&marks).Error; err != nil {
		log.Printf("Failed to find marks for student with ID %d: %v\n", student.ID, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Marks not found"})
	}

	// Determine the pass status
	passStatus := "pass"
	for _, mark := range marks {
		if mark.Status != "pass" {
			passStatus = "fail"
			break
		}
	}
	// Calculate totalMarksObtained marks
	totalMarksObtained := 0
	for _, mark := range marks {
		totalMarksObtained += mark.TotalMarks
	}
	// TotalFull Marks
	totalFullMarks := 0
	for _, fullmark := range marks {
		if fullmark.Course.SemesterTotalMarks != 0 {
			totalFullMarks += fullmark.Course.SemesterTotalMarks
		}
		if fullmark.Course.AssistantTotalMarks != nil {
			totalFullMarks += *fullmark.Course.AssistantTotalMarks
		}
		if fullmark.Course.PracticalTotalMarks != nil {
			totalFullMarks += *fullmark.Course.PracticalTotalMarks
		}
	}

	// Calculate totalMarksTotal marks
	// var totalMarksTotal int
	// if err := initializers.DB.Model(&models.Mark{}).Where("course_id IN (?) AND student_id =?", student.Batch.CourseIDs, student.ID).
	//     Joins("JOIN courses ON marks.course_id = courses.id").
	//     Sum(&totalMarksTotal).Error; err!= nil {
	//     log.Printf("Failed to calculate total marks: %v\n", err)
	//     return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error"})
	// }

	// // Calculate percentage
	// percentage := float64(totalMarksObtained) / float64(totalMarksTotal) * 10

	// Return user profile
	return c.Render("users/profile", fiber.Map{
		"Users":    user,
		"Students": student,
		"Marks":    marks,
		// "Courses":            courses,
		"PassStatus":         passStatus,
		"totalMarksObtained": totalMarksObtained,
		"totalFullMarks":     totalFullMarks,
	})
}
