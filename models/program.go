package models

import "gorm.io/gorm"

type Program struct {
	gorm.Model
	Name      string     `gorm:"not null" json:"name" form:"name"`
	Semesters []Semester `gorm:"foreignKey:ProgramID"`
}
