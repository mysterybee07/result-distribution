package validation

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
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
		if err := initializers.DB.Where("registration = ? AND batch_id = ? AND program_id = ?", student.Registration, student.BatchID, student.ProgramID).First(&existingStudent).Error; err == nil {
			log.Println("Registration Number already taken:", student.Registration)
			return fiber.NewError(fiber.StatusBadRequest, "Registration Number is already taken for the specified batch and program")
		}
	} else {
		// In update mode, ensure unique constraints exclude the current student
		if err := initializers.DB.Where("id <> ? AND symbol_number = ? AND batch_id = ? AND program_id = ?", student.ID, student.SymbolNumber, student.BatchID, student.ProgramID).First(&existingStudent).Error; err == nil {
			log.Println("Symbol Number already taken:", student.SymbolNumber)
			return fiber.NewError(fiber.StatusBadRequest, "Symbol Number is already taken for the specified batch and program")
		}
		if err := initializers.DB.Where("id <> ? AND registration = ? AND batch_id = ? AND program_id = ?", student.ID, student.Registration, student.BatchID, student.ProgramID).First(&existingStudent).Error; err == nil {
			log.Println("Registration Number already taken:", student.Registration)
			return fiber.NewError(fiber.StatusBadRequest, "Registration Number is already taken for the specified batch and program")
		}
	}

	return nil
}
