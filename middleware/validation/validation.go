package validation

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
	"github.com/mysterybee07/result-distribution-system/utils"
)

func ValidateStudent(student *models.Student, isUpdate bool) error {
	// Check if Batch and Program exist
	var batch models.Batch
	var program models.Program

	if err := initializers.DB.First(&batch, student.BatchID).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Batch not found")
	}

	if err := initializers.DB.First(&program, student.ProgramID).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Program not found")
	}

	// Check if Symbol Number and Registration Number already exists
	var existingStudent models.Student
	if !isUpdate {
		if err := initializers.DB.Where("symbol_number = ? AND batch_id = ? AND program_id = ?", student.SymbolNumber, student.BatchID, student.ProgramID).First(&existingStudent).Error; err == nil {
			log.Println("Symbol Number already taken:", student.SymbolNumber)
			return fiber.NewError(fiber.StatusBadRequest, "Symbol Number is already taken for the specified batch and program")
		}
		if err := initializers.DB.Where("registration_number = ? AND batch_id = ? AND program_id = ?", student.RegistrationNumber, student.BatchID, student.ProgramID).First(&existingStudent).Error; err == nil {
			log.Println("Registration Number already taken:", student.RegistrationNumber)
			return fiber.NewError(fiber.StatusBadRequest, "Registration Number is already taken for the specified batch and program")
		}
	} else {
		// In update mode, ensure unique constraints exclude the current student
		if err := initializers.DB.Where("id <> ? AND symbol_number = ? AND batch_id = ? AND program_id = ?", student.ID, student.SymbolNumber, student.BatchID, student.ProgramID).First(&existingStudent).Error; err == nil {
			log.Println("Symbol Number already taken:", student.SymbolNumber)
			return fiber.NewError(fiber.StatusBadRequest, "Symbol Number is already taken for the specified batch and program")
		}
		if err := initializers.DB.Where("id <> ? AND registration_number = ? AND batch_id = ? AND program_id = ?", student.ID, student.RegistrationNumber, student.BatchID, student.ProgramID).First(&existingStudent).Error; err == nil {
			log.Println("Registration Number already taken:", student.RegistrationNumber)
			return fiber.NewError(fiber.StatusBadRequest, "Registration Number is already taken for the specified batch and program")
		}
	}

	return nil
}

// VAlidation for creating marks

func ValidateMarksInput(input *models.MarksPayload, isUpdate bool) error {
	// Check if the program exists
	var program models.Program
	if err := initializers.DB.First(&program, input.ProgramID).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Program not found")
	}

	// Check if the semester exists
	var semester models.Semester
	if err := initializers.DB.First(&semester, input.SemesterID).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Semester not found")
	}

	// Check if the course exists for the given program and semester
	var course models.Course
	if err := initializers.DB.Where("id = ? AND program_id = ? AND semester_id = ?", input.CourseID, input.ProgramID, input.SemesterID).First(&course).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Course not found for the given batch, program and semester")
	}

	// Ensure no duplicate mark entries for students
	for _, markEntry := range input.Marks {
		var existingMark models.Mark
		err := initializers.DB.Where("batch_id = ? AND program_id = ? AND semester_id = ? AND course_id = ? AND student_id = ?",
			input.BatchID, input.ProgramID, input.SemesterID, input.CourseID, markEntry.StudentID).First(&existingMark).Error

		if isUpdate {
			// For updates, ensure that the mark entry exists
			if err != nil {
				return fiber.NewError(fiber.StatusNotFound, "Mark entry not found for the student")
			}
		} else {
			// For creation, ensure that the mark entry does not already exist
			if err == nil {
				return fiber.NewError(fiber.StatusConflict, "Mark entry already exists for the student")
			}
		}
	}

	return nil
}

func ValidateUser(data *models.User, isUpdate bool) error {
	// Validate email format
	if !utils.ValidateEmail(data.Email) {
		log.Println("Invalid email format:", data.Email)
		return errors.New("invalid email format")
	}

	// Validate password length
	if !isUpdate || (isUpdate && len(data.Password) > 0) {
		if len(data.Password) < 8 {
			log.Println("Password too short:", data.Password)
			return errors.New("password must be at least 8 characters long")
		}
	}

	if isUpdate {
		// Logging for update validation
		log.Println("Update detected. Skipping batch, program, symbol checks.")
		if data.Email == "" && data.Password == "" && data.ImageURL == "" {
			log.Println("No update fields provided")
			return errors.New("at least one of email, password, or image URL must be provided for update")
		}
		var existingUser models.User
		if err := initializers.DB.Where("email = ?", data.Email).First(&existingUser).Error; err == nil && existingUser.ID != data.ID {
			log.Println("Email is already taken:", data.Email)
			return errors.New("email is already taken")
		}
		return nil
	}

	// New user validation
	if data.BatchID == nil || data.ProgramID == nil || data.SymbolNumber == "" || data.RegistrationNumber == "" || data.Email == "" || data.Password == "" {
		log.Println("Missing fields - BatchID:", data.BatchID, "ProgramID:", data.ProgramID, "SymbolNumber:", data.SymbolNumber, "RegistrationNumber:", data.RegistrationNumber)
		return errors.New("batch ID, program ID, symbol, registration, email, and password are required for regular users")
	}

	// Check if symbol and registration exist in the students table for the given batch and program
	var student models.Student
	if err := initializers.DB.Where("symbol_number = ? AND registration_number = ? AND batch_id = ? AND program_id = ?", data.SymbolNumber, data.RegistrationNumber, *data.BatchID, *data.ProgramID).First(&student).Error; err != nil {
		log.Println("Invalid symbol/registration for batch and program:", data.SymbolNumber, data.RegistrationNumber)
		return errors.New("invalid symbol or registration for the specified batch and program")
	}

	// Check for unique email
	var existingUser models.User
	if err := initializers.DB.Where("email = ?", data.Email).First(&existingUser).Error; err == nil && existingUser.ID != data.ID {
		log.Println("Email already taken:", data.Email)
		return errors.New("email is already taken")
	}

	// Check if symbol number is unique for the batch and program
	if err := initializers.DB.Where("symbol_number = ? AND batch_id = ? AND program_id = ?", data.SymbolNumber, data.BatchID, data.ProgramID).First(&existingUser).Error; err == nil && existingUser.ID != data.ID {
		log.Println("Symbol number already taken:", data.SymbolNumber)
		return errors.New("symbol number is already taken for the specified batch and program")
	}

	// Check if registration number is unique for the batch and program
	if err := initializers.DB.Where("registration_number = ? AND batch_id = ? AND program_id = ?", data.RegistrationNumber, data.BatchID, data.ProgramID).First(&existingUser).Error; err == nil && existingUser.ID != data.ID {
		log.Println("Registration number already taken:", data.RegistrationNumber)
		return errors.New("registration number is already taken for the specified batch and program")
	}

	return nil
}
