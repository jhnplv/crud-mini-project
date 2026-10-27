package models

import (
	"gorm.io/gorm"
)

// Materi 17 & 18: Struct representasi tabel products untuk modul CRUD kita
type Product struct {
	gorm.Model
	Name        string `json:"name" gorm:"size:255;not null"`
	Description string `json:"description" gorm:"type:text"`
	Price       int    `json:"price" gorm:"not null"`
	Stock       int    `json:"stock" gorm:"not null;default:0"`
	Category    string `json:"category" gorm:"size:100"`
}
