package repository

import (
	accessrepository "project_pos_app/repository/access_repository"
	authrepository "project_pos_app/repository/auth_repository"
	examplerepository "project_pos_app/repository/example_repository"
	"project_pos_app/repository/notification"
	orderrepository "project_pos_app/repository/order_repository"
	productrepository "project_pos_app/repository/product"
	profilesuperadmin "project_pos_app/repository/profile_superadmin"
	revenuerepository "project_pos_app/repository/revenue_repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AllRepository struct {
	Example    examplerepository.ExampleRepository
	Auth       authrepository.AuthRepoInterface
	Notif      notification.NotifRepoInterface
	Revenue    revenuerepository.RevenueRepositoryInterface
	Product    productrepository.ProductRepo
	Order      orderrepository.OrderRepository
	Superadmin profilesuperadmin.SuperadminRepo
	Access     accessrepository.AccessRepository
}

func NewAllRepo(DB *gorm.DB, Log *zap.Logger) *AllRepository {
	return &AllRepository{
		Example:    examplerepository.NewExampleRepo(DB, Log),
		Auth:       authrepository.NewManagementVoucherRepo(DB, Log),
		Notif:      notification.NewNotifRepo(DB, Log),
		Revenue:    revenuerepository.NewRevenueRepository(DB, Log),
		Product:    productrepository.NewProductRepo(DB, Log),
		Order:      orderrepository.NewOrderRepo(DB, Log),
		Superadmin: profilesuperadmin.NewSuperadmin(DB, Log),
		Access:     accessrepository.NewAccessRepository(DB, Log),
	}
}
