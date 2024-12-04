package models

import (
	"time"
)

type Points struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Points    uint      `json:"points"`
	User      User      `gorm:"foreignKey:UserID"` // Relasi ke model User
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
