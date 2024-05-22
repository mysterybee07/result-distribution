package models

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	Name       string   `gorm:"not null" json:"name"`
	CourseCode string   `gorm:"course_code" json:"course_code"`
	ProgramID  uint     `gorm:"not null" json:"program_id"`
	SemesterID uint     `gorm:"not null" json:"semester_id"`
	Program    Program  `gorm:"foreignKey:ProgramID"`
	Semester   Semester `gorm:"foreignKey:SemesterID"`
}
