package controllers

import (
	"campus-lost-and-found/internal/middleware"
	"campus-lost-and-found/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	Service *services.NotificationService
}

func NewNotificationController(service *services.NotificationService) *NotificationController {
	return &NotificationController{Service: service}
}

// GetNotifications godoc
// @Summary Get user notifications
// @Description Get notifications for the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []models.Notification
// @Router /notifications [get]
func (ctrl *NotificationController) GetNotifications(c *gin.Context) {
	userID := middleware.GetUserID(c)
	notifs, err := ctrl.Service.GetUserNotifications(userID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notifs)
}

// MarkAsRead godoc
// @Summary Mark notification as read
// @Description Mark a notification as read
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]string
// @Router /notifications/{id}/read [put]
func (ctrl *NotificationController) MarkAsRead(c *gin.Context) {
	id := c.Param("id")
	err := ctrl.Service.MarkAsRead(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}
