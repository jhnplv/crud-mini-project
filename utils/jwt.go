package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// [Analogi: Stempel basah resmi pimpinan perusahaan yang disimpan rapat-rapat]
// Variabel kunci rahasia untuk menandatangani dan memverifikasi keaslian token digital
var jwtKey = []byte("rahasia_super_aman_sekali")

// [Analogi: Mencetak kartu akses ID Card khusus bagi karyawan yang sah]
// Materi 20: Membuat token JWT berisi klaim data (email) dan masa aktif selama 24 jam
func GenerateToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
