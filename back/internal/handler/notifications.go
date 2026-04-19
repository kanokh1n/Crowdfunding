package handler

import (
	"net/http"

	"crowdfunding/back/internal/model"

	"github.com/gin-gonic/gin"
)

// GET /api/notifications
func (h *Handler) ListNotifications(c *gin.Context) {
	uid := h.currentUserID(c)
	var notifs []model.Notification
	h.db.Where("user_id = ?", uid).Order("created_at DESC").Find(&notifs)
	c.JSON(http.StatusOK, notifs)
}

// PATCH /api/notifications/:id/read
func (h *Handler) MarkNotificationRead(c *gin.Context) {
	uid := h.currentUserID(c)
	var notif model.Notification
	if err := h.db.Where("id = ? AND user_id = ?", c.Param("id"), uid).First(&notif).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	h.db.Model(&notif).Update("is_read", true)
	c.JSON(http.StatusOK, notif)
}
