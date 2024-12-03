package controllers

import (
	"Backend-Recything/config"
	"Backend-Recything/helper"
	"Backend-Recything/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Struct untuk input artikel
type CreateArticleInput struct {
	Judul     string `json:"judul" validate:"required"`
	Author    string `json:"author" validate:"required"`
	Konten    string `json:"konten" validate:"required"`
	LinkFoto  string `json:"link_foto" validate:"required,url"` // Foto wajib, URL
	LinkVideo string `json:"link_video" validate:"omitempty,url"`
}

func CreateArticle(c echo.Context) error {
	var input CreateArticleInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Invalid input", http.StatusBadRequest, "error", nil))
	}

	// Validasi input
	if err := c.Validate(&input); err != nil {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Validation error", http.StatusBadRequest, "error", helper.FormatValidationError(err)))
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
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to create article", http.StatusInternalServerError, "error", nil))
	}

	return c.JSON(http.StatusOK, helper.APIResponse("Article created successfully", http.StatusOK, "success", article))
}
