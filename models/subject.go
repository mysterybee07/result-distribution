package models

import "gorm.io/gorm"

type Subject struct {
	gorm.Model
	Name       string   `gorm:"not null" json:"name"`
	ProgramID  uint     `gorm:"not null" json:"program_id"`
	SemesterID uint     `gorm:"not null" json:"semester_id"`
	Program    Program  `gorm:"foreignKey:ProgramID"`
	Semester   Semester `gorm:"foreignKey:SemesterID"`
}
