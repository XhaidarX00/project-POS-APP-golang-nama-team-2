package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	TableID       uint            `json:"table_id" binding:"required"`
	CustomerName  string          `json:"customer_name,omitempty" binding:"required"`
	Status        string          `json:"status"`
	TotalAmount   float64         `json:"total_amount"`
	Tax           float64         `json:"tax"`
	PaymentMethod uint            `json:"payment_method"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggerignore:"true"`
	OrderProducts []OrderProduct  `gorm:"-" json:"order_products" `
}

type OrderResponse struct {
	ID           uint                   `json:"id"`
	CustomerName string                 `json:"customer_name"`
	TableID      int                    `json:"table_id"`
	Status       string                 `json:"status"`
	OrderDate    time.Time              `json:"order_date"`
	SubTotal     int                    `json:"sub_total"`
	OrderProduct []OrderProductResponse `json:"order_products" gorm:"-"`
}

type OrderProductResponse struct {
	OrderID int     `json:"order_id"`
	Qty     int     `json:"qty"`
	Item    string  `json:"item"`
	Price   float64 `json:"price"`
}

func SeedOrders() []Order {
	return []Order{
		{TableID: 1, CustomerName: "Juned", Status: "Completed", TotalAmount: 500.00, Tax: 50.00, PaymentMethod: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 2, CustomerName: "Bob", Status: "ready", TotalAmount: 300.00, Tax: 30.00, PaymentMethod: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 3, CustomerName: "Deni", Status: "in progres", TotalAmount: 0.00, Tax: 0.00, PaymentMethod: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 4, CustomerName: "Diana", Status: "Completed", TotalAmount: 450.00, Tax: 45.00, PaymentMethod: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 5, CustomerName: "Adam", Status: "In Progress", TotalAmount: 150.00, Tax: 15.00, PaymentMethod: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 6, CustomerName: "Fiona", Status: "Completed", TotalAmount: 700.00, Tax: 70.00, PaymentMethod: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 7, CustomerName: "Fina", Status: "ready", TotalAmount: 200.00, Tax: 20.00, PaymentMethod: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 8, CustomerName: "Helen", Status: "Completed", TotalAmount: 350.00, Tax: 35.00, PaymentMethod: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 9, CustomerName: "Sule", Status: "ready", TotalAmount: 0.00, Tax: 0.00, PaymentMethod: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 10, CustomerName: "Jack", Status: "In Progress", TotalAmount: 400.00, Tax: 40.00, PaymentMethod: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 11, CustomerName: "Kipli", Status: "Completed", TotalAmount: 600.00, Tax: 60.00, PaymentMethod: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 12, CustomerName: "Leo", Status: "ready", TotalAmount: 250.00, Tax: 25.00, PaymentMethod: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 13, CustomerName: "Jaenab", Status: "Cancelled", TotalAmount: 0.00, Tax: 0.00, PaymentMethod: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 14, CustomerName: "Nina", Status: "Completed", TotalAmount: 550.00, Tax: 55.00, PaymentMethod: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 15, CustomerName: "Siti", Status: "In Progress", TotalAmount: 300.00, Tax: 30.00, PaymentMethod: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 16, CustomerName: "Paula", Status: "Completed", TotalAmount: 800.00, Tax: 80.00, PaymentMethod: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 17, CustomerName: "Maemunah", Status: "ready", TotalAmount: 180.00, Tax: 18.00, PaymentMethod: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 18, CustomerName: "Rachel", Status: "Completed", TotalAmount: 420.00, Tax: 42.00, PaymentMethod: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 19, CustomerName: "Tukiman", Status: "ready", TotalAmount: 0.00, Tax: 0.00, PaymentMethod: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 20, CustomerName: "Tina", Status: "In Progress", TotalAmount: 500.00, Tax: 50.00, PaymentMethod: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
}
