package notifcontroller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	notifcontroller "project_pos_app/controller/notif_controller"
	"project_pos_app/helper"
	"project_pos_app/model"
	"project_pos_app/repository/notification"
	"project_pos_app/service"
	notifservice "project_pos_app/service/notif_service"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func notifBase() (notifcontroller.NotifController, *notification.MockDB) {
	log := *zap.NewNop()

	mockRepo, _ := helper.InitService()
	serviceNotif := notifservice.NewMockNotifService(mockRepo, &log)
	service := service.AllService{
		Notif: serviceNotif,
	}

	return notifcontroller.NewNotifController(&service, &log), mockRepo
}

func TestCreateNotification(t *testing.T) {
	handler, mockService := notifBase()

	t.Run("Successfully create a notification", func(t *testing.T) {
		r := gin.Default()
		r.POST("/api/notification", handler.CreateNotifications)

		newNotification := model.Notification{
			Title:   "Testing",
			Message: "Test notification",
		}
		mockService.On("Create", newNotification).Once().Return(nil)
		body, _ := json.Marshal(newNotification)
		req := httptest.NewRequest(http.MethodPost, "/api/notification", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertCalled(t, "Create", newNotification)

		var actualResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, "Create notification successfully", actualResponse["Message"])
		assert.Nil(t, actualResponse["data"])
	})

	t.Run("Failed to create a notification - Invalid Payload", func(t *testing.T) {
		r := gin.Default()
		r.POST("/api/notification", handler.CreateNotifications)

		invalidBody := `{"Title": "New Notification"` // Invalid JSON
		req := httptest.NewRequest(http.MethodPost, "/api/notification", bytes.NewBufferString(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var actualResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)

		assert.Contains(t, actualResponse["Message"], "Invalid Payload")
		assert.Nil(t, actualResponse["data"])
	})
}

func TestGetAllNotifications(t *testing.T) {
	handler, mockRepo := notifBase()

	t.Run("Successfully retrieve all notifications", func(t *testing.T) {
		// Setup router
		r := gin.Default()
		r.GET("/api/notification", handler.GetAllNotifications)

		// Mock data
		now := time.Now()
		notifications := []model.Notification{
			{
				ID:        1,
				Title:     "Testing",
				Message:   "Test notification",
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        2,
				Title:     "Testing2",
				Message:   "Test notification2",
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		// Mock repository behavior
		mockRepo.On("GetAll").Once().Return(notifications, nil)

		// Create request and response recorder
		req := httptest.NewRequest(http.MethodGet, "/api/notification", nil)
		w := httptest.NewRecorder()

		// Serve request
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertCalled(t, "GetAll")

		// Parse response body
		var actualResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)

		// Validate "message" field
		message, ok := actualResponse["Message"].(string)
		assert.True(t, ok, "expected 'Message' field in the response to be a string")
		assert.Equal(t, "Get all notification successfully", message)

		// Validate "Status" field
		status, ok := actualResponse["Status"].(float64)
		assert.True(t, ok, "expected 'Status' field in the response to be a number")
		assert.Equal(t, float64(200), status)

		// Validate "data" field
		data, ok := actualResponse["data"].([]interface{})
		assert.True(t, ok, "expected 'data' field in the response to be a slice")
		assert.Len(t, data, 2)

		// Validate individual notification fields
		firstNotification, ok := data[0].(map[string]interface{})
		assert.True(t, ok, "expected each item in 'data' to be a map")
		assert.Equal(t, float64(1), firstNotification["id"])
		assert.Equal(t, "Testing", firstNotification["title"])
		assert.Equal(t, "Test notification", firstNotification["message"])
		assert.Equal(t, "", firstNotification["status"])

		secondNotification, ok := data[1].(map[string]interface{})
		assert.True(t, ok, "expected each item in 'data' to be a map")
		assert.Equal(t, float64(2), secondNotification["id"])
		assert.Equal(t, "Testing2", secondNotification["title"])
		assert.Equal(t, "Test notification2", secondNotification["message"])
		assert.Equal(t, "", secondNotification["status"])
	})

	t.Run("Failed to retrieve notifications", func(t *testing.T) {
		// Setup router
		r := gin.Default()
		r.GET("/api/notification", handler.GetAllNotifications)

		// Mock repository behavior
		mockRepo.On("GetAll").Once().Return(nil, fmt.Errorf("database error"))

		// Create request and response recorder
		req := httptest.NewRequest(http.MethodGet, "/api/notification", nil)
		w := httptest.NewRecorder()

		// Serve request
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockRepo.AssertCalled(t, "GetAll")

		// Parse response body
		var actualResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)

		// Validate "message" field
		message, ok := actualResponse["Message"].(string)
		assert.True(t, ok, "expected 'Message' field in the response to be a string")
		assert.Equal(t, "Failed to fetch notifications", message)

		// Validate "Status" field
		status, ok := actualResponse["Status"].(float64)
		assert.True(t, ok, "expected 'Status' field in the response to be a number")
		assert.Equal(t, float64(500), status)
	})
}

func TestGetNotificationByID(t *testing.T) {
	handler, mockService := notifBase()

	t.Run("Successfully retrieve notification by ID", func(t *testing.T) {
		r := gin.Default()
		r.GET("/api/notification/:id", handler.GetNotificationByID)

		now := time.Now()
		expectedNotification := model.Notification{
			ID:        1,
			Title:     "Testing",
			Message:   "Test notification",
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockService.On("FindByID", 1).Once().Return(expectedNotification, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/notification/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertCalled(t, "FindByID", 1)

		var actualResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, "Notification retrieved successfully", actualResponse["Message"])

		data, ok := actualResponse["data"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(1), data["id"])
		assert.Equal(t, "Testing", data["title"])
		assert.Equal(t, "Test notification", data["message"])
	})

	t.Run("Failed to retrieve notification - Invalid ID", func(t *testing.T) {
		r := gin.Default()
		r.GET("/api/notification/:id", handler.GetNotificationByID)

		req := httptest.NewRequest(http.MethodGet, "/api/notification/invalid", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var actualResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, "Invalid ID format", actualResponse["Message"])
	})

	t.Run("Failed to retrieve notification - Database Error", func(t *testing.T) {
		r := gin.Default()
		r.GET("/api/notification/:id", handler.GetNotificationByID)

		mockService.On("FindByID", 999).Once().Return(model.Notification{}, fmt.Errorf("database error"))

		req := httptest.NewRequest(http.MethodGet, "/api/notification/999", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertCalled(t, "FindByID", 999)

		var actualResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, "Failed to fetch notification", actualResponse["Message"])
	})
}

func TestUpdateNotification(t *testing.T) {
	handler, mockService := notifBase()

	t.Run("Successfully update notification", func(t *testing.T) {
		r := gin.Default()
		r.PUT("/api/notification/:id", handler.UpdateNotification)

		mockService.On("Update", mock.Anything, 1).Once().Return(nil)

		req := httptest.NewRequest(http.MethodPut, "/api/notification/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertCalled(t, "Update", mock.Anything, 1)

		var actualResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, "Notification updated successfully", actualResponse["Message"])
	})

	t.Run("Failed to update notification - Invalid ID", func(t *testing.T) {
		r := gin.Default()
		r.PUT("/api/notification/:id", handler.UpdateNotification)

		req := httptest.NewRequest(http.MethodPut, "/api/notification/invalid", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var actualResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, "Invalid ID format", actualResponse["Message"])
	})

	t.Run("Failed to update notification - Database Error", func(t *testing.T) {
		r := gin.Default()
		r.PUT("/api/notification/:id", handler.UpdateNotification)

		mockService.On("Update", mock.Anything, 999).Once().Return(fmt.Errorf("database error"))

		req := httptest.NewRequest(http.MethodPut, "/api/notification/999", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		// mockService.AssertCalled(t, "Update", mock.Anything, 999)

		// var actualResponse map[string]interface{}
		// err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		// assert.NoError(t, err)

		// assert.Equal(t, "Failed to update notification", actualResponse["Message"])
	})
}

func TestDeleteNotification(t *testing.T) {
	handler, mockService := notifBase()

	t.Run("Successfully delete notification", func(t *testing.T) {
		r := gin.Default()
		r.DELETE("/api/notification/:id", handler.DeleteNotification)
		now := time.Now()
		notif := model.Notification{
			ID:        1,
			Title:     "Testing",
			Message:   "Test notification",
			Status:    "new",
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockService.On("FindByID", 1).Return(notif, nil)
		mockService.On("Delete", 1).Once().Return(&gorm.DB{Error: nil})

		req := httptest.NewRequest(http.MethodDelete, "/api/notification/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertCalled(t, "Delete", 1)

		var actualResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, "Notification deleted successfully", actualResponse["Message"])
	})

	t.Run("Failed to delete notification - Invalid ID", func(t *testing.T) {
		r := gin.Default()
		r.DELETE("/api/notification/:id", handler.DeleteNotification)

		req := httptest.NewRequest(http.MethodDelete, "/api/notification/invalid", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var actualResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, "Invalid ID format", actualResponse["Message"])
	})

	t.Run("Failed to delete notification - Database Error", func(t *testing.T) {
		r := gin.Default()
		r.DELETE("/api/notification/:id", handler.DeleteNotification)
		mockService.On("Delete", 999).Once().Return(&gorm.DB{Error: nil})

		req := httptest.NewRequest(http.MethodDelete, "/api/notification/999", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertNotCalled(t, "Delete", 999)

		// var actualResponse map[string]interface{}
		// err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		// assert.NoError(t, err)

		// assert.Equal(t, "Failed to delete notification", actualResponse["Message"])
	})
}

func TestMarkAllNotificationsAsRead(t *testing.T) {
	handler, mockService := notifBase()

	t.Run("Successfully mark all notifications as read", func(t *testing.T) {
		r := gin.Default()
		r.PUT("/api/notification/mark-read", handler.MarkAllNotificationsAsRead)
		notifications := []model.Notification{}
		mockService.On("MarkAllAsRead").Once().Return(notifications, nil)

		req := httptest.NewRequest(http.MethodPut, "/api/notification/mark-read", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertCalled(t, "MarkAllAsRead")

		var actualResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, "All notifications marked as read", actualResponse["Message"])
	})

	t.Run("Failed to mark all notifications as read - Database Error", func(t *testing.T) {
		r := gin.Default()
		r.PUT("/api/notification/mark-read", handler.MarkAllNotificationsAsRead)

		mockService.On("MarkAllAsRead").Once().Return(fmt.Errorf("database error"))

		req := httptest.NewRequest(http.MethodPut, "/api/notification/mark-read", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertCalled(t, "MarkAllAsRead")

		// var actualResponse map[string]interface{}
		// err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		// assert.NoError(t, err)

		// assert.Equal(t, "Failed to mark all notifications as read", actualResponse["Message"])
	})
}
