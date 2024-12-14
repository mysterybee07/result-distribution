package models

import (
	"errors"

	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	SymbolNumber       string   `gorm:"not null" json:"symbol_number"`
	RegistrationNumber string   `gorm:"not null" json:"registration_number"`
	Fullname           string   `gorm:"not null" json:"fullname"`
	BatchID            uint     `gorm:"not null" json:"batch_id"`
	ProgramID          uint     `gorm:"not null" json:"program_id"`
	CollegeID          uint     `gorm:"not null" json:"college_id"`
	CurrentSemester    uint     `gorm:"not null;default:1" json:"current_semester"`
	Status             string   `gorm:"not null;default:Active" json:"status"`
	College            College  `gorm:"foreignKey:CollegeID"`
	Batch              Batch    `gorm:"foreignKey:BatchID"`
	Program            Program  `gorm:"foreignKey:ProgramID"`
	Semester           Semester `gorm:"foreignKey:CurrentSemester"`
}

func (s *Student) AfterCreate(tx *gorm.DB) error {
	// Check if the CapacityAndCount entry exists
	var capacityAndCount CapacityAndCount
	err := tx.Where("college_id = ? AND batch_id = ? AND program_id = ?", s.CollegeID, s.BatchID, s.ProgramID).First(&capacityAndCount).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// If it doesn't exist, create a new CapacityAndCount entry
		capacityAndCount = CapacityAndCount{
			CollegeID:     s.CollegeID,
			BatchID:       s.BatchID,
			ProgramID:     s.ProgramID,
			StudentsCount: 1, // Initialize count to 1 for the first student
			Capacity:      0, // Set capacity as needed
		}
		return tx.Create(&capacityAndCount).Error
	}

	// If it exists, increment the StudentsCount
	return tx.Model(&capacityAndCount).
		Update("students_count", gorm.Expr("students_count + ?", 1)).
		Error
}
