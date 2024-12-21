package model

type Summary struct {
	DailySales   int `json:"dailySales"`
	MonthlySales int `json:"monthlySales"`
	TotalTables  int `json:"totalTables"`
}

type ReportExcel struct {
	No           int
	OrderId      uint
	CustomerName string
	ProductName  string
	Price        float64
	Qty          int
	Status       string
	CreatedAt    string
}
