package controllers

import (
	"errors"
	"log"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
	"github.com/mysterybee07/result-distribution-system/utils"
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
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "All compulsory courses and exactly one optional course are required to be marked for each student",
			})
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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update semester",
			})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save result",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Results published and semesters updated",
	})
}

// To get only pass students by semester and rank them by total marks they obtained in all subject
func PassingStudentsBySemester(c *fiber.Ctx) error {
	programID := c.Query("program_id")
	batchID := c.Query("batch_id")
	semesterID := c.Query("semester_id")

	// Retrieve students based on program, batch, and semester
	var students []models.Student
	query := initializers.DB.Preload("Program")

	// Apply filters based on query parameters
	if programID != "" {
		query = query.Where("program_id = ?", programID)
	}
	if batchID != "" {
		query = query.Where("batch_id = ?", batchID)
	}
	if semesterID != "" {
		query = query.Where("semester_id = ?", semesterID)
	}

	if err := query.Find(&students).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching students")
	}

	// Get pass status for the semester
	passStatus, err := utils.GetPassStatusBySemester(semesterID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching pass status")
	}

	// Prepare data to pass to the response
	var studentsData []struct {
		models.Student
		TotalMarks int
		PassStatus string
		Rank       int
	}

	// Calculate total marks for each student
	studentMarks := make(map[uint]int)

	var marks []models.Mark
	if err := initializers.DB.Find(&marks).Error; err != nil {
		log.Printf("Failed to fetch marks: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching marks")
	}

	for _, mark := range marks {
		studentMarks[mark.StudentID] += mark.TotalMarks
	}

	// Populate the studentsData slice
	for _, student := range students {
		totalMarks := studentMarks[student.ID]
		status := passStatus[student.ID]
		if status == "fail" {
			continue // Skip students with a "fail" status
		}
		studentsData = append(studentsData, struct {
			models.Student
			TotalMarks int
			PassStatus string
			Rank       int
		}{
			Student:    student,
			TotalMarks: totalMarks,
			PassStatus: status,
		})
	}

	// Sort studentsData by TotalMarks in descending order
	sort.Slice(studentsData, func(i, j int) bool {
		return studentsData[i].TotalMarks > studentsData[j].TotalMarks
	})

	// Assign ranks
	rank := 1
	previousMarks := -1
	for i, student := range studentsData {
		if previousMarks == -1 || student.TotalMarks != previousMarks {
			rank = i + 1
		}
		studentsData[i].Rank = rank
		previousMarks = student.TotalMarks
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Rank": studentsData,
	})
}

func FailedStudentsByCourse(c *fiber.Ctx) error {
	batchID := c.Query("batch_id")
	programID := c.Query("program_id")
	semesterID := c.Query("semester_id")
	courseID := c.Query("course_id")

	// Get the pass status for the semester
	passStatus, err := utils.GetPassStatusBySemester(semesterID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching pass status",
		})
	}

	// Prepare the main query with joins and necessary conditions
	var failedStudents []struct {
		models.Student
		BatchID    uint `json:"batch_id"`
		ProgramID  uint `json:"program_id"`
		SemesterID uint `json:"semester_id"`
		CourseID   uint `json:"course_id"`
	}

	// Query for students who have marks and join with the marks table
	query := initializers.DB.
		Table("students").
		Select("students.*, marks.batch_id, marks.program_id, marks.semester_id, marks.course_id").
		Joins("JOIN marks ON students.id = marks.student_id").
		Where("marks.course_id = ?", courseID)

	// Add filters based on query parameters
	if batchID != "" {
		query = query.Where("marks.batch_id = ?", batchID)
	}
	if programID != "" {
		query = query.Where("marks.program_id = ?", programID)
	}
	if semesterID != "" {
		query = query.Where("marks.semester_id = ?", semesterID)
	}

	// Execute the query and retrieve students who have marks for the specific course
	if err := query.Find(&failedStudents).Error; err != nil {
		log.Printf("Failed to fetch failed students: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching failed students",
		})
	}

	// Filter out students who have failed in the specified course
	var result []struct {
		models.Student
		BatchID    uint `json:"batch_id"`
		ProgramID  uint `json:"program_id"`
		SemesterID uint `json:"semester_id"`
		CourseID   uint `json:"course_id"`
	}

	for _, student := range failedStudents {
		if status, ok := passStatus[student.ID]; ok && status == "fail" {
			result = append(result, student)
		}
	}

	// Return the failed students data in the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"FailedStudents": result,
	})
}

// func PassingStudentsBySemester(c *fiber.Ctx) error {
// 	batchID := c.Query("batch_id")
// 	programID := c.Query("program_id")
// 	semesterID := c.Query("semester_id")

// 	// Get the pass status for the semester
// 	passStatus, err := utils.GetPassStatusBySemester(semesterID)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Error fetching pass status",
// 		})
// 	}

// 	// Prepare the main query with necessary conditions
// 	var passingStudents []struct {
// 		models.Student
// 		BatchID    uint `json:"batch_id"`
// 		ProgramID  uint `json:"program_id"`
// 		SemesterID uint `json:"semester_id"`
// 	}

// 	query := initializers.DB.
// 		Table("students").
// 		Select("students.*, marks.batch_id, marks.program_id, marks.semester_id").
// 		Joins("JOIN marks ON students.id = marks.student_id").
// 		Where("marks.status = ?", "pass")

// 	// Add filters based on query parameters
// 	if semesterID != "" {
// 		query = query.Where("marks.semester_id = ?", semesterID)
// 	}
// 	if batchID != "" {
// 		query = query.Where("marks.batch_id = ?", batchID)
// 	}
// 	if programID != "" {
// 		query = query.Where("marks.program_id = ?", programID)
// 	}

// 	// Execute the query and retrieve students
// 	if err := query.Find(&passingStudents).Error; err != nil {
// 		log.Printf("Failed to fetch passing students: %v\n", err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Error fetching passing students",
// 		})
// 	}

// 	// Filter out students who failed
// 	var result []struct {
// 		models.Student
// 		BatchID    uint `json:"batch_id"`
// 		ProgramID  uint `json:"program_id"`
// 		SemesterID uint `json:"semester_id"`
// 	}

// 	for _, student := range passingStudents {
// 		if status, ok := passStatus[student.ID]; ok && status == "pass" {
// 			result = append(result, student)
// 		}
// 	}

// 	// Return the passing students data in the response
// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"PassingStudents": result,
// 	})
// }
