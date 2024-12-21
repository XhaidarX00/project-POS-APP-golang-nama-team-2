package dashboardcontroller

import (
	"net/http"
	"project_pos_app/helper"
	"project_pos_app/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ControllerDashboard struct {
	Service *service.AllService
	Log     *zap.Logger
}

func NewControllerDashboard(service *service.AllService, log *zap.Logger) ControllerDashboard {
	return ControllerDashboard{Service: service, Log: log}
}

func (ctrl *ControllerDashboard) GetPopularProduct(ctx *gin.Context) {
	data, err := ctrl.Service.Dashboard.GetPopularProduct()
	if err != nil {
		helper.Responses(ctx, http.StatusInternalServerError, err.Error(), nil)
		ctx.Abort()
		return
	}
	helper.Responses(ctx, http.StatusOK, "Get Reservations success", data)
}
