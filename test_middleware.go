package main1

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"crud-mini-project/config"
	"crud-mini-project/controllers"
	"crud-mini-project/middleware"
)

func main() {
	// 1. Inisialisasi koneksi database
	config.ConnectDB()

	// 2. Endpoint publik (Bisa diakses tanpa token)
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/login", controllers.Login)

	// 3. Endpoint privat / dilindungi (Wajib membawa token JWT yang valid)
	// Kita membungkus handler menggunakan middleware.AuthMiddleware
	http.HandleFunc("/protected-product", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "sukses",
			"message": "Selamat! Anda berhasil mengakses endpoint rahasia yang dilindungi middleware.",
		})
	}))

	fmt.Println("Server Middleware berjalan di port :8080...")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
