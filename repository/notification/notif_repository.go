package notification

import (
	"fmt"
	"project_pos_app/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NotifRepoInterface interface {
	Create(data model.Notification) error
	GetAll(data *[]model.Notification, status string) error
	FindByID(id int) (model.Notification, error)
	Update(data *model.Notification, id int) error
	Delete(id int) error
	MarkAllAsRead() error
}

type notifRepo struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewNotifRepo(db *gorm.DB, log *zap.Logger) NotifRepoInterface {
	return &notifRepo{DB: db, Log: log}
}

func (r *notifRepo) Create(data model.Notification) error {
	if err := r.DB.Create(&data).Error; err != nil {
		return fmt.Errorf("failed to create notification")
	}

	return nil
}

func (r *notifRepo) GetAll(data *[]model.Notification, status string) error {
	query := r.DB
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("created_at desc").Find(&data).Error; err != nil {
		r.Log.Error("Error GetnotifAll : ", zap.Error(err))
		return fmt.Errorf("error get list notif : %s", err.Error())
	}

	return nil
}

func (r *notifRepo) FindByID(id int) (model.Notification, error) {
	var data model.Notification
	err := r.DB.First(&data, id).Error
	if err != nil {
		r.Log.Error("Failed get data", zap.Error(err))
		return model.Notification{}, fmt.Errorf("notification not found")
	}
	return data, nil
}

func (r *notifRepo) Update(data *model.Notification, id int) error {
	var err error
	*data, err = r.FindByID(id)
	if err != nil {
		r.Log.Error("Error UpdateNotif : ", zap.Error(err))
		return fmt.Errorf("error update notif : %s", err.Error())
	}

	data.Status = "readed"
	if err := r.DB.Save(&data).Error; err != nil {
		r.Log.Error("Error UpdateNotif : ", zap.Error(err))
		return fmt.Errorf("error update notif : %s", err.Error())
	}

	return nil
}

func (r *notifRepo) Delete(id int) error {
	if err := r.DB.Delete(&model.Notification{}, id).Error; err != nil {
		r.Log.Error("Error Delete Notif : ", zap.Error(err))
		return fmt.Errorf("error delete notif : %s", err.Error())
	}

	return nil
}

func (r *notifRepo) MarkAllAsRead() error {
	if err := r.DB.Model(&model.Notification{}).Where("status = ?", "new").Update("status", "readed").Error; err != nil {
		r.Log.Error("Error Update Status All Notif : ", zap.Error(err))
		return fmt.Errorf("error update status all notif")
	}

	return nil
}
