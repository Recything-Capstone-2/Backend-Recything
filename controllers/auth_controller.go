package controllers

import (
	"Backend-Recything/config"
	"Backend-Recything/helper"
	"Backend-Recything/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Struct untuk response login
type LoginResponseData struct {
	IDUser       uint   `json:"id_user"`
	NamaLengkap  string `json:"nama_lengkap"`
	TanggalLahir string `json:"tanggal_lahir"`
	NoTelepon    string `json:"no_telepon"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	Role         string `json:"role"`
	Photo        string `json:"photo"`
}

// Struct untuk validasi input login
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
	IDUser       uint   `json:"id_user"`
	NamaLengkap  string `json:"nama_lengkap"`
	TanggalLahir string `json:"tanggal_lahir"`
	NoTelepon    string `json:"no_telepon"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Photo        string `json:"photo"`
}

// Struct untuk validasi input registrasi
type RegisterInput struct {
	NamaLengkap  string `json:"nama_lengkap" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=6"`
	TanggalLahir string `json:"tanggal_lahir" validate:"required"`
	NoTelepon    string `json:"no_telepon" validate:"required"`
	Role         string `json:"role" validate:"oneof=admin user"` // Validasi untuk admin/user
	Photo        string `json:"photo" validate:"required,url"`
}

type RegisterResponse struct {
	IDUser       uint   `json:"id_user"`
	NamaLengkap  string `json:"nama_lengkap"`
	TanggalLahir string `json:"tanggal_lahir"`
	NoTelepon    string `json:"no_telepon"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Photo        string `json:"photo"`
}

type UpdateUserDataInput struct {
	NamaLengkap  string `json:"nama_lengkap" validate:"required"`
	TanggalLahir string `json:"tanggal_lahir" validate:"required"`
	NoTelepon    string `json:"no_telepon" validate:"required"`
	Email        string `json:"email" validate:"email"`
	OldPassword  string `json:"old_password"` // Optional: Password lama
	NewPassword  string `json:"new_password"` // Optional: Password baru
}

// Struct untuk JWT Claims
type jwtCustomClaims struct {
	Name   string `json:"name"`
	UserID uint   `json:"userID"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// LoginHandler menangani proses login
func LoginHandler(c echo.Context) error {
	var input LoginInput
	if err := c.Bind(&input); err != nil {
		response := helper.APIResponse("Invalid request", http.StatusBadRequest, "error", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	// Validasi input
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		response := helper.APIResponse("Validation error", http.StatusBadRequest, "error", err.Error())
		return c.JSON(http.StatusBadRequest, response)
	}

	// Cari user berdasarkan email
	var user models.User
	result := config.DB.First(&user, "email = ?", input.Email)
	if result.Error != nil || user.ID == 0 {
		response := helper.APIResponse("Email not found", http.StatusUnauthorized, "error", nil)
		return c.JSON(http.StatusUnauthorized, response)
	}

	// Cek password
	if !CheckPasswordHash(input.Password, user.Password) {
		response := helper.APIResponse("Incorrect password", http.StatusUnauthorized, "error", nil)
		return c.JSON(http.StatusUnauthorized, response)
	}

	// Generate token JWT
	token, err := GenerateJWT(user.ID, user.NamaLengkap, user.Role)
	if err != nil {
		response := helper.APIResponse("Failed to generate token", http.StatusInternalServerError, "error", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Response data dengan field photo
	data := LoginResponseData{
		IDUser:       user.ID,
		NamaLengkap:  user.NamaLengkap,
		Email:        user.Email,
		NoTelepon:    user.NoTelepon,
		TanggalLahir: user.TanggalLahir.Format("2006-01-02"),
		Token:        token,
		Role:         user.Role,
		Photo:        user.Photo, // Tambahkan photo ke respons
	}

	response := helper.APIResponse("Login successful", http.StatusOK, "success", data)
	return c.JSON(http.StatusOK, response)
}

// RegisterHandler menangani proses registrasi
// RegisterHandler menangani proses registrasi
func RegisterHandler(c echo.Context) error {
	var input RegisterInput
	if err := c.Bind(&input); err != nil {
		response := helper.APIResponse("Invalid request", http.StatusBadRequest, "error", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	// Jika role kosong, isi default menjadi "user"
	if input.Role == "" {
		input.Role = "user"
	}

	// Validasi input
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		response := helper.APIResponse("Validation error", http.StatusBadRequest, "error", err.Error())
		return c.JSON(http.StatusBadRequest, response)
	}

	// Hash password
	hash, err := HashPassword(input.Password)
	if err != nil {
		response := helper.APIResponse("Failed to hash password", http.StatusInternalServerError, "error", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Parse tanggal lahir
	tanggalLahir, err := time.Parse("2006-01-02", input.TanggalLahir)
	if err != nil {
		response := helper.APIResponse("Invalid birth date format", http.StatusBadRequest, "error", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	// Set TanggalLahir to nil if it's zero value
	if tanggalLahir.IsZero() {
		tanggalLahir = time.Time{} // Set to zero value (January 1, 0001)
	}

	// Membuat user baru
	user := models.User{
		NamaLengkap:  input.NamaLengkap,
		Email:        input.Email,
		NoTelepon:    input.NoTelepon,
		Password:     hash,
		TanggalLahir: tanggalLahir, // Assign the parsed date
		Role:         input.Role,
		Photo:        input.Photo, // Simpan URL foto langsung
	}

	// Simpan ke database
	result := config.DB.Create(&user)
	if result.Error != nil {
		response := helper.APIResponse("Failed to register", http.StatusInternalServerError, "error", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Format respons
	responseData := RegisterResponse{
		IDUser:       user.ID,
		NamaLengkap:  user.NamaLengkap,
		TanggalLahir: user.TanggalLahir.Format("2006-01-02"),
		NoTelepon:    user.NoTelepon,
		Email:        user.Email,
		Role:         user.Role,
		Photo:        user.Photo,
	}

	response := helper.APIResponse("Registration successful", http.StatusOK, "success", responseData)
	return c.JSON(http.StatusOK, response)
}

// GetAllUsers mengembalikan daftar semua pengguna
func GetAllUsers(c echo.Context) error {
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

	// Query database untuk menghitung total data
	var totalItems int64
	if err := config.DB.Model(&models.User{}).Count(&totalItems).Error; err != nil {
		response := helper.APIResponse("Failed to count users", http.StatusInternalServerError, "error", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Query database untuk mengambil data dengan paginasi
	var users []models.User
	if err := config.DB.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		response := helper.APIResponse("Failed to retrieve users", http.StatusInternalServerError, "error", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Format data untuk menghapus field yang tidak diperlukan
	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			IDUser:       user.ID,
			NamaLengkap:  user.NamaLengkap,
			TanggalLahir: user.TanggalLahir.Format("2006-01-02"),
			NoTelepon:    user.NoTelepon,
			Email:        user.Email,
			Role:         user.Role,
			Photo:        user.Photo, // Menambahkan photo ke respons
		})
	}

	// Hitung total halaman
	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))

	// Buat respons dengan data dan informasi paginasi
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"items": userResponses,
			"pagination": map[string]interface{}{
				"current_page": page,
				"per_page":     limit,
				"total_users":  totalItems,
				"total_pages":  totalPages,
			},
		},
		"error": nil,
	}

	return c.JSON(http.StatusOK, helper.APIResponse("Users retrieved successfully", http.StatusOK, "success", response))
}

// GetUserByID mengembalikan data pengguna berdasarkan ID
func GetUserByID(c echo.Context) error {
	id := c.Param("id")

	// Cari user berdasarkan ID
	var user models.User
	result := config.DB.First(&user, "id = ?", id)
	if result.Error != nil || user.ID == 0 {
		// User not found
		response := helper.APIResponse("User not found", http.StatusNotFound, "error", nil)
		return c.JSON(http.StatusNotFound, response)
	}

	// Ambil laporan untuk user
	var reports []models.ReportRubbish

	// Ambil semua laporan untuk user tanpa pagination
	config.DB.Where("user_id = ?", id).Find(&reports)

	// Format laporan ke dalam respons
	var reportsResponse []ReportResponse
	for _, report := range reports {
		// Append each report including its status
		reportsResponse = append(reportsResponse, ReportResponse{
			ID:             report.ID,
			UserID:         report.UserID,
			Category:       report.Category,
			TanggalLaporan: report.TanggalLaporan.Format("2006-01-02"),
			Location:       report.Location,
			Description:    report.Description,
			Photo:          report.Photo,
			Status:         report.Status, // This field is now included in the response
			Longitude:      report.Longitude,
			Latitude:       report.Latitude,
		})
	}

	// Format data untuk respons
	userResponse := struct {
		IDUser       uint             `json:"id_user"`
		NamaLengkap  string           `json:"nama_lengkap"`
		TanggalLahir string           `json:"tanggal_lahir"`
		NoTelepon    string           `json:"no_telepon"`
		Email        string           `json:"email"`
		Role         string           `json:"role"`
		Photo        string           `json:"photo"`
		Reports      []ReportResponse `json:"reports"` // Include reports here
	}{
		IDUser:       user.ID,
		NamaLengkap:  user.NamaLengkap,
		TanggalLahir: user.TanggalLahir.Format("2006-01-02"),
		NoTelepon:    user.NoTelepon,
		Email:        user.Email,
		Role:         user.Role,
		Photo:        user.Photo,
		Reports:      reportsResponse,
	}

	// Response berhasil
	response := helper.APIResponse("User retrieved successfully", http.StatusOK, "success", userResponse)
	return c.JSON(http.StatusOK, response)
}

// GenerateJWT membuat token JWT
func GenerateJWT(userID uint, name string, role string) (string, error) {
	claims := &jwtCustomClaims{
		Name:   name,
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

// HashPassword mengenkripsi password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash mencocokkan password dengan hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// UpdatePhotoHandler untuk menangani update foto pengguna
func UpdateUserPhoto(c echo.Context) error {
	// Ambil user ID dari parameter
	id := c.Param("id")

	// Ambil file dari request
	file, err := c.FormFile("photo")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Failed to retrieve photo file",
		})
	}

	// Buka file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to open photo file",
		})
	}
	defer src.Close()

	// Upload ke Cloudinary
	cld, err := config.InitCloudinary()
	if err != nil {
		log.Printf("Cloudinary initialization error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Cloudinary initialization failed",
		})
	}

	// Upload file ke Cloudinary
	uploadResult, err := cld.Upload.Upload(c.Request().Context(), src, uploader.UploadParams{
		Folder: "user_photos", // Menyimpan foto dalam folder khusus
	})
	if err != nil {
		log.Printf("Cloudinary upload error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to upload photo to Cloudinary",
		})
	}

	// Update foto pengguna di database
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "User not found",
		})
	}

	user.Photo = uploadResult.SecureURL
	if err := config.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to save user photo",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "User photo updated successfully",
		"photo":   user.Photo,
	})
}

func Logout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Berhasil Logout",
	})
}

// Fungsi untuk mengupdate data diri dan mengganti password pengguna
func UpdateUserData(c echo.Context) error {
	// Ambil ID user dari parameter
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32) // Konversi id (string) ke uint
	if err != nil {
		response := helper.APIResponse("Invalid user ID", http.StatusBadRequest, "error", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	// Ambil userID dari JWT token
	tokenUserID := c.Get("userID").(uint)
	if tokenUserID != uint(userID) {
		response := helper.APIResponse("You can only update your own data", http.StatusForbidden, "error", nil)
		return c.JSON(http.StatusForbidden, response)
	}

	// Bind input JSON ke dalam struct
	var input UpdateUserDataInput
	if err := c.Bind(&input); err != nil {
		response := helper.APIResponse("Invalid request", http.StatusBadRequest, "error", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	// Cari user berdasarkan ID
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		response := helper.APIResponse("User not found", http.StatusNotFound, "error", nil)
		return c.JSON(http.StatusNotFound, response)
	}

	// Jika password lama ada, maka verifikasi dan update password
	if input.OldPassword != "" && input.NewPassword != "" {
		if !CheckPasswordHash(input.OldPassword, user.Password) {
			response := helper.APIResponse("Old password is incorrect", http.StatusUnauthorized, "error", nil)
			return c.JSON(http.StatusUnauthorized, response)
		}
		hashedNewPassword, err := HashPassword(input.NewPassword)
		if err != nil {
			response := helper.APIResponse("Failed to hash new password", http.StatusInternalServerError, "error", nil)
			return c.JSON(http.StatusInternalServerError, response)
		}
		user.Password = hashedNewPassword
	}

	// Periksa setiap field input apakah diisi, jika iya, baru update
	if input.NamaLengkap != "" {
		user.NamaLengkap = input.NamaLengkap
	}

	if input.TanggalLahir != "" {
		parsedDate, err := time.Parse("2006-01-02", input.TanggalLahir)
		if err != nil {
			response := helper.APIResponse("Invalid date format", http.StatusBadRequest, "error", nil)
			return c.JSON(http.StatusBadRequest, response)
		}
		user.TanggalLahir = parsedDate
	}

	if input.NoTelepon != "" {
		user.NoTelepon = input.NoTelepon
	}

	if input.Email != "" {
		user.Email = input.Email
	}

	// Simpan perubahan ke database
	if err := config.DB.Save(&user).Error; err != nil {
		response := helper.APIResponse("Failed to update user data", http.StatusInternalServerError, "error", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Format response data
	userResponse := UserResponse{
		IDUser:       user.ID,
		NamaLengkap:  user.NamaLengkap,
		TanggalLahir: user.TanggalLahir.Format("2006-01-02"),
		NoTelepon:    user.NoTelepon,
		Email:        user.Email,
		Role:         user.Role,
		Photo:        user.Photo,
	}

	response := helper.APIResponse("User data updated successfully", http.StatusOK, "success", userResponse)
	return c.JSON(http.StatusOK, response)
}
