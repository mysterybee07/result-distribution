package controllers

import (
	"log"
	"sort"

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
		Rank       int
	}

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

	err := c.Render("dashboard/index", fiber.Map{
		"StudentsData": studentsData,
	})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}

func FailStudents(c *fiber.Ctx) error {
	batchID := c.Query("batch_id")
	programID := c.Query("program_id")
	semesterID := c.Query("semester_id")
	courseID := c.Query("course_id")

	var students []models.Student
	if err := initializers.DB.Preload("Program").Find(&students).Error; err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching students"})
		return err
	}

	var marks []models.Mark
	query := initializers.DB
	if batchID != "" {
		query = query.Where("batch_id = ?", batchID)
	}
	if programID != "" {
		query = query.Where("program_id = ?", programID)
	}
	if semesterID != "" {
		query = query.Where("semester_id = ?", semesterID)
	}
	if courseID != "" {
		query = query.Where("course_id = ?", courseID)
	}
	if err := query.Find(&marks).Error; err != nil {
		log.Printf("Failed to fetch marks: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching marks"})
	}

	// Filter students who have failed
	failStudentsMap := make(map[uint]bool)
	for _, mark := range marks {
		if mark.Status != "pass" {
			failStudentsMap[mark.StudentID] = true
		}
	}

	// Prepare data to pass as JSON response
	var failStudentsData []struct {
		models.Student
		BatchID    uint
		ProgramID  uint
		SemesterID uint
		CourseID   uint
	}

	for _, student := range students {
		if failStudentsMap[student.ID] {
			// Find the batch, program, semester, and course from the marks
			var studentBatchID, studentProgramID, studentSemesterID, studentCourseID uint
			for _, mark := range marks {
				if mark.StudentID == student.ID {
					studentBatchID = mark.BatchID
					studentProgramID = mark.ProgramID
					studentSemesterID = mark.SemesterID
					studentCourseID = mark.CourseID
					break
				}
			}
			failStudentsData = append(failStudentsData, struct {
				models.Student
				BatchID    uint
				ProgramID  uint
				SemesterID uint
				CourseID   uint
			}{
				Student:    student,
				BatchID:    studentBatchID,
				ProgramID:  studentProgramID,
				SemesterID: studentSemesterID,
				CourseID:   studentCourseID,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(failStudentsData)
}

// func AddStudent(c *fiber.Ctx) error {

// }
