package models

import (
	"github.com/mysterybee07/result-distribution-system/models"
	"gorm.io/gorm"
)

type Result struct {
	*gorm.Model
	BatchID    uint            `gorm:"not null" json:"batch_id"`
	ProgramID  uint            `gorm:"not null" json:"program_id"`
	SemesterID uint            `gorm:"not null" json:"semester_id"`
	Batch      models.Batch    `gorm:"foreignKey:BatchID"`
	Program    models.Program  `gorm:"foreignKey:ProgramID"`
	Semester   models.Semester `gorm:"foreignKey:SemesterID"`
	Status     string          `gorm:"not null; default:NotPublished" json:"status"`
}
