package handler

import (
	"net/http"

	"crowdfunding/back/internal/model"

	"github.com/gin-gonic/gin"
)

// GET /api/categories
func (h *Handler) ListCategories(c *gin.Context) {
	var categories []model.Category
	if err := h.db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch categories"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// POST /api/admin/categories
func (h *Handler) AdminCreateCategory(c *gin.Context) {
	var input struct {
		Title string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat := model.Category{Title: input.Title}
	if err := h.db.Create(&cat).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "category already exists"})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

// DELETE /api/admin/categories/:id
func (h *Handler) AdminDeleteCategory(c *gin.Context) {
	if err := h.db.Delete(&model.Category{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
