package service

import (
	"project_pos_app/repository"
	authservice "project_pos_app/service/auth_service"
	exampleservice "project_pos_app/service/example_service"

	"go.uber.org/zap"
)

type AllService struct {
	Example exampleservice.ExampleService
	Auth    authservice.AuthService
}

func NewAllService(repo *repository.AllRepository, log *zap.Logger) *AllService {
	return &AllService{
		Example: exampleservice.NewExampleService(repo, log),
		Auth:    authservice.NewManagementVoucherService(repo, log),
	}
}
