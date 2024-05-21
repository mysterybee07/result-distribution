package models

import "gorm.io/gorm"

type Program struct {
	gorm.Model
	Name    string `gorm:"not null" json:"name"`
	BatchID uint   `gorm:"not null" json:"batch_id"`
	Batch   Batch  `gorm:"foreignKey:BatchID"`
}
