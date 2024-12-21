package ordercontroller

import (
	"net/http"
	"project_pos_app/helper"

	"github.com/gin-gonic/gin"
)

// GetAllTable godoc
// @Summary Retrieve all tables
// @Description Get a list of all tables
// @Tags Tables
// @Accept json
// @Produce json
// @Security Authentication
// @Success 200 {object} model.SuccessResponse{data=[]model.Table} "Tables successfully retrieved"
// @Failure 404 {object} model.ErrorResponse "Tables not found"
// @Router /order/table [get]
func (oc *orderController) GetAllTable(c *gin.Context) {

	tables, err := oc.service.Order.GetAllTable()
	if err != nil {
		helper.Responses(c, http.StatusNotFound, "Error: "+err.Error(), nil)
		return
	}

	helper.Responses(c, http.StatusOK, "Table Successfully Retrieved", tables)
}
