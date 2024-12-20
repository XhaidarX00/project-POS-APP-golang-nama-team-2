package router

import (
	"project_pos_app/infra"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRoutes(ctx *infra.IntegrationContext) *gin.Engine {

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/login", ctx.Ctl.Auth.Login)

	NotificationRoutes(r, ctx)

	order := r.Group("/order")
	{
		order.GET("/", ctx.Ctl.Order.GetAllOrder)
		order.GET("/table", ctx.Ctl.Order.GetAllTable)
		order.GET("/payment", ctx.Ctl.Order.GetAllPayment)
		order.POST("/", ctx.Ctl.Order.CreateOrder)
		order.PUT("/:id", ctx.Ctl.Order.UpdateOrder)
		order.DELETE("/:id", ctx.Ctl.Order.DeleteOrder)
	}

	return r
}

func NotificationRoutes(r *gin.Engine, ctx *infra.IntegrationContext) {
	notifRoute := r.Group("/api")
	{
		notifRoute.POST("/notifications", ctx.Ctl.Notif.CreateNotifications)
		notifRoute.GET("/notifications", ctx.Ctl.Notif.GetAllNotifications)
		notifRoute.GET("/notifications/:id", ctx.Ctl.Notif.GetNotificationByID)
		notifRoute.PUT("/notifications/:id", ctx.Ctl.Notif.UpdateNotification)
		notifRoute.DELETE("/notifications/:id", ctx.Ctl.Notif.DeleteNotification)
		notifRoute.PUT("/notifications/mark-all-read", ctx.Ctl.Notif.MarkAllNotificationsAsRead)
	}
}
