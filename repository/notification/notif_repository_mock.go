package notification

import (
	"project_pos_app/model"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return m.Called(query, args).Get(0).(*gorm.DB)
}

func (m *MockDB) Order(value interface{}) *gorm.DB {
	return m.Called(value).Get(0).(*gorm.DB)
}

func (m *MockDB) Find(dest interface{}) *gorm.DB {
	return m.Called(dest).Get(0).(*gorm.DB)
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return m.Called(dest, conds).Get(0).(*gorm.DB)
}

func (m *MockDB) Save(value interface{}) *gorm.DB {
	return m.Called(value).Get(0).(*gorm.DB)
}

func (m *MockDB) Model(value interface{}) *gorm.DB {
	return m.Called(value).Get(0).(*gorm.DB)
}

func (m *MockDB) Create(data model.Notification) error {
	args := m.Called(data)

	if err := args.Error(0); err != nil {
		return err
	}

	return nil
}

func (m *MockDB) GetAll(data *[]model.Notification, status string) error {
	args := m.Called()

	if notification := args.Get(0); notification != nil {
		*data = notification.([]model.Notification)
		return nil
	}

	return args.Error(1)
}

func (m *MockDB) FindByID(id int) (model.Notification, error) {
	args := m.Called(id) // id = 9999
	if args.Get(0) != nil {
		return args.Get(0).(model.Notification), args.Error(1)
	}

	return model.Notification{}, args.Error(1)
}

func (m *MockDB) Update(data *model.Notification, id int) error {
	args := m.Called(data, id)

	if args.Get(0) != nil {
		updatedNotif := args.Get(0).(*model.Notification)
		updatedNotif.Status = "reade"
		*data = *updatedNotif
	}

	return args.Error(0)
}

func (m *MockDB) Delete(id int) error {
	data, err := m.FindByID(id)
	if err != nil {
		return err
	}

	args := m.Called(data.ID)
	if notification := args.Get(0); notification != nil {
		return nil
	}

	return args.Error(1)
}

func (m *MockDB) MarkAllAsRead() error {
	args := m.Called()

	if notification := args.Get(0); notification != nil {
		for _, v := range notification.([]model.Notification) {
			if v.Status == "new" {
				v.Status = "readed"
			}
		}

		return nil
	}

	return args.Error(1)
}
