package models

import (
	"time"
)

type Article struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Judul     string    `gorm:"type:varchar(255);not null" json:"judul"`
	Author    string    `gorm:"type:varchar(255);not null" json:"author"`
	Konten    string    `gorm:"type:text;not null" json:"konten"`
	LinkFoto  string    `gorm:"type:varchar(255);not null" json:"link_foto"` // Foto wajib
	LinkVideo string    `gorm:"type:varchar(255)" json:"link_video"`         // Video opsional
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
