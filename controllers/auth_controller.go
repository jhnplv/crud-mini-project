package controllers

import (
	"crud-mini-project/config"
	"crud-mini-project/models"
	"crud-mini-project/utils"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// [Analogi: Meja resepsionis pendaftaran anggota baru yang mencatat dan mengamankan data]
// Materi 14 & 20: Menangani proses registrasi akun dengan parsing JSON dan hashing password menggunakan bcrypt
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metode HTTP tidak diizinkan.", http.StatusMethodNotAllowed)
		return
	}

	var userInput models.User

	// Mendecode payload JSON dari body request ke dalam struktur data User
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		http.Error(w, "Format data tidak valid.", http.StatusBadRequest)
		return
	}

	// [Analogi: Memasukkan password asli ke mesin penghancur dokumen rahasia menjadi kode acak satu arah]
	// Mengenkripsi password mentah menjadi hash menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Gagal mengenkripsi password.", http.StatusInternalServerError)
		return
	}

	newUser := models.User{
		Email:    userInput.Email,
		Password: string(hashedPassword),
	}

	// [Analogi: Menyimpan data ke dalam lemari arsip pusat database]
	// Materi 15 & 18: Menyimpan data pengguna baru ke database PostgreSQL melalui GORM
	result := config.DB.Create(&newUser)
	if result.Error != nil {
		http.Error(w, "Email sudah terdaftar atau terjadi kesalahan database.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status": "sukses", "message":"Registrasi akun berhasil!"}`))
}

// [Analogi: Pos pemeriksaan tiket masuk yang mencocokkan identitas dan memberikan gelang akses]
// Materi 16 & 20: Menangani proses login, verifikasi password hash, dan penerbitan token JWT
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metode HTTP tidak diizinkan.", http.StatusMethodNotAllowed)
		return
	}

	var userInput models.User
	json.NewDecoder(r.Body).Decode(&userInput)

	var dbUser models.User

	// Mencari data pengguna berdasarkan email di dalam database
	err := config.DB.Where("email = ?", userInput.Email).First(&dbUser).Error
	if err != nil {
		http.Error(w, "Email atau password salah.", http.StatusUnauthorized)
		return
	}

	// [Analogi: Mencocokkan bentuk fisik kunci dengan cetakan lubang brankas]
	// Membandingkan password input dengan hash password yang tersimpan di database
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(userInput.Password))
	if err != nil {
		http.Error(w, "Email atau password salah.", http.StatusUnauthorized)
		return
	}

	// Menghasilkan token JWT jika proses autentikasi berhasil
	tokenString, err := utils.GenerateToken(dbUser.Email)
	if err != nil {
		http.Error(w, "Gagal membuat token akses.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Mengembalikan respons JSON berisi token akses
	json.NewEncoder(w).Encode(map[string]string{
		"status": "sukses",
		"token":  tokenString,
	})
}
