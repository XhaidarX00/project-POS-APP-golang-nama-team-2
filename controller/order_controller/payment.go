package ordercontroller

import (
	"net/http"
	"project_pos_app/helper"

	"github.com/gin-gonic/gin"
)

func (oc *orderController) GetAllPayment(c *gin.Context) {

	payments, err := oc.service.Order.GetAllPayment()
	if err != nil {
		helper.Responses(c, http.StatusNotFound, "Error: "+err.Error(), nil)
		return
	}

	helper.Responses(c, http.StatusOK, "payment successfully Retrieved", payments)
}
