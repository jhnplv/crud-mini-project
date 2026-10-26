package models

import (
	"gorm.io/gorm"
)

// Materi 17 & 18: Struct representasi tabel products untuk modul CRUD kita
type Product struct {
	gorm.Model
	Name        string `gorm:"size:255;not null"`
	Description string `gorm:"type:text"`
	Price       int    `gorm:"not null"`
	Stock       int    `gorm:"not null;default:0"`
}
