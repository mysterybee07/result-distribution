package models

import (
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	CourseCode          string   `gorm:"unique" json:"course_code"`
	Name                string   `gorm:"not null" json:"name"`
	SemesterPassMarks   int      `json:"semester_pass_marks" validate:"required"`
	PracticalPassMarks  *int     `json:"practical_pass_marks,omitempty"`
	AssistantPassMarks  *int     `json:"assistant_pass_marks,omitempty"`
	SemesterTotalMarks  int      `json:"semester_total_marks" validate:"required"`
	PracticalTotalMarks *int     `json:"practical_total_marks,omitempty"`
	AssistantTotalMarks *int     `json:"assistant_total_marks,omitempty"`
	ProgramID           uint     `gorm:"not null" json:"program_id"`
	SemesterID          uint     `gorm:"not null" json:"semester_id"`
	IsCompulsory        bool     `gorm:"not null" json:"is_compulsory" default:"false"`
	Program             Program  `gorm:"foreignKey:ProgramID"`
	Semester            Semester `gorm:"foreignKey:SemesterID"`
}

type CoursesPayload struct {
	ProgramID  uint     `json:"program_id" validate:"required"`
	SemesterID uint     `json:"semester_id" validate:"required"`
	Courses    []Course `json:"courses" validate:"required"`
}
