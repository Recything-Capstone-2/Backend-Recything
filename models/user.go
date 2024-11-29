package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id_user"`
	NamaLengkap  string    `gorm:"type:varchar(255)" json:"nama_lengkap"`
	TanggalLahir time.Time `gorm:"type:datetime" json:"tanggal_lahir"`
	NoTelepon    string    `gorm:"type:varchar(15)" json:"no_telepon" validate:"regexp"`
	Password     string    `gorm:"type:varchar(255)" json:"password"`
	Email        string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Role         string    `gorm:"type:varchar(50);default:'user'" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
