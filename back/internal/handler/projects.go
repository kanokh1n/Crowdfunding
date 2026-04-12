package handler

import (
	"net/http"
	"time"

	"crowdfunding/back/internal/model"
	"crowdfunding/back/internal/moderation"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GET /api/projects
func (h *Handler) ListProjects(c *gin.Context) {
	page, limit := paginate(c)
	offset := (page - 1) * limit

	query := h.db.Model(&model.Project{}).Preload("User").Preload("Categories")

	// По умолчанию публичный список — только одобренные проекты
	status := c.Query("status")
	if status == "" {
		query = query.Where("projects.status = ?", model.StatusActive)
	} else {
		query = query.Where("projects.status = ?", status)
	}

	if cat := c.Query("category_id"); cat != "" {
		query = query.Joins("JOIN project_categories pc ON pc.project_id = projects.id").
			Where("pc.category_id = ?", cat)
	}
	if search := c.Query("search"); search != "" {
		query = query.Where("projects.title ILIKE ?", "%"+search+"%")
	}

	sort := c.Query("sort")
	switch sort {
	case "current_amount":
		query = query.Order("projects.current_amount DESC")
	case "likes_count":
		query = query.
			Joins("LEFT JOIN likes l ON l.project_id = projects.id").
			Group("projects.id").
			Order("COUNT(l.id) DESC")
	default:
		query = query.Order("projects.created_at DESC")
	}

	var total int64
	query.Count(&total)

	var projects []model.Project
	if err := query.Offset(offset).Limit(limit).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch projects"})
		return
	}

	// attach likes count
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

// GET /api/projects/:id
func (h *Handler) GetProject(c *gin.Context) {
	var project model.Project
	err := h.db.Preload("User").Preload("Categories").First(&project, c.Param("id")).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	h.db.Model(&model.Like{}).Where("project_id = ?", project.ID).Count(&project.LikesCount)
	c.JSON(http.StatusOK, project)
}

// POST /api/projects
func (h *Handler) CreateProject(c *gin.Context) {
	var input struct {
		Title        string     `json:"title"         binding:"required"`
		Description  string     `json:"description"`
		GoalAmount   float64    `json:"goal_amount"   binding:"required,gt=0"`
		EndDate      *time.Time `json:"end_date"`
		CategoryIDs  []uint     `json:"category_ids"`
		ProjectImg   string     `json:"project_img"`
		LinkTelegram string     `json:"link_telegram" binding:"omitempty,url"`
		LinkGithub   string     `json:"link_github"   binding:"omitempty,url"`
		LinkLinkedin string     `json:"link_linkedin" binding:"omitempty,url"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project := model.Project{
		UserID:       h.currentUserID(c),
		Title:        input.Title,
		Description:  input.Description,
		GoalAmount:   input.GoalAmount,
		EndDate:      input.EndDate,
		ProjectImg:   input.ProjectImg,
		LinkTelegram: input.LinkTelegram,
		LinkGithub:   input.LinkGithub,
		LinkLinkedin: input.LinkLinkedin,
		Status:       model.StatusPendingAI,
	}

	if len(input.CategoryIDs) > 0 {
		var categories []model.Category
		h.db.Find(&categories, input.CategoryIDs)
		project.Categories = categories
	}

	if err := h.db.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create project"})
		return
	}

	// Run AI check synchronously (replace with async job when needed)
	h.runAICheck(&project)

	h.db.Preload("User").Preload("Categories").Preload("Moderation").First(&project, project.ID)
	c.JSON(http.StatusCreated, project)
}

// PATCH /api/projects/:id
func (h *Handler) UpdateProject(c *gin.Context) {
	var project model.Project
	if err := h.db.First(&project, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	uid := h.currentUserID(c)
	role := h.currentRole(c)
	if project.UserID != uid && role != model.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var input struct {
		Title        *string              `json:"title"`
		Description  *string              `json:"description"`
		GoalAmount   *float64             `json:"goal_amount"   binding:"omitempty,gt=0"`
		EndDate      *time.Time           `json:"end_date"`
		ProjectImg   *string              `json:"project_img"`
		Status       *model.ProjectStatus `json:"status"`
		CategoryIDs  []uint               `json:"category_ids"`
		LinkTelegram *string              `json:"link_telegram" binding:"omitempty,url"`
		LinkGithub   *string              `json:"link_github"   binding:"omitempty,url"`
		LinkLinkedin *string              `json:"link_linkedin" binding:"omitempty,url"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]any{}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.GoalAmount != nil {
		updates["goal_amount"] = *input.GoalAmount
	}
	if input.EndDate != nil {
		updates["end_date"] = *input.EndDate
	}
	if input.ProjectImg != nil {
		updates["project_img"] = *input.ProjectImg
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.LinkTelegram != nil {
		updates["link_telegram"] = *input.LinkTelegram
	}
	if input.LinkGithub != nil {
		updates["link_github"] = *input.LinkGithub
	}
	if input.LinkLinkedin != nil {
		updates["link_linkedin"] = *input.LinkLinkedin
	}

	if len(updates) > 0 {
		h.db.Model(&project).Updates(updates)
	}

	if input.CategoryIDs != nil {
		var categories []model.Category
		h.db.Find(&categories, input.CategoryIDs)
		h.db.Model(&project).Association("Categories").Replace(categories)
	}

	h.db.Preload("User").Preload("Categories").First(&project, project.ID)
	h.db.Model(&model.Like{}).Where("project_id = ?", project.ID).Count(&project.LikesCount)
	c.JSON(http.StatusOK, project)
}

// DELETE /api/projects/:id
func (h *Handler) DeleteProject(c *gin.Context) {
	var project model.Project
	if err := h.db.First(&project, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	uid := h.currentUserID(c)
	role := h.currentRole(c)
	if project.UserID != uid && role != model.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	h.db.Delete(&project)
	c.Status(http.StatusNoContent)
}

// GET /api/users/me/projects
func (h *Handler) MyProjects(c *gin.Context) {
	page, limit := paginate(c)
	offset := (page - 1) * limit

	uid := h.currentUserID(c)
	var total int64
	h.db.Model(&model.Project{}).Where("user_id = ?", uid).Count(&total)

	var projects []model.Project
	h.db.Where("user_id = ?", uid).
		Preload("Categories").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&projects)

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

// helper for admin
func (h *Handler) findProjectByID(id string) (*model.Project, error) {
	var p model.Project
	err := h.db.Unscoped().Preload("User").Preload("Categories").First(&p, id).Error
	return &p, err
}

var _ = gorm.ErrRecordNotFound // keep gorm import

// runAICheck runs the stub AI moderation and updates project status.
func (h *Handler) runAICheck(project *model.Project) {
	result := moderation.RunAICheck(project.Title, project.Description)

	now := result.CheckedAt
	aiStatus := model.AIStatusPassed
	projectStatus := model.StatusPendingHuman
	if !result.Passed {
		aiStatus = model.AIStatusFailed
		projectStatus = model.StatusRejectedAI
	}

	mod := model.ProjectModeration{
		ProjectID:   project.ID,
		AIStatus:    aiStatus,
		AIScore:     result.Score,
		AIFlags:     moderation.FlagsToJSON(result.Flags),
		AICheckedAt: &now,
		HumanStatus: model.HumanStatusPending,
	}
	h.db.Create(&mod)
	h.db.Model(project).Update("status", projectStatus)
	project.Status = projectStatus
}
