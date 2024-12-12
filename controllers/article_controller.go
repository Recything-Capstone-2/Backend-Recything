package controllers

import (
	"Backend-Recything/config"
	"Backend-Recything/helper"
	"Backend-Recything/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// Struct untuk input artikel
type InputBikinArtikel struct {
	Judul     string `json:"judul" validate:"required"`
	Author    string `json:"author" validate:"required"`
	Konten    string `json:"konten" validate:"required"`
	LinkFoto  string `json:"link_foto" validate:"required,url"` // Foto wajib, URL
	LinkVideo string `json:"link_video" validate:"omitempty,url"`
}

// Struct untuk input update artikel
type InputUpdateArtikel struct {
	Judul     string `json:"judul" validate:"required"`
	Author    string `json:"author" validate:"required"`
	Konten    string `json:"konten" validate:"required"`
	LinkFoto  string `json:"link_foto" validate:"required,url"`
	LinkVideo string `json:"link_video" validate:"omitempty,url"`
}

func BikinArtikel(c echo.Context) error {
	var input InputBikinArtikel
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Input invalid", http.StatusBadRequest, "error", nil))
	}

	// Validasi input
	if err := c.Validate(&input); err != nil {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Validasi error", http.StatusBadRequest, "error", helper.FormatValidationError(err)))
	}

	// Buat artikel baru
	article := models.Article{
		Judul:     input.Judul,
		Author:    input.Author,
		Konten:    input.Konten,
		LinkFoto:  input.LinkFoto,
		LinkVideo: input.LinkVideo,
	}

	// Simpan artikel ke database
	if err := config.DB.Create(&article).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Gagal untuk membuat artikel", http.StatusInternalServerError, "error", nil))
	}

	// Format data respons tanpa `created_at` dan `updated_at`
	articleResponse := map[string]interface{}{
		"id":         article.ID,
		"judul":      article.Judul,
		"author":     article.Author,
		"konten":     article.Konten,
		"link_foto":  article.LinkFoto,
		"link_video": article.LinkVideo,
	}

	return c.JSON(http.StatusOK, helper.APIResponse("Artikel sukses terbuat", http.StatusOK, "success", articleResponse))
}

func AmbilSemuaArtikel(c echo.Context) error {
	// Ambil parameter query untuk paginasi
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	// Default nilai untuk paginasi
	page := 1
	limit := 10

	// Parse parameter jika ada
	if pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}
	if limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}

	// Hitung offset berdasarkan page dan limit
	offset := (page - 1) * limit

	// Hitung total artikel
	var totalItems int64
	if err := config.DB.Model(&models.Article{}).Count(&totalItems).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Gagal menghitung artikel", http.StatusInternalServerError, "error", nil))
	}

	// Ambil data artikel dengan paginasi
	var articles []models.Article
	if err := config.DB.Offset(offset).Limit(limit).Order("created_at DESC").Find(&articles).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Gagal mengambil artikel", http.StatusInternalServerError, "error", nil))
	}

	// Mapping artikel ke dalam respons
	var articleResponses []map[string]interface{}
	for _, article := range articles {
		articleResponses = append(articleResponses, map[string]interface{}{
			"id":         article.ID,
			"judul":      article.Judul,
			"author":     article.Author,
			"konten":     article.Konten,
			"link_foto":  article.LinkFoto,
			"link_video": article.LinkVideo,
			"created_at": article.CreatedAt.Format("2006-01-02"), // Format waktu
		})
	}

	// Hitung total halaman
	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))

	// Format respons dengan paginasi
	response := map[string]interface{}{
		"items": articleResponses,
		"pagination": map[string]interface{}{
			"current_page": page,
			"per_page":     limit,
			"total_items":  totalItems,
			"total_pages":  totalPages,
		},
	}

	return c.JSON(http.StatusOK, helper.APIResponse("Artikel sukses diambil", http.StatusOK, "success", response))
}

func AmbilArtikelByID(c echo.Context) error {
	// Ambil ID dari parameter URL
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		// Respon untuk ID tidak valid
		return c.JSON(http.StatusBadRequest, helper.APIResponse("ID tidak valid", http.StatusBadRequest, "error", nil))
	}

	// Cari artikel berdasarkan ID
	var artikel models.Article
	if err := config.DB.First(&artikel, id).Error; err != nil {
		// Respon jika artikel tidak ditemukan
		return c.JSON(http.StatusNotFound, helper.APIResponse("Artikel tidak ditemukan", http.StatusNotFound, "error", nil))
	}

	// Format respons artikel
	artikelResponse := map[string]interface{}{
		"id":         artikel.ID,
		"judul":      artikel.Judul,
		"author":     artikel.Author,
		"konten":     artikel.Konten,
		"link_foto":  artikel.LinkFoto,
		"link_video": artikel.LinkVideo,
		"created_at": artikel.CreatedAt.Format("2006-01-02"),
	}
	// Respon sukses
	return c.JSON(http.StatusOK, helper.APIResponse("Artikel berhasil ditemukan", http.StatusOK, "success", artikelResponse))
}

func UpdateArtikel(c echo.Context) error {
	// Ambil ID dari parameter URL
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("ID tidak valid", http.StatusBadRequest, "error", nil))
	}

	// Cari artikel berdasarkan ID
	var artikel models.Article
	if err := config.DB.First(&artikel, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, helper.APIResponse("Artikel tidak ditemukan", http.StatusNotFound, "error", nil))
	}

	// Bind input dari request body
	var input InputUpdateArtikel
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Input invalid", http.StatusBadRequest, "error", nil))
	}

	// Validasi input
	if err := c.Validate(&input); err != nil {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Validasi error", http.StatusBadRequest, "error", helper.FormatValidationError(err)))
	}

	// Update artikel
	artikel.Judul = input.Judul
	artikel.Author = input.Author
	artikel.Konten = input.Konten
	artikel.LinkFoto = input.LinkFoto
	artikel.LinkVideo = input.LinkVideo

	// Simpan perubahan ke database
	if err := config.DB.Save(&artikel).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Gagal mengupdate artikel", http.StatusInternalServerError, "error", nil))
	}

	// Format respons
	artikelResponse := map[string]interface{}{
		"id":         artikel.ID,
		"judul":      artikel.Judul,
		"author":     artikel.Author,
		"konten":     artikel.Konten,
		"link_foto":  artikel.LinkFoto,
		"link_video": artikel.LinkVideo,
	}

	return c.JSON(http.StatusOK, helper.APIResponse("Artikel berhasil diupdate", http.StatusOK, "success", artikelResponse))
}

func DeleteArtikel(c echo.Context) error {
	// Ambil ID dari parameter URL
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		// Respon untuk ID tidak valid
		return c.JSON(http.StatusBadRequest, helper.APIResponse("ID tidak valid", http.StatusBadRequest, "error", nil))
	}

	// Cari artikel berdasarkan ID
	var artikel models.Article
	if err := config.DB.First(&artikel, id).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			// Respon jika artikel tidak ditemukan
			return c.JSON(http.StatusNotFound, helper.APIResponse("Artikel tidak ditemukan", http.StatusNotFound, "error", nil))
		}
		// Respon jika ada kesalahan lain
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Gagal mencari artikel", http.StatusInternalServerError, "error", nil))
	}

	// Hapus artikel dari database
	if err := config.DB.Delete(&artikel).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Gagal menghapus artikel", http.StatusInternalServerError, "error", nil))
	}

	// Respon sukses
	return c.JSON(http.StatusOK, helper.APIResponse("Artikel berhasil dihapus", http.StatusOK, "success", nil))
}
