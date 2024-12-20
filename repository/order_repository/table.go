package orderrepository

import (
	"errors"
	"fmt"
	"project_pos_app/model"

	"gorm.io/gorm"
)

func (or *orderRepository) GetAllTable() ([]*model.Table, error) {

	table := []*model.Table{}
	if err := or.DB.Find(&table).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("table not found")
		}
		return nil, err
	}

	return table, nil
}

func (or *orderRepository) findTable(id int) error {

	table := model.Table{}

	if err := or.DB.First(&table, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("table %d does not exist", id)
		}

		return err
	}

	if table.IsBook {
		return fmt.Errorf("table %d is already booked", id)
	}

	return nil
}
