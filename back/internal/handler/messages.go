package handler

import (
	"net/http"

	"crowdfunding/back/internal/model"

	"github.com/gin-gonic/gin"
)

// GET /api/messages
func (h *Handler) ListMessages(c *gin.Context) {
	page, limit := paginate(c)
	offset := (page - 1) * limit
	uid := h.currentUserID(c)

	var total int64
	h.db.Model(&model.Message{}).
		Where("sender_id = ? OR recipient_id = ?", uid, uid).
		Count(&total)

	var messages []model.Message
	h.db.Where("sender_id = ? OR recipient_id = ?", uid, uid).
		Preload("Sender").
		Preload("Recipient").
		Preload("Project").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&messages)

	c.JSON(http.StatusOK, gin.H{
		"data":  messages,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// POST /api/messages
func (h *Handler) SendMessage(c *gin.Context) {
	var input struct {
		ProjectID uint   `json:"project_id" binding:"required"`
		Title     string `json:"title"      binding:"required"`
		Content   string `json:"content"    binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var project model.Project
	if err := h.db.First(&project, input.ProjectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	senderID := h.currentUserID(c)
	if project.UserID == senderID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot send message to yourself"})
		return
	}

	msg := model.Message{
		SenderID:    senderID,
		RecipientID: project.UserID,
		ProjectID:   project.ID,
		Title:       input.Title,
		Content:     input.Content,
	}
	if err := h.db.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send message"})
		return
	}

	h.db.Preload("Sender").Preload("Recipient").Preload("Project").First(&msg, msg.ID)
	c.JSON(http.StatusCreated, msg)
}

// PATCH /api/messages/:id/read
func (h *Handler) MarkMessageRead(c *gin.Context) {
	var msg model.Message
	if err := h.db.First(&msg, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
		return
	}

	if msg.RecipientID != h.currentUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	h.db.Model(&msg).Update("is_read", true)
	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}
