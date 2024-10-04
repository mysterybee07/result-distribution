package controllers

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
	"gorm.io/gorm"
)

func Result(c *fiber.Ctx) error {

	var batches []models.Batch
	if err := initializers.DB.Find(&batches).Error; err != nil {
		log.Printf("Failed to fetch batches: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching batches")
	}

	var programs []models.Program
	if err := initializers.DB.Preload("Semesters").Find(&programs).Error; err != nil {
		log.Printf("Failed to fetch programs: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch programs")
	}

	var semesters []models.Semester
	if err := initializers.DB.Find(&semesters).Error; err != nil {
		log.Printf("Failed to fetch semesters: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch semesters")
	}

	var results []models.Result
	if err := initializers.DB.Preload("Batch").Preload("Program").Preload("Semester").Find(&results).Error; err != nil {
		log.Printf("Failed to fetch results: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch results")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Batches":   batches,
		"Programs":  programs,
		"Semesters": semesters,
		"Results":   results,
	})
}

func PublishResults(c *fiber.Ctx) error {
	// Get batch and semester from request body or query params
	type PublishRequest struct {
		BatchID    uint `json:"batch_id" form:"batch_id"`
		ProgramID  uint `json:"program_id" form:"program_id"`
		SemesterID uint `json:"semester_id" form:"semester_id"`
	}

	var req PublishRequest
	if err := c.BodyParser(&req); err != nil {
		log.Println("Unable to parse form data:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid form data"})
	}

	// Check if results are already published for the given batch, program, and semester
	var existingResult models.Result
	if err := initializers.DB.Where("batch_id = ? AND program_id = ? AND semester_id = ?", req.BatchID, req.ProgramID, req.SemesterID).First(&existingResult).Error; err == nil {
		// Result already exists, return an error
		log.Printf("Result already published for batch %d, program %d, and semester %d\n", req.BatchID, req.ProgramID, req.SemesterID)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Result already published for the given semester with batch and program"})
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Other database error
		log.Printf("Failed to fetch existing result: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check existing result"})
	}

	var students []models.Student
	if err := initializers.DB.Where("status = ? AND batch_id = ? AND program_id = ?", "active", req.BatchID, req.ProgramID).Find(&students).Error; err != nil {
		log.Printf("Failed to fetch students: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch students"})
	}

	// Check if all students have marks for all courses in the current semester
	for _, student := range students {
		var marks []models.Mark
		if err := initializers.DB.Where("student_id = ? AND semester_id = ?", student.ID, req.SemesterID).Find(&marks).Error; err != nil {
			log.Printf("Failed to fetch marks for student %d: %v\n", student.ID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch marks"})
		}

		var courses []models.Course
		if err := initializers.DB.Where("program_id = ?", student.ProgramID).Find(&courses).Error; err != nil {
			log.Printf("Failed to fetch courses for student %d: %v\n", student.ID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch courses"})
		}

		courseIDMap := make(map[uint]bool)
		for _, course := range courses {
			courseIDMap[course.ID] = false
		}

		for _, mark := range marks {
			courseIDMap[mark.CourseID] = true
		}

		allCoursesMarked := true
		for _, marked := range courseIDMap {
			if !marked {
				allCoursesMarked = false
				break
			}
		}

		if !allCoursesMarked {
			log.Printf("Not all courses are marked for student %d in semester %d\n", student.ID, req.SemesterID)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Not all courses are marked for all students"})
		}
	}

	// If all students have marks for all courses, proceed to publish results
	for _, student := range students {
		// Increment the semester
		if student.CurrentSemester < 8 {
			student.CurrentSemester++
		} else {
			student.Status = "Graduated"
		}

		if err := initializers.DB.Save(&student).Error; err != nil {
			log.Printf("Failed to update semester for student %d: %v\n", student.ID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update semester"})
		}
	}

	// Save the result in the database
	newResult := models.Result{
		BatchID:    req.BatchID,
		ProgramID:  req.ProgramID,
		SemesterID: req.SemesterID,
		Status:     "Published",
	}
	if err := initializers.DB.Create(&newResult).Error; err != nil {
		log.Printf("Failed to save result: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save result"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Results published and semesters updated",
	})
}
