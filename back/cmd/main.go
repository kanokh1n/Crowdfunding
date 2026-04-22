package main

import (
	"log"

	"crowdfunding/back/internal/config"
	"crowdfunding/back/internal/database"
	"crowdfunding/back/internal/handler"
	"crowdfunding/back/internal/middleware"
	"crowdfunding/back/internal/model"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db := database.NewPostgres(cfg.DBUrl)
	rdb := database.NewRedis(cfg.RedisUrl)

	// Auto migrate all models
	if err := db.AutoMigrate(
		&model.User{},
		&model.EmailToken{},
		&model.Category{},
		&model.Project{},
		&model.ProjectModeration{},
		&model.Pledge{},
		&model.Comment{},
		&model.Like{},
		&model.Message{},
		&model.Notification{},
		&model.ProjectImage{},
	); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	h := handler.New(db, rdb, cfg)

	r := gin.Default()
	r.Static("/uploads", "./uploads")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	// ── Public ──────────────────────────────────────────────────────────────
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/verify-email", h.VerifyEmail)
		auth.POST("/refresh", h.RefreshToken)

		api.GET("/projects", h.ListProjects)
		api.GET("/projects/:id", h.GetProject)
		api.GET("/projects/:id/comments", h.ListComments)
		api.GET("/categories", h.ListCategories)
	}

	// ── Authenticated ────────────────────────────────────────────────────────
	secured := api.Group("")
	secured.Use(middleware.Auth(cfg.JWTSecret))
	{
		// users
		secured.GET("/users/me", h.GetMe)
		secured.PATCH("/users/me", h.UpdateMe)
		secured.GET("/users/me/projects", h.MyProjects)
		secured.GET("/users/:id", h.GetUser)

		// projects
		secured.POST("/projects", h.CreateProject)
		secured.PATCH("/projects/:id", h.UpdateProject)
		secured.DELETE("/projects/:id", h.DeleteProject)

		// pledges
		secured.POST("/projects/:id/pledges", h.CreatePledge)
		secured.GET("/projects/:id/pledges", h.ListPledges)

		// comments
		secured.POST("/projects/:id/comments", h.CreateComment)
		secured.PATCH("/comments/:id", h.UpdateComment)
		secured.DELETE("/comments/:id", h.DeleteComment)

		// likes
		secured.POST("/projects/:id/like", h.LikeProject)
		secured.DELETE("/projects/:id/like", h.UnlikeProject)

		// messages
		secured.GET("/messages", h.ListMessages)
		secured.POST("/messages", h.SendMessage)
		secured.PATCH("/messages/:id/read", h.MarkMessageRead)

		// notifications
		secured.GET("/notifications", h.ListNotifications)
		secured.PATCH("/notifications/:id/read", h.MarkNotificationRead)

		// upload
		secured.POST("/upload", h.UploadImage)
		secured.POST("/upload/document", h.UploadDocument)
	}

	// ── Admin ────────────────────────────────────────────────────────────────
	admin := api.Group("/admin")
	admin.Use(middleware.Auth(cfg.JWTSecret), middleware.RequireAdmin())
	{
		admin.GET("/users", h.AdminListUsers)
		admin.PATCH("/users/:id", h.AdminUpdateUser)

		admin.GET("/projects", h.AdminListProjects)
		admin.PATCH("/projects/:id", h.AdminUpdateProject)
		admin.DELETE("/projects/:id", h.AdminDeleteProject)

		admin.POST("/categories", h.AdminCreateCategory)
		admin.DELETE("/categories/:id", h.AdminDeleteCategory)

		admin.GET("/moderation", h.AdminModerationList)
		admin.GET("/moderation/:project_id", h.AdminModerationGet)
		admin.PATCH("/moderation/:project_id", h.AdminModerationDecide)
		admin.POST("/moderation/:project_id/invite", h.AdminModerationInvite)
		admin.POST("/moderation/:project_id/recheck", h.AdminModerationRecheck)
	}

	log.Printf("Server starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
