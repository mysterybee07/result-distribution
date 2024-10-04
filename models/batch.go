package models

import "gorm.io/gorm"

type Batch struct {
	gorm.Model
	Batch uint `gorm:"not null" json:"batch"`
	// Programs []Program `gorm:"foreignKey:BatchID"`
}
