// models/mark.go

package models

import (
	"gorm.io/gorm"
)

type Mark struct {
	gorm.Model
	BatchID        uint     `gorm:"not null" json:"batch_id"`
	ProgramID      uint     `gorm:"not null" json:"program_id"`
	SemesterID     uint     `gorm:"not null" json:"semester_id"`
	SubjectID      uint     `gorm:"not null" json:"subject_id"`
	StudentID      uint     `gorm:"not null" json:"student_id"`
	SemesterMarks  int      `gorm:"not null" json:"semester_marks"`
	AssistantMarks int      `gorm:"not null" json:"assistant_marks"`
	PracticalMarks int      `gorm:"not null" json:"practical_marks"`
	TotalMarks     int      `gorm:"->;type:int GENERATED ALWAYS AS (semester_marks + assistant_marks + practical_marks) STORED" json:"total_marks"`
	Batch          Batch    `gorm:"foreignkey:BatchID"`
	Program        Program  `gorm:"foreignkey:ProgramID"`
	Semester       Semester `gorm:"foreignkey:SemesterID"`
	Subject        Subject  `gorm:"foreignkey:SubjectID"`
	Student        Student  `gorm:"foreignkey:StudentID"`
}
