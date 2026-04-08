package handler

import (
	"fmt"
	"time"

	"crowdfunding/back/internal/config"
	"crowdfunding/back/internal/middleware"
	"crowdfunding/back/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	db  *gorm.DB
	rdb *redis.Client
	cfg *config.Config
}

func New(db *gorm.DB, rdb *redis.Client, cfg *config.Config) *Handler {
	return &Handler{db: db, rdb: rdb, cfg: cfg}
}

func (h *Handler) currentUserID(c *gin.Context) uint {
	return c.MustGet("claims").(*middleware.Claims).UserID
}

func (h *Handler) currentRole(c *gin.Context) model.Role {
	return c.MustGet("claims").(*middleware.Claims).Role
}

func (h *Handler) generateAccessToken(userID uint, role model.Role) (string, error) {
	claims := &middleware.Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.cfg.JWTSecret))
}

func paginate(c *gin.Context) (page, limit int) {
	page = intQuery(c, "page", 1)
	limit = intQuery(c, "limit", 10)
	if limit > 100 {
		limit = 100
	}
	return
}

func intQuery(c *gin.Context, key string, def int) int {
	var v int
	if _, err := fmt.Sscan(c.Query(key), &v); err != nil || v < 1 {
		return def
	}
	return v
}
