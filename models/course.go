package models

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	CourseCode          string   `gorm:"course_code" json:"course_code" form:"course_code"`
	Name                string   `gorm:"not null" json:"name" form:"name"`
	SemesterPassMarks   int      `json:"semester_pass_marks" validate:"required" form:"semester_pass_marks"`
	PracticalPassMarks  *int     `json:"practical_pass_marks,omitempty" form:"practical_pass_marks"`
	AssistantPassMarks  *int     `json:"assistant_pass_marks,omitempty" form:"assistant_pass_marks"`
	SemesterTotalMarks  int      `json:"semester_total_marks" validate:"required" form:"semester_total_marks"`
	PracticalTotalMarks *int     `json:"practical_total_marks,omitempty" form:"practical_total_marks"`
	AssistantTotalMarks *int     `json:"assistant_total_marks,omitempty" form:"assistant_total_marks"`
	ProgramID           uint     `gorm:"not null" json:"program_id" form:"program_id"`
	SemesterID          uint     `gorm:"not null" json:"semester_id" form:"semester_id"`
	Program             Program  `gorm:"foreignKey:ProgramID"`
	Semester            Semester `gorm:"foreignKey:SemesterID"`
}
