package notifcontroller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	notifcontroller "project_pos_app/controller/notif_controller"
	"project_pos_app/helper"
	mocktesting "project_pos_app/mock_testing"
	"project_pos_app/model"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NotifControllerSuite struct {
	suite.Suite
	controller notifcontroller.NotifController
	mockDB     *mocktesting.MockDB
	ctx        *gin.Context
	writer     *httptest.ResponseRecorder
}

func TestNotifControllerSuite(t *testing.T) {
	suite.Run(t, new(NotifControllerSuite))
}

func (suite *NotifControllerSuite) SetupTest() {
	mockLogger := zap.NewNop()
	mockDB, mockService := helper.InitService()

	suite.controller = notifcontroller.NewNotifController(mockService, mockLogger)
	suite.mockDB = mockDB
	gin.SetMode(gin.TestMode)
	suite.writer = httptest.NewRecorder()
	suite.ctx, _ = gin.CreateTestContext(suite.writer)
}

func (suite *NotifControllerSuite) TestCreateNotifications() {
	now := time.Now()
	dynamicDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	payload := model.Notification{
		Title:     "New Notification",
		Message:   "This is a test notification",
		Status:    "new",
		CreatedAt: dynamicDate,
		UpdatedAt: dynamicDate,
	}

	suite.Run("Success Created", func() {
		// Mock service response
		suite.mockDB.On("Create", payload).Once().Return(nil)

		// Prepare request
		body, _ := json.Marshal(payload)
		suite.ctx.Request = httptest.NewRequest(http.MethodPost, "/api/notifications", bytes.NewBuffer(body))
		suite.ctx.Request.Header.Set("Content-Type", "application/json")

		// Call controller
		suite.controller.CreateNotifications(suite.ctx)

		expected := model.SuccessResponse{
			Status:  http.StatusCreated,
			Message: "Create notification successfully",
			Data:    nil,
		}

		var apiResponse model.SuccessResponse
		err := json.Unmarshal(suite.writer.Body.Bytes(), &apiResponse)

		suite.NoError(err)
		suite.Equal(expected, apiResponse)
		suite.Equal(http.StatusCreated, suite.writer.Code)
		suite.Contains(suite.writer.Body.String(), "Create notification successfully")
		suite.mockDB.AssertCalled(suite.T(), "Create", payload)
	})

	suite.Run("Invalid Payload", func() {
		// Reset setup data
		suite.SetupTest()

		// Mock service response
		suite.mockDB.On("Create", payload).Return(errors.New("invalid payload"))

		// Invalid JSON payload
		suite.ctx.Request = httptest.NewRequest(http.MethodPost, "/api/notifications", bytes.NewBuffer([]byte("{invalid_json")))
		suite.ctx.Request.Header.Set("Content-Type", "application/json")

		// Call controller
		suite.controller.CreateNotifications(suite.ctx)

		suite.Equal(http.StatusBadRequest, suite.writer.Code)
		suite.Contains(suite.writer.Body.String(), "Invalid Payload")
	})

	suite.Run("Database Error", func() {
		// Reset setup data
		suite.SetupTest()

		now := time.Now()
		dynamicDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		payload := model.Notification{
			Title:     "New Notification",
			Message:   "This is a test notification",
			Status:    "new",
			CreatedAt: dynamicDate,
			UpdatedAt: dynamicDate,
		}

		// Mock service response with error
		suite.mockDB.On("Create", payload).Return(errors.New("database error"))

		// Prepare request
		body, _ := json.Marshal(payload)
		suite.ctx.Request = httptest.NewRequest(http.MethodPost, "/api/notifications", bytes.NewBuffer(body))
		suite.ctx.Request.Header.Set("Content-Type", "application/json")

		// Call controller
		suite.controller.CreateNotifications(suite.ctx)

		suite.Equal(http.StatusInternalServerError, suite.writer.Code)
		suite.Contains(suite.writer.Body.String(), "Failed to fetch notifications")
		suite.mockDB.AssertCalled(suite.T(), "Create", payload)
	})
}

func (suite *NotifControllerSuite) TestGetAllNotifications() {
	suite.Run("Success Fetch", func() {
		notifications := []model.Notification{
			{ID: 1, Title: "Notif 1", Message: "Message 1", Status: "new"},
			{ID: 2, Title: "Notif 2", Message: "Message 2", Status: "read"},
		}

		suite.mockDB.On("GetAll").Return(notifications, nil)

		suite.ctx.Request = httptest.NewRequest(http.MethodGet, "/api/notifications", nil)

		suite.controller.GetAllNotifications(suite.ctx)

		mapNotifications, err := helper.StructToMapSlice(&notifications)
		suite.NoError(err)

		expected := model.SuccessResponse{
			Status:  http.StatusOK,
			Message: "Get all notification successfully",
			Data:    mapNotifications,
		}

		var apiResponse model.SuccessResponse
		err = json.Unmarshal(suite.writer.Body.Bytes(), &apiResponse)
		suite.NoError(err)
		apiResponse.Data, err = helper.ConvertFieldInData(apiResponse.Data, "id", "int")
		suite.NoError(err)

		suite.Equal(expected, apiResponse)
		suite.Equal(http.StatusOK, suite.writer.Code)
		suite.mockDB.AssertCalled(suite.T(), "GetAll")
	})

	suite.Run("Database Error", func() {
		suite.SetupTest()
		suite.mockDB.On("GetAll").Return(nil, errors.New("database error"))

		suite.ctx.Request = httptest.NewRequest(http.MethodGet, "/api/notifications", nil)

		suite.controller.GetAllNotifications(suite.ctx)

		suite.Equal(http.StatusInternalServerError, suite.writer.Code)
		suite.Contains(suite.writer.Body.String(), "Failed to fetch notifications")
		suite.mockDB.AssertCalled(suite.T(), "GetAll")
	})
}

func (suite *NotifControllerSuite) TestGetNotificationByID() {
	suite.Run("Success Fetch", func() {
		notification := model.Notification{ID: 1, Title: "Notif 1", Message: "Message 1", Status: "new"}

		suite.mockDB.On("FindByID", 1).Return(&notification, nil)

		suite.ctx.Request = httptest.NewRequest(http.MethodGet, "/api/notifications/1", nil)
		suite.ctx.Params = append(suite.ctx.Params, gin.Param{Key: "id", Value: "1"})

		suite.controller.GetNotificationByID(suite.ctx)
		mapNotifications, err := helper.StructToMap(&notification)
		suite.NoError(err)
		expected := model.SuccessResponse{
			Status:  http.StatusOK,
			Message: "Notification retrieved successfully",
			Data:    mapNotifications,
		}

		var apiResponse model.SuccessResponse
		err = json.Unmarshal(suite.writer.Body.Bytes(), &apiResponse)
		suite.NoError(err)

		apiResponse.Data, err = helper.ConvertFieldInMap(apiResponse.Data, "id", "int")
		suite.NoError(err)

		suite.Equal(expected, apiResponse)
		suite.Equal(http.StatusOK, suite.writer.Code)
		suite.mockDB.AssertCalled(suite.T(), "FindByID", 1)
	})

	suite.Run("Invalid ID Format", func() {
		suite.SetupTest()
		suite.ctx.Request = httptest.NewRequest(http.MethodGet, "/api/notifications/invalid", nil)
		suite.ctx.Params = append(suite.ctx.Params, gin.Param{Key: "id", Value: "invalid"})

		suite.controller.GetNotificationByID(suite.ctx)

		suite.Equal(http.StatusBadRequest, suite.writer.Code)
		suite.Contains(suite.writer.Body.String(), "Invalid ID format")
	})

	suite.Run("Database Error", func() {
		suite.SetupTest()
		suite.mockDB.On("FindByID", 1).Return(nil, errors.New("database error"))

		suite.ctx.Request = httptest.NewRequest(http.MethodGet, "/api/notifications/1", nil)
		suite.ctx.Params = append(suite.ctx.Params, gin.Param{Key: "id", Value: "1"})

		suite.controller.GetNotificationByID(suite.ctx)

		suite.Equal(http.StatusInternalServerError, suite.writer.Code)
		suite.Contains(suite.writer.Body.String(), "Failed to fetch notification")
		suite.mockDB.AssertCalled(suite.T(), "FindByID", 1)
	})
}

func (suite *NotifControllerSuite) TestUpdateNotification() {
	now := time.Now()
	dynamicDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	suite.Run("Success Update", func() {
		notif := &model.Notification{
			ID:      1,
			Title:   "Testing",
			Message: "Test notification",
			Status:  "new",
		}

		suite.mockDB.On("FindByID", notif.ID).Return(notif, nil)
		suite.mockDB.On("Update", notif, 1).Return(notif, nil)

		notif.UpdatedAt = dynamicDate
		suite.ctx.Request = httptest.NewRequest(http.MethodPut, "/api/notifications/1", nil)
		suite.ctx.Params = append(suite.ctx.Params, gin.Param{Key: "id", Value: "1"})

		suite.controller.UpdateNotification(suite.ctx)

		suite.Equal(http.StatusOK, suite.writer.Code)
		suite.Contains(suite.writer.Body.String(), "Notification updated successfully")
		suite.mockDB.AssertCalled(suite.T(), "Update", notif, 1)
	})

	suite.Run("Invalid ID Format", func() {
		suite.SetupTest()
		suite.ctx.Request = httptest.NewRequest(http.MethodPut, "/api/notifications/invalid", nil)
		suite.ctx.Params = append(suite.ctx.Params, gin.Param{Key: "id", Value: "invalid"})

		suite.controller.UpdateNotification(suite.ctx)

		suite.Equal(http.StatusBadRequest, suite.writer.Code)
		suite.Contains(suite.writer.Body.String(), "Invalid ID format")
	})

	suite.Run("Database Error", func() {
		suite.SetupTest()
		notif := &model.Notification{
			ID:      9999,
			Title:   "Testing",
			Message: "Test notification",
			Status:  "new",
		}

		suite.mockDB.On("FindByID", notif.ID).Return(nil, errors.New("record not found"))
		suite.mockDB.On("Update", notif, notif.ID).Return(nil, errors.New("database error"))

		suite.ctx.Request = httptest.NewRequest(http.MethodPut, "/api/notifications/9999", nil)
		suite.ctx.Params = append(suite.ctx.Params, gin.Param{Key: "id", Value: "9999"})

		suite.controller.UpdateNotification(suite.ctx)

		suite.Equal(http.StatusInternalServerError, suite.writer.Code)
		suite.Contains(suite.writer.Body.String(), "Failed to update notification")
	})
}

func (suite *NotifControllerSuite) TestDeleteNotification() {
	suite.Run("Success Delete", func() {
		notif := &model.Notification{
			ID:      1,
			Title:   "Testing",
			Message: "Test notification",
			Status:  "new",
		}

		suite.mockDB.On("FindByID", notif.ID).Return(notif, nil)
		suite.mockDB.On("Delete", notif.ID).Return(&gorm.DB{Error: nil})

		suite.ctx.Request = httptest.NewRequest(http.MethodDelete, "/api/notifications/1", nil)
		suite.ctx.Params = append(suite.ctx.Params, gin.Param{Key: "id", Value: "1"})

		suite.controller.DeleteNotification(suite.ctx)

		suite.Equal(http.StatusOK, suite.writer.Code)
		suite.Contains(suite.writer.Body.String(), "Notification deleted successfully")
		suite.mockDB.AssertCalled(suite.T(), "Delete", 1)
	})

	suite.Run("Invalid ID Format", func() {
		suite.SetupTest()
		suite.ctx.Request = httptest.NewRequest(http.MethodDelete, "/api/notifications/invalid", nil)
		suite.ctx.Params = append(suite.ctx.Params, gin.Param{Key: "id", Value: "invalid"})

		suite.controller.DeleteNotification(suite.ctx)

		suite.Equal(http.StatusBadRequest, suite.writer.Code)
		suite.Contains(suite.writer.Body.String(), "Invalid ID format")
	})

	suite.Run("Database Error", func() {
		suite.SetupTest()
		notif := &model.Notification{
			ID: 9999,
		}

		suite.mockDB.On("FindByID", notif.ID).Return(nil, errors.New("record not found"))
		suite.mockDB.On("Delete", notif.ID).Return(errors.New("database error"))

		suite.ctx.Request = httptest.NewRequest(http.MethodDelete, "/api/notifications/9999", nil)
		suite.ctx.Params = append(suite.ctx.Params, gin.Param{Key: "id", Value: "9999"})

		suite.controller.DeleteNotification(suite.ctx)

		suite.Equal(http.StatusInternalServerError, suite.writer.Code)
		suite.Contains(suite.writer.Body.String(), "Failed to delete notification")
	})
}

func (suite *NotifControllerSuite) TestMarkAllNotificationsAsRead() {
	suite.Run("Success Mark as Read", func() {
		now := time.Now()
		dynamicDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		notifications := []model.Notification{
			{
				ID:        1,
				Title:     "Testing",
				Message:   "Test notification",
				Status:    "new",
				CreatedAt: dynamicDate,
				UpdatedAt: dynamicDate,
			},
			{
				ID:        2,
				Title:     "Testing2",
				Message:   "Test notification2",
				Status:    "new",
				CreatedAt: dynamicDate,
				UpdatedAt: dynamicDate,
			},
		}

		suite.mockDB.On("MarkAllAsRead").Once().Return(notifications, nil)

		suite.ctx.Request = httptest.NewRequest(http.MethodPut, "/api/notifications/mark-as-read", nil)

		suite.controller.MarkAllNotificationsAsRead(suite.ctx)

		suite.Equal(http.StatusOK, suite.writer.Code)
		suite.Contains(suite.writer.Body.String(), "All notifications marked as read")
		suite.mockDB.AssertCalled(suite.T(), "MarkAllAsRead")
	})

	suite.Run("Database Error", func() {
		suite.SetupTest()
		suite.mockDB.On("MarkAllAsRead").Return(nil, errors.New("database error"))

		suite.ctx.Request = httptest.NewRequest(http.MethodPut, "/api/notifications/mark-as-read", nil)

		suite.controller.MarkAllNotificationsAsRead(suite.ctx)

		suite.Equal(http.StatusInternalServerError, suite.writer.Code)
		suite.Contains(suite.writer.Body.String(), "Failed to mark all notifications as read")
		suite.mockDB.AssertCalled(suite.T(), "MarkAllAsRead")
	})
}
