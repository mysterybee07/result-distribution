package models

import (
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
