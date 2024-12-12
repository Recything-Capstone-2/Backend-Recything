package controllers

import (
	"Backend-Recything/config"
	"Backend-Recything/helper"
	"Backend-Recything/models"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"encoding/json"
	"net/url"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/labstack/echo/v4"
)

// Struct untuk input laporan
type ReportInput struct {
	Category       string `form:"category" validate:"required,oneof=report_rubbish report_littering"`
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
	Category       string       `json:"category"`
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

	// Bind input data from the request
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Invalid input", http.StatusBadRequest, "error", nil))
	}

	// Get user ID from context
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, helper.APIResponse("Invalid user ID from token", http.StatusUnauthorized, "error", nil))
	}

	// Handle photo upload if provided
	var photoURL string
	file, err := c.FormFile("photo")
	if file != nil {
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to open photo file", http.StatusInternalServerError, "error", nil))
		}
		defer src.Close()

		if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
			return c.JSON(http.StatusBadRequest, helper.APIResponse("Invalid file type, only images are allowed", http.StatusBadRequest, "error", nil))
		}

		cld, err := config.InitCloudinary()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.APIResponse("Cloudinary initialization failed", http.StatusInternalServerError, "error", nil))
		}

		uploadResult, err := cld.Upload.Upload(c.Request().Context(), src, uploader.UploadParams{
			Folder: "report_rubbish",
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to upload photo", http.StatusInternalServerError, "error", nil))
		}
		photoURL = uploadResult.SecureURL
	}

	// Set status to "process" if required fields are provided
	status := "rejected"
	if input.Location != "" && input.Description != "" && photoURL != "" && input.TanggalLaporan != "" {
		status = "process"
	}

	// Get coordinates if location is provided
	var latitude, longitude float64
	if input.Location != "" {
		latitude, longitude, err = getCoordinatesFromAddress(input.Location)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to get coordinates from address", http.StatusInternalServerError, "error", err.Error()))
		}
	}

	// Parse TanggalLaporan into time.Time
	tanggalLaporan, err := time.Parse("2006-01-02", input.TanggalLaporan)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Invalid date format. Please use YYYY-MM-DD.", http.StatusBadRequest, "error", nil))
	}

	// Create the report
	report := models.ReportRubbish{
		UserID:         userID,
		Category:       input.Category,
		Location:       input.Location,
		Description:    input.Description,
		Photo:          photoURL,
		Status:         status,
		Longitude:      longitude,
		Latitude:       latitude,
		TanggalLaporan: tanggalLaporan, // Store as time.Time
	}

	// Save the report to the database
	if err := config.DB.Create(&report).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to create report", http.StatusInternalServerError, "error", nil))
	}

	// Load the report with associated user data
	var reportWithUser models.ReportRubbish
	if err := config.DB.Preload("User").First(&reportWithUser, report.ID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to load report with user", http.StatusInternalServerError, "error", nil))
	}

	// Prepare the response
	response := ReportResponse{
		ID:             reportWithUser.ID,
		UserID:         report.UserID,
		Category:       report.Category,
		TanggalLaporan: report.TanggalLaporan.Format("2006-01-02"), // Return as formatted string
		Location:       report.Location,
		Description:    report.Description,
		Photo:          report.Photo,
		Status:         report.Status,
		Longitude:      report.Longitude,
		Latitude:       report.Latitude,
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

	// Return success response
	return c.JSON(http.StatusOK, helper.APIResponse("Report created successfully", http.StatusOK, "success", response))
}

// Fungsi untuk memperbarui status laporan dan memberikan poin jika laporan disetujui
func UpdateReportStatus(c echo.Context) error {
	// Mendapatkan ID laporan dari parameter URL
	id := c.Param("id")

	// Mendapatkan input status
	input := struct {
		Status string `json:"status" validate:"required,oneof=pending approved rejected"`
	}{}

	// Bind dan validasi input
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Invalid input format", http.StatusBadRequest, "error", nil))
	}

	// Ambil laporan berdasarkan ID
	var report models.ReportRubbish
	if err := config.DB.First(&report, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, helper.APIResponse("Report not found", http.StatusNotFound, "error", nil))
	}

	// Update status laporan
	report.Status = input.Status
	if err := config.DB.Save(&report).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to update report status", http.StatusInternalServerError, "error", nil))
	}

	// Jika status laporan adalah "approved", beri poin ke user
	if report.Status == "approved" {
		var user models.User
		if err := config.DB.First(&user, report.UserID).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to find user", http.StatusInternalServerError, "error", nil))
		}

		// Poin yang akan diberikan
		points := uint(1000)

		// Tambahkan poin ke user
		user.Points += points
		if err := config.DB.Save(&user).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to update user points", http.StatusInternalServerError, "error", nil))
		}

		// Update atau buat data poin di tabel Points
		var userPoints models.Points
		err := config.DB.Where("user_id = ?", user.ID).First(&userPoints).Error
		if err != nil {
			// Jika tidak ada data poin, buat data baru
			if err := config.DB.Create(&models.Points{
				UserID: user.ID,
				Points: points,
			}).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to add points", http.StatusInternalServerError, "error", nil))
			}
		} else {
			// Jika data poin sudah ada, update
			userPoints.Points += points
			if err := config.DB.Save(&userPoints).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to update points", http.StatusInternalServerError, "error", nil))
			}
		}
	}

	// Siapkan respons dengan metadata dan data yang relevan
	responseData := struct {
		ID     uint   `json:"id"`
		Status string `json:"status"`
		UserID uint   `json:"user_id"`
	}{
		ID:     report.ID,
		Status: report.Status,
		UserID: report.UserID,
	}

	// Mengembalikan respons sukses dengan metadata dan data yang relevan
	response := helper.APIResponse("Report status updated successfully", http.StatusOK, "success", responseData)
	return c.JSON(http.StatusOK, response)
}

func GetAllReportRubbish(c echo.Context) error {
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

	// Ambil kategori, status, dan sorting dari query parameter
	category := c.QueryParam("category")
	status := c.QueryParam("status") // New status filter
	sortOrder := c.QueryParam("sort")

	// Validasi nilai sorting (default: "asc" jika tidak diisi)
	if sortOrder != "asc" && sortOrder != "desc" && sortOrder != "" {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Invalid sort order. Use 'asc' or 'desc'.", http.StatusBadRequest, "error", nil))
	}
	if sortOrder == "" {
		sortOrder = "asc" // Default sorting: ascending
	}

	// Query database dengan filter kategori dan status
	var totalItems int64
	db := config.DB.Model(&models.ReportRubbish{})

	// Filter by category
	if category != "" {
		db = db.Where("category = ?", category)
	}

	// Filter by status if provided
	if status != "" {
		db = db.Where("status = ?", status)
	}

	// Hitung total data sebelum paginasi
	if err := db.Count(&totalItems).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to count reports", http.StatusInternalServerError, "error", nil))
	}

	// Ambil data dengan paginasi
	var reports []models.ReportRubbish
	if err := db.Order(fmt.Sprintf("tanggal_laporan %s", sortOrder)).
		Offset(offset).Limit(limit).Preload("User").Find(&reports).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to load reports", http.StatusInternalServerError, "error", nil))
	}

	// Mapping hasil ke response
	var reportResponses []ReportResponse
	for _, report := range reports {
		reportResponses = append(reportResponses, ReportResponse{
			ID:             report.ID,
			UserID:         report.UserID,
			Category:       report.Category,
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
				TanggalLahir: report.User.TanggalLahir.Format("2006-01-02"),
				NoTelepon:    report.User.NoTelepon,
				Email:        report.User.Email,
				Role:         report.User.Role,
				Photo:        report.User.Photo,
			},
		})
	}

	// Hitung total halaman
	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))

	// Respons yang diinginkan
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"items": reportResponses,
			"pagination": map[string]interface{}{
				"current_page": page,
				"per_page":     limit,
				"total_report": totalItems,
				"total_pages":  totalPages,
			},
		},
		"error": nil,
	}

	return c.JSON(http.StatusOK, helper.APIResponse("Report history retrieved successfully", http.StatusOK, "success", response))
}

func GetReportHistoryByUser(c echo.Context) error {
	// Mendapatkan userID dari konteks (token JWT)
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, helper.APIResponse("Invalid user ID from token", http.StatusUnauthorized, "error", nil))
	}

	// Query laporan berdasarkan userID
	var reports []models.ReportRubbish
	if err := config.DB.Preload("User").Where("user_id = ?", userID).Find(&reports).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to retrieve report history", http.StatusInternalServerError, "error", nil))
	}

	// Mapping hasil query ke struktur response
	var reportResponses []ReportResponse
	for _, report := range reports {
		reportResponses = append(reportResponses, ReportResponse{
			ID:             report.ID,
			UserID:         report.UserID,
			Category:       report.Category,
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
				TanggalLahir: report.User.TanggalLahir.Format("2006-01-02"),
				NoTelepon:    report.User.NoTelepon,
				Email:        report.User.Email,
				Role:         report.User.Role,
				Photo:        report.User.Photo,
			},
		})
	}

	// Mengembalikan respons sukses
	return c.JSON(http.StatusOK, helper.APIResponse("Report history retrieved successfully", http.StatusOK, "success", reportResponses))
}

func AddPointsToUser(c echo.Context) error {
	// Ambil userID dari token (misalnya menggunakan JWT)
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "Invalid user ID from token")
	}

	// Ambil jumlah poin yang ingin ditambahkan (misalnya dari input request)
	points := uint(1000) // Ganti dengan poin yang Anda inginkan

	// Dapatkan user berdasarkan userID
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	// Tambahkan poin ke user
	user.Points += points

	// Simpan perubahan pada user
	if err := config.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to update user points")
	}

	// Cek apakah data poin untuk user sudah ada
	var userPoints models.Points
	err := config.DB.Where("user_id = ?", user.ID).First(&userPoints).Error
	if err != nil {
		// Jika tidak ada, buat data poin baru
		if err := config.DB.Create(&models.Points{
			UserID: user.ID,
			Points: points,
		}).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to create points record")
		}
	} else {
		// Jika sudah ada, update jumlah poin
		userPoints.Points += points
		if err := config.DB.Save(&userPoints).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to update points")
		}
	}

	return c.JSON(http.StatusOK, "Points added successfully")
}

func GetLatestReports(c echo.Context) error {
	// Inisialisasi slice untuk menampung hasil query
	var reports []models.ReportRubbish

	// Query database: Ambil 10 laporan terbaru berdasarkan TanggalLaporan
	if err := config.DB.Order("tanggal_laporan DESC").Limit(10).Preload("User").Find(&reports).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to retrieve latest reports", http.StatusInternalServerError, "error", nil))
	}

	// Mapping hasil query ke response
	var reportResponses []ReportResponse
	for _, report := range reports {
		reportResponses = append(reportResponses, ReportResponse{
			ID:             report.ID,
			UserID:         report.UserID,
			Category:       report.Category,
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
				TanggalLahir: report.User.TanggalLahir.Format("2006-01-02"),
				NoTelepon:    report.User.NoTelepon,
				Email:        report.User.Email,
				Role:         report.User.Role,
				Photo:        report.User.Photo,
			},
		})
	}

	// Return response
	return c.JSON(http.StatusOK, helper.APIResponse("Latest reports retrieved successfully", http.StatusOK, "success", reportResponses))
}

func DeductPointsFromUser(c echo.Context) error {
	// Ambil input dari request
	input := struct {
		UserID uint `json:"user_id" validate:"required"`
		Points int  `json:"points" validate:"required,min=1"` // Pastikan poin yang dikurangi minimal 1
	}{}

	// Bind dan validasi input
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Invalid input format", http.StatusBadRequest, "error", nil))
	}

	// Cari pengguna berdasarkan ID
	var user models.User
	if err := config.DB.First(&user, input.UserID).Error; err != nil {
		return c.JSON(http.StatusNotFound, helper.APIResponse("User not found", http.StatusNotFound, "error", nil))
	}

	// Pastikan poin tidak menjadi negatif
	if int(user.Points) < input.Points {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Insufficient points", http.StatusBadRequest, "error", nil))
	}

	// Kurangi poin pengguna
	user.Points -= uint(input.Points)
	if err := config.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to update user points", http.StatusInternalServerError, "error", nil))
	}

	// Update atau buat record di tabel Points
	var userPoints models.Points
	err := config.DB.Where("user_id = ?", user.ID).First(&userPoints).Error
	if err != nil {
		// Jika tidak ada record, buat baru
		if err := config.DB.Create(&models.Points{
			UserID: user.ID,
			Points: uint(user.Points), // Simpan poin terbaru
		}).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to update points record", http.StatusInternalServerError, "error", nil))
		}
	} else {
		// Jika ada, perbarui
		userPoints.Points = uint(user.Points)
		if err := config.DB.Save(&userPoints).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to save points record", http.StatusInternalServerError, "error", nil))
		}
	}

	// Siapkan respons dengan tambahan nama dan email
	responseData := struct {
		UserID          uint   `json:"user_id"`
		NamaLengkap     string `json:"nama_lengkap"`
		Email           string `json:"email"`
		NoTelepone      string `json:"no_telepone"`
		RemainingPoints uint   `json:"remaining_points"`
	}{
		UserID:          user.ID,
		NamaLengkap:     user.NamaLengkap,
		Email:           user.Email,
		NoTelepone:      user.NoTelepon,
		RemainingPoints: user.Points,
	}

	return c.JSON(http.StatusOK, helper.APIResponse("Points deducted successfully", http.StatusOK, "success", responseData))
}

func DeleteReportByID(c echo.Context) error {
	// Mendapatkan ID laporan dari parameter URL
	id := c.Param("id")

	// Validasi apakah ID valid (berbentuk angka)
	reportID, err := strconv.Atoi(id)
	if err != nil || reportID <= 0 {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Invalid report ID", http.StatusBadRequest, "error", nil))
	}

	// Ambil laporan berdasarkan ID
	var report models.ReportRubbish
	if err := config.DB.First(&report, reportID).Error; err != nil {
		if err.Error() == "record not found" {
			return c.JSON(http.StatusNotFound, helper.APIResponse("Report not found", http.StatusNotFound, "error", nil))
		}
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to retrieve report", http.StatusInternalServerError, "error", nil))
	}

	// Hapus laporan dari database
	if err := config.DB.Delete(&report).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to delete report", http.StatusInternalServerError, "error", nil))
	}

	// Kembalikan respons sukses
	return c.JSON(http.StatusOK, helper.APIResponse("Report deleted successfully", http.StatusOK, "success", nil))
}
func GetReportByID(c echo.Context) error {
	// Mendapatkan ID laporan dari parameter URL
	id := c.Param("id")

	// Validasi apakah ID valid (berbentuk angka)
	reportID, err := strconv.Atoi(id)
	if err != nil || reportID <= 0 {
		return c.JSON(http.StatusBadRequest, helper.APIResponse("Invalid report ID", http.StatusBadRequest, "error", nil))
	}

	// Ambil laporan berdasarkan ID
	var report models.ReportRubbish
	if err := config.DB.Preload("User").First(&report, reportID).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return c.JSON(http.StatusNotFound, helper.APIResponse("Report not found", http.StatusNotFound, "error", nil))
		}
		return c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to retrieve report", http.StatusInternalServerError, "error", nil))
	}

	// Mapping hasil ke response
	reportResponse := ReportResponse{
		ID:             report.ID,
		UserID:         report.UserID,
		Category:       report.Category,
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
			TanggalLahir: report.User.TanggalLahir.Format("2006-01-02"),
			NoTelepon:    report.User.NoTelepon,
			Email:        report.User.Email,
			Role:         report.User.Role,
			Photo:        report.User.Photo,
		},
	}

	// Kembalikan respons sukses
	return c.JSON(http.StatusOK, helper.APIResponse("Report retrieved successfully", http.StatusOK, "success", reportResponse))
}
