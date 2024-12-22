package revenueservice

import (
	mocktesting "project_pos_app/mock_testing"
	"project_pos_app/model"

	"go.uber.org/zap"
)

// MockRevenueServiceInterface mendefinisikan kontrak untuk mock service revenue
type MockRevenueServiceInterface interface {
	FetchTotalRevenueByStatus() (map[string]float64, error)
	FetchMonthlyRevenue() (map[string]float64, error)
	FetchProductRevenues() ([]model.ProductRevenue, error)
	SaveOrderRevenue(order model.OrderRevenue) error
	CalculateOrderRevenue() ([]model.OrderRevenue, error)
	SaveProductRevenue(product model.ProductRevenue) error
	CalculateProductRevenue() ([]model.ProductRevenue, error)
	GetLowStockProducts(threshold int) ([]model.Product, error)
}

// MockRevenueService implementasi mock service untuk revenue
type MockRevenueService struct {
	Repo *mocktesting.MockDB
	Log  *zap.Logger
}

// NewMockRevenueService membuat instance baru dari mock service revenue
func NewMockRevenueService(repo *mocktesting.MockDB, log *zap.Logger) MockRevenueServiceInterface {
	return &MockRevenueService{
		Repo: repo,
		Log:  log,
	}
}

// FetchTotalRevenueByStatus mengambil total revenue berdasarkan status
func (m *MockRevenueService) FetchTotalRevenueByStatus() (map[string]float64, error) {
	m.Log.Info("Fetching total revenue by status")
	result, err := m.Repo.GetTotalRevenueByStatus()
	if err != nil {
		m.Log.Error("Failed to fetch total revenue by status", zap.Error(err))
		return nil, err
	}
	return result, nil
}

// FetchMonthlyRevenue mengambil revenue bulanan
func (m *MockRevenueService) FetchMonthlyRevenue() (map[string]float64, error) {
	m.Log.Info("Fetching monthly revenue")
	result, err := m.Repo.GetMonthlyRevenue()
	if err != nil {
		m.Log.Error("Failed to fetch monthly revenue", zap.Error(err))
		return nil, err
	}
	return result, nil
}

// FetchProductRevenues mengambil revenue per produk
func (m *MockRevenueService) FetchProductRevenues() ([]model.ProductRevenue, error) {
	m.Log.Info("Fetching product revenues")
	var result []model.ProductRevenue
	result, err := m.Repo.GetProductRevenues()
	if err != nil {
		m.Log.Error("Failed to fetch product revenues", zap.Error(err))
		return nil, err
	}
	return result, nil
}

// SaveOrderRevenue menyimpan revenue pesanan
func (m *MockRevenueService) SaveOrderRevenue(order model.OrderRevenue) error {
	m.Log.Info("Saving order revenue", zap.Any("order", order))
	return m.Repo.SaveOrderRevenue(order)
}

// CalculateOrderRevenue menghitung revenue pesanan
func (m *MockRevenueService) CalculateOrderRevenue() ([]model.OrderRevenue, error) {
	m.Log.Info("Calculating order revenue")
	result, err := m.Repo.CalculateOrderRevenue()
	if err != nil {
		m.Log.Error("Failed to calculate order revenue", zap.Error(err))
		return nil, err
	}
	return result, nil
}

// SaveProductRevenue menyimpan revenue produk
func (m *MockRevenueService) SaveProductRevenue(product model.ProductRevenue) error {
	m.Log.Info("Saving product revenue", zap.Any("product", product))
	return m.Repo.SaveProductRevenue(product)
}

// CalculateProductRevenue menghitung revenue produk
func (m *MockRevenueService) CalculateProductRevenue() ([]model.ProductRevenue, error) {
	m.Log.Info("Calculating product revenue")
	result, err := m.Repo.CalculateProductRevenue()
	if err != nil {
		m.Log.Error("Failed to calculate product revenue", zap.Error(err))
		return nil, err
	}
	return result, nil
}

// GetLowStockProducts mendapatkan produk dengan stok rendah
func (m *MockRevenueService) GetLowStockProducts(threshold int) ([]model.Product, error) {
	m.Log.Info("Fetching low stock products", zap.Int("threshold", threshold))
	result, err := m.Repo.FindLowStockProducts(threshold)
	if err != nil {
		m.Log.Error("Failed to fetch low stock products", zap.Error(err))
		return nil, err
	}
	return result, nil
}
