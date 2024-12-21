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
	RevenueRoutes(r, ctx)
	ProductRoutes(r, ctx)
	ReservationRoutes(r, ctx)
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

func RevenueRoutes(r *gin.Engine, ctx *infra.IntegrationContext) {
	revenueRoute := r.Group("/api")
	{
		revenueRoute.GET("/revenue/month", ctx.Ctl.Revenue.GetMonthlyRevenue)
		revenueRoute.GET("/revenue/products", ctx.Ctl.Revenue.GetProductRevenues)
		revenueRoute.GET("/revenue/status", ctx.Ctl.Revenue.GetTotalRevenueByStatus)
	}
}

func ProductRoutes(r *gin.Engine, ctx *infra.IntegrationContext) {
	productRoute := r.Group("/api")
	{
		productRoute.GET("/products", ctx.Ctl.Product.GetAllProducts)
		productRoute.GET("/products/:id", ctx.Ctl.Product.GetProductByID)
		productRoute.POST("/products", ctx.Ctl.Product.CreateProduct)
		productRoute.PUT("/products/:id", ctx.Ctl.Product.UpdateProduct)
		productRoute.DELETE("/product/:id", ctx.Ctl.Product.DeleteProduct)
	}
}
func ReservationRoutes(r *gin.Engine, ctx *infra.IntegrationContext) {
	reservationRoute := r.Group("/api")
	{
		reservationRoute.GET("/reservation", ctx.Ctl.Reservation.GetAll)
		reservationRoute.GET("/reservation/:id", ctx.Ctl.Reservation.GetById)
		reservationRoute.POST("/reservation", ctx.Ctl.Reservation.Create)
		reservationRoute.PUT("/reservation/:id", ctx.Ctl.Reservation.Edit)
	}
}
