package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// [Analogi: Petugas pos pemeriksaan paspor di gerbang utama kawasan khusus]
// Materi 17 & 19: Middleware untuk memvalidasi keberadaan dan keaslian token JWT pada header request
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Mengambil nilai header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Token akses tidak ditemukan.", http.StatusUnauthorized)
			return
		}

		// Memastikan format header berawalan "Bearer "
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Format token tidak valid.", http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]
		jwtKey := []byte("rahasia_super_aman_sekali")

		// [Analogi: Menerawang keaslian stempel hologram pada kartu akses]
		// Memverifikasi token JWT menggunakan kunci rahasia yang sama saat diterbitkan
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Token tidak sah atau sudah kadaluarsa.", http.StatusUnauthorized)
			return
		}

		// Jika token valid, lanjutkan permintaan ke fungsi handler berikutnya
		next(w, r)
	}
}
