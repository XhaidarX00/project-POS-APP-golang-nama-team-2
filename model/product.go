package model

import (
	"time"
)

type Product struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	ImageURL   string     `json:"image_url" form:"image_url"`
	Name       string     `json:"name" form:"name"`
	ItemID     string     `json:"item_id" form:"item_id"`
	Stock      string     `json:"stock" form:"stock"`
	CategoryID uint       `json:"category_id" form:"category_id"`
	Qty        int        `json:"qty" form:"qty"`
	Price      float64    `json:"price" form:"price"`
	Status     string     `json:"status" form:"status"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func SeedProducts() []Product {
	now := time.Now()
	return []Product{
		{ImageURL: "https://example.com/image1.jpg", Name: "Product 1", ItemID: "P001", Stock: "in stock", CategoryID: 1, Qty: 10, Price: 100.50, Status: "Active", CreatedAt: now, UpdatedAt: now},
		{ImageURL: "https://example.com/image2.jpg", Name: "Product 2", ItemID: "P002", Stock: "out of stock", CategoryID: 2, Qty: 5, Price: 150.00, Status: "Inactive", CreatedAt: now, UpdatedAt: now},
		{ImageURL: "https://example.com/image3.jpg", Name: "Product 3", ItemID: "P003", Stock: "in stock", CategoryID: 3, Qty: 8, Price: 75.75, Status: "Active", CreatedAt: now, UpdatedAt: now},
		{ImageURL: "https://example.com/image4.jpg", Name: "Product 4", ItemID: "P004", Stock: "limited stock", CategoryID: 4, Qty: 2, Price: 50.00, Status: "Inactive", CreatedAt: now, UpdatedAt: now},
		{ImageURL: "https://example.com/image5.jpg", Name: "Product 5", ItemID: "P005", Stock: "in stock", CategoryID: 5, Qty: 12, Price: 120.00, Status: "Active", CreatedAt: now, UpdatedAt: now},
		{ImageURL: "https://example.com/image6.jpg", Name: "Product 6", ItemID: "P006", Stock: "out of stock", CategoryID: 1, Qty: 3, Price: 60.75, Status: "Inactive", CreatedAt: now, UpdatedAt: now},
		{ImageURL: "https://example.com/image7.jpg", Name: "Product 7", ItemID: "P007", Stock: "limited stock", CategoryID: 2, Qty: 6, Price: 85.00, Status: "Active", CreatedAt: now, UpdatedAt: now},
		{ImageURL: "https://example.com/image8.jpg", Name: "Product 8", ItemID: "P008", Stock: "in stock", CategoryID: 3, Qty: 10, Price: 95.99, Status: "Active", CreatedAt: now, UpdatedAt: now},
		{ImageURL: "https://example.com/image9.jpg", Name: "Product 9", ItemID: "P009", Stock: "out of stock", CategoryID: 4, Qty: 4, Price: 40.00, Status: "Inactive", CreatedAt: now, UpdatedAt: now},
		{ImageURL: "https://example.com/image10.jpg", Name: "Product 10", ItemID: "P010", Stock: "in stock", CategoryID: 5, Qty: 15, Price: 130.25, Status: "Active", CreatedAt: now, UpdatedAt: now},
		{ImageURL: "https://example.com/image11.jpg", Name: "Product 11", ItemID: "P011", Stock: "in stock", CategoryID: 1, Qty: 18, Price: 200.00, Status: "Active", CreatedAt: now, UpdatedAt: now},
		{ImageURL: "https://example.com/image12.jpg", Name: "Product 12", ItemID: "P012", Stock: "out of stock", CategoryID: 2, Qty: 8, Price: 140.00, Status: "Inactive", CreatedAt: now, UpdatedAt: now},
		{ImageURL: "https://example.com/image13.jpg", Name: "Product 13", ItemID: "P013", Stock: "limited stock", CategoryID: 3, Qty: 6, Price: 80.50, Status: "Active", CreatedAt: now, UpdatedAt: now},
		{ImageURL: "https://example.com/image14.jpg", Name: "Product 14", ItemID: "P014", Stock: "limited stock", CategoryID: 4, Qty: 3, Price: 30.00, Status: "Inactive", CreatedAt: now, UpdatedAt: now},
		{ImageURL: "https://example.com/image15.jpg", Name: "Product 15", ItemID: "P015", Stock: "in stock", CategoryID: 5, Qty: 12, Price: 125.75, Status: "Active", CreatedAt: now, UpdatedAt: now},
		{ImageURL: "https://example.com/image16.jpg", Name: "Product 16", ItemID: "P016", Stock: "out of stock", CategoryID: 1, Qty: 2, Price: 35.00, Status: "Inactive", CreatedAt: now, UpdatedAt: now},
		{ImageURL: "https://example.com/image17.jpg", Name: "Product 17", ItemID: "P017", Stock: "in stock", CategoryID: 2, Qty: 10, Price: 90.00, Status: "Active", CreatedAt: now, UpdatedAt: now},
		{ImageURL: "https://example.com/image18.jpg", Name: "Product 18", ItemID: "P018", Stock: "in stock", CategoryID: 3, Qty: 9, Price: 115.99, Status: "Active", CreatedAt: now, UpdatedAt: now},
		{ImageURL: "https://example.com/image19.jpg", Name: "Product 19", ItemID: "P019", Stock: "limited stock", CategoryID: 4, Qty: 1, Price: 20.00, Status: "Inactive", CreatedAt: now, UpdatedAt: now},
		{ImageURL: "https://example.com/image20.jpg", Name: "Product 20", ItemID: "P020", Stock: "in stock", CategoryID: 5, Qty: 14, Price: 150.00, Status: "Active", CreatedAt: now, UpdatedAt: now},
	}
}
