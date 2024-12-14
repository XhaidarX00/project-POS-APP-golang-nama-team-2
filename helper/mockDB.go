package helper

import (
	"project_pos_app/repository"
	"project_pos_app/repository/notification"
	notifservice "project_pos_app/service/notif_service"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})
	gormDB, _ := gorm.Open(dialector, &gorm.Config{})
	return gormDB, mock
}

func InitService() (*notification.MockDB, notifservice.NotifServiceInterface) {
	mockDB := new(notification.MockDB)
	MockRepo := &repository.AllRepository{
		Notif: mockDB,
	}
	mockLogger := zap.NewNop()
	service := notifservice.NewNotifService(MockRepo, mockLogger)

	return mockDB, service
}
