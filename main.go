package main

import (
	"Backend-Recything/config"
	"Backend-Recything/controllers"
	"Backend-Recything/middlewares"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Memuat file .env
	loadEnv()

	// Inisialisasi database
	config.InitDB()

	// Inisialisasi Echo
	e := echo.New()

	// Set validator untuk Echo
	e.Validator = &middlewares.CustomValidator{Validator: validator.New()}

	// Middleware global
	e.Use(middleware.Logger())  // Logging request dan response
	e.Use(middleware.Recover()) // Menangani panic agar server tidak crash
	e.Use(middleware.CORS())    // Mendukung CORS untuk client-side apps

	// Rute publik (tanpa autentikasi)
	publicRoutes(e)

	// Rute dengan autentikasi
	protectedRoutes(e)

	// Menjalankan server pada port 8000
	if err := e.Start(":8000"); err != nil {
		e.Logger.Fatal("Failed to start server: ", err)
	}
}

// Rute publik (tanpa autentikasi)
func publicRoutes(e *echo.Echo) {
	e.POST("/api/v1/register", controllers.RegisterHandler) // Registrasi user baru
	e.POST("/api/v1/login", controllers.LoginHandler)       // Login user
	e.Static("/uploads", "uploads")                         // Akses file statis
}

// Rute dengan autentikasi (hanya untuk user login)
func protectedRoutes(e *echo.Echo) {
	authGroup := e.Group("/api/v1")
	authGroup.Use(middlewares.AuthMiddleware) // Middleware untuk validasi token JWT

	// Rute untuk user
	authGroup.GET("/logout", controllers.Logout)                  // Logout user
	authGroup.PUT("/user/photo/:id", controllers.UpdateUserPhoto) // Update foto user
	authGroup.GET("/users/points", controllers.GetUserPoints)
	authGroup.PUT("/user/data/:id", controllers.UpdateUserData) // Update data diri user

	// Rute laporan sampah
	authGroup.POST("/report-rubbish", controllers.CreateReportRubbish) // Membuat laporan
	authGroup.GET("/report-rubbish", controllers.GetAllReportRubbish)
	authGroup.GET("/report-rubbish/history", controllers.GetReportHistoryByUser)

	// Rute khusus admin (misalnya untuk memvalidasi laporan)
	authGroup.PUT("/report-rubbish/:id/status", middlewares.RoleMiddleware("admin")(controllers.UpdateReportStatus))

	adminGroup := authGroup.Group("/admin", middlewares.RoleMiddleware("admin"))
	adminGroup.GET("/users/points", controllers.GetAllUserPoints)
	adminGroup.GET("/users", controllers.GetAllUsers)
	adminGroup.GET("/users/:id", controllers.GetUserByID) // Mendapatkan user berdasarkan ID

	// Rute Artikel Edukasi
	adminGroup.POST("/articles", controllers.BikinArtikel)
	authGroup.GET("/articles", controllers.AmbilSemuaArtikel)
}

// loadEnv memuat variabel environment dari file .env
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Gagal memuat file .env")
	}
}