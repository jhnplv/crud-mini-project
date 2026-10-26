package main

import (
	"crud-mini-project/config"
	// "crud-mini-project/models"
)

// Materi 18: Titik eksekusi utama program untuk menguji koneksi database dan migrasi
func main() {
	// Memanggil fungsi koneksi dan migrasi GORM ke PostgreSQL
	config.ConnectDB()
}