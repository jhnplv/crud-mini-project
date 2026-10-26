package models

import (
	"gorm.io/gorm"
)

// Materi 17 & 18: Struct representasi tabel users dengan GORM Model (ID, CreatedAt, UpdatedAt, DeletedAt)
type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null;size:255"`
	Password string `gorm:"not null;size:255"`
}
