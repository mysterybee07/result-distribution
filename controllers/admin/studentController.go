package controllers

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/middleware/validation"
	"github.com/mysterybee07/result-distribution-system/models"
)

func Student(c *fiber.Ctx) error {
	var programs []models.Program
	if err := initializers.DB.Find(&programs).Error; err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error fetching programs")
		return err
	}
	var batches []models.Batch
	if err := initializers.DB.Find(&batches).Error; err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error fetching batches")
		return err
	}
	err := c.Render("dashboard/students/addstudent", fiber.Map{
		"Programs": programs,
		"Batches":  batches,
	})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}

func CreateStudents(c *fiber.Ctx) error {
	// Check if a file is uploaded
	file, err := c.FormFile("file")
	if err == nil {
		// Parse TSV file
		tsvFile, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not open uploaded file",
			})
		}
		defer tsvFile.Close()

		// BatchID and ProgramID might be required as query params or passed another way
		batchID, _ := strconv.ParseUint(c.FormValue("batch_id"), 10, 32)
		programID, _ := strconv.ParseUint(c.FormValue("program_id"), 10, 32)

		var students []models.Student
		scanner := bufio.NewScanner(tsvFile)
		firstLine := true

		for scanner.Scan() {
			if firstLine {
				firstLine = false // Skip the header row
				continue
			}

			line := scanner.Text()
			fields := strings.Split(line, "\t")

			if len(fields) < 4 { // Ensure minimum required columns (fullname, symbol number, registration number, college_id)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid TSV format",
				})
			}

			symbolNumber := fields[0]
			registrationNumber := fields[1]
			fullname := fields[2]
			collegeIdentifier := fields[3] // Can be either an ID or name

			var collegeID uint
			if id, err := strconv.ParseUint(collegeIdentifier, 10, 32); err == nil {
				collegeID = uint(id) // CollegeID as a number
			} else {
				// CollegeID as a string (college name)
				var college models.College
				if err := initializers.DB.Where("college_name = ?", collegeIdentifier).First(&college).Error; err != nil {
					return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
						"error": fmt.Sprintf("College not found for name: %s", collegeIdentifier),
					})
				}
				collegeID = college.ID
			}

			student := models.Student{
				SymbolNumber:       symbolNumber,
				RegistrationNumber: registrationNumber,
				Fullname:           fullname,
				BatchID:            uint(batchID),
				ProgramID:          uint(programID),
				CollegeID:          collegeID,
			}

			if err := validation.ValidateStudent(&student, false); err != nil {
				return c.Status(err.(*fiber.Error).Code).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			students = append(students, student)
		}

		// Bulk insert students
		if err := initializers.DB.Create(&students).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not add students",
			})
		}

		return c.JSON(fiber.Map{
			"message":  "Students added successfully from TSV file",
			"students": students,
		})
	}

	// JSON input parsing
	var input struct {
		BatchID   uint `json:"batch_id"`
		ProgramID uint `json:"program_id"`
		Students  []struct {
			Fullname           string      `json:"fullname"`
			SymbolNumber       string      `json:"symbol_number"`
			RegistrationNumber string      `json:"registration_number"`
			CollegeID          interface{} `json:"college_id"` // Use interface{} to allow for both uint and string
		} `json:"students"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to parse request body",
		})
	}

	var students []models.Student
	for _, s := range input.Students {
		var collegeID uint

		switch v := s.CollegeID.(type) {
		case float64: // CollegeID as a number (ID)
			collegeID = uint(v)
		case string: // CollegeID as a string (college name)
			var college models.College
			if err := initializers.DB.Where("college_name = ?", v).First(&college).Error; err != nil {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": fmt.Sprintf("College not found for name: %s", v),
				})
			}
			collegeID = college.ID
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid college identifier",
			})
		}

		student := models.Student{
			SymbolNumber:       s.SymbolNumber,
			RegistrationNumber: s.RegistrationNumber,
			Fullname:           s.Fullname,
			BatchID:            input.BatchID,
			ProgramID:          input.ProgramID,
			CollegeID:          collegeID,
		}

		if err := validation.ValidateStudent(&student, false); err != nil {
			return c.Status(err.(*fiber.Error).Code).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		students = append(students, student)
	}

	// Bulk insert students
	if err := initializers.DB.Create(&students).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not add students",
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Students added successfully",
		"students": students,
	})
}

func EditStudent(c *fiber.Ctx) error {
	err := c.Render("Students/edit", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}

func UpdateStudent(c *fiber.Ctx) error {
	id := c.Params("id") // Extract student ID from the URL parameter

	var input struct {
		Fullname           string `json:"fullname"`
		SymbolNumber       string `json:"symbol_number"`
		RegistrationNumber string `json:"registration_number"`
		BatchID            uint   `json:"batch_id"`
		ProgramID          uint   `json:"program_id"`
		CollegeID          uint   `json:"college_id"`
	}

	// Parse the JSON body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to parse request body",
		})
	}

	// Find the student by ID
	var student models.Student
	if err := initializers.DB.First(&student, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Student not found",
		})
	}

	// Update the fields
	student.Fullname = input.Fullname
	student.SymbolNumber = input.SymbolNumber
	student.RegistrationNumber = input.RegistrationNumber
	student.BatchID = input.BatchID
	student.ProgramID = input.ProgramID
	student.CollegeID = input.CollegeID

	// Validate updated student data
	if err := validation.ValidateStudent(&student, true); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Save the updated student data
	if err := initializers.DB.Save(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not update student",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Student updated successfully",
		"student": student,
	})
}

func GetStudents(c *fiber.Ctx) error {
	var students []models.Student

	// Use Preload to include Batch and Program associations
	result := initializers.DB.Preload("Batch").Preload("Program").Find(&students)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve students",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Students retrieved successfully",
		"students": students,
	})
}

func GetStudentById(c *fiber.Ctx) error {
	id := c.Params("id")
	var student models.Student
	if err := initializers.DB.Preload("Batch").Preload("Program").Preload("Semester").First(&student, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Student not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Student retrieved successfully",
		"student": student,
	})
}

func GetFilteredStudents(c *fiber.Ctx) error {
	batchID := c.Query("batch_id")
	programID := c.Query("program_id")
	semesterID := c.Query("semester_id")

	var students []models.Student
	if err := initializers.DB.Preload("Batch").Preload("Program").Preload("Semester").
		Where("batch_id = ? AND program_id = ? AND current_semester = ?", batchID, programID, semesterID).
		Find(&students).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching students",
		})
	}

	return c.JSON(fiber.Map{
		"students": students,
	})
}

func DeleteStudent(c *fiber.Ctx) error {
	id := c.Params("id")
	var student models.Student
	if err := initializers.DB.First(&student, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Student not found",
		})
	}
	if err := initializers.DB.Delete(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not delete student",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Student deleted successfully",
		"student": student,
	})
}
