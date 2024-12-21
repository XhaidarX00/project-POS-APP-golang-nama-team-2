package model

import "time"

type RevenueByStatus struct {
	Status  string  `gorm:"type:varchar(50)" json:"status" binding:"required" example:"Completed"`
	Revenue float64 `gorm:"type:decimal(10,2)" json:"revenue" binding:"required" example:"100.50"`
}

type MonthlyRevenue struct {
	Month   string  `gorm:"type:varchar(20)" json:"month" binding:"required" example:"January"`
	Revenue float64 `gorm:"type:decimal(10,2)" json:"revenue" binding:"required" example:"100.50"`
}

// OrderRevenue represents the structure for orders
type OrderRevenue struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	Status    string    `gorm:"type:varchar(50)" json:"status" binding:"required" example:"Completed"`
	Revenue   float64   `gorm:"type:decimal(10,2)" json:"revenue" binding:"required" example:"100.50"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" example:"2024-12-01T00:00:00Z"`
	ProductID uint      `json:"product_id"`
}

// ProductRevenue represents the revenue details for products
type ProductRevenue struct {
	ID           uint      `gorm:"primaryKey" json:"id" example:"1"`
	ProductName  string    `gorm:"type:varchar(100)" json:"product_name" binding:"required" example:"Chicken Parmesan"`
	SellPrice    float64   `gorm:"type:decimal(10,2)" json:"sell_price" binding:"required" example:"55.00"`
	Profit       float64   `gorm:"type:decimal(10,2)" json:"profit" binding:"required" example:"7985.00"`
	ProfitMargin float64   `gorm:"type:decimal(5,2)" json:"profit_margin" binding:"required" example:"15.00"`
	TotalRevenue float64   `gorm:"type:decimal(10,2)" json:"total_revenue" binding:"required" example:"8000.00"`
	RevenueDate  time.Time `gorm:"type:date" json:"revenue_date" binding:"required" example:"2024-03-28"`
}

func RevenueSeedProduct() []ProductRevenue {
	return []ProductRevenue{
		{ProductName: "Chicken Parmesan", SellPrice: 55.00, Profit: 7985.00, ProfitMargin: 15.00, TotalRevenue: 8000.00, RevenueDate: time.Now()},
		{ProductName: "Grilled Salmon", SellPrice: 70.00, Profit: 8900.00, ProfitMargin: 20.00, TotalRevenue: 10000.00, RevenueDate: time.Now().AddDate(0, 0, -1)},
		{ProductName: "Vegetarian Pizza", SellPrice: 40.00, Profit: 5000.00, ProfitMargin: 12.50, TotalRevenue: 5600.00, RevenueDate: time.Now().AddDate(0, 0, -2)},
		{ProductName: "Beef Burger", SellPrice: 30.00, Profit: 3000.00, ProfitMargin: 10.00, TotalRevenue: 3500.00, RevenueDate: time.Now().AddDate(0, 0, -3)},
		{ProductName: "Pasta Primavera", SellPrice: 45.00, Profit: 7000.00, ProfitMargin: 18.00, TotalRevenue: 8500.00, RevenueDate: time.Now().AddDate(0, 0, -4)},
		{ProductName: "Chicken Salad", SellPrice: 20.00, Profit: 4000.00, ProfitMargin: 15.00, TotalRevenue: 4800.00, RevenueDate: time.Now().AddDate(0, 0, -5)},
		{ProductName: "Margherita Pizza", SellPrice: 25.00, Profit: 5500.00, ProfitMargin: 16.00, TotalRevenue: 6000.00, RevenueDate: time.Now().AddDate(0, 0, -6)},
		{ProductName: "Fish Tacos", SellPrice: 35.00, Profit: 6200.00, ProfitMargin: 14.50, TotalRevenue: 7000.00, RevenueDate: time.Now().AddDate(0, 0, -7)},
		{ProductName: "Lasagna", SellPrice: 50.00, Profit: 7500.00, ProfitMargin: 17.00, TotalRevenue: 8500.00, RevenueDate: time.Now().AddDate(0, 0, -8)},
		{ProductName: "Steak and Fries", SellPrice: 60.00, Profit: 8000.00, ProfitMargin: 22.00, TotalRevenue: 9000.00, RevenueDate: time.Now().AddDate(0, 0, -9)},
	}
}

// RevenueSeedOrder generates dummy data for OrderRevenue
func RevenueSeedOrder() []OrderRevenue {
	return []OrderRevenue{
		{Status: "Completed", Revenue: 100.50, CreatedAt: time.Now(), ProductID: 1},
		{Status: "In Progress", Revenue: 150.00, CreatedAt: time.Now().AddDate(0, 0, -1), ProductID: 2},
		{Status: "Completed", Revenue: 200.00, CreatedAt: time.Now().AddDate(0, 0, -2), ProductID: 3},
		{Status: "Cancelled", Revenue: 50.00, CreatedAt: time.Now().AddDate(0, 0, -3), ProductID: 4},
		{Status: "Completed", Revenue: 75.00, CreatedAt: time.Now().AddDate(0, 0, -4), ProductID: 5},
		{Status: "In Progress", Revenue: 120.00, CreatedAt: time.Now().AddDate(0, 0, -5), ProductID: 6},
		{Status: "Completed", Revenue: 95.00, CreatedAt: time.Now().AddDate(0, 0, -6), ProductID: 7},
		{Status: "Cancelled", Revenue: 40.00, CreatedAt: time.Now().AddDate(0, 0, -7), ProductID: 8},
		{Status: "Completed", Revenue: 110.00, CreatedAt: time.Now().AddDate(0, 0, -8), ProductID: 9},
		{Status: "In Progress", Revenue: 130.00, CreatedAt: time.Now().AddDate(0, 0, -9), ProductID: 10},
		{Status: "Completed", Revenue: 220.00, CreatedAt: time.Now().AddDate(0, -1, 0), ProductID: 8},
		{Status: "In Progress", Revenue: 180.00, CreatedAt: time.Now().AddDate(0, -2, 0), ProductID: 10},
		{Status: "Cancelled", Revenue: 60.00, CreatedAt: time.Now().AddDate(0, -3, 0), ProductID: 9},
		{Status: "Completed", Revenue: 310.00, CreatedAt: time.Now().AddDate(0, -4, 0), ProductID: 7},
		{Status: "In Progress", Revenue: 190.00, CreatedAt: time.Now().AddDate(0, -5, 0), ProductID: 15},
		{Status: "Failed", Revenue: 0.00, CreatedAt: time.Now().AddDate(0, 0, -1), ProductID: 2},
		{Status: "Failed", Revenue: 0.00, CreatedAt: time.Now().AddDate(0, 0, -1), ProductID: 3},
		{Status: "Failed", Revenue: 0.00, CreatedAt: time.Now().AddDate(0, -1, -1), ProductID: 4},
		{Status: "Failed", Revenue: 0.00, CreatedAt: time.Now().AddDate(0, -2, -1), ProductID: 5},
	}
}
