package controllers

import (
	"Backend-Recything/config"
	"Backend-Recything/helper"
	"Backend-Recything/models"
	"net/http"

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
		LinkFoto:  input.LinkFoto,  // Ini link foto wajib (URL)
		LinkVideo: input.LinkVideo, // Ini link video opsional (URL)
	}

	// simpen artikel ke database
	if err := config.DB.Create(&article).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Gagal untuk membuat artikel", http.StatusInternalServerError, "error", nil))
	}

	return c.JSON(http.StatusOK, helper.APIResponse("Artikel sukses terbuat", http.StatusOK, "success", article))
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
