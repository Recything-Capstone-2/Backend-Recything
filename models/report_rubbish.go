package models

import (
	"time"
)

type ReportRubbish struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `json:"user_id"`
	Location       string    `json:"location"`
	Description    string    `json:"description"`
	Photo          string    `json:"photo"`
	Status         string    `json:"status"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	TanggalLaporan time.Time `json:"tanggal_laporan"`
	Category       string    `gorm:"type:varchar(50);not null"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	User           User      `json:"user" gorm:"foreignKey:UserID;references:ID"`
}
