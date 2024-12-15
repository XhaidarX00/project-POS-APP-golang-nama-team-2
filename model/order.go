package model

import (
	"time"
)

type Order struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	TableID       uint           `json:"table_id"`
	ReservationID uint           `json:"reservation_id"`
	CustomerName  string         `json:"customer_name"`
	Status        string         `json:"status"`
	TotalAmount   float64        `json:"total_amount"`
	Tax           float64        `json:"tax"`
	PaymentMethod string         `json:"payment_method"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     *time.Time     `gorm:"index" json:"deleted_at"`
	OrderProducts []OrderProduct `gorm:"foreignKey:OrderID" json:"order_products"`
}

func SeedOrders() []Order {
	return []Order{
		{TableID: 1, ReservationID: 1, CustomerName: "", Status: "Completed", TotalAmount: 500.00, Tax: 50.00, PaymentMethod: "Credit Card", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 2, ReservationID: 0, CustomerName: "Bob", Status: "ready", TotalAmount: 300.00, Tax: 30.00, PaymentMethod: "Cash", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 3, ReservationID: 3, CustomerName: "", Status: "in progres", TotalAmount: 0.00, Tax: 0.00, PaymentMethod: "None", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 4, ReservationID: 0, CustomerName: "Diana", Status: "Completed", TotalAmount: 450.00, Tax: 45.00, PaymentMethod: "Debit Card", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 5, ReservationID: 5, CustomerName: "", Status: "In Progress", TotalAmount: 150.00, Tax: 15.00, PaymentMethod: "Mobile Payment", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 6, ReservationID: 0, CustomerName: "Fiona", Status: "Completed", TotalAmount: 700.00, Tax: 70.00, PaymentMethod: "Credit Card", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 7, ReservationID: 7, CustomerName: "", Status: "ready", TotalAmount: 200.00, Tax: 20.00, PaymentMethod: "Cash", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 8, ReservationID: 0, CustomerName: "Helen", Status: "Completed", TotalAmount: 350.00, Tax: 35.00, PaymentMethod: "Debit Card", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 9, ReservationID: 9, CustomerName: "", Status: "ready", TotalAmount: 0.00, Tax: 0.00, PaymentMethod: "None", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 10, ReservationID: 0, CustomerName: "Jack", Status: "In Progress", TotalAmount: 400.00, Tax: 40.00, PaymentMethod: "Mobile Payment", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 11, ReservationID: 11, CustomerName: "", Status: "Completed", TotalAmount: 600.00, Tax: 60.00, PaymentMethod: "Credit Card", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 12, ReservationID: 0, CustomerName: "Leo", Status: "ready", TotalAmount: 250.00, Tax: 25.00, PaymentMethod: "Cash", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 13, ReservationID: 13, CustomerName: "", Status: "Cancelled", TotalAmount: 0.00, Tax: 0.00, PaymentMethod: "None", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 14, ReservationID: 0, CustomerName: "Nina", Status: "Completed", TotalAmount: 550.00, Tax: 55.00, PaymentMethod: "Debit Card", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 15, ReservationID: 15, CustomerName: "", Status: "In Progress", TotalAmount: 300.00, Tax: 30.00, PaymentMethod: "Mobile Payment", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 16, ReservationID: 0, CustomerName: "Paula", Status: "Completed", TotalAmount: 800.00, Tax: 80.00, PaymentMethod: "Credit Card", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 17, ReservationID: 17, CustomerName: "", Status: "ready", TotalAmount: 180.00, Tax: 18.00, PaymentMethod: "Cash", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 18, ReservationID: 0, CustomerName: "Rachel", Status: "Completed", TotalAmount: 420.00, Tax: 42.00, PaymentMethod: "Debit Card", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 19, ReservationID: 19, CustomerName: "", Status: "ready", TotalAmount: 0.00, Tax: 0.00, PaymentMethod: "None", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{TableID: 20, ReservationID: 0, CustomerName: "Tina", Status: "In Progress", TotalAmount: 500.00, Tax: 50.00, PaymentMethod: "Mobile Payment", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
}
