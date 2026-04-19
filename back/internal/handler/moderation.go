package handler

import (
	"net/http"
	"time"

	"crowdfunding/back/internal/model"

	"github.com/gin-gonic/gin"
)

// GET /api/admin/moderation
// Список проектов, ожидающих проверки человеком (pending_human).
func (h *Handler) AdminModerationList(c *gin.Context) {
	page, limit := paginate(c)
	offset := (page - 1) * limit

	var total int64
	h.db.Model(&model.Project{}).Where("status = ?", model.StatusPendingHuman).Count(&total)

	var projects []model.Project
	h.db.Where("status = ?", model.StatusPendingHuman).
		Preload("User").
		Preload("Categories").
		Preload("Moderation").
		Order("created_at ASC"). // старые сначала
		Offset(offset).Limit(limit).
		Find(&projects)

	c.JSON(http.StatusOK, gin.H{
		"data":  projects,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GET /api/admin/moderation/:project_id
// Детали конкретного проекта + полная информация о модерации.
func (h *Handler) AdminModerationGet(c *gin.Context) {
	var project model.Project
	err := h.db.
		Preload("User").
		Preload("Categories").
		Preload("Moderation.Moderator").
		First(&project, c.Param("project_id")).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	h.db.Model(&model.Like{}).Where("project_id = ?", project.ID).Count(&project.LikesCount)
	c.JSON(http.StatusOK, project)
}

// PATCH /api/admin/moderation/:project_id
// Решение модератора: approve или reject.
func (h *Handler) AdminModerationDecide(c *gin.Context) {
	var project model.Project
	if err := h.db.Preload("Moderation").First(&project, c.Param("project_id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	if project.Status != model.StatusPendingHuman {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project is not awaiting human moderation"})
		return
	}

	var input struct {
		Decision      string `json:"decision"       binding:"required,oneof=approve reject"`
		ModeratorNote string `json:"moderator_note"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	moderatorID := h.currentUserID(c)
	now := time.Now()

	var newProjectStatus model.ProjectStatus
	var newHumanStatus model.HumanStatus

	if input.Decision == "approve" {
		newProjectStatus = model.StatusActive
		newHumanStatus = model.HumanStatusApproved
	} else {
		newProjectStatus = model.StatusRejected
		newHumanStatus = model.HumanStatusRejected
	}

	// Update moderation record
	if project.Moderation != nil {
		h.db.Model(project.Moderation).Updates(map[string]any{
			"human_status":       newHumanStatus,
			"moderator_id":       moderatorID,
			"moderator_note":     input.ModeratorNote,
			"human_moderated_at": now,
		})
	}

	// Update project status
	h.db.Model(&project).Update("status", newProjectStatus)

	h.db.Preload("User").Preload("Categories").Preload("Moderation.Moderator").First(&project, project.ID)
	c.JSON(http.StatusOK, project)
}

// POST /api/admin/moderation/:project_id/invite
// Отправляет стартаперу приглашение на модерацию с произвольным сообщением от админа.
func (h *Handler) AdminModerationInvite(c *gin.Context) {
	var project model.Project
	if err := h.db.First(&project, c.Param("project_id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	var input struct {
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notif := model.Notification{
		UserID:    project.UserID,
		ProjectID: &project.ID,
		Type:      model.NotifInvite,
		Title:     "Приглашение на модерацию",
		Body:      input.Message,
	}
	h.db.Create(&notif)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// POST /api/admin/moderation/:project_id/recheck
// Повторно прогнать AI-проверку (например, после редактирования проекта).
func (h *Handler) AdminModerationRecheck(c *gin.Context) {
	var project model.Project
	if err := h.db.Preload("Moderation").First(&project, c.Param("project_id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	// Reset moderation state
	if project.Moderation != nil {
		h.db.Delete(project.Moderation)
	}
	h.db.Model(&project).Update("status", model.StatusPendingAI)
	project.Status = model.StatusPendingAI

	// Re-run AI
	h.runAICheck(&project)

	h.db.Preload("User").Preload("Categories").Preload("Moderation").First(&project, project.ID)
	c.JSON(http.StatusOK, project)
}
