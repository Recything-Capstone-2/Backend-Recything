package controllers

import (
	"Backend-Recything/config"
	"Backend-Recything/helper"
	"Backend-Recything/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Fungsi untuk mendapatkan poin dari pengguna berdasarkan userID
func GetUserPoints(c echo.Context) error {
	// Mendapatkan userID dari token JWT
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "Invalid user ID from token")
	}

	// Mengambil data poin dari tabel Points dan preload data User
	var userPoints models.Points
	if err := config.DB.Preload("User").Where("user_id = ?", userID).First(&userPoints).Error; err != nil {
		// Jika tidak ditemukan, kembalikan poin 0
		if err.Error() == "record not found" {
			userPoints.Points = 0
		} else {
			return c.JSON(http.StatusInternalServerError, "Failed to retrieve points")
		}
	}

	// Menyiapkan respons dengan menambahkan data pengguna (User)
	response := struct {
		ID          uint   `json:"id"`
		UserID      uint   `json:"user_id"`
		Points      uint   `json:"points"`
		NamaLengkap string `json:"nama_lengkap"`
		Email       string `json:"email"`
	}{
		ID:          userPoints.ID,
		UserID:      userPoints.UserID,
		Points:      userPoints.Points,
		NamaLengkap: userPoints.User.NamaLengkap, // Akses data User
		Email:       userPoints.User.Email,       // Akses data User
	}

	// Mengembalikan response sukses dengan data yang diinginkan
	return c.JSON(http.StatusOK, helper.APIResponse("User points retrieved successfully", http.StatusOK, "success", response))
}
