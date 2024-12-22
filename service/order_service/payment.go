package orderservice

import "project_pos_app/model"

func (os *orderService) GetAllPayment() ([]*model.Payment, error) {

	payments, err := os.Repo.Order.GetAllPayment()
	if err != nil {
		return nil, err
	}

	return payments, nil
}
