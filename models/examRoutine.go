package models

import (
	"time"

	"gorm.io/gorm"
)

type ExamRoutine struct {
	gorm.Model
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	BatchID    uint      `json:"batch_id"`
	ProgramID  uint      `json:"program_id"`
	SemesterID uint      `json:"semester_id"`
	Status     bool      `gorm:"not null; default:false" json:"status"`

	// Foreign key associations
	Batch    Batch    `gorm:"foreignKey:BatchID"`
	Program  Program  `gorm:"foreignKey:ProgramID"`
	Semester Semester `gorm:"foreignKey:SemesterID"`
}

type ExamSchedules struct {
	gorm.Model
	CourseID      uint      `json:"course_id"`
	ExamRoutineID uint      `json:"exam_routine_id"`
	ExamDate      time.Time `json:"exam_date"`

	Course      Course      `gorm:"foreignKey:CourseID"`
	ExamRoutine ExamRoutine `gorm:"foreignKey:ExamRoutineID"`
}

type ExamRoutineRequest struct {
	BatchID    uint      `json:"batch_id"`
	ProgramID  uint      `json:"program_id"`
	SemesterID uint      `json:"semester_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
}
