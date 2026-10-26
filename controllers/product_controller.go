package controllers

import (
	"crud-mini-project/config"
	"crud-mini-project/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// [Analogi: Petugas pencatat barang masuk yang menerima dan menyimpan data produk baru ke dalam sistem]
// [Materi 18: Menangani proses pembuatan data produk baru ke database menggunakan GORM]
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metode HTTP tidak diizinkan.", http.StatusMethodNotAllowed)
		return
	}

	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Format data JSON tidak valid", http.StatusBadRequest)
		return
	}

	// [Analogi: Memasukan data barang fisik ke lemari arsip pusat]
	result := config.DB.Create(&product)
	if result.Error != nil {
		http.Error(w, "Gagal menyimpan produk ke database.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "sukses",
		"message": "Produk berhasil ditambahkan!",
		"data":    product,
	})
}

// [Analogi: Petugas gudang yang memperbarui informasi data barang yang sudah ada berdasarkan ID]
// Materi 18: Menangani proses pembaruan data produk berdasarkan ID di database
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Metode HTTP tidak diizinkan.", http.StatusMethodNotAllowed)
		return
	}

	// Mengambil ID dari URL path (contoh: /products/1 -> ambil angka 1)
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "ID produk tidak valid.", http.StatusBadRequest)
		return
	}

	idStr := pathParts[len(pathParts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID produk harus berupa angka.", http.StatusBadRequest)
		return
	}

	var existingProduct models.Product
	if err := config.DB.First(&existingProduct, id).Error; err != nil {
		http.Error(w, "Produk tidak ditemukan.", http.StatusNotFound)
		return
	}

	var updateData models.Product
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Format data JSON tidak valid.", http.StatusBadRequest)
		return
	}

	existingProduct.Name = updateData.Name
	existingProduct.Description = updateData.Description
	existingProduct.Price = updateData.Price
	existingProduct.Stock = updateData.Stock

	config.DB.Save(&existingProduct)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "sukses",
		"message": "Produk berhasil diperbarui.",
		"date":    existingProduct,
	})
}

// [Analogi: Petugas penghapus arsip yang mengeluarkan data barang dari sistem gudang secara permanen]
// Materi 18: Menangani proses penghapusan data produk berdasarkan ID dari database
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Metode HTTP tidak diizinkan.", http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "ID produk tidak valid.", http.StatusBadRequest)
		return
	}

	idStr := pathParts[len(pathParts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID produk harus berupa angka.", http.StatusBadRequest)
		return
	}

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		http.Error(w, "Produk tidak ditemukan.", http.StatusNotFound)
		return
	}

	config.DB.Delete(&product)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "sukses",
		"message": "Produk berhasil dihapus!",
	})
}
