package config

import (
	"fmt"
	"log"

	"crud-mini-project/models"

	//"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Materi 15 & 18: Memuat konfigurasi environment dan menyambungkan aplikasi Go ke PostgreSQL
func ConnectDB() {
	// Memuat file .env agar variabel lingkungan terbaca oleh os.Getenv
	//_ = godotenv.Load()

	// [Analogi: Membaca paspor dan kunci rahasia dari brankas file lingkungan (.env)]
	//host := os.Getenv("DB_HOST")
	//user := os.Getenv("DB_USER")
	//password := os.Getenv("DB_PASSWORD")
	//dbname := os.Getenv("DB_NAME")
	//port := os.Getenv("DB_PORT")
	host := "localhost"
	user := "root"
	password := "secretpassword"
	dbname := "golang_crud_db"
	port := "5432"

	// Merangkai string DSN (Data Source Name) secara dinamis dari variabel environment
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbname, port)

	var err error

	// [Analogi: Menghubungkan kabel utama antara server aplikasi Go ke pusat penyimpanan data PostgreSQL]
	// Materi 18: Membuka koneksi GORM ke database relasional menggunakan driver Postgres
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung dengan database: %v", err)
	}

	fmt.Println("Berhasil terhubung ke database PostgreSQL via Environment Variables!")

	// [Analogi: Kontraktor bangunan otomatis yang membaca cetak biru struct lalu merakit tabel secara fisik]
	// Materi 18: Menjalankan GORM AutoMigrate untuk membuat tabel users dan products secara otomatis
	err = DB.AutoMigrate(&models.User{}, &models.Product{})
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}

	fmt.Println("Migrasi tabel database berhasil dijalankan!")
}
