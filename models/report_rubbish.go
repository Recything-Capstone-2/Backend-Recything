package models

import (
	"time"
)

type ReportRubbish struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `json:"user_id"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	Photo       string    `json:"photo"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        User      `json:"user" gorm:"foreignKey:UserID;references:ID"`
}
