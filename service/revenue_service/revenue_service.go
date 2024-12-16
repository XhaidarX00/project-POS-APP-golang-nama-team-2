package revenueservice

import (
	"errors"
	"project_pos_app/model"
	"project_pos_app/repository"

	"go.uber.org/zap"
)

type RevenueServiceInterface interface {
	FetchTotalRevenueByStatus() (map[string]float64, error)
	FetchMonthlyRevenue() (map[string]float64, error)
	FetchProductRevenues() ([]model.ProductRevenue, error)
	SaveOrderRevenue(order model.OrderRevenue) error
	CalculateOrderRevenue() ([]model.OrderRevenue, error)
	SaveProductRevenue(product model.ProductRevenue) error
	CalculateProductRevenue() ([]model.ProductRevenue, error)
	GetLowStockProducts(threshold int) ([]model.ProductRevenue, error)
}

type revenueService struct {
	Repo *repository.AllRepository
	Log  *zap.Logger
}

func NewRevenueService(repo *repository.AllRepository, log *zap.Logger) RevenueServiceInterface {
	return &revenueService{
		Repo: repo,
		Log:  log,
	}
}

func (s *revenueService) FetchTotalRevenueByStatus() (map[string]float64, error) {
	return s.Repo.Revenue.GetTotalRevenueByStatus()
}

func (s *revenueService) FetchMonthlyRevenue() (map[string]float64, error) {
	return s.Repo.Revenue.GetMonthlyRevenue()
}

func (s *revenueService) FetchProductRevenues() ([]model.ProductRevenue, error) {
	return s.Repo.Revenue.GetProductRevenues()
}

func (s *revenueService) GetLowStockProducts(threshold int) ([]model.ProductRevenue, error) {
	if threshold <= 0 {
		return nil, errors.New("threshold must be a positive number")
	}
	return s.Repo.Revenue.FindLowStockProducts(threshold)
}

func (s *revenueService) CalculateProductRevenue() ([]model.ProductRevenue, error) {
	return s.Repo.Revenue.CalculateProductRevenue()
}

func (s *revenueService) SaveProductRevenue(product model.ProductRevenue) error {
	if product.ProductName == "" {
		return errors.New("product name cannot be empty")
	}
	return s.Repo.Revenue.SaveProductRevenue(product)
}

func (s *revenueService) CalculateOrderRevenue() ([]model.OrderRevenue, error) {
	return s.Repo.Revenue.CalculateOrderRevenue()
}

func (s *revenueService) SaveOrderRevenue(order model.OrderRevenue) error {
	if order.Status == "" {
		return errors.New("order status cannot be empty")
	}
	return s.Repo.Revenue.SaveOrderRevenue(order)
}
