package repository

import (
	authrepository "project_pos_app/repository/auth_repository"
	examplerepository "project_pos_app/repository/example_repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AllRepository struct {
	Example examplerepository.ExampleRepository
	Auth    authrepository.AuthRepoInterface
}

func NewAllRepo(DB *gorm.DB, Log *zap.Logger) *AllRepository {
	return &AllRepository{
		Example: examplerepository.NewExampleRepo(DB, Log),
		Auth:    authrepository.NewManagementVoucherRepo(DB, Log),
	}
}
