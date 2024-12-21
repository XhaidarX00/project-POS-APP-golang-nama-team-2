package dashboardservice

import (
	"project_pos_app/model"
	"project_pos_app/repository"
	"time"

	"go.uber.org/zap"
)

type ServiceDashboard interface {
	GetPopularProduct() ([]model.Product, error)
	GetAll(date string) ([]model.Reservation, error)
}

type serviceDashboard struct {
	Repo *repository.AllRepository
	Log  *zap.Logger
}

func NewRevenueService(repo *repository.AllRepository, log *zap.Logger) ServiceDashboard {
	return &serviceDashboard{
		Repo: repo,
		Log:  log,
	}
}

func (s *serviceDashboard) GetPopularProduct() ([]model.Product, error) {
	// products, err := s.Repo.Dashboard.FindPopularProduct()
	// if err != nil {
	// 	return nil, err
	// }
	// return products, nil
	return s.Repo.Dashboard.FindPopularProduct()
}
func (s *serviceDashboard) GetAll(date string) ([]model.Reservation, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	reservations, err := s.Repo.Reservation.FindReservations(date)
	if err != nil {
		return nil, err
	}
	for i, _ := range reservations {
		reservations[i].Date = reservations[i].ReservationDate.Format("2006-01-02")
		reservations[i].Time = reservations[i].ReservationDate.Format("15:04:05")
	}
	return reservations, nil
}
