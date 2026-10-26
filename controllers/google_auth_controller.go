package controllers

import (
	"crud-mini-project/config"
	"crud-mini-project/models"
	"crud-mini-project/utils"
	"encoding/json"
	"net/http"
)

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	googleConfig := utils.GetGoogleOAuthConfig()
	url := googleConfig.AuthCodeURL("state-pseudo-random")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Kode otorisasi Google tidak ditemukan.", http.StatusBadRequest)
		return
	}

	googleConfig := utils.GetGoogleOAuthConfig()

	token, err := googleConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Gagal menukar token Google.", http.StatusInternalServerError)
		return
	}

	client := googleConfig.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Gagal mengambil data pengguna dari Google.", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		http.Error(w, "Gagal membaca respons profil Google.", http.StatusInternalServerError)
		return
	}

	var dbUser models.User
	result := config.DB.Where("email = ?", googleUser.Email).First(&dbUser)

	if result.Error != nil {
		newUser := models.User{
			Email:    googleUser.Email,
			Password: "",
		}
		if err := config.DB.Create(&newUser).Error; err != nil {
			http.Error(w, "Gagal menyimpan akun Google ke database.", http.StatusInternalServerError)
			return
		}
		dbUser = newUser
	}

	tokenString, err := utils.GenerateToken(dbUser.Email)
	if err != nil {
		http.Error(w, "Gagal membuat token akses.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "sukses",
		"message": "Login dengan Google berhasil.",
		"token":   tokenString,
	})
}
