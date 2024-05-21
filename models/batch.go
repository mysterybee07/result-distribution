package models

import "gorm.io/gorm"

type Batch struct {
	gorm.Model
	Year     string    `gorm:"not null" json:"year"`
	Programs []Program `gorm:"foreignKey:BatchID"`
}
