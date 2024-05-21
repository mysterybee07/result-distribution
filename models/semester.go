package models

import "gorm.io/gorm"

type Semester struct {
	gorm.Model
	Name      string  `gorm:"not null" json:"name"`
	ProgramID uint    `gorm:"not null" json:"program_id"`
	Program   Program `gorm:"foreignKey:ProgramID"`
}
