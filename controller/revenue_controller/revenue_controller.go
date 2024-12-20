package revenuecontroller

import (
	"net/http"
	"project_pos_app/helper"
	"project_pos_app/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RevenueController struct {
	Service *service.AllService
	Log     *zap.Logger
}

func NewRevenueController(service *service.AllService, log *zap.Logger) RevenueController {
	return RevenueController{Service: service, Log: log}
}

// GetTotalRevenueByStatus godoc
// @Summary Fetch total revenue by status
// @Description Get total revenue grouped by order status
// @Tags Revenue
// @Produce json
// @Success 200 {object} model.SuccessResponse{data=[]model.RevenueByStatus} "Fetch total revenue by status successfully"
// @Failure 500 {object} model.ErrorResponse "Failed to fetch total revenue by status"
// @Router /api/revenue/status [get]
func (ctrl *RevenueController) GetTotalRevenueByStatus(ctx *gin.Context) {
	data, err := ctrl.Service.Revenue.FetchTotalRevenueByStatus()
	if err != nil {
		ctrl.Log.Error("Failed to fetch total revenue by status", zap.Error(err))
		helper.Responses(ctx, http.StatusInternalServerError, "Failed to fetch total revenue by status: "+err.Error(), nil)
		ctx.Abort()
		return
	}

	helper.Responses(ctx, http.StatusOK, "Fetch total revenue by status successfully", data)
}

// GetMonthlyRevenue godoc
// @Summary Fetch monthly revenue
// @Description Get total revenue grouped by month
// @Tags Revenue
// @Produce json
// @Success 200 {object} model.SuccessResponse{data=[]model.MonthlyRevenue} "Fetch monthly revenue successfully"
// @Failure 500 {object} model.ErrorResponse "Failed to fetch monthly revenue"
// @Router /api/revenue/month [get]
func (ctrl *RevenueController) GetMonthlyRevenue(ctx *gin.Context) {
	data, err := ctrl.Service.Revenue.FetchMonthlyRevenue()
	if err != nil {
		ctrl.Log.Error("Failed to fetch monthly revenue", zap.Error(err))
		helper.Responses(ctx, http.StatusInternalServerError, "Failed to fetch monthly revenue: "+err.Error(), nil)
		ctx.Abort()
		return
	}

	helper.Responses(ctx, http.StatusOK, "Fetch monthly revenue successfully", data)
}

// GetProductRevenues godoc
// @Summary Fetch product revenues
// @Description Get revenue details for all products
// @Tags Revenue
// @Produce json
// @Success 200 {object} model.SuccessResponse{data=[]model.ProductRevenue} "Fetch product revenues successfully"
// @Failure 500 {object} model.ErrorResponse "Failed to fetch product revenues"
// @Router /api/revenue/products [get]
func (ctrl *RevenueController) GetProductRevenues(ctx *gin.Context) {
	data, err := ctrl.Service.Revenue.FetchProductRevenues()
	if err != nil {
		ctrl.Log.Error("Failed to fetch product revenues", zap.Error(err))
		helper.Responses(ctx, http.StatusInternalServerError, "Failed to fetch product revenues: "+err.Error(), nil)
		ctx.Abort()
		return
	}

	helper.Responses(ctx, http.StatusOK, "Fetch product revenues successfully", data)
}
