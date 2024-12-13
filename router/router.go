package router

import (
	"project_pos_app/infra"

	"github.com/gin-gonic/gin"
)

func NewRoutes(ctx infra.IntegrationContext) *gin.Engine {

	r := gin.Default()

	r.POST("/login", ctx.Ctl.Auth.Login)

	return r
}
