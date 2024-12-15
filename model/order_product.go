package model

type OrderProduct struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID"`
	Qty       int     `json:"qty"`
}

func SeedOrderProducts() []OrderProduct {
	return []OrderProduct{
		{OrderID: 1, ProductID: 1, Qty: 2},
		{OrderID: 1, ProductID: 2, Qty: 1},
		{OrderID: 2, ProductID: 3, Qty: 3},
		{OrderID: 2, ProductID: 4, Qty: 2},
		{OrderID: 3, ProductID: 5, Qty: 1},
		{OrderID: 4, ProductID: 6, Qty: 4},
		{OrderID: 5, ProductID: 7, Qty: 1},
		{OrderID: 5, ProductID: 8, Qty: 2},
		{OrderID: 6, ProductID: 9, Qty: 5},
		{OrderID: 7, ProductID: 10, Qty: 2},
		{OrderID: 8, ProductID: 11, Qty: 3},
		{OrderID: 9, ProductID: 12, Qty: 2},
		{OrderID: 10, ProductID: 13, Qty: 1},
		{OrderID: 10, ProductID: 14, Qty: 3},
		{OrderID: 11, ProductID: 15, Qty: 4},
		{OrderID: 12, ProductID: 16, Qty: 2},
		{OrderID: 13, ProductID: 17, Qty: 1},
		{OrderID: 14, ProductID: 18, Qty: 2},
		{OrderID: 15, ProductID: 19, Qty: 3},
		{OrderID: 16, ProductID: 20, Qty: 1},
	}
}
