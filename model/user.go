package model

import (
	"time"

	"gorm.io/gorm"
)

type Login struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID        int             `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string          `grom:"type:varchar(100)" json:"name" binding:"required"`
	Email     string          `grom:"type:varchar(255);unique" json:"email" binding:"required,email"`
	Password  string          `grom:"type:varchar(50)" json:"password" binding:"required,min=8"`
	CreatedAt time.Time       `gorm:"autoCreateTime"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime"`
	DeletedAt *gorm.DeletedAt `gorm:"index"`
}

type Session struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int       `gorm:"type:int"`
	Token        string    `gorm:"not null"`
	IpAddress    string    `gorm:"not null"`
	LastActivity time.Time `gorm:"not null"`
}
