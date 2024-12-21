package superadmincontroller

import (
	"net/http"
	"project_pos_app/helper"
	"project_pos_app/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SuperadminController interface {
	ListDataAdmin(c *gin.Context)
}

type superadminController struct {
	service *service.AllService
	log     *zap.Logger
}

func NewSuperadminController(service *service.AllService, log *zap.Logger) SuperadminController {
	return &superadminController{service, log}
}

func (sc *superadminController) ListDataAdmin(c *gin.Context) {

	admins, err := sc.service.Superadmin.ListDataAdmin()
	if err != nil {
		helper.Responses(c, http.StatusNotFound, "Error: "+err.Error(), nil)
		return
	}

	helper.Responses(c, http.StatusOK, "Succesfully Retrieved data", admins)

}
