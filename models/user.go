package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id_user"`
	Nama      string    `gorm:"type:varchar(255)"json:"nama"`
	Username  string    `gorm:"type:varchar(255);unique" json:"username"`
	Password  string    `gorm:"type:varchar(255)" json:"password"`
	Email     string    `gorm:"type:varchar(255);unique;not null"`
	CreatedAt time.Time `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime" json:"updated_at"`
}
