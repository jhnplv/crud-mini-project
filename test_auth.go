package main3

import (
	"fmt"
	"log"
	"net/http"

	"crud-mini-project/config"
	controller "crud-mini-project/controllers"
)

func main() {
	// Memuat koneksi database dan migrasi dari Task 3
	config.ConnectDB()

	// Mendaftarkan endpoint auth
	http.HandleFunc("/register", controller.Register)
	http.HandleFunc("/login", controller.Login)

	fmt.Println("Server Auth berjalan di port :9000...")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
