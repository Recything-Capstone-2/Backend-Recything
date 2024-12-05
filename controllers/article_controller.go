package controllers

import (
	"Backend-Recything/config"
	"Backend-Recything/helper"
	"Backend-Recything/models"
	"net/http"
	"strconv"

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
	var articles []models.Article

	if err := config.DB.Find(&articles).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Gagal mengambil artikel", http.StatusInternalServerError, "error", nil))
	}
	var articleResponses []map[string]interface{}
	for _, article := range articles {
		articleResponses = append(articleResponses, map[string]interface{}{
			"id":         article.ID,
			"judul":      article.Judul,
			"author":     article.Author,
			"konten":     article.Konten,
			"link_foto":  article.LinkFoto,
			"link_video": article.LinkVideo,
			"created_at": article.CreatedAt,
			"updated_at": article.UpdatedAt,
		})
	}

	return c.JSON(http.StatusOK, helper.APIResponse("Artikel sukses diambil", http.StatusOK, "success", articleResponses))
}
func AmbilArtikelByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
	}

	var artikel models.Article
	if err := config.DB.First(&artikel, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Article tidak ditemukan"})
	}

	return c.JSON(http.StatusOK, artikel)
}
