package model

import (
	"log"
	"project_pos_app/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Login struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID        int             `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string          `gorm:"type:varchar(255);unique" json:"email" binding:"required,email"`
	Password  string          `gorm:"type:varchar(255)" json:"password" binding:"required,min=8"`
	Role      string          `gorm:"type:varchar(255)" json:"role" binding:"required"`
	CreatedAt time.Time       `gorm:"autoCreateTime"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime"`
	DeletedAt *gorm.DeletedAt `gorm:"index"`
}

type Session struct {
	ID           int `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int `gorm:"type:int"`
	Token        string
	IpAddress    string    `gorm:"not null"`
	LastActivity time.Time `gorm:"not null"`
}

func SeedUsers() []User {

	users := []struct {
		Email    string
		Password string
		Role     string
	}{
		{"superadmin@example.com", "superadmin123", "super_admin"},
		{"admin@example.com", "admin123", "admin"},
		{"admin1@example.com", "admin123", "admin"},
		{"admin2@example.com", "admin123", "admin"},
		{"admin3@example.com", "admin123", "admin"},
		{"staff@example.com", "staff123", "staff"},
		{"staff1@example.com", "staff123", "staff"},
		{"staff2@example.com", "staff123", "staff"},
		{"staff3@example.com", "staff123", "staff"},
		{"staff4@example.com", "staff123", "staff"},
	}

	var seededUsers []User
	for _, user := range users {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			log.Fatalf("Error hashing password for user %s: %v", user.Email, err)
		}
		seededUsers = append(seededUsers, User{
			Email:     user.Email,
			Password:  hashedPassword,
			Role:      user.Role,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	return seededUsers
}

func SeedSessions() []Session {
	return []Session{
		{
			UserID:       1,
			Token:        uuid.New().String(),
			IpAddress:    "192.168.1.1",
			LastActivity: time.Now(),
		},
		{
			UserID:       2,
			Token:        uuid.New().String(),
			IpAddress:    "192.168.1.2",
			LastActivity: time.Now(),
		},
		{
			UserID:       3,
			Token:        uuid.New().String(),
			IpAddress:    "192.168.1.3",
			LastActivity: time.Now(),
		},
		{
			UserID:       4,
			Token:        uuid.New().String(),
			IpAddress:    "192.168.1.4",
			LastActivity: time.Now(),
		},
		{
			UserID:       5,
			Token:        uuid.New().String(),
			IpAddress:    "192.168.1.5",
			LastActivity: time.Now(),
		},
		{
			UserID:       6,
			Token:        uuid.New().String(),
			IpAddress:    "192.168.1.6",
			LastActivity: time.Now(),
		},
		{
			UserID:       7,
			Token:        uuid.New().String(),
			IpAddress:    "192.168.1.7",
			LastActivity: time.Now(),
		},
		{
			UserID:       8,
			Token:        uuid.New().String(),
			IpAddress:    "192.168.1.8",
			LastActivity: time.Now(),
		},
		{
			UserID:       9,
			Token:        uuid.New().String(),
			IpAddress:    "192.168.1.9",
			LastActivity: time.Now(),
		},
		{
			UserID:       10,
			Token:        uuid.New().String(),
			IpAddress:    "192.168.1.10",
			LastActivity: time.Now(),
		},
	}
}
