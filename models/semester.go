package models

import "gorm.io/gorm"

type Semester struct {
	gorm.Model
	SemesterName uint    `gorm:"not null" json:"semester_name"`
	ProgramID    uint    `gorm:"not null" json:"program_id"`
	Program      Program `gorm:"foreignKey:ProgramID"`
}
