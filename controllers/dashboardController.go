package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

func Index(c *fiber.Ctx) error {
	var students []models.Student
	if err := initializers.DB.Preload("Program").Find(&students).Error; err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error fetching students")
		return err
	}

	var marks []models.Mark
	if err := initializers.DB.Find(&marks).Error; err != nil {
		log.Printf("Failed to fetch marks: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching marks")
	}

	// Calculate total marks for each student
	studentMarks := make(map[uint]int)
	passStatus := make(map[uint]string)

	for _, mark := range marks {
		studentMarks[mark.StudentID] += mark.TotalMarks
		if mark.Status != "pass" {
			passStatus[mark.StudentID] = "fail"
		} else if _, ok := passStatus[mark.StudentID]; !ok {
			passStatus[mark.StudentID] = "pass"
		}
	}

	// Prepare data to pass to the template
	var studentsData []struct {
		models.Student
		TotalMarks int
		PassStatus string
	}

	for _, student := range students {
		totalMarks := studentMarks[student.ID]
		status := passStatus[student.ID]
		studentsData = append(studentsData, struct {
			models.Student
			TotalMarks int
			PassStatus string
		}{
			Student:    student,
			TotalMarks: totalMarks,
			PassStatus: status,
		})
	}

	err := c.Render("dashboard/index", fiber.Map{
		"StudentsData": studentsData,
	})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}

// func AddStudent(c *fiber.Ctx) error {

// }
