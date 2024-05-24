package models

import (
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	SymbolNumber string  `gorm:"not null" json:"symbol_number"`
	Registration string  `gorm:"not null" json:"registration"`
	Fullname     string  `gorm:"not null" json:"fullname"`
	BatchID      uint    `gorm:"not null" json:"batch_id"`
	ProgramID    uint    `gorm:"not null" json:"program_id"`
	Batch        Batch   `gorm:"foreignKey:BatchID"`
	Program      Program `gorm:"foreignKey:ProgramID"`
}
