package ordercontroller

import (
	"net/http"
	"project_pos_app/helper"

	"github.com/gin-gonic/gin"
)

// GetAllPayment godoc
// @Summary Retrieve all payment methods
// @Description Get a list of all payment methods
// @Tags Payments
// @Accept json
// @Produce json
// @Security Authentication
// @Success 200 {object} model.SuccessResponse{data=[]model.Payment} "Payment methods successfully retrieved"
// @Failure 404 {object} model.ErrorResponse "Payment methods not found"
// @Router /order/payment [get]
func (oc *orderController) GetAllPayment(c *gin.Context) {

	payments, err := oc.service.Order.GetAllPayment()
	if err != nil {
		helper.Responses(c, http.StatusNotFound, "Error: "+err.Error(), nil)
		return
	}

	helper.Responses(c, http.StatusOK, "payment successfully Retrieved", payments)
}
