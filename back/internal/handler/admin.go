package handler

import (
	"net/http"

	"crowdfunding/back/internal/model"

	"github.com/gin-gonic/gin"
)

// GET /api/admin/users
func (h *Handler) AdminListUsers(c *gin.Context) {
	page, limit := paginate(c)
	offset := (page - 1) * limit

	var total int64
	h.db.Model(&model.User{}).Count(&total)

	var users []model.User
	h.db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"data":  users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// PATCH /api/admin/users/:id
func (h *Handler) AdminUpdateUser(c *gin.Context) {
	var user model.User
	if err := h.db.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var input struct {
		Role       *model.Role `json:"role"`
		IsVerified *bool       `json:"is_verified"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]any{}
	if input.Role != nil {
		updates["role"] = *input.Role
	}
	if input.IsVerified != nil {
		updates["is_verified"] = *input.IsVerified
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	h.db.Model(&user).Updates(updates)
	h.db.First(&user, user.ID)
	c.JSON(http.StatusOK, user)
}

// GET /api/admin/projects
func (h *Handler) AdminListProjects(c *gin.Context) {
	page, limit := paginate(c)
	offset := (page - 1) * limit

	query := h.db.Unscoped().Model(&model.Project{}).Preload("User").Preload("Categories")

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var projects []model.Project
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&projects)

	for i := range projects {
		h.db.Model(&model.Like{}).Where("project_id = ?", projects[i].ID).Count(&projects[i].LikesCount)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  projects,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// PATCH /api/admin/projects/:id
func (h *Handler) AdminUpdateProject(c *gin.Context) {
	var project model.Project
	if err := h.db.Unscoped().First(&project, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	var input struct {
		Status *model.ProjectStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Status != nil {
		h.db.Unscoped().Model(&project).Update("status", *input.Status)
	}

	h.db.Unscoped().Preload("User").Preload("Categories").First(&project, project.ID)
	c.JSON(http.StatusOK, project)
}

// DELETE /api/admin/projects/:id  (hard delete)
func (h *Handler) AdminDeleteProject(c *gin.Context) {
	result := h.db.Unscoped().Delete(&model.Project{}, c.Param("id"))
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
