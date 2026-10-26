package main

import (
	"fmt"
	"log"
	"net/http"

	"crud-mini-project/config"
	"crud-mini-project/controllers"

	"github.com/joho/godotenv"
)

func main() {
	// Load file .env untuk membaca Client ID & Secret
	godotenv.Load()

	// Hubungkan database
	config.ConnectDB()

	// Rute Google OAuth
	http.HandleFunc("/auth/google/login", controllers.GoogleLogin)
	http.HandleFunc("/auth/google/callback", controllers.GoogleCallback)

	fmt.Println("Server Google Auth berjalan di port :9000...")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
