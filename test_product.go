package main

import (
	"fmt"
	"log"
	"net/http"

	"crud-mini-project/config"
	"crud-mini-project/controllers"
	"crud-mini-project/middleware"
)

func main() {
	// 1. Hubungkan database
	config.ConnectDB()

	// 2. Rute Publik (Auth)
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/login", controllers.Login)

	// 3. Rute Privat / Terproteksi (CRUD Produk) dilindungi Middleware JWT
	http.HandleFunc("/products", middleware.AuthMiddleware(controllers.CreateProduct))

	// Untuk Update & Delete, kita menggunakan endpoint dasar (contoh sederhana routing ID)
	http.HandleFunc("/products/update", middleware.AuthMiddleware(controllers.UpdateProduct))
	http.HandleFunc("/products/delete", middleware.AuthMiddleware(controllers.DeleteProduct))

	fmt.Println("Server CRUD Produk berjalan di port :9000...")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
