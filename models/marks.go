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
type MarksPayload struct {
	BatchID    uint `json:"batch_id" validate:"required"`
	ProgramID  uint `json:"program_id" validate:"required"`
	SemesterID uint `json:"semester_id" validate:"required"`
	CourseID   uint `json:"course_id" validate:"required"`
	Marks      []struct {
		StudentID      uint `json:"student_id" validate:"required"`
		SemesterMarks  int  `json:"semester_marks" validate:"required"`
		AssistantMarks int  `json:"assistant_marks"`
		PracticalMarks int  `json:"practical_marks"`
	} `json:"marks" validate:"required,dive"`
}

func (m *Mark) BeforeSave(tx *gorm.DB) (err error) {
	var course Course
	if err := tx.First(&course, m.CourseID).Error; err != nil {
		return err
	}

	if m.SemesterMarks < course.SemesterPassMarks ||
		(course.PracticalPassMarks != nil && m.PracticalMarks < *course.PracticalPassMarks) ||
		(course.AssistantPassMarks != nil && m.AssistantMarks < *course.AssistantPassMarks) {
		m.Status = "failed"
	} else {
		m.Status = "pass"
	}

	return
}
