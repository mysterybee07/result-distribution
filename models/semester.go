package models

import "gorm.io/gorm"

type Semester struct {
	gorm.Model
	Name      uint    `gorm:"not null" json:"name" form:"name"`
	ProgramID uint    `gorm:"not null" json:"program_id" form:"program_id"`
	Program   Program `gorm:"foreignKey:ProgramID"`
}
