package controllers

import (
	"crud-mini-project/config"
	"crud-mini-project/models"
	"crud-mini-project/utils"
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

	//err := json.NewDecoder(r.Body).Decode(&product)
	//if err != nil {
	//	http.Error(w, "Format data JSON tidak valid", http.StatusBadRequest)
	//	return
	//}

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Status json tidak valid.",
		})
		return
	}

	//TASK 9: Penambahan validasi
	if utils.IsEmpty(product.Name) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Validasi gagal: Nama produk tidak boleh kosong.",
		})
		return
	}

	if product.Price <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Validasi gagal: Harga produk harus lebih dari 0.",
		})
		return
	}

	if product.Stock < 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Validasi gagal: Stok produk tidak boleh bernilai negatif.",
		})
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

// [Analogi: Rak arsip khusus untuk mengambil satu berkas spesifik berdasarkan nomor ID dokumen]
// Materi Task 8: Mengambil detail produk berdasarkan ID
func GetProductByID(w http.ResponseWriter, r *http.Request) {
	// Mengambil parameter id dari query string (contoh: /products/detail?id=1)
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Parameter ID produk diperlukan", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Format tidak valid.", http.StatusBadRequest)
		return
	}

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		http.Error(w, "Produk tidak ditemukan.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "sukses",
		"data":   product,
	})
}

// [Analogi: Buku katalog toko yang dilengkapi halaman (pagination) dan indeks pencarian nama (search)]
// Materi Task 8: Mengambil daftar produk dengan fitur Pagination dan Search
func GetProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	//query := config.DB.Model(&models.Product{})

	//Fitur search by name
	//search := r.URL.Query().Get("search")
	//if search != "" {
	//	query = query.Where("name ILIKE ?", "%"+search+"%")
	//}

	// 1. Parameter Pagination (Default: page 1, limit 10)
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	offset := (page - 1) * limit

	// 2. Parameter Search (Pencarian berdasarkan nama produk)
	search := r.URL.Query().Get("search")

	// Query dasar GORM
	query := config.DB.Model(&models.Product{})

	if search != "" {
		// Menggunakan ILike atau LIKE untuk pencarian case-insensitive pada nama produk
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	//Filter
	category := r.URL.Query().Get("category")
	if category != "" {
		query = query.Where("category ILIKE ?", "%"+category+"%")
	}

	minPrice := r.URL.Query().Get("min_price")
	if minPrice != "" {
		query = query.Where("price >= ?", minPrice)
	}

	maxPrice := r.URL.Query().Get("max_price")
	if maxPrice != "" {
		query = query.Where("price <= ?", maxPrice)
	}

	// Terapkan limit, offset untuk pagination, lalu ambil data
	if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		http.Error(w, "Gagal mengambil data produk.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "sukses",
		"page":   page,
		"limit":  limit,
		"data":   products,
	})

}
