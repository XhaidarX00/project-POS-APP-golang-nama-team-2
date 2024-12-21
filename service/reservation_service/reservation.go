package reservationservice

import (
	"errors"
	"project_pos_app/model"
	"project_pos_app/repository"
	"time"

	"go.uber.org/zap"
)

type ServiceReservation interface {
	GetAll(date string) ([]model.Reservation, error)
	GetById(reservation *model.Reservation) error
	Create(reservation *model.Reservation) error
	Edit(reservation *model.Reservation, form model.FormUpdate) error
}

type serviceReservation struct {
	Repo *repository.AllRepository
	Log  *zap.Logger
}

func NewRevenueService(repo *repository.AllRepository, log *zap.Logger) ServiceReservation {
	return &serviceReservation{
		Repo: repo,
		Log:  log,
	}
}

func (s *serviceReservation) GetAll(date string) ([]model.Reservation, error) {
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
func (s *serviceReservation) GetById(reservation *model.Reservation) error {
	err := s.Repo.Reservation.FindReservation(reservation)
	if err != nil {
		return err
	}
	reservation.Date = reservation.ReservationDate.Format("2006-01-02")
	reservation.Time = reservation.ReservationDate.Format("15:04:05")
	return nil
}
func (s *serviceReservation) Create(reservation *model.Reservation) error {
	parsedTime, err := time.Parse("2006-01-02 15:04:05", reservation.Date+" "+reservation.Time)
	if err != nil {
		s.Log.Error("Failed to parsed time", zap.Error(err))
		return errors.New(" Bad Request")
	}
	reservations, err := s.Repo.Reservation.FindReservationsByTable(reservation.Date, reservation.TableNumber)
	if err != nil {
		s.Log.Error("Failed to get data reservation", zap.Error(err))
		return errors.New(" Bad Request")
	}
	reservation.ReservationDate = parsedTime
	for _, v := range reservations {
		if v.ReservationDate.Hour() == reservation.ReservationDate.Hour() {
			s.Log.Error("Failed to Create Reservation", zap.String("Error", "Already Booked"))
			return errors.New(" Already Booked")
		}
	}
	return s.Repo.Reservation.Insert(reservation)
}
func (s *serviceReservation) Edit(reservation *model.Reservation, form model.FormUpdate) error {
	err := s.Repo.Reservation.FindReservation(reservation)
	if err != nil {
		return err
	}
	date := reservation.ReservationDate.Format("2006-01-02")
	reservations, err := s.Repo.Reservation.FindReservationsByTable(date, form.TableNumber)
	if err != nil {
		s.Log.Error("Failed to get data reservation", zap.Error(err))
		return errors.New(" Bad Request")
	}
	for _, v := range reservations {
		if v.ReservationDate.Hour() == reservation.ReservationDate.Hour() {
			s.Log.Error("Failed to Update Reservation", zap.String("Error", "Already Booked"))
			return errors.New(" Already Booked")
		}
	}
	reservation.TableNumber = form.TableNumber
	if form.Status != "" {
		reservation.Status = form.Status
	}
	return s.Repo.Reservation.Update(reservation)
}
