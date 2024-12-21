package model

import (
	"time"
)

type Superadmin struct {
	ID              uint   `gorm:"primaryKey"`
	UserID          uint   `json:"user_id"` // Relasi ke tabel User
	FullName        string `gorm:"type:varchar(100);not null"`
	Address         string `gorm:"type:varchar(255)"`
	NewPassword     string `gorm:"-"`
	ConfirmPassword string `gorm:"-"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func SeedSuperadmins() []Superadmin {
	return []Superadmin{
		{
			UserID:    1,
			FullName:  "John Doe",
			Address:   "123 Street USA, Chicago",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}
