package orderrepository

import (
	"errors"
	"project_pos_app/model"

	"gorm.io/gorm"
)

func (or *orderRepository) GetAllPayment() ([]*model.Payment, error) {

	payment := []*model.Payment{}
	if err := or.DB.Find(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	return payment, nil
}
