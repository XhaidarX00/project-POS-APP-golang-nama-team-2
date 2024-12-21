package ordercontroller

import (
	"net/http"
	"project_pos_app/helper"

	"github.com/gin-gonic/gin"
)

func (oc *orderController) GetAllTable(c *gin.Context) {

	tables, err := oc.service.Order.GetAllTable()
	if err != nil {
		helper.Responses(c, http.StatusNotFound, "Error: "+err.Error(), nil)
		return
	}

	helper.Responses(c, http.StatusOK, "Table Successfully Retrieved", tables)
}
