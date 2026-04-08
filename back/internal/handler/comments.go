package handler

import (
	"net/http"

	"crowdfunding/back/internal/model"

	"github.com/gin-gonic/gin"
)

// GET /api/projects/:id/comments
func (h *Handler) ListComments(c *gin.Context) {
	page, limit := paginate(c)
	offset := (page - 1) * limit

	var total int64
	h.db.Model(&model.Comment{}).Where("project_id = ?", c.Param("id")).Count(&total)

	var comments []model.Comment
	h.db.Where("project_id = ?", c.Param("id")).
		Preload("User").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&comments)

	c.JSON(http.StatusOK, gin.H{
		"data":  comments,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// POST /api/projects/:id/comments
func (h *Handler) CreateComment(c *gin.Context) {
	var project model.Project
	if err := h.db.First(&project, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required,min=1,max=1500"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := model.Comment{
		UserID:    h.currentUserID(c),
		ProjectID: project.ID,
		Content:   input.Content,
	}
	if err := h.db.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
		return
	}

	h.db.Preload("User").First(&comment, comment.ID)
	c.JSON(http.StatusCreated, comment)
}

// DELETE /api/comments/:id
func (h *Handler) DeleteComment(c *gin.Context) {
	var comment model.Comment
	if err := h.db.First(&comment, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}

	uid := h.currentUserID(c)
	role := h.currentRole(c)
	if comment.UserID != uid && role != model.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	h.db.Delete(&comment)
	c.Status(http.StatusNoContent)
}
