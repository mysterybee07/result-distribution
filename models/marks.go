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
	CourseID       uint     `gorm:"not null" json:"course_id"`
	StudentID      uint     `gorm:"not null" json:"student_id"`
	SemesterMarks  int      `gorm:"not null" json:"semester_marks"`
	AssistantMarks int      `gorm:"not null" json:"assistant_marks"`
	PracticalMarks int      `gorm:"not null" json:"practical_marks"`
	TotalMarks     int      `gorm:"->;type:int GENERATED ALWAYS AS (semester_marks + assistant_marks + practical_marks) STORED" json:"total_marks"`
	Status         string   `gorm:"default:pass" json:"status"`
	Batch          Batch    `gorm:"foreignkey:BatchID"`
	Program        Program  `gorm:"foreignkey:ProgramID"`
	Semester       Semester `gorm:"foreignkey:SemesterID"`
	Course         Course   `gorm:"foreignkey:CourseID"`
	Student        Student  `gorm:"foreignkey:StudentID"`
}

// BeforeSave hook to set status based on marks
func (m *Mark) BeforeSave(tx *gorm.DB) (err error) {
	if m.SemesterMarks < 24 || m.AssistantMarks < 8 || m.PracticalMarks < 8 {
		m.Status = "failed"
	} else {
		m.Status = "pass"
	}
	return
}
