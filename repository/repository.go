package repository

import (
	authrepository "project_pos_app/repository/auth_repository"
	examplerepository "project_pos_app/repository/example_repository"
	"project_pos_app/repository/notification"
	revenuerepository "project_pos_app/repository/revenue_repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AllRepository struct {
	Example examplerepository.ExampleRepository
	Auth    authrepository.AuthRepoInterface
	Notif   notification.NotifRepoInterface
	Revenue revenuerepository.RevenueRepositoryInterface
}

func NewAllRepo(DB *gorm.DB, Log *zap.Logger) *AllRepository {
	return &AllRepository{
		Example: examplerepository.NewExampleRepo(DB, Log),
		Auth:    authrepository.NewManagementVoucherRepo(DB, Log),
		Notif:   notification.NewNotifRepo(DB, Log),
		Revenue: revenuerepository.NewRevenueRepository(DB, Log),
	}
}
