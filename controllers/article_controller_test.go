package controllers

import (
	"Backend-Recything/config"
	"Backend-Recything/helper"
	"Backend-Recything/models"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

var e *echo.Echo

func init() {
	// Inisialisasi instance Echo
	e = echo.New()
	// Inisialisasi database
	loadEnv()
	config.InitDB()
}

func TestAmbilArtikelByID(t *testing.T) {
	article := models.Article{
		Judul:     "Artikel Tes",
		Author:    "Penulis Tes",
		Konten:    "Konten Tes",
		LinkFoto:  "http://contoh.com/foto.jpg",
		LinkVideo: "http://contoh.com/video.mp4",
	}
	config.DB.Create(&article)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/articles/"+strconv.Itoa(int(article.ID)), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(article.ID)))

	// Panggil fungsi controller
	if err := AmbilArtikelByID(c); err != nil {
		t.Fatalf("Gagal mengambil artikel berdasarkan ID: %v", err)
	}

	// Verifikasi status kode
	if rec.Code != http.StatusOK {
		t.Errorf("Status code seharusnya %d, tapi dapat %d", http.StatusOK, rec.Code)
	}

	// Verifikasi response
	var response helper.Response
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Gagal membaca response: %v", err)
	}

	if response.Meta.Status != "success" {
		t.Errorf("Response status seharusnya success, tapi dapat %s", response.Meta.Status)
	}
}

func TestAmbilSemuaArtikel(t *testing.T) {
	article := models.Article{
		Judul:     "Tes Artikel",
		Author:    "Tes Penulis",
		Konten:    "Ini tes",
		LinkFoto:  "http://contoh.com/poto.jpg",
		LinkVideo: "http://contoh.com/pidio.mp4",
	}
	config.DB.Create(&article)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/articles", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := AmbilSemuaArtikel(c); err != nil {
		t.Fatalf("Gagal mendapatkan artikel: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Stats kode harunya %d, tapi dapet %d", http.StatusOK, rec.Code)
	}

	var response helper.Response
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Gagal batalin respon: %v", err)
	}

	if response.Meta.Status != "success" {
		t.Errorf("Status yang diharapkan adalah sukses, tapi malah dapat %s", response.Meta.Status)
	}
}
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Gagal memuat file .env")
	}
}
