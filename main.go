package main

import (
	"Backend-Recything/config"
	"Backend-Recything/controllers"
	"log"

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

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS()) // Tambahkan CORS jika aplikasi digunakan oleh client-side

	// Rute autentikasi
	e.POST("/api/v1/register", controllers.RegisterHandler)
	e.POST("/api/v1/login", controllers.LoginHandler)

	// Menjalankan server8
	if err := e.Start(":8090"); err != nil {
		e.Logger.Fatal("Failed to start server: ", err)
	}
}

// loadEnv memuat variabel environment dari file .env
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Gagal memuat file .env")
	}
}
