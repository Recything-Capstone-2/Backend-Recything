package controllers

import (
	"Backend-Recything/config"
	"Backend-Recything/helper"
	"Backend-Recything/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

func TestGetUserPoints(t *testing.T) {
	// Membuat instance Echo
	e := echo.New()

	// Membuat data user dan poin untuk testing
	user := models.User{
		NamaLengkap:  "User Test",
		Email:        "usertest@example.com",
		NoTelepon:    "081234567890",
		TanggalLahir: time.Now(), // Menggunakan tanggal valid
	}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("Gagal membuat user: %v", err)
	}

	points := models.Points{
		UserID: user.ID,
		Points: 50,
	}
	if err := config.DB.Create(&points).Error; err != nil {
		t.Fatalf("Gagal membuat points: %v", err)
	}

	// Membuat request dan recorder
	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/points", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", user.ID)

	// Memanggil controller
	if err := GetUserPoints(c); err != nil {
		t.Fatalf("Gagal mendapatkan user points: %v", err)
	}

	// Memeriksa status kode
	if rec.Code != http.StatusOK {
		t.Errorf("Status code seharusnya %d, tapi dapat %d", http.StatusOK, rec.Code)
	}

	// Memeriksa response
	var response helper.Response
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Gagal membaca response: %v", err)
	}

	if response.Meta.Status != "success" {
		t.Errorf("Response status seharusnya success, tapi dapat %s", response.Meta.Status)
	}

	data := response.Data.(map[string]interface{})
	if data["points"] != float64(50) {
		t.Errorf("Points seharusnya 50, tapi dapat %v", data["points"])
	}
	config.DB.Where("user_id = ?", user.ID).Delete(&models.Points{})
	config.DB.Where("email = ?", "usertest@example.com").Delete(&models.User{})
}

func TestGetAllUserPoints(t *testing.T) {
	// Membuat instance Echo
	e := echo.New()

	// Membuat data user dan poin untuk testing
	user1 := models.User{
		NamaLengkap:  "User Test 1",
		Email:        "user1@example.com",
		NoTelepon:    "081234567891",
		TanggalLahir: time.Now(), // Menggunakan tanggal valid
	}
	if err := config.DB.Create(&user1).Error; err != nil {
		t.Fatalf("Gagal membuat user1: %v", err)
	}

	user2 := models.User{
		NamaLengkap:  "User Test 2",
		Email:        "user2@example.com",
		NoTelepon:    "081234567892",
		TanggalLahir: time.Now(), // Menggunakan tanggal valid
	}
	if err := config.DB.Create(&user2).Error; err != nil {
		t.Fatalf("Gagal membuat user2: %v", err)
	}

	points1 := models.Points{
		UserID: user1.ID, // Menggunakan ID user1 yang benar
		Points: 100,
	}
	if err := config.DB.Create(&points1).Error; err != nil {
		t.Fatalf("Gagal membuat points1: %v", err)
	}

	points2 := models.Points{
		UserID: user2.ID, // Menggunakan ID user2 yang benar
		Points: 200,
	}
	if err := config.DB.Create(&points2).Error; err != nil {
		t.Fatalf("Gagal membuat points2: %v", err)
	}

	// Membuat request dan recorder
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/points", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Memanggil controller
	if err := GetAllUserPoints(c); err != nil {
		t.Fatalf("Gagal mendapatkan semua user points: %v", err)
	}

	// Memeriksa status kode
	if rec.Code != http.StatusOK {
		t.Errorf("Status code seharusnya %d, tapi dapat %d", http.StatusOK, rec.Code)
	}

	// Memeriksa response
	var response helper.Response
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Gagal membaca response: %v", err)
	}

	if response.Meta.Status != "success" {
		t.Errorf("Response status seharusnya success, tapi dapat %s", response.Meta.Status)
	}

	// Memastikan data berisi dua user
	data := response.Data.([]interface{})
	if len(data) != 2 {
		t.Errorf("Jumlah data seharusnya 2, tapi dapat %d", len(data))
	}

	// Memeriksa jika setiap user memiliki points yang benar
	for _, item := range data {
		userData := item.(map[string]interface{})
		points := userData["points"].(float64)

		// Periksa points sesuai dengan yang diharapkan
		if userData["email"] == "user1@example.com" && points != 100 {
			t.Errorf("Points untuk user1 seharusnya 100, tapi dapat %v", points)
		}
		if userData["email"] == "user2@example.com" && points != 200 {
			t.Errorf("Points untuk user2 seharusnya 200, tapi dapat %v", points)
		}
	}

	// Penghapusan data setelah test
	config.DB.Where("user_id = ?", user1.ID).Delete(&models.Points{})
	config.DB.Where("user_id = ?", user2.ID).Delete(&models.Points{})

	config.DB.Where("email = ?", "user1@example.com").Delete(&models.User{})
	config.DB.Where("email = ?", "user2@example.com").Delete(&models.User{})

}
