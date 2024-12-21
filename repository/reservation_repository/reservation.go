package reservationrepository

import (
	"errors"
	"project_pos_app/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RepositoryReservation interface {
	FindReservations(date string) ([]model.Reservation, error)
	FindReservationsByTable(date string, tableNumber uint) ([]model.Reservation, error)
	FindReservation(reservation *model.Reservation) error
	Insert(reservation *model.Reservation) error
	Update(reservation *model.Reservation) error
}

type repositoryReservation struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewReservationRepository(db *gorm.DB, log *zap.Logger) RepositoryReservation {
	return &repositoryReservation{
		DB:  db,
		Log: log,
	}
}

func (r *repositoryReservation) FindReservations(date string) ([]model.Reservation, error) {
	var reservations []model.Reservation
	err := r.DB.Where("DATE(reservation_date) = ?", date).Find(&reservations).Error
	if err != nil {
		r.Log.Error("Failed to find all reservation", zap.Error(err))
		return nil, errors.New(" Internal Server Error")
	}
	return reservations, nil
}
func (r *repositoryReservation) FindReservationsByTable(date string, tableNumber uint) ([]model.Reservation, error) {
	var reservations []model.Reservation
	err := r.DB.Where("DATE(reservation_date) = ? and table_number = ?", date, tableNumber).Find(&reservations).Error
	if err != nil {
		r.Log.Error("Failed to find all reservation", zap.Error(err))
		return nil, errors.New(" Internal Server Error")
	}
	return reservations, nil
}
func (r *repositoryReservation) FindReservation(reservation *model.Reservation) error {
	err := r.DB.Where("id = ?", reservation.ID).Find(&reservation).Error
	if err != nil {
		r.Log.Error("Failed to find detail reservation", zap.Error(err))
		return errors.New(" Internal Server Error")
	}
	return nil
}
func (r *repositoryReservation) Insert(reservation *model.Reservation) error {
	err := r.DB.Create(&reservation).Error
	if err != nil {
		r.Log.Error("Failed to find insert reservation", zap.Error(err))
		return errors.New(" Bad Request")
	}
	return nil
}
func (r *repositoryReservation) Update(reservation *model.Reservation) error {
	err := r.DB.Save(&reservation).Error
	if err != nil {
		r.Log.Error("Failed to find update reservation", zap.Error(err))
		return errors.New(" Bad Request")
	}
	return nil
}
