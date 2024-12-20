package orderservice

import "project_pos_app/model"

func (os *orderService) GetAllTable() ([]*model.Table, error) {

	table, err := os.Repo.Order.GetAllTable()
	if err != nil {
		return nil, err
	}

	return table, nil
}
