package model

import (
	"time"
)

type Superadmin struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `json:"user_id"`
	User      User   `gorm:"foreignKey:UserID" binding:"required"`
	FullName  string `gorm:"type:varchar(100);not null" binding:"required,min=3,max=100"`
	Address   string `gorm:"type:varchar(255)" binding:"omitempty,max=255"`
	Image     string `json:"image" binding:"required,imagefile"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func SeedSuperadmins() []Superadmin {
	return []Superadmin{
		{
			UserID:    1,
			FullName:  "John Doe",
			Address:   "123 Street USA, Chicago",
			Image:     "http://imageitem.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}
