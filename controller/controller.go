package controller

import (
	authcontroller "project_pos_app/controller/auth_controller"
	examplecontroller "project_pos_app/controller/example_controller"
	"project_pos_app/database"
	"project_pos_app/service"

	"go.uber.org/zap"
)

type AllController struct {
	Example examplecontroller.ExampleController
	Auth    authcontroller.AuthHadler
}

func NewAllController(service *service.AllService, log *zap.Logger, cfg *database.Cache) AllController {
	return AllController{
		Example: examplecontroller.NewExampleController(service, log),
		Auth:    authcontroller.NewUserHandler(service, log, cfg),
	}
}
