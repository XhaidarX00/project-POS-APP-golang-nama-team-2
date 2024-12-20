package orderservice

import (
	"project_pos_app/model"
	"project_pos_app/repository"

	"go.uber.org/zap"
)

type OrderService interface {
	GetAllOrder(search, status string) ([]*model.OrderResponse, error)
	CreateOrder(order *model.Order) error
	UpdateOrder(id int, order *model.Order) error
	GetAllTable() ([]*model.Table, error)
	GetAllPayment() ([]*model.Payment, error)
	DeleteOrder(id int) error
}

type orderService struct {
	Repo *repository.AllRepository
	Log  *zap.Logger
}

func NewOrderService(Repo *repository.AllRepository, Log *zap.Logger) OrderService {
	return &orderService{Repo, Log}
}

func (os *orderService) GetAllOrder(search, status string) ([]*model.OrderResponse, error) {

	orders, err := os.Repo.Order.GetAllOrder(search, status)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (os *orderService) CreateOrder(order *model.Order) error {

	order.Tax = 12
	order.Status = "In Process"

	if err := os.Repo.Order.CreateOrder(order); err != nil {
		return err
	}

	return nil
}

func (os *orderService) UpdateOrder(id int, order *model.Order) error {

	order.Tax = 12

	if order.PaymentMethod != "" && order.Status != "cancelled" {
		order.Status = "completed"
	}

	if err := os.Repo.Order.UpdateOrder(id, order); err != nil {
		return err
	}

	return nil
}

func (os *orderService) DeleteOrder(id int) error {

	if err := os.Repo.Order.DeleteOrder(id); err != nil {
		return err
	}

	return nil
}
