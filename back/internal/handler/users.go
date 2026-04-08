package handler

import (
	"net/http"

	"crowdfunding/back/internal/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GET /api/users/me
func (h *Handler) GetMe(c *gin.Context) {
	var user model.User
	if err := h.db.First(&user, h.currentUserID(c)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// PATCH /api/users/me
func (h *Handler) UpdateMe(c *gin.Context) {
	var input struct {
		FIO         *string `json:"fio"`
		Description *string `json:"description"`
		Phone       *string `json:"phone"`
		ProfileImg  *string `json:"profile_img"`
		Password    *string `json:"password" binding:"omitempty,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]any{}
	if input.FIO != nil {
		updates["fio"] = *input.FIO
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Phone != nil {
		updates["phone"] = *input.Phone
	}
	if input.ProfileImg != nil {
		updates["profile_img"] = *input.ProfileImg
	}
	if input.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
		updates["password"] = string(hash)
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	id := h.currentUserID(c)
	if err := h.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	var user model.User
	h.db.First(&user, id)
	c.JSON(http.StatusOK, user)
}

// GET /api/users/:id
func (h *Handler) GetUser(c *gin.Context) {
	var user model.User
	if err := h.db.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
