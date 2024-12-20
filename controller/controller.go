package controller

import (
	authcontroller "project_pos_app/controller/auth_controller"
	examplecontroller "project_pos_app/controller/example_controller"
	notifcontroller "project_pos_app/controller/notif_controller"
	productcontroller "project_pos_app/controller/product_controller"
	revenuecontroller "project_pos_app/controller/revenue_controller"

	// productcontroller "project_pos_app/controller/product_controller"
	"project_pos_app/database"
	"project_pos_app/service"

	"go.uber.org/zap"
)

type AllController struct {
	Example examplecontroller.ExampleController
	Auth    authcontroller.AuthHadler
	Notif   notifcontroller.NotifController
	Revenue revenuecontroller.RevenueController
	Product productcontroller.ProductController
}

func NewAllController(service *service.AllService, log *zap.Logger, cfg *database.Cache) AllController {
	return AllController{
		Example: examplecontroller.NewExampleController(service, log),
		Auth:    authcontroller.NewUserHandler(service, log, cfg),
		Notif:   notifcontroller.NewNotifController(service, log),
		Revenue: revenuecontroller.NewRevenueController(service, log),
		Product: *productcontroller.NewProductController(service, log),
	}
}
