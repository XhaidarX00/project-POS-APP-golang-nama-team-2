package notifservice

import (
	"fmt"
	"project_pos_app/model"
	"project_pos_app/repository"
	"time"

	"go.uber.org/zap"
)

type NotifServiceInterface interface {
	CreateNotification(data model.Notification) error
	GetAllNotifications(status string) ([]model.Notification, error)
	GetNotificationByID(id int) (model.Notification, error)
	UpdateNotification(id int) error
	DeleteNotification(id int) error
	MarkAllNotificationsAsRead() error
}

type notifService struct {
	Repo *repository.AllRepository
	Log  *zap.Logger
}

func NewNotifService(repo *repository.AllRepository, log *zap.Logger) NotifServiceInterface {
	return &notifService{
		Repo: repo,
		Log:  log,
	}
}

func (s *notifService) CreateNotification(data model.Notification) error {
	now := time.Now()
	if data.Status == "" {
		data.Status = "new"
	}
	data.CreatedAt = now
	data.UpdatedAt = now

	if data.Title == "" {
		return fmt.Errorf("title should not none")
	}

	return s.Repo.Notif.Create(data)
}

func (s *notifService) GetAllNotifications(status string) ([]model.Notification, error) {
	var notifications []model.Notification
	if err := s.Repo.Notif.GetAll(&notifications, status); err != nil {
		s.Log.Error("Failed to fetch notifications", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch notifications: %w", err)
	}
	return notifications, nil
}

func (s *notifService) GetNotificationByID(id int) (model.Notification, error) {
	notification, err := s.Repo.Notif.FindByID(id)
	if err != nil {
		s.Log.Error("Failed to fetch notification by ID", zap.Error(err))
		return model.Notification{}, err
	}

	if notification.Title == "" {
		return model.Notification{}, fmt.Errorf("title is reqired")
	}

	return notification, nil
}

func (s *notifService) UpdateNotification(id int) error {
	notification, err := s.Repo.Notif.FindByID(id)
	if err != nil {
		s.Log.Error("Failed to find notification for update", zap.Error(err))
		return fmt.Errorf("failed to find notification for update: %w", err)
	}

	if notification.Title == "" {
		return fmt.Errorf("title is reqired")
	}

	notification.Status = "readed"
	if err := s.Repo.Notif.Update(&notification, id); err != nil {
		s.Log.Error("Failed to update notification", zap.Error(err))
		return fmt.Errorf("failed to update notification: %w", err)
	}

	return nil
}

func (s *notifService) DeleteNotification(id int) error {
	notif, err := s.Repo.Notif.FindByID(id)
	if err != nil {
		s.Log.Error("Failed to find notification for update", zap.Error(err))
		return fmt.Errorf("failed to find notification for update: %w", err)
	}

	if notif.Title == "" {
		return fmt.Errorf("title is reqired")
	}

	if err := s.Repo.Notif.Delete(id); err != nil {
		s.Log.Error("Failed to delete notification", zap.Error(err))
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

func (s *notifService) MarkAllNotificationsAsRead() error {
	if err := s.Repo.Notif.MarkAllAsRead(); err != nil {
		s.Log.Error("Failed to mark all notifications as read", zap.Error(err))
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}
	return nil
}
