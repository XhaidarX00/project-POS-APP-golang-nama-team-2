package helper

import (
	mocktesting "project_pos_app/mock_testing"
	"project_pos_app/repository"
	"project_pos_app/service"
	notifservice "project_pos_app/service/notif_service"
	revenueservice "project_pos_app/service/revenue_service"

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

// type ServiceMock struct {
// 	Notif   notifservice.NotifServiceInterface
// 	Revenue revenueservice.RevenueServiceInterface
// }

func InitService() (*mocktesting.MockDB, *service.AllService) {
	mockDB := new(mocktesting.MockDB)
	MockRepo := &repository.AllRepository{
		Notif:   mockDB,
		Revenue: mockDB,
	}
	mockLogger := zap.NewNop()
	serviceNotif := notifservice.NewNotifService(MockRepo, mockLogger)
	serviceRevenue := revenueservice.NewRevenueService(MockRepo, mockLogger)

	var service service.AllService
	service.Notif = serviceNotif
	service.Revenue = serviceRevenue

	return mockDB, &service
}
