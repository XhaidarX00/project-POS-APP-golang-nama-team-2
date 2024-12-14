package notifcontroller

import (
	"net/http"
	"project_pos_app/helper"
	"project_pos_app/model"
	"project_pos_app/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type NotifController struct {
	Service *service.AllService
	Log     *zap.Logger
}

func NewNotifController(service *service.AllService, log *zap.Logger) NotifController {
	return NotifController{
		Service: service,
		Log:     log,
	}
}

// CreateNotifications godoc
// @Summary Create a new notification
// @Description Create a new notification
// @Tags Notification
// @Accept json
// @Produce json
// @Param notification body model.Notification true "Notification payload"
// @Success 201 {object} model.SuccessResponse{data=model.Notification} "Notification created successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid payload"
// @Failure 500 {object} model.ErrorResponse "Failed to create notification"
// @Router /api/notifications [post]
func (c *NotifController) CreateNotifications(ctx *gin.Context) {
	var data model.Notification
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		c.Log.Error("Invalid payload", zap.Error(err))
		helper.Responses(ctx, http.StatusInternalServerError, "Invalid Payload: "+err.Error(), nil)
		ctx.Abort()
		return
	}

	err = c.Service.Notif.CreateNotification(data)
	if err != nil {
		c.Log.Error("Failed to get all notifications", zap.Error(err))
		helper.Responses(ctx, http.StatusInternalServerError, "Failed to fetch notifications", nil)
		ctx.Abort()
		return
	}

	helper.Responses(ctx, http.StatusCreated, "Create notification successfully", nil)
}

// GetAllNotifications godoc
// @Summary Get all notifications
// @Description Retrieve all notifications, optionally filtered by status
// @Tags Notification
// @Accept json
// @Produce json
// @Param status query string false "Notification status (e.g., unread)"
// @Success 200 {object} model.SuccessResponse{data=[]model.Notification} "List of notifications retrieved successfully"
// @Failure 500 {object} model.ErrorResponse "Failed to fetch notifications"
// @Router /api/notifications [get]
func (c *NotifController) GetAllNotifications(ctx *gin.Context) {
	status := ctx.Query("status")
	notifications, err := c.Service.Notif.GetAllNotifications(status)
	if err != nil {
		c.Log.Error("Failed to get all notifications", zap.Error(err))
		helper.Responses(ctx, http.StatusInternalServerError, "Failed to fetch notifications", nil)
		ctx.Abort()
		return
	}
	helper.Responses(ctx, http.StatusOK, "Get all notification successfully", notifications)
}

// GetNotificationByID godoc
// @Summary Get a notification by ID
// @Description Retrieve a notification by its ID
// @Tags Notification
// @Accept json
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} model.SuccessResponse{data=model.Notification} "Notification retrieved successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid ID format"
// @Failure 500 {object} model.ErrorResponse "Failed to fetch notification"
// @Router /api/notifications/{id} [get]
func (c *NotifController) GetNotificationByID(ctx *gin.Context) {
	id := ctx.Param("id")
	notifID, err := strconv.Atoi(id)
	if err != nil {
		c.Log.Error("Invalid notification ID", zap.Error(err))
		helper.Responses(ctx, http.StatusBadRequest, "Invalid ID format", nil)
		ctx.Abort()
		return
	}

	notification, err := c.Service.Notif.GetNotificationByID(notifID)
	if err != nil {
		c.Log.Error("Failed to get notification by ID", zap.Error(err))
		helper.Responses(ctx, http.StatusInternalServerError, "Failed to fetch notification", nil)
		ctx.Abort()
		return
	}
	helper.Responses(ctx, http.StatusOK, "Notification retrieved successfully", notification)
}

// UpdateNotification godoc
// @Summary Update a notification by ID
// @Description Update the status or details of a notification
// @Tags Notification
// @Accept json
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} model.SuccessResponse "Notification updated successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid ID format"
// @Failure 500 {object} model.ErrorResponse "Failed to update notification"
// @Router /api/notifications/{id} [put]
func (c *NotifController) UpdateNotification(ctx *gin.Context) {
	id := ctx.Param("id")
	notificationID, err := strconv.Atoi(id)
	if err != nil {
		c.Log.Error("Invalid notification ID", zap.Error(err))
		helper.Responses(ctx, http.StatusBadRequest, "Invalid ID format", nil)
		ctx.Abort()
		return
	}

	if err := c.Service.Notif.UpdateNotification(notificationID); err != nil {
		c.Log.Error("Failed to update notification", zap.Error(err))
		helper.Responses(ctx, http.StatusInternalServerError, "Failed to update notification", nil)
		ctx.Abort()
		return
	}
	helper.Responses(ctx, http.StatusOK, "Notification updated successfully", nil)
}

// DeleteNotification godoc
// @Summary Delete a notification by ID
// @Description Delete a notification using its ID
// @Tags Notification
// @Accept json
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} model.SuccessResponse "Notification deleted successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid ID format"
// @Failure 500 {object} model.ErrorResponse "Failed to delete notification"
// @Router /api/notifications/{id} [delete]
func (c *NotifController) DeleteNotification(ctx *gin.Context) {
	id := ctx.Param("id")
	notificationID, err := strconv.Atoi(id)
	if err != nil {
		c.Log.Error("Invalid notification ID", zap.Error(err))
		helper.Responses(ctx, http.StatusBadRequest, "Invalid ID format", nil)
		ctx.Abort()
		return
	}

	if err := c.Service.Notif.DeleteNotification(notificationID); err != nil {
		c.Log.Error("Failed to delete notification", zap.Error(err))
		helper.Responses(ctx, http.StatusInternalServerError, "Failed to delete notification", nil)
		ctx.Abort()
		return
	}
	helper.Responses(ctx, http.StatusOK, "Notification deleted successfully", nil)
}

// MarkAllNotificationsAsRead godoc
// @Summary Mark all notifications as read
// @Description Mark all notifications as read
// @Tags Notification
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessResponse "All notifications marked as read successfully"
// @Failure 500 {object} model.ErrorResponse "Failed to mark notifications as read"
// @Router /api/notifications/mark-as-read [put]
func (c *NotifController) MarkAllNotificationsAsRead(ctx *gin.Context) {
	if err := c.Service.Notif.MarkAllNotificationsAsRead(); err != nil {
		c.Log.Error("Failed to mark all notifications as read", zap.Error(err))
		helper.Responses(ctx, http.StatusInternalServerError, "Failed to mark all notifications as read", nil)
		ctx.Abort()
		return
	}
	helper.Responses(ctx, http.StatusOK, "All notifications marked as read", nil)
}
