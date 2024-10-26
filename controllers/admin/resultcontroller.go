package controllers

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
	exam "github.com/mysterybee07/result-distribution-system/models/exam"
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

	var results []exam.Result
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
	var existingResult exam.Result
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
	// Check if all students have marks for all compulsory courses and exactly one optional course in the current semester
	for _, student := range students {
		var marks []models.Mark
		if err := initializers.DB.Where("student_id = ? AND semester_id = ?", student.ID, req.SemesterID).Find(&marks).Error; err != nil {
			log.Printf("Failed to fetch marks for student %d: %v\n", student.ID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch marks"})
		}

		var courses []models.Course
		if err := initializers.DB.Where("program_id = ? AND semester_id = ?", student.ProgramID, req.SemesterID).Find(&courses).Error; err != nil {
			log.Printf("Failed to fetch courses for student %d: %v\n", student.ID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch courses"})
		}

		// Separate compulsory and optional courses
		compulsoryCourses := make(map[uint]bool)
		optionalCourses := make(map[uint]bool)

		for _, course := range courses {
			if course.IsCompulsory {
				compulsoryCourses[course.ID] = false
			} else {
				optionalCourses[course.ID] = false
			}
		}

		// Mark courses as completed based on the student's marks
		for _, mark := range marks {
			if _, exists := compulsoryCourses[mark.CourseID]; exists {
				compulsoryCourses[mark.CourseID] = true
			} else if _, exists := optionalCourses[mark.CourseID]; exists {
				optionalCourses[mark.CourseID] = true
			}
		}

		// Check if all compulsory courses are marked
		allCompulsoryMarked := true
		for _, marked := range compulsoryCourses {
			if !marked {
				allCompulsoryMarked = false
				break
			}
		}

		// Count the number of marked optional courses
		optionalMarkedCount := 0
		for _, marked := range optionalCourses {
			if marked {
				optionalMarkedCount++
			}
		}

		// Validate that all compulsory courses are marked and exactly one optional course is marked
		if !allCompulsoryMarked || optionalMarkedCount != 1 {
			log.Printf("All compulsory courses or exactly one optional course are not marked for student %d in semester %d\n", student.ID, req.SemesterID)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "All compulsory courses and exactly one optional course are required to be marked for each student"})
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
	newResult := exam.Result{
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
