package service

import (
	"project_pos_app/repository"
	authservice "project_pos_app/service/auth_service"
	categoryservice "project_pos_app/service/category_service"
	exampleservice "project_pos_app/service/example_service"
	notifservice "project_pos_app/service/notif_service"
	orderservice "project_pos_app/service/order_service"
	productservice "project_pos_app/service/product_service"
	revenueservice "project_pos_app/service/revenue_service"
	superadminservice "project_pos_app/service/superadmin_service"

	"go.uber.org/zap"
)

type AllService struct {
	Example    exampleservice.ExampleService
	Auth       authservice.AuthService
	Notif      notifservice.NotifServiceInterface
	Revenue    revenueservice.RevenueServiceInterface
	Product    productservice.ProductService
	Order      orderservice.OrderService
	Superadmin superadminservice.SuperadminService
	Category   categoryservice.CategoryService
}

func NewAllService(repo *repository.AllRepository, log *zap.Logger) *AllService {
	return &AllService{
		Example:    exampleservice.NewExampleService(repo, log),
		Auth:       authservice.NewManagementVoucherService(repo, log),
		Notif:      notifservice.NewNotifService(repo, log),
		Revenue:    revenueservice.NewRevenueService(repo, log),
		Product:    productservice.NewProductService(repo, log),
		Order:      orderservice.NewOrderService(repo, log),
		Superadmin: superadminservice.NewSuperadminService(repo, log),
		Category:   categoryservice.NewCategoryService(repo, log),
	}
}
