package handler

import (
	"errors"
	"net/http"

	"crowdfunding/back/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// POST /api/projects/:id/like
func (h *Handler) LikeProject(c *gin.Context) {
	var project model.Project
	if err := h.db.First(&project, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	like := model.Like{
		UserID:    h.currentUserID(c),
		ProjectID: project.ID,
	}

	if err := h.db.Create(&like).Error; err != nil {
		// unique constraint violation means already liked
		c.JSON(http.StatusConflict, gin.H{"error": "already liked"})
		return
	}

	var count int64
	h.db.Model(&model.Like{}).Where("project_id = ?", project.ID).Count(&count)
	c.JSON(http.StatusOK, gin.H{"likes_count": count})
}

// DELETE /api/projects/:id/like
func (h *Handler) UnlikeProject(c *gin.Context) {
	uid := h.currentUserID(c)

	result := h.db.Where("user_id = ? AND project_id = ?", uid, c.Param("id")).Delete(&model.Like{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unlike"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "like not found"})
		return
	}

	var count int64
	h.db.Model(&model.Like{}).Where("project_id = ?", c.Param("id")).Count(&count)
	c.JSON(http.StatusOK, gin.H{"likes_count": count})
}

// helper — keep gorm import
var _ = errors.New
var _ *gorm.DB
