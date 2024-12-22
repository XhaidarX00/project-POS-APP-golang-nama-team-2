package model

import (
	"strconv"
	"time"
)

type Category struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	IconURL     string     `json:"icon_url"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
}

func SeedCategories() []Category {
	categories := []Category{}
	for i := 1; i <= 20; i++ {
		categories = append(categories, Category{
			IconURL:     "https://example.com/icon" + strconv.Itoa(i) + ".png",
			Name:        "Category " + strconv.Itoa(i),
			Description: "Description for Category " + strconv.Itoa(i),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}

	return categories
}
