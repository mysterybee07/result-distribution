package models

import "gorm.io/gorm"

type Program struct {
	gorm.Model
	ProgramName string     `gorm:"not null" json:"program_name"`
	Semesters   []Semester `gorm:"foreignKey:ProgramID"`
}
