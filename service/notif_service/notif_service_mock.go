package notifservice

import (
	mocktesting "project_pos_app/mock_testing"
	"project_pos_app/model"
	"time"

	"go.uber.org/zap"
)

// MockNotifServiceInterface mendefinisikan kontrak untuk mock service notifikasi
type MockNotifServiceInterface interface {
	CreateNotification(data model.Notification) error
	GetAllNotifications(status string) ([]model.Notification, error)
	GetNotificationByID(id int) (*model.Notification, error)
	DeleteNotification(id int) error
	MarkAllNotificationsAsRead() error
	UpdateNotification(int) error
}

// MockNotifService implementasi mock service untuk notifikasi
type MockNotifService struct {
	Repo *mocktesting.MockDB
	Log  *zap.Logger
}

// NewMockNotifService membuat instance baru dari mock service notifikasi
func NewMockNotifService(repo *mocktesting.MockDB, log *zap.Logger) MockNotifServiceInterface {
	return &MockNotifService{
		Repo: repo,
		Log:  log,
	}
}

// CreateNotification membuat notifikasi baru
func (m *MockNotifService) CreateNotification(data model.Notification) error {
	m.Log.Info("Creating Mock notification", zap.Any("notification", data))
	args := m.Repo.Called(data)
	return args.Error(0)
}

// GetAllNotifications mengambil semua notifikasi berdasarkan status
func (m *MockNotifService) GetAllNotifications(status string) ([]model.Notification, error) {
	now := time.Now()
	defaultData := []model.Notification{
		{
			ID:        1,
			Title:     "Notification 1",
			Message:   "First test notification",
			CreatedAt: now,
			UpdatedAt: now,
			Status:    status,
		},
		{
			ID:        2,
			Title:     "Notification 2",
			Message:   "Second test notification",
			CreatedAt: now,
			UpdatedAt: now,
			Status:    status,
		},
	}

	m.Log.Info("Fetching notifications", zap.String("status", status))

	err := m.Repo.GetAll(&defaultData, status)
	if err != nil {
		m.Log.Error("Failed to get notifications", zap.Error(err))
		return nil, err
	}

	return defaultData, nil
}

// GetNotificationByID mengambil notifikasi berdasarkan ID
func (m *MockNotifService) GetNotificationByID(id int) (*model.Notification, error) {
	m.Log.Info("Fetching notification by ID", zap.Int("id", id))
	return m.Repo.FindByID(id)
}

// UpdateNotification memperbarui notifikasi
func (m *MockNotifService) UpdateNotification(id int) error {
	data := &model.Notification{
		ID:      1,
		Title:   "Testing",
		Message: "Test notification",
		Status:  "new",
	}
	m.Log.Info("Updating notification", zap.Int("id", id))
	return m.Repo.Update(data, id)
}

// DeleteNotification menghapus notifikasi
func (m *MockNotifService) DeleteNotification(id int) error {
	m.Log.Info("Deleting notification", zap.Int("id", id))
	return m.Repo.Delete(id)
}

// MarkAllNotificationsAsRead menandai semua notifikasi sebagai sudah dibaca
func (m *MockNotifService) MarkAllNotificationsAsRead() error {
	m.Log.Info("Marking all notifications as read")
	return m.Repo.MarkAllAsRead()
}
