package controllers

import (
	"Backend-Recything/config"
	"Backend-Recything/helper"
	"Backend-Recything/models"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

	// Membuat user baru
	user := models.User{
		NamaLengkap:  input.NamaLengkap,
		Email:        input.Email,
		NoTelepon:    input.NoTelepon,
		Password:     hash,
		TanggalLahir: tanggalLahir,
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
	var users []models.User
	result := config.DB.Find(&users)
	if result.Error != nil {
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

	response := helper.APIResponse("Users retrieved successfully", http.StatusOK, "success", userResponses)
	return c.JSON(http.StatusOK, response)
}

// GetUserByID mengembalikan data pengguna berdasarkan ID
func GetUserByID(c echo.Context) error {
	id := c.Param("id")

	var user models.User
	result := config.DB.First(&user, "id = ?", id)
	if result.Error != nil || user.ID == 0 {
		response := helper.APIResponse("User not found", http.StatusNotFound, "error", nil)
		return c.JSON(http.StatusNotFound, response)
	}

	// Format data untuk respons
	userResponse := UserResponse{
		IDUser:       user.ID,
		NamaLengkap:  user.NamaLengkap,
		TanggalLahir: user.TanggalLahir.Format("2006-01-02"),
		NoTelepon:    user.NoTelepon,
		Email:        user.Email,
		Role:         user.Role,
		Photo:        user.Photo, // Menambahkan photo ke respons
	}

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
func UpdatePhotoHandler(c echo.Context) error {
	id := c.Param("id")

	// Mengambil file yang diupload
	file, err := c.FormFile("photo")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "No file uploaded",
		})
	}

	// Validasi tipe file
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid file type. Only images are allowed.",
		})
	}

	// Membuka file upload
	src, err := file.Open()
	if err != nil {
		log.Printf("Failed to open file: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to open uploaded file",
		})
	}
	defer src.Close()

	// Path absolut untuk menyimpan file
	uploadDir := "/var/www/recythingtech/uploads/"
	photoPath := uploadDir + id + "_" + file.Filename

	// Membuka file tujuan
	dst, err := os.Create(photoPath)
	if err != nil {
		log.Printf("Failed to create file: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to save uploaded file",
		})
	}
	defer dst.Close()

	// Menyalin data file
	if _, err := io.Copy(dst, src); err != nil {
		log.Printf("Failed to copy file: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to save uploaded file",
		})
	}

	// Cari pengguna berdasarkan ID
	var user models.User
	result := config.DB.First(&user, id)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "User not found",
		})
	}

	// Update path foto di database
	user.Photo = id + "_" + file.Filename
	result = config.DB.Save(&user)
	if result.Error != nil {
		log.Printf("Failed to update database: %v\n", result.Error)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to update photo path in database",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Photo uploaded and updated successfully",
	})
}
