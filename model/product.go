package model

import (
	"time"
)

type Product struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	ImageURL   string     `json:"image_url"`
	Name       string     `json:"name"`
	CategoryID uint       `json:"category_id"`
	Qty        int        `json:"qty"`
	Price      float64    `json:"price"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at"`
}

func SeedProducts() []Product {
	return []Product{
		{ImageURL: "https://example.com/image1.jpg", Name: "Product 1", CategoryID: 1, Qty: 10, Price: 100.50, Status: "Available", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ImageURL: "https://example.com/image2.jpg", Name: "Product 2", CategoryID: 2, Qty: 15, Price: 200.00, Status: "Available", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ImageURL: "https://example.com/image3.jpg", Name: "Product 3", CategoryID: 3, Qty: 20, Price: 150.75, Status: "Out of Stock", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ImageURL: "https://example.com/image4.jpg", Name: "Product 4", CategoryID: 4, Qty: 25, Price: 125.00, Status: "Available", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ImageURL: "https://example.com/image5.jpg", Name: "Product 5", CategoryID: 5, Qty: 30, Price: 75.90, Status: "Available", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ImageURL: "https://example.com/image6.jpg", Name: "Product 6", CategoryID: 1, Qty: 12, Price: 80.20, Status: "Out of Stock", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ImageURL: "https://example.com/image7.jpg", Name: "Product 7", CategoryID: 2, Qty: 50, Price: 250.00, Status: "Available", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ImageURL: "https://example.com/image8.jpg", Name: "Product 8", CategoryID: 3, Qty: 45, Price: 120.60, Status: "Out of Stock", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ImageURL: "https://example.com/image9.jpg", Name: "Product 9", CategoryID: 4, Qty: 22, Price: 180.30, Status: "Available", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ImageURL: "https://example.com/image10.jpg", Name: "Product 10", CategoryID: 5, Qty: 28, Price: 95.75, Status: "Available", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ImageURL: "https://example.com/image11.jpg", Name: "Product 11", CategoryID: 1, Qty: 11, Price: 110.00, Status: "Available", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ImageURL: "https://example.com/image12.jpg", Name: "Product 12", CategoryID: 2, Qty: 32, Price: 200.90, Status: "Out of Stock", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ImageURL: "https://example.com/image13.jpg", Name: "Product 13", CategoryID: 3, Qty: 18, Price: 140.00, Status: "Available", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ImageURL: "https://example.com/image14.jpg", Name: "Product 14", CategoryID: 4, Qty: 35, Price: 130.25, Status: "Out of Stock", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ImageURL: "https://example.com/image15.jpg", Name: "Product 15", CategoryID: 5, Qty: 40, Price: 85.40, Status: "Available", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ImageURL: "https://example.com/image16.jpg", Name: "Product 16", CategoryID: 1, Qty: 14, Price: 90.00, Status: "Available", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ImageURL: "https://example.com/image17.jpg", Name: "Product 17", CategoryID: 2, Qty: 55, Price: 220.00, Status: "Out of Stock", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ImageURL: "https://example.com/image18.jpg", Name: "Product 18", CategoryID: 3, Qty: 60, Price: 150.40, Status: "Available", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ImageURL: "https://example.com/image19.jpg", Name: "Product 19", CategoryID: 4, Qty: 45, Price: 115.00, Status: "Available", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ImageURL: "https://example.com/image20.jpg", Name: "Product 20", CategoryID: 5, Qty: 50, Price: 99.99, Status: "Available", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

}
