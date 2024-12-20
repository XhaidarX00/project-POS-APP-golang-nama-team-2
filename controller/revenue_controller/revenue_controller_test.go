package revenuecontroller_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	revenuecontroller "project_pos_app/controller/revenue_controller"
	"project_pos_app/helper"
	mocktesting "project_pos_app/mock_testing"
	"project_pos_app/model"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type NotifControllerSuite struct {
	suite.Suite
	controller revenuecontroller.RevenueController
	mockDB     *mocktesting.MockDB
	ctx        *gin.Context
	writer     *httptest.ResponseRecorder
}

func TestNotifControllerSuite(t *testing.T) {
	suite.Run(t, new(NotifControllerSuite))
}

func (suite *NotifControllerSuite) SetupTest() {
	mockLogger := zap.NewNop()
	mockDB, mockService := helper.InitService()

	suite.controller = revenuecontroller.NewRevenueController(mockService, mockLogger)
	suite.mockDB = mockDB
	gin.SetMode(gin.TestMode)
	suite.writer = httptest.NewRecorder()
	suite.ctx, _ = gin.CreateTestContext(suite.writer)
}

func (suite *NotifControllerSuite) TestGetTotalRevenueByStatus() {
	suite.Run("Successfully fetch total revenue by status", func() {
		expectedResult := map[string]float64{
			"cancelled": 150,
			"confirmed": 1110.5,
			"failed":    0,
			"pending":   770,
		}

		suite.mockDB.On("GetTotalRevenueByStatus").Once().Return(expectedResult, nil)

		// Prepare request
		suite.ctx.Request = httptest.NewRequest(http.MethodPost, "/api/notifications", nil)
		suite.ctx.Request.Header.Set("Content-Type", "application/json")

		// Call controller
		suite.controller.GetTotalRevenueByStatus(suite.ctx)

		Result, err := helper.ConvertToMap(expectedResult)
		suite.NoError(err)
		expected := model.SuccessResponse{
			Status:  http.StatusOK,
			Message: "Fetch total revenue by status successfully",
			Data:    Result,
		}

		var apiResponse model.SuccessResponse
		err = json.Unmarshal(suite.writer.Body.Bytes(), &apiResponse)
		suite.NoError(err)

		suite.Equal(expected, apiResponse)
		suite.Equal(http.StatusOK, suite.writer.Code)
		suite.mockDB.AssertCalled(suite.T(), "GetTotalRevenueByStatus")
	})

	suite.Run("Failed to fetch total revenue by status", func() {
		suite.SetupTest()
		suite.mockDB.On("GetTotalRevenueByStatus").Once().Return(nil, errors.New("database error"))

		// Prepare request
		suite.ctx.Request = httptest.NewRequest(http.MethodPost, "/api/notifications", nil)
		suite.ctx.Request.Header.Set("Content-Type", "application/json")

		// Call controller
		suite.controller.GetTotalRevenueByStatus(suite.ctx)

		suite.Equal(http.StatusInternalServerError, suite.writer.Code)
		suite.Contains(suite.writer.Body.String(), "Failed to fetch total revenue by status")
	})
}

func (suite *NotifControllerSuite) TestGetMonthlyRevenue() {
	suite.Run("Successfully fetch monthly revenue", func() {
		expectedResult := map[string]float64{"January": 1000.0, "February": 1500.0}
		suite.mockDB.On("GetMonthlyRevenue").Once().Return(expectedResult, nil)

		suite.ctx.Request = httptest.NewRequest(http.MethodGet, "/api/monthly_revenue", nil)
		suite.ctx.Request.Header.Set("Content-Type", "application/json")

		// Call controller
		suite.controller.GetMonthlyRevenue(suite.ctx)

		Result, err := helper.ConvertToMap(expectedResult)
		suite.NoError(err)
		expected := model.SuccessResponse{
			Status:  http.StatusOK,
			Message: "Fetch monthly revenue successfully",
			Data:    Result,
		}

		var apiResponse model.SuccessResponse
		err = json.Unmarshal(suite.writer.Body.Bytes(), &apiResponse)
		suite.NoError(err)

		suite.Equal(expected, apiResponse)
		suite.Equal(http.StatusOK, suite.writer.Code)
		suite.Contains(suite.writer.Body.String(), "Fetch monthly revenue successfully")
	})

	suite.Run("Failed to fetch monthly revenue", func() {
		suite.SetupTest()
		suite.mockDB.On("GetMonthlyRevenue").Return(nil, errors.New("failed to fetch data"))
		suite.ctx.Request = httptest.NewRequest(http.MethodGet, "/api/monthly_revenue", nil)
		suite.ctx.Request.Header.Set("Content-Type", "application/json")

		// Call controller
		suite.controller.GetMonthlyRevenue(suite.ctx)

		suite.Equal(http.StatusInternalServerError, suite.writer.Code)
		suite.Contains(suite.writer.Body.String(), "Failed to fetch monthly revenue")
	})
}

func (suite *NotifControllerSuite) TestGetProductRevenues() {
	suite.Run("Successfully fetch product revenues", func() {
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

		// Call mock
		suite.mockDB.On("GetProductRevenues").Once().Return(existingProductRevenue, nil)

		// Call request
		suite.ctx.Request = httptest.NewRequest(http.MethodGet, "/api/product_revenues", nil)
		suite.ctx.Request.Header.Set("Content-Type", "application/json")

		// Call controller
		suite.controller.GetProductRevenues(suite.ctx)
		suite.Equal(http.StatusOK, suite.writer.Code)
		suite.Contains(suite.writer.Body.String(), "Fetch product revenues successfully")
	})

	suite.Run("Failed to fetch product revenues", func() {
		suite.SetupTest()
		suite.mockDB.On("GetProductRevenues").Return(nil, errors.New("repository error"))

		suite.ctx.Request = httptest.NewRequest(http.MethodGet, "/api/product_revenues", nil)
		suite.ctx.Request.Header.Set("Content-Type", "application/json")

		// Call controller
		suite.controller.GetProductRevenues(suite.ctx)

		suite.Equal(http.StatusInternalServerError, suite.writer.Code)
		suite.Contains(suite.writer.Body.String(), "Failed to fetch product revenues")
	})
}
