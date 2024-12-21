package dashboardrepository

import (
	"errors"
	"project_pos_app/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RepositoryDashboard interface {
	FindPopularProduct() ([]model.Product, error)
	FindSummary(date string) ([]model.Reservation, error)
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
		Select("product_id, SUM(qty) AS total_qty").
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
func (r *repositoryDashboard) FindNewProduct(date string) ([]model.Product, error) {
	var products []model.Product
	err := r.DB.Where("DATE(reservation_date) = ?", date).Find(&products).Error
	if err != nil {
		r.Log.Error("Failed to find all reservation", zap.Error(err))
		return nil, errors.New(" Internal Server Error")
	}
	return products, nil
}
func (r *repositoryDashboard) FindSummary(date string) ([]model.Reservation, error) {
	var reservations []model.Reservation
	err := r.DB.Where("DATE(reservation_date) = ?", date).Find(&reservations).Error
	if err != nil {
		r.Log.Error("Failed to find all reservation", zap.Error(err))
		return nil, errors.New(" Internal Server Error")
	}
	return reservations, nil
}
