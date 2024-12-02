package controllers

import (
	"Backend-Recything/config"
	"Backend-Recything/helper"
	"Backend-Recything/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"encoding/json"
	"net/url"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/labstack/echo/v4"
)

// Struct untuk input laporan
type ReportInput struct {
	Location       string `form:"location" validate:"required"` // Menggunakan alamat untuk mendapatkan latitude dan longitude
	Description    string `form:"description" validate:"required"`
	Photo          string `form:"photo"`
	Status         string `form:"status"`
	TanggalLaporan string `form:"tanggal_laporan" validate:"required"`
}

// Struct untuk respons laporan
type ReportResponse struct {
	ID             uint         `json:"id"`
	UserID         uint         `json:"user_id"`
	TanggalLaporan string       `json:"tanggal_laporan"`
	Location       string       `json:"location"`
	Description    string       `json:"description"`
	Photo          string       `json:"photo"`
	Status         string       `json:"status"`
	Longitude      float64      `json:"longitude"`
	Latitude       float64      `json:"latitude"`
	User           UserResponse `json:"user"`
}

// Struktur untuk respons dari HERE API
type HereGeocodeResponse struct {
	Items []struct {
		Position struct {
			Latitude  float64 `json:"lat"`
			Longitude float64 `json:"lng"`
		} `json:"position"`
	} `json:"items"`
}

// Fungsi untuk mendapatkan koordinat (latitude, longitude) dari alamat menggunakan HERE API
func getCoordinatesFromAddress(address string) (float64, float64, error) {
	apiKey := os.Getenv("HERE_API_KEY")
	baseURL := os.Getenv("HERE_BASE_URL")

	// Membuat URL request dengan parameter alamat dan API key
	requestURL := fmt.Sprintf("%s?q=%s&apiKey=%s", baseURL, url.QueryEscape(address), apiKey)

	// Mengirimkan request ke HERE API
	resp, err := http.Get(requestURL)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to make request to HERE API: %v", err)
	}
	defer resp.Body.Close()

	// Mengecek jika status kode bukan 200 OK
	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("error: received status code %d from HERE API", resp.StatusCode)
	}

	// Parsing JSON response dari HERE API
	var geocodeResponse HereGeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&geocodeResponse); err != nil {
		return 0, 0, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	// Memeriksa apakah ada hasil yang ditemukan
	if len(geocodeResponse.Items) == 0 {
		return 0, 0, fmt.Errorf("no results found for the address")
	}

	// Mendapatkan latitude dan longitude dari hasil pertama
	latitude := geocodeResponse.Items[0].Position.Latitude
	longitude := geocodeResponse.Items[0].Position.Longitude

	return latitude, longitude, nil
}

// Fungsi untuk membuat laporan baru
func CreateReportRubbish(c echo.Context) error {
	var input ReportInput

	// Mengikat data dari form-data
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Invalid input", http.StatusBadRequest, "error", nil))
	}

	// Validasi input
	if err := c.Validate(input); err != nil {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Validation error", http.StatusBadRequest, "error", err.Error()))
	}

	// Mengambil userID dari token yang sudah diset di context
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, helper.APIResponse("Invalid user ID from token", http.StatusUnauthorized, "error", nil))
	}

	// Mengambil file foto dari form-data
	file, err := c.FormFile("photo")
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Failed to retrieve photo file", http.StatusBadRequest, "error", nil))
	}

	// Membuka file foto
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to open photo file", http.StatusInternalServerError, "error", nil))
	}
	defer src.Close()

	// Memeriksa jenis file foto yang di-upload
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Invalid file type, only images are allowed", http.StatusBadRequest, "error", nil))
	}

	// Menginisialisasi Cloudinary untuk upload foto
	cld, err := config.InitCloudinary()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Cloudinary initialization failed", http.StatusInternalServerError, "error", nil))
	}

	// Upload foto ke Cloudinary
	uploadResult, err := cld.Upload.Upload(c.Request().Context(), src, uploader.UploadParams{
		Folder: "report_rubbish", // Menyimpan foto di folder "report_rubbish" di Cloudinary
	})
	if err != nil {
		log.Printf("Cloudinary upload error: %v", err)
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to upload photo", http.StatusInternalServerError, "error", nil))
	}

	// Mendapatkan koordinat (latitude, longitude) dari alamat
	latitude, longitude, err := getCoordinatesFromAddress(input.Location)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to get coordinates from address", http.StatusInternalServerError, "error", err.Error()))
	}

	// Konversi tanggal_laporan string menjadi time.Time
	tanggalLaporan, err := time.Parse("2006-01-02", input.TanggalLaporan)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Invalid date format, must be YYYY-MM-DD", http.StatusBadRequest, "error", nil))
	}

	// Membuat laporan baru
	report := models.ReportRubbish{
		UserID:         userID,
		Location:       input.Location,
		Description:    input.Description,
		Photo:          uploadResult.SecureURL, // Menyimpan URL foto yang diupload
		Status:         "pending",              // Status default
		Longitude:      longitude,
		Latitude:       latitude,
		TanggalLaporan: tanggalLaporan, // Menyimpan tanggal pelaporan dalam format time.Time
	}

	// Menyimpan laporan ke database
	if err := config.DB.Create(&report).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to create report", http.StatusInternalServerError, "error", nil))
	}

	// Mengambil laporan dengan data pengguna
	var reportWithUser models.ReportRubbish
	if err := config.DB.Preload("User").First(&reportWithUser, report.ID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to load report with user", http.StatusInternalServerError, "error", nil))
	}

	// Membentuk respons dengan field baru
	response := ReportResponse{
		ID:             reportWithUser.ID,
		UserID:         reportWithUser.UserID,
		Location:       reportWithUser.Location,
		Description:    reportWithUser.Description,
		Photo:          reportWithUser.Photo,
		Status:         reportWithUser.Status,
		Longitude:      reportWithUser.Longitude,
		Latitude:       reportWithUser.Latitude,
		TanggalLaporan: reportWithUser.TanggalLaporan.Format("2006-01-02"), // Format tanggal untuk respons
		User: UserResponse{
			IDUser:       reportWithUser.User.ID,
			NamaLengkap:  reportWithUser.User.NamaLengkap,
			TanggalLahir: reportWithUser.User.TanggalLahir.Format("2006-01-02"),
			NoTelepon:    reportWithUser.User.NoTelepon,
			Email:        reportWithUser.User.Email,
			Role:         reportWithUser.User.Role,
			Photo:        reportWithUser.User.Photo,
		},
	}

	// Mengembalikan respons sukses
	return c.JSON(http.StatusOK, helper.APIResponse("Report created successfully", http.StatusOK, "success", response))
}

func UpdateReportStatus(c echo.Context) error {
	// Ambil ID laporan dari parameter
	id := c.Param("id")

	// Struktur untuk menerima input JSON
	input := struct {
		Status string `json:"status" validate:"required,oneof=pending approved rejected"`
	}{}

	// Parse dan validasi input
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Invalid input format", http.StatusBadRequest, "error", nil))
	}

	if err := c.Validate(&input); err != nil {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Validation error: "+err.Error(), http.StatusBadRequest, "error", nil))
	}

	// Akses database untuk memperbarui status
	db := config.DB                  // Pastikan `config.DB` adalah instance database yang diinisialisasi sebelumnya
	report := models.ReportRubbish{} // Pastikan model `Report` sesuai dengan database Anda

	// Cari laporan berdasarkan ID
	if err := db.First(&report, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, helper.APIResponse("Report not found", http.StatusNotFound, "error", nil))
	}

	// Perbarui status laporan
	report.Status = input.Status
	if err := db.Save(&report).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to update report status", http.StatusInternalServerError, "error", nil))
	}

	// Berikan respons sukses
	return c.JSON(http.StatusOK, helper.APIResponse("Report status updated successfully", http.StatusOK, "success", map[string]interface{}{
		"id":     report.ID,
		"status": report.Status,
	}))
}

func GetAllReportRubbish(c echo.Context) error {
	// Mendapatkan semua laporan dari database
	var reports []models.ReportRubbish
	if err := config.DB.Preload("User").Find(&reports).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to load reports", http.StatusInternalServerError, "error", nil))
	}

	// Menyiapkan response
	var reportResponses []ReportResponse
	for _, report := range reports {
		reportResponses = append(reportResponses, ReportResponse{
			ID:             report.ID,
			UserID:         report.UserID,
			TanggalLaporan: report.TanggalLaporan.Format("2006-01-02"),
			Location:       report.Location,
			Description:    report.Description,
			Photo:          report.Photo,
			Status:         report.Status,
			Longitude:      report.Longitude,
			Latitude:       report.Latitude,
			User: UserResponse{
				IDUser:       report.User.ID,
				NamaLengkap:  report.User.NamaLengkap,
				TanggalLahir: report.User.TanggalLahir.Format("2006-01-02"), // Format sesuai kebutuhan
				NoTelepon:    report.User.NoTelepon,
				Email:        report.User.Email,
				Role:         report.User.Role,
				Photo:        report.User.Photo,
			},
		})
	}

	// Mengembalikan response dengan semua laporan
	return c.JSON(http.StatusOK, helper.APIResponse("Reports retrieved successfully", http.StatusOK, "success", reportResponses))
}