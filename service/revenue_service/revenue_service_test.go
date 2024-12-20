package revenueservice_test

import (
	"errors"
	"project_pos_app/helper"
	"project_pos_app/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFetchTotalRevenueByStatus(t *testing.T) {
	mockDB, service := helper.InitService()

	if mockDB == nil {
		t.Errorf("mockdb bernilai nil")
	}

	if service.Revenue == nil {
		t.Errorf("service revenue bernilai nil")
	}

	t.Run("Successfully fetch total revenue by status", func(t *testing.T) {
		expectedResult := map[string]float64{
			"cancelled": 150,
			"confirmed": 1110.5,
			"failed":    0,
			"pending":   770,
		}
		mockDB.On("GetTotalRevenueByStatus").Once().Return(expectedResult, nil)

		result, err := service.Revenue.FetchTotalRevenueByStatus()
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("Failed to fetch total revenue by status", func(t *testing.T) {
		mockDB.On("GetTotalRevenueByStatus").Return(nil, errors.New("failed to fetch data"))

		result, err := service.Revenue.FetchTotalRevenueByStatus()

		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
	})
}

func TestFetchMonthlyRevenue(t *testing.T) {
	mockDB, service := helper.InitService()

	t.Run("Successfully fetch monthly revenue", func(t *testing.T) {
		expectedResult := map[string]float64{"January": 1000.0, "February": 1500.0}
		mockDB.On("GetMonthlyRevenue").Once().Return(expectedResult, nil)

		result, err := service.Revenue.FetchMonthlyRevenue()

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("Failed to fetch monthly revenue", func(t *testing.T) {
		mockDB.On("GetMonthlyRevenue").Return(nil, errors.New("failed to fetch data"))
		result, err := service.Revenue.FetchMonthlyRevenue()

		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
	})
}

func TestGetLowStockProducts(t *testing.T) {
	mockDB, service := helper.InitService()

	t.Run("Successfully fetch low stock products", func(t *testing.T) {
		expectedProducts := []model.Product{
			{ID: 1, Name: "Product A", Qty: 5},
			{ID: 2, Name: "Product B", Qty: 2},
		}
		mockDB.On("FindLowStockProducts", 10).Once().Return(expectedProducts, nil)

		result, err := service.Revenue.GetLowStockProducts(10)

		assert.NoError(t, err)
		assert.Equal(t, expectedProducts, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("Failed to fetch low stock products with invalid threshold", func(t *testing.T) {
		result, err := service.Revenue.GetLowStockProducts(-5)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "threshold must be a positive number")
	})

	t.Run("Failed to fetch low stock products due to repository error", func(t *testing.T) {
		mockDB.On("FindLowStockProducts", 10).Return(nil, errors.New("repository error"))

		result, err := service.Revenue.GetLowStockProducts(10)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
	})
}

func TestSaveProductRevenue(t *testing.T) {
	mockDB, service := helper.InitService()

	t.Run("Successfully save product revenue", func(t *testing.T) {
		productRevenue := model.ProductRevenue{
			ProductName:  "Product A",
			TotalRevenue: 100.0,
		}
		mockDB.On("SaveProductRevenue", productRevenue).Once().Return(nil)

		err := service.Revenue.SaveProductRevenue(productRevenue)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("Failed to save product revenue with empty product name", func(t *testing.T) {
		productRevenue := model.ProductRevenue{TotalRevenue: 100.0}

		err := service.Revenue.SaveProductRevenue(productRevenue)

		assert.Error(t, err)
		assert.EqualError(t, err, "product name cannot be empty")
	})
}

func TestSaveOrderRevenue(t *testing.T) {
	mockDB, service := helper.InitService()

	t.Run("Successfully save order revenue", func(t *testing.T) {
		orderRevenue := model.OrderRevenue{
			Status:  "completed",
			Revenue: 500.0,
		}
		mockDB.On("SaveOrderRevenue", orderRevenue).Once().Return(nil)

		err := service.Revenue.SaveOrderRevenue(orderRevenue)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("Failed to save order revenue with empty status", func(t *testing.T) {
		orderRevenue := model.OrderRevenue{Revenue: 500.0}

		err := service.Revenue.SaveOrderRevenue(orderRevenue)

		assert.Error(t, err)
		assert.EqualError(t, err, "order status cannot be empty")
	})
}

func TestFetchProductRevenue(t *testing.T) {
	mockDB, service := helper.InitService()

	t.Run("Successfully fetch products revenue", func(t *testing.T) {
		revenueDate := time.Date(2024, 12, 18, 0, 0, 0, 0, time.UTC)
		existingProductRevenue := []model.ProductRevenue{{
			ID:           1,
			ProductName:  "Product A",
			SellPrice:    80,
			Profit:       1500,
			ProfitMargin: 12,
			TotalRevenue: 1800,
			RevenueDate:  revenueDate,
		},
			{
				ProductName:  "Product A",
				SellPrice:    100,
				Profit:       1900,
				ProfitMargin: 15,
				TotalRevenue: 2000,
				RevenueDate:  revenueDate,
			},
		}

		mockDB.On("GetProductRevenues").Once().Return(existingProductRevenue, nil)

		result, err := service.Revenue.FetchProductRevenues()

		assert.NoError(t, err)
		assert.Equal(t, existingProductRevenue, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("Failed to fetch products revenue due to repository error", func(t *testing.T) {
		mockDB.On("GetProductRevenues").Return(nil, errors.New("repository error"))

		result, err := service.Revenue.FetchProductRevenues()

		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
	})
}
