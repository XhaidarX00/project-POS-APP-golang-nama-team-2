package dashboardrepository

import (
	"errors"
	"project_pos_app/model"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RepositoryDashboard interface {
	FindPopularProduct() ([]model.Product, error)
	FindNewProduct() ([]model.Product, error)
	FindSummary(summary *model.Summary) error
	FindReport(report *[]model.ReportExcel) error
}

type repositoryDashboard struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewReservationRepository(db *gorm.DB, log *zap.Logger) RepositoryDashboard {
	return &repositoryDashboard{
		DB:  db,
		Log: log,
	}
}

func (r *repositoryDashboard) FindPopularProduct() ([]model.Product, error) {
	var avgOrder int
	err := r.DB.Raw(`
    SELECT ROUND(AVG(sub.qty), 0) AS avg_order
    FROM (
        SELECT "product_id", SUM("qty") AS qty
        FROM "order_products"
        GROUP BY "product_id"
    ) sub
`).Scan(&avgOrder).Error

	if err != nil {
		r.Log.Error("Failed to find avg order", zap.Error(err))
		return nil, errors.New(" Internal Server Error")
	}
	var productId []uint
	err = r.DB.Model(&model.OrderProduct{}).
		Select("product_id").
		Group("product_id").
		Having("SUM(qty) > ?", avgOrder).
		Scan(&productId).Error
	if err != nil {
		r.Log.Error("Failed to find popular productId", zap.Error(err))
		return nil, errors.New(" Internal Server Error")
	}
	var products []model.Product
	for _, v := range productId {
		var product model.Product
		err := r.DB.Find(&product, v).Error
		if err != nil {
			r.Log.Error("Failed to find popular product", zap.Error(err))
			return nil, errors.New(" Internal Server Error")
		}
		products = append(products, product)
	}
	return products, nil
}
func (r *repositoryDashboard) FindNewProduct() ([]model.Product, error) {
	var products []model.Product
	err := r.DB.Where("created_at >= CURRENT_DATE - INTERVAL '30 days'").Find(&products).Error
	if err != nil {
		r.Log.Error("Failed to find new product", zap.Error(err))
		return nil, errors.New(" Internal Server Error")
	}
	return products, nil
}
func (r *repositoryDashboard) FindSummary(summary *model.Summary) error {
	date := time.Now().Format("2006-01-02")
	year := time.Now().Year()
	month := int(time.Now().Month())
	var dailySales float64
	var monthlySales float64
	var count int64
	err := r.DB.Model(&model.Order{}).
		Where("status ILIKE ?", "Completed").
		Where("DATE(created_at) = ?", date).
		Select("COALESCE(SUM(total_amount), 0) AS daily_sales").
		Group("DATE(created_at)").
		Scan(&dailySales).Error
	if err != nil {
		r.Log.Error("Failed to find daily sales", zap.Error(err))
		return errors.New(" Internal Server Error")
	}
	err = r.DB.Model(&model.Order{}).
		Where("status ILIKE ?", "Completed").
		Where("EXTRACT(YEAR FROM created_at) = ?", year).
		Where("EXTRACT(MONTH FROM created_at) = ?", month).
		Select("COALESCE(SUM(total_amount), 0) AS monthly_sales").
		Group("EXTRACT(YEAR FROM created_at), EXTRACT(MONTH FROM created_at)").
		Scan(&monthlySales).Error
	if err != nil {
		r.Log.Error("Failed to find monthly sales", zap.Error(err))
		return errors.New(" Internal Server Error")
	}
	err = r.DB.Model(&model.Table{}).Count(&count).Error
	if err != nil {
		r.Log.Error("Failed to find total table", zap.Error(err))
		return errors.New(" Internal Server Error")
	}
	// fmt.Println("MASUK FIND SUMMARY REPO", date, month, year, dailySales, monthlySales)
	summary.DailySales = int(dailySales)
	summary.MonthlySales = int(monthlySales)
	summary.TotalTables = int(count)
	return nil
}
func (r *repositoryDashboard) FindReport(report *[]model.ReportExcel) error {
	err := r.DB.Table("order_products as op").
		Select("o.id as order_id, o.customer_name, p.name as product_name, p.price, op.qty, o.status, o.created_at").
		Joins("join orders as o on op.order_id = o.id").
		Joins("join products as p on op.product_id = p.id").
		Scan(&report).
		Error
	if err != nil {
		r.Log.Error("Failed to find monthly sales", zap.Error(err))
		return errors.New(" Internal Server Error")
	}
	return nil
}
