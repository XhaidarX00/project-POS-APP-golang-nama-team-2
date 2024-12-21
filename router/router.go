package router

import (
	"project_pos_app/helper"
	"project_pos_app/infra"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRoutes(ctx *infra.IntegrationContext) *gin.Engine {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("imagefile", helper.ImageFile)
	}

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/login", ctx.Ctl.Auth.Login)
	r.PATCH("/logout", ctx.Ctl.Superadmin.Logout)

	NotificationRoutes(r, ctx)
	RevenueRoutes(r, ctx)
	ProductRoutes(r, ctx)
	CategoryRoutes(r, ctx)

	ReservationRoutes(r, ctx)
	OrderRoutes(r, ctx)
	SuperAdmin(r, ctx)

	return r
}

func NotificationRoutes(r *gin.Engine, ctx *infra.IntegrationContext) {
	notifRoute := r.Group("/notification")
	{
		notifRoute.Use(ctx.Middleware.Access.AccessMiddleware())
		notifRoute.POST("/", ctx.Ctl.Notif.CreateNotifications)
		notifRoute.GET("/", ctx.Ctl.Notif.GetAllNotifications)
		notifRoute.GET("/:id", ctx.Ctl.Notif.GetNotificationByID)
		notifRoute.PUT("/:id", ctx.Ctl.Notif.UpdateNotification)
		notifRoute.DELETE("/:id", ctx.Ctl.Notif.DeleteNotification)
		notifRoute.PUT("/mark-all-read", ctx.Ctl.Notif.MarkAllNotificationsAsRead)
	}
}

func RevenueRoutes(r *gin.Engine, ctx *infra.IntegrationContext) {
	revenueRoute := r.Group("/revenue")
	{
		revenueRoute.Use(ctx.Middleware.Access.AccessMiddleware())
		revenueRoute.GET("/month", ctx.Ctl.Revenue.GetMonthlyRevenue)
		revenueRoute.GET("/products", ctx.Ctl.Revenue.GetProductRevenues)
		revenueRoute.GET("/status", ctx.Ctl.Revenue.GetTotalRevenueByStatus)
	}
}

func ProductRoutes(r *gin.Engine, ctx *infra.IntegrationContext) {
	productRoute := r.Group("/product")
	{
		productRoute.Use(ctx.Middleware.Access.AccessMiddleware())
		productRoute.GET("/", ctx.Ctl.Product.GetAllProducts)
		productRoute.GET("/:id", ctx.Ctl.Product.GetProductByID)
		productRoute.POST("/", ctx.Ctl.Product.CreateProduct)
		productRoute.PUT("/:id", ctx.Ctl.Product.UpdateProduct)
		productRoute.DELETE("/:id", ctx.Ctl.Product.DeleteProduct)
	}
}

func ReservationRoutes(r *gin.Engine, ctx *infra.IntegrationContext) {
	reservationRoute := r.Group("/reservation")
	{
		reservationRoute.Use(ctx.Middleware.Access.AccessMiddleware())
		reservationRoute.GET("/", ctx.Ctl.Reservation.GetAll)
		reservationRoute.GET("/:id", ctx.Ctl.Reservation.GetById)
		reservationRoute.POST("/", ctx.Ctl.Reservation.Create)
		reservationRoute.PUT("/:id", ctx.Ctl.Reservation.Edit)
	}
}

func OrderRoutes(r *gin.Engine, ctx *infra.IntegrationContext) {
	order := r.Group("/order")
	{
		order.Use(ctx.Middleware.Access.AccessMiddleware())
		order.GET("/", ctx.Ctl.Order.GetAllOrder)
		order.GET("/table", ctx.Ctl.Order.GetAllTable)
		order.GET("/payment", ctx.Ctl.Order.GetAllPayment)
		order.POST("/", ctx.Ctl.Order.CreateOrder)
		order.PUT("/:id", ctx.Ctl.Order.UpdateOrder)
		order.DELETE("/:id", ctx.Ctl.Order.DeleteOrder)
	}
}

func SuperAdmin(r *gin.Engine, ctx *infra.IntegrationContext) {
	superadmin := r.Group("/superadmin")
	{
		superadmin.Use(ctx.Middleware.Access.SuperAdminOnly())
		superadmin.GET("/", ctx.Ctl.Superadmin.ListDataAdmin)
		superadmin.PUT("/", ctx.Ctl.Superadmin.UpdateSuperadmin)
		superadmin.PUT("/:id", ctx.Ctl.Superadmin.UpdateAccessUser)
	}
}

func CategoryRoutes(r *gin.Engine, ctx *infra.IntegrationContext) {
	categoryRoute := r.Group("/api")
	{
		categoryRoute.GET("/categories", ctx.Ctl.Category.GetAllCategory)
		categoryRoute.GET("/categories/products", func(c *gin.Context) {
			ctx.Ctl.Product.GetAllProducts(c)
		})
		categoryRoute.GET("/categories/:id", ctx.Ctl.Category.GetCategoryByID)
		categoryRoute.POST("/categories", ctx.Ctl.Category.CreateCategory)
		categoryRoute.PUT("/categories/:id", ctx.Ctl.Category.UpdateCategory)
	}
}
