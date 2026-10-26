package main

import (
	"crud-mini-project/config"
	"crud-mini-project/controllers"
	"crud-mini-project/middleware"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.ConnectDB()

	// [Analogi: Menaruh papan tombol bel di meja informasi restoran bertuliskan "Cek Kesehatan Dapur"]
	// Materi 16: Mendaftarkan rute endpoint /healthcheck menggunakan http.HandleFunc
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		// [Analogi: Menyiapkan piring saji (ResponseWriter) dengan label menu JSON]
		w.Header().Set("Content-Type", "application/json")

		// Mengirimkan kode status HTTP 200 OK (tanda sukses)
		w.WriteHeader(http.StatusOK)

		// [Analogi: Menyajikan makanan ke atas piring berupa pesan status sistem]
		w.Write([]byte(`{"status": "sukses", "message": "Server backend Go siap melayani!"}`))
	})

	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/auth/google/login", controllers.GoogleLogin)
	http.HandleFunc("/auth/google/callback", controllers.GoogleCallback)

	http.HandleFunc("/products", middleware.AuthMiddleware(controllers.CreateProduct))
	http.HandleFunc("/products/update", middleware.AuthMiddleware(controllers.UpdateProduct))
	http.HandleFunc("/products/delete", middleware.AuthMiddleware(controllers.DeleteProduct))

	port := ":9000"
	fmt.Printf("Server menyala dan mendengarkan pada port %s\n", port)

	// [Analogi: Membuka pintu gerbang utama toko secara fisik dan berdiri menunggu tamu datang]
	// Materi 9 & 16: Menyalakan HTTP server dan menangani error jika gagal berjalan
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Server gagal dijalankan: %v", err)
	}

}
