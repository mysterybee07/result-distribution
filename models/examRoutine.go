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
}
