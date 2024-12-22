package revenuerepository

import (
	"errors"
	"log"
	"project_pos_app/model"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RevenueRepositoryInterface interface {
	GetTotalRevenueByStatus() (map[string]float64, error)
	GetMonthlyRevenue() (map[string]float64, error)
	GetProductRevenues() ([]model.ProductRevenue, error)
	SaveOrderRevenue(order model.OrderRevenue) error
	CalculateOrderRevenue() ([]model.OrderRevenue, error)
	SaveProductRevenue(product model.ProductRevenue) error
	CalculateProductRevenue() ([]model.ProductRevenue, error)
	FindLowStockProducts(threshold int) ([]model.Product, error)
}

type RevenueRepository struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewRevenueRepository(db *gorm.DB, log *zap.Logger) RevenueRepositoryInterface {
	return &RevenueRepository{
		DB:  db,
		Log: log,
	}
}

func (r *RevenueRepository) FindLowStockProducts(threshold int) ([]model.Product, error) {
	var products []model.Product
	result := r.DB.Where("qty < ?", threshold).Find(&products)
	return products, result.Error
}

func (r *RevenueRepository) GetProductRevenues() ([]model.ProductRevenue, error) {
	var productRevenues []model.ProductRevenue
	err := r.DB.Table("product_revenues").
		Order("total_revenue DESC").
		Find(&productRevenues).Error
	if err != nil {
		return nil, err
	}

	return productRevenues, nil
}

func (r *RevenueRepository) GetTotalRevenueByStatus() (map[string]float64, error) {
	var results []model.RevenueByStatus

	totalRevenue := make(map[string]float64)

	err := r.DB.Model(&model.OrderRevenue{}).
		Select("status, SUM(revenue) as revenue").
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	for _, res := range results {
		totalRevenue[res.Status] = res.Revenue
	}

	return totalRevenue, nil
}

func (r *RevenueRepository) GetMonthlyRevenue() (map[string]float64, error) {
	var results []model.MonthlyRevenue

	monthlyRevenue := make(map[string]float64)

	// Dapatkan tahun sekarang
	currentYear := time.Now().Year()

	// Tambahkan filter untuk tahun sekarang
	err := r.DB.Model(&model.OrderRevenue{}).
		Select("TO_CHAR(created_at, 'YYYY-MM') as month, SUM(revenue) as revenue").
		Where("EXTRACT(YEAR FROM created_at) = ?", currentYear). // Filter data berdasarkan tahun sekarang
		Group("month").
		Order("month").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	for _, res := range results {
		monthlyRevenue[res.Month] = res.Revenue
	}

	return monthlyRevenue, nil
}

func (r *RevenueRepository) CalculateOrderRevenue() ([]model.OrderRevenue, error) {
	var revenues []model.OrderRevenue

	// Query untuk menghitung revenue dari tabel Order dan OrderProduct
	err := r.DB.Table("orders").
		Select(`
			orders.status AS status,
			SUM(orders.total_amount) AS revenue,
			CURRENT_DATE AS created_at
		`).
		Joins(`
			LEFT JOIN order_products 
			ON orders.id = order_products.order_id
		`).
		Group("orders.status").
		Scan(&revenues).Error

	return revenues, err
}

// func (r *RevenueRepository) CalculateOrderRevenue() ([]model.OrderRevenue, error) {
// 	var orders []model.OrderRevenue

// 	err := r.DB.Model(&model.OrderRevenue{}).
// 		Select("status, SUM(revenue) as revenue, CURRENT_DATE as created_at").
// 		Group("status").
// 		Scan(&orders).Error

// 	return orders, err
// }

func (r *RevenueRepository) SaveOrderRevenue(order model.OrderRevenue) error {
	// Validasi input
	if order.Status == "" {
		return errors.New("order status cannot be empty")
	}
	if order.Revenue < 0 {
		return errors.New("revenue cannot be negative")
	}
	if order.CreatedAt.IsZero() {
		return errors.New("created_at cannot be empty")
	}

	// Mulai transaksi untuk memastikan konsistensi data
	tx := r.DB.Begin()
	if tx.Error != nil {
		return errors.New("database error")
		// return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // Rollback jika terjadi panic
		}
	}()

	// Proses pencarian order yang sudah ada
	var existingOrder model.OrderRevenue
	result := tx.Where("id = ? AND created_at = ?", order.ID, order.CreatedAt).First(&existingOrder)

	if result.Error == nil {
		// Jika order sudah ada, lakukan pembaruan (update) data order
		if err := tx.Model(&existingOrder).Updates(order).Error; err != nil {
			tx.Rollback()
			return errors.New("database error")
			// return err
		}
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Jika order belum ada, buat data baru
		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			return errors.New("database error")
			// return err
		}
	} else {
		tx.Rollback()
		// return errors.New("database error")
		// return err
	}

	// Komit transaksi jika tidak ada error
	tx.Commit()
	return nil
}

// func (r *RevenueRepository) CalculateProductRevenue() ([]model.ProductRevenue, error) {
// 	var products []model.ProductRevenue

// 	err := r.DB.Table("products").
// 		Select("products.name AS product_name, products.price AS sell_price, SUM(order_products.qty * products.price) AS total_revenue, 15.00 AS profit_margin, CURRENT_DATE AS revenue_date").
// 		Joins("JOIN order_products ON products.id = order_products.product_id").
// 		Joins("JOIN orders ON order_products.order_id = orders.id").
// 		Where("orders.status = ?", "confirmed").
// 		Group("products.name, products.price").
// 		Scan(&products).Error

// 	return products, err
// }

// CalculateProductRevenue calculates revenue details for all products
func (r *RevenueRepository) CalculateProductRevenue() ([]model.ProductRevenue, error) {
	var products []model.ProductRevenue

	err := r.DB.Table("products").
		Select(`
			products.name AS product_name, 
			products.price AS sell_price, 
			SUM(order_products.qty * products.price) AS total_revenue, 
			CURRENT_DATE AS revenue_date
		`).
		Joins("JOIN order_products ON products.id = order_products.product_id").
		Joins("JOIN orders ON order_products.order_id = orders.id").
		Where("orders.status = ?", "Completed").
		Group("products.name, products.price").
		Scan(&products).Error

	if err != nil {
		return nil, errors.New("failed to calculate product revenue: " + err.Error())
	}

	// Calculate profit and profit margin for each product
	ProfitMargin := viper.GetFloat64("PROFIT_MARGIN")
	for i := range products {
		products[i].Profit = calculateProfit(products[i].TotalRevenue, ProfitMargin)
		products[i].ProfitMargin = calculateProfitMargin(products[i].TotalRevenue, products[i].Profit)
	}

	log.Printf("Data Product : %v\n", products)
	return products, nil
}

// Helper function to calculate profit
func calculateProfit(totalRevenue, profitMargin float64) float64 {
	return totalRevenue * (profitMargin / 100)
}

// Helper function to calculate profit margin
func calculateProfitMargin(totalRevenue, profit float64) float64 {
	if totalRevenue == 0 {
		return 0
	}
	return (profit / totalRevenue) * 100
}

func (r *RevenueRepository) SaveProductRevenue(product model.ProductRevenue) error {
	// Validasi input
	if product.ProductName == "" {
		return errors.New("product name cannot be empty")
	}
	if product.SellPrice <= 0 {
		return errors.New("sell price must be positive")
	}
	if product.RevenueDate.IsZero() {
		return errors.New("revenue date cannot be empty")
	}

	// Mulai transaksi untuk memastikan konsistensi data
	tx := r.DB.Begin()
	if tx.Error != nil {
		// return errors.New("database error")
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // Rollback jika terjadi panic
		}
	}()

	// Proses pencarian produk yang sudah ada
	var existingRevenue model.ProductRevenue
	result := tx.Where("product_name = ? AND revenue_date = ?", product.ProductName, product.RevenueDate).First(&existingRevenue)

	if result.Error == nil {
		// Jika produk sudah ada, lakukan pembaruan (update) data produk
		if err := tx.Model(&existingRevenue).Updates(product).Error; err != nil {
			tx.Rollback()
			// return errors.New("database error")
			return err
		}
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Jika produk belum ada, buat data baru
		if err := tx.Create(&product).Error; err != nil {
			tx.Rollback()
			// return errors.New("database error")
			return err
		}
	} else {
		// Jika error lain terjadi
		tx.Rollback()
		// return errors.New("database error")
	}

	// Komit transaksi jika tidak ada error
	if err := tx.Commit().Error; err != nil {
		return errors.New("failed to commit transaction")
	}

	return nil
}
