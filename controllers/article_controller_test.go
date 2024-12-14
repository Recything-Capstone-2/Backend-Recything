package controllers

import (
	"Backend-Recything/config"
	"Backend-Recything/helper"
	"Backend-Recything/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

var (
	e *echo.Echo
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading file .env")
	}

	e = echo.New()

	if err := config.InitDB(); err != nil {
		log.Fatalf("inisialisasi database gagal: %v", err)
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
