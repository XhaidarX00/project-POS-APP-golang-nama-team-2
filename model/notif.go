package model

import (
	"fmt"
	"time"
)

type CreateNotification struct {
	Title     string    `json:"title" example:"New Message"`
	Message   string    `json:"message" example:"You have a new message"`
	Status    string    `json:"status" example:"new"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Notification struct {
	ID        int       `gorm:"primaryKey" json:"id" example:"1"`
	Title     string    `json:"title" example:"New Message"`
	Message   string    `json:"message" example:"You have a new message"`
	Status    string    `json:"status" example:"new"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func NotifStock(productName string) Notification {
	return Notification{
		Title:     "Stock Alert",
		Message:   fmt.Sprintf("products %s are out of stock. Please restock immediately.", productName),
		Status:    "new",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NotificationSeed() []Notification {
	return []Notification{
		{
			Title:     "Welcome",
			Message:   "Thank you for joining our platform!",
			Status:    "new",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Title:     "Discount Offer",
			Message:   "Get 20% off on your next purchase.",
			Status:    "new",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Title:     "System Update",
			Message:   "Our system will undergo maintenance tonight from 12 AM to 3 AM.",
			Status:    "readed",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Title:     "Stock Alert",
			Message:   "One or more products are out of stock. Please restock immediately.",
			Status:    "new",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}
