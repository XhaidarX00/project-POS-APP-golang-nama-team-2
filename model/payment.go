package model

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Name      string          `json:"name"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggerignore:"true"`
}

func SeedPayments() []Payment {

	payments := []Payment{
		{
			Name:      "Cash",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Credit Card",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Bank Transfer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "E-Wallet",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	return payments
}
