package orderrepository

import (
	"errors"
	"fmt"
	"project_pos_app/model"
	"strconv"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderRepository interface {
	GetAllOrder(search, status string) ([]*model.OrderResponse, error)
	CreateOrder(order *model.Order) error
	UpdateOrder(id int, order *model.Order) error
	GetAllTable() ([]*model.Table, error)
	GetAllPayment() ([]*model.Payment, error)
	DeleteOrder(id int) error
}

type orderRepository struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewOrderRepo(DB *gorm.DB, Log *zap.Logger) OrderRepository {
	return &orderRepository{DB, Log}
}

func (or *orderRepository) GetAllOrder(search, status string) ([]*model.OrderResponse, error) {
	orders := []*model.OrderResponse{}

	result := or.DB.Table("orders as o").
		Select("o.id, o.table_id, o.customer_name, o.status, o.created_at as order_date").
		Where("o.deleted_at IS NULL")

	if search != "" {
		if _, err := strconv.Atoi(search); err == nil {
			result = result.Where("o.id = ? OR o.customer_name ILIKE ?", search, "%"+search+"%")
		} else {
			result = result.Where("o.customer_name ILIKE ?", "%"+search+"%")
		}
	}

	if status != "" {
		result.Where("o.status ILIKE ?", "%"+status+"%")
	}

	if err := result.Scan(&orders).Error; err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return []*model.OrderResponse{}, nil
	}

	orderIDs := []int{}
	for _, order := range orders {
		orderIDs = append(orderIDs, int(order.ID))
	}

	orderProducts := []*model.OrderProductResponse{}
	if err := or.DB.Table("order_products as po").
		Select("po.order_id, po.qty, p.name AS item, p.price").
		Joins("JOIN products as p ON p.id = po.product_id").
		Where("po.order_id IN ?", orderIDs).
		Scan(&orderProducts).Error; err != nil {
		return nil, err
	}

	for _, order := range orders {
		order.SubTotal = 0
		order.OrderProduct = []model.OrderProductResponse{}

		for _, op := range orderProducts {
			if op.OrderID == int(order.ID) {
				order.SubTotal += int(op.Price * float64(op.Qty))
				order.OrderProduct = append(order.OrderProduct, *op)
			}
		}
	}

	return orders, nil
}

func (or *orderRepository) CreateOrder(order *model.Order) error {
	return or.DB.Transaction(func(tx *gorm.DB) error {

		if err := or.findTable(int(order.TableID)); err != nil {
			return err
		}

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Table{}).
			Where("id = ?", order.TableID).
			Update("is_book", true).Error; err != nil {
			return err
		}

		if len(order.OrderProducts) == 0 {
			return fmt.Errorf("order has no products")
		}

		for _, op := range order.OrderProducts {

			product := model.Product{}

			if err := tx.First(&product, "id = ?", op.ProductID).Error; err != nil {
				return fmt.Errorf("failed to fetch product with id %d: %v", op.ProductID, err)
			}

			if product.Qty <= 0 {
				return fmt.Errorf("product with id %d stock habis", op.ProductID)
			}

			if product.Qty < int(op.Qty) {
				return fmt.Errorf("stock product with id %d less than qty", op.ProductID)
			}

			if err := tx.Model(&model.Product{}).
				Where("id = ?", op.ProductID).
				Update("qty", gorm.Expr("qty - ?", op.Qty)).Error; err != nil {
				return err
			}

			op.OrderID = order.ID
			if err := tx.Create(&op).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (or *orderRepository) UpdateOrder(id int, order *model.Order) error {
	return or.DB.Transaction(func(tx *gorm.DB) error {

		existingOrder := model.Order{}
		if err := tx.First(&existingOrder, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("order with id %d does not exist", id)
			}
			return err
		}

		if existingOrder.TableID != order.TableID {
			if err := or.findTable(int(order.TableID)); err != nil {
				return err
			}

			if err := tx.Model(&model.Table{}).Where("id = ?", existingOrder.TableID).Update("is_book", false).Error; err != nil {
				return fmt.Errorf("failed to release old table: %v", err)
			}

			if err := tx.Model(&model.Table{}).Where("id = ?", order.TableID).Update("is_book", true).Error; err != nil {
				return fmt.Errorf("failed to book new table: %v", err)
			}
		}

		existingOrderProducts := []model.OrderProduct{}
		if err := tx.Where("order_id = ?", id).Find(&existingOrderProducts).Error; err != nil {
			return fmt.Errorf("failed to retrieve existing order products: %v", err)
		}

		for _, existingOrderProduct := range existingOrderProducts {
			if err := updateStock(tx, int(existingOrderProduct.ProductID), existingOrderProduct.Qty); err != nil {
				return err
			}
		}

		if order.Status != "canceled" {
			if err := tx.Where("order_id = ?", id).Delete(&model.OrderProduct{}).Error; err != nil {
				return fmt.Errorf("failed to delete order products: %v", err)
			}

			var totalAmount float64
			for _, orderProduct := range order.OrderProducts {
				product := model.Product{}
				if err := tx.First(&product, "id = ?", orderProduct.ProductID).Error; err != nil {
					return fmt.Errorf("failed to find product with id %d: %v", orderProduct.ProductID, err)
				}

				subtotal := float64(orderProduct.Qty) * product.Price
				totalAmount += subtotal

				if err := tx.Model(&model.Product{}).
					Where("id = ?", orderProduct.ProductID).
					Update("qty", gorm.Expr("qty - ?", orderProduct.Qty)).Error; err != nil {
					return fmt.Errorf("failed to deduct stock for product %d: %v", orderProduct.ProductID, err)
				}

				orderProduct.ID = 0
				orderProduct.OrderID = uint(id)

				if err := tx.Create(&orderProduct).Error; err != nil {
					return fmt.Errorf("failed to recreate order product: %v", err)
				}
			}

			order.TotalAmount = totalAmount + (totalAmount * order.Tax / 100)

		}

		if order.Status == "completed" || order.Status == "canceled" {
			tx.Model(&model.Table{}).Where("id = ?", order.TableID).Update("is_book", false)
		}

		if err := tx.Model(&model.Order{}).Where("id = ?", id).Updates(&order).Error; err != nil {
			return fmt.Errorf("failed to update order: %v", err)
		}

		return nil
	})
}

func (or *orderRepository) DeleteOrder(id int) error {
	result := or.DB.Where("id = ?", id).Delete(&model.Order{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete order: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("order with ID %d not found", id)
	}

	return nil
}

func updateStock(tx *gorm.DB, id, qty int) error {

	if err := tx.Model(&model.Product{}).
		Where("id = ?", id).
		Update("qty", gorm.Expr("qty + ?", qty)).Error; err != nil {
		return fmt.Errorf("failed to restore stock for product %d: %v", id, err)
	}

	return nil
}
