package controller

import (
	authcontroller "project_pos_app/controller/auth_controller"
	categorycontroller "project_pos_app/controller/category_controller"
	dashboardcontroller "project_pos_app/controller/dashboard_controller"
	notifcontroller "project_pos_app/controller/notif_controller"
	productcontroller "project_pos_app/controller/product_controller"
	reservationcontroller "project_pos_app/controller/reservation_controller"
	revenuecontroller "project_pos_app/controller/revenue_controller"
	superadmincontroller "project_pos_app/controller/superadmin_controller"

	// productcontroller "project_pos_app/controller/product_controller"
	ordercontroller "project_pos_app/controller/order_controller"
	"project_pos_app/database"
	"project_pos_app/service"

	"go.uber.org/zap"
)

type AllController struct {
	Auth        authcontroller.AuthHadler
	Notif       notifcontroller.NotifController
	Revenue     revenuecontroller.RevenueController
	Product     productcontroller.ProductController
	Order       ordercontroller.OrderController
	Superadmin  superadmincontroller.SuperadminController
	Category    categorycontroller.CategoryController
	Reservation reservationcontroller.ControllerReservation
	Dashboard   dashboardcontroller.ControllerDashboard
}

func NewAllController(service *service.AllService, log *zap.Logger, cfg *database.Cache) AllController {
	return AllController{
		Auth:        authcontroller.NewUserHandler(service, log, cfg),
		Notif:       notifcontroller.NewNotifController(service, log),
		Revenue:     revenuecontroller.NewRevenueController(service, log),
		Product:     *productcontroller.NewProductController(service, log),
		Order:       ordercontroller.NewOrderController(service, log),
		Superadmin:  superadmincontroller.NewSuperadminController(service, log),
		Category:    *categorycontroller.NewCategoryController(service, log),
		Reservation: reservationcontroller.NewControllerReservation(service, log),
		Dashboard:   dashboardcontroller.NewControllerDashboard(service, log),
	}
}
