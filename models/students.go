package models

import (
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	SymbolNumber       string   `gorm:"not null" json:"symbol_number" form:"symbol_number"`
	RegistrationNumber string   `gorm:"not null" json:"registration_number" form:"registration"`
	Fullname           string   `gorm:"not null" json:"fullname" form:"fullname"`
	BatchID            uint     `gorm:"not null" json:"batch_id" form:"batch_id"`
	ProgramID          uint     `gorm:"not null" json:"program_id" form:"program_id"`
	CurrentSemester    uint     `gorm:"not null;default:1" json:"current_semester" form:"current_semester"`
	Status             string   `gorm:"not null;default:Active" json:"status" form:"status"`
	Batch              Batch    `gorm:"foreignKey:BatchID"`
	Program            Program  `gorm:"foreignKey:ProgramID"`
	Semester           Semester `gorm:"foreignKey:CurrentSemester"`
}
