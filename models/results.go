package models

import (
	"gorm.io/gorm"
)

type Result struct {
	*gorm.Model
	BatchID    uint     `gorm:"not null" json:"batch_id"`
	ProgramID  uint     `gorm:"not null" json:"program_id"`
	SemesterID uint     `gorm:"not null" json:"semester_id"`
	Batch      Batch    `gorm:"foreignKey:BatchID"`
	Program    Program  `gorm:"foreignKey:ProgramID"`
	Semester   Semester `gorm:"foreignKey:SemesterID"`
	Status     string   `gorm:"not null; default:NotPublished" json:"status"`
}
