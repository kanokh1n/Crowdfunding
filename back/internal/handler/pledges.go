package handler

import (
	"net/http"

	"crowdfunding/back/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// POST /api/projects/:id/pledges
func (h *Handler) CreatePledge(c *gin.Context) {
	var project model.Project
	if err := h.db.First(&project, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	if project.Status != model.StatusActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project is not accepting pledges"})
		return
	}

	var input struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pledge := model.Pledge{
		UserID:    h.currentUserID(c),
		ProjectID: project.ID,
		Amount:    input.Amount,
	}

	tx := h.db.Begin()
	if err := tx.Create(&pledge).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create pledge"})
		return
	}

	if err := tx.Model(&model.Project{}).Where("id = ?", project.ID).
		UpdateColumn("current_amount", gorm.Expr("current_amount + ?", input.Amount)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update project amount"})
		return
	}
	tx.Commit()

	h.db.Preload("User").First(&pledge, pledge.ID)
	c.JSON(http.StatusCreated, pledge)
}

// GET /api/projects/:id/pledges
func (h *Handler) ListPledges(c *gin.Context) {
	page, limit := paginate(c)
	offset := (page - 1) * limit

	var total int64
	h.db.Model(&model.Pledge{}).Where("project_id = ?", c.Param("id")).Count(&total)

	var pledges []model.Pledge
	h.db.Where("project_id = ?", c.Param("id")).
		Preload("User").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&pledges)

	c.JSON(http.StatusOK, gin.H{
		"data":  pledges,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}
