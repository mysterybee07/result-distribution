package models

import "gorm.io/gorm"

type Notice struct {
	*gorm.Model
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ProgramID   uint     `gorm:"not null" json:"program_id"`
	BatchID     *uint    `gorm:"" json:"batch_id"`
	SemesterID  *uint    `gorm:"" json:"semester_id"`
	FilePath    string   `json:"file_path"`
	Batch       Batch    `gorm:"foreignKey:BatchID" json:"-"`
	Program     Program  `gorm:"foreignKey:ProgramID" json:"-"`
	Semester    Semester `gorm:"foreignKey:SemesterID" json:"-"`
	Status      string   `gorm:"not null; default:NotPublished" json:"status"`
}

type NoticeInput struct {
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
	ProgramID   uint   `form:"program_id" json:"program_id"`
	BatchID     *uint  `form:"batch_id" json:"batch_id"`
	SemesterID  *uint  `form:"semester_id" json:"semester_id"`
}
