package controllers

import (
	"Backend-Recything/config"
	"Backend-Recything/helper"
	"Backend-Recything/models"
	"net/http"
	"os"
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
}

// Struct untuk validasi input login
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// Struct untuk validasi input registrasi
type RegisterInput struct {
	NamaLengkap  string `json:"nama_lengkap" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=6"`
	TanggalLahir string `json:"tanggal_lahir" validate:"required"`
	NoTelepon    string `json:"no_telepon" validate:"required"`
}

// Struct untuk JWT Claims
type jwtCustomClaims struct {
	Name   string `json:"name"`
	UserID uint   `json:"userID"`
	jwt.RegisteredClaims
}

// LoginHandler menangani proses login
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
	token, err := GenerateJWT(user.ID, user.NamaLengkap)
	if err != nil {
		response := helper.APIResponse("Failed to generate token", http.StatusInternalServerError, "error", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Response data
	data := LoginResponseData{
		IDUser:       user.ID,
		NamaLengkap:  user.NamaLengkap,
		Email:        user.Email,
		NoTelepon:    user.NoTelepon,                         // Menambahkan NoTelepon di response login
		TanggalLahir: user.TanggalLahir.Format("2006-01-02"), // Menyertakan TanggalLahir jika perlu
		Token:        token,
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

	// Validasi input
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		response := helper.APIResponse("Validation error", http.StatusBadRequest, "error", err.Error())
		return c.JSON(http.StatusBadRequest, response)
	}

	// Hash password
	hash, _ := HashPassword(input.Password)

	// Parse TanggalLahir ke time.Time
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
	}

	// Simpan ke database
	result := config.DB.Create(&user)
	if result.Error != nil {
		response := helper.APIResponse("Failed to register", http.StatusInternalServerError, "error", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Generate token JWT
	token, err := GenerateJWT(user.ID, user.NamaLengkap)
	if err != nil {
		response := helper.APIResponse("Failed to generate token", http.StatusInternalServerError, "error", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Response data
	data := LoginResponseData{
		IDUser:       user.ID,
		NamaLengkap:  user.NamaLengkap,
		Email:        user.Email,
		NoTelepon:    user.NoTelepon,                         // Pastikan NoTelepon ada dalam response
		TanggalLahir: user.TanggalLahir.Format("2006-01-02"), // Format tanggal jika diperlukan
		Token:        token,
	}

	response := helper.APIResponse("Registration successful", http.StatusOK, "success", data)
	return c.JSON(http.StatusOK, response)
}

// GenerateJWT membuat token JWT
func GenerateJWT(userID uint, name string) (string, error) {
	claims := &jwtCustomClaims{
		Name:   name,
		UserID: userID,
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
