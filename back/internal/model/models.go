package model

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type ProjectStatus string

const (
	StatusPendingAI    ProjectStatus = "pending_ai"    // только создан, ждёт AI-проверки
	StatusRejectedAI   ProjectStatus = "rejected_ai"   // не прошёл AI
	StatusPendingHuman ProjectStatus = "pending_human" // прошёл AI, ждёт модератора
	StatusActive       ProjectStatus = "active"        // одобрен модератором
	StatusRejected     ProjectStatus = "rejected"      // отклонён модератором
	StatusCompleted    ProjectStatus = "completed"     // сбор завершён
	StatusCancelled    ProjectStatus = "cancelled"     // отменён автором
)

// ─── Moderation ───────────────────────────────────────────────────────────────

type AIStatus string

const (
	AIStatusPending AIStatus = "pending"
	AIStatusPassed  AIStatus = "passed"
	AIStatusFailed  AIStatus = "failed"
)

type HumanStatus string

const (
	HumanStatusPending  HumanStatus = "pending"
	HumanStatusApproved HumanStatus = "approved"
	HumanStatusRejected HumanStatus = "rejected"
)

type ProjectModeration struct {
	ID        uint    `gorm:"primaryKey"              json:"id"`
	ProjectID uint    `gorm:"uniqueIndex;not null"    json:"project_id"`
	Project   Project `gorm:"foreignKey:ProjectID"    json:"project,omitempty"`

	// AI stage
	AIStatus    AIStatus   `gorm:"default:'pending'"       json:"ai_status"`
	AIScore     float64    `gorm:"default:0"               json:"ai_score"`
	AIFlags     string     `gorm:"type:text;default:'[]'"  json:"ai_flags"` // JSON array
	AICheckedAt *time.Time `                               json:"ai_checked_at"`

	// Human stage
	HumanStatus      HumanStatus `gorm:"default:'pending'"   json:"human_status"`
	ModeratorID      *uint       `                           json:"moderator_id"`
	Moderator        *User       `gorm:"foreignKey:ModeratorID" json:"moderator,omitempty"`
	ModeratorNote    string      `                           json:"moderator_note"`
	HumanModeratedAt *time.Time  `                           json:"human_moderated_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ─── User ────────────────────────────────────────────────────────────────────

type User struct {
	ID          uint      `gorm:"primaryKey"                   json:"id"`
	Email       string    `gorm:"uniqueIndex;not null"         json:"email"`
	Password    string    `gorm:"not null"                     json:"-"`
	FIO         string    `                                    json:"fio"`
	Description string    `                                    json:"description"`
	ProfileImg  string    `                                    json:"profile_img"`
	Phone       string    `                                    json:"phone"`
	Role        Role      `gorm:"default:'user'"               json:"role"`
	IsVerified  bool      `gorm:"default:false"                json:"is_verified"`
	CreatedAt   time.Time `                                    json:"created_at"`
	UpdatedAt   time.Time `                                    json:"updated_at"`
}

// ─── Email verification token ─────────────────────────────────────────────────

type EmailToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

// ─── Category ─────────────────────────────────────────────────────────────────

type Category struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Title string `gorm:"uniqueIndex;not null" json:"title"`
}

// ─── Project ──────────────────────────────────────────────────────────────────

type ProjectImage struct {
	ID        uint   `gorm:"primaryKey"         json:"id"`
	ProjectID uint   `gorm:"not null;index"     json:"project_id"`
	URL       string `gorm:"not null"           json:"url"`
	Position  int    `gorm:"default:0"          json:"position"`
}

type Project struct {
	ID               uint               `gorm:"primaryKey"                   json:"id"`
	UserID           uint               `gorm:"not null"                     json:"user_id"`
	User             User               `gorm:"foreignKey:UserID"            json:"user,omitempty"`
	Title            string             `gorm:"not null"                     json:"title"`
	ShortDescription string             `                                    json:"short_description"`
	Description      string             `                                    json:"description"`
	GoalAmount       float64            `gorm:"type:decimal(10,2);not null"  json:"goal_amount"`
	CurrentAmount    float64            `gorm:"type:decimal(10,2);default:0" json:"current_amount"`
	ProjectImg       string             `                                    json:"project_img"`
	RoadmapFile      string             `                                    json:"roadmap_file"`
	LinkTelegram     string             `                                    json:"link_telegram"`
	LinkGithub       string             `                                    json:"link_github"`
	LinkLinkedin     string             `                                    json:"link_linkedin"`
	Status           ProjectStatus      `gorm:"default:'pending_ai'"         json:"status"`
	Images           []ProjectImage     `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"images,omitempty"`
	Moderation       *ProjectModeration `gorm:"foreignKey:ProjectID"     json:"moderation,omitempty"`
	EndDate          *time.Time         `                                    json:"end_date"`
	Categories       []Category         `gorm:"many2many:project_categories" json:"categories,omitempty"`
	LikesCount       int64              `gorm:"-"                            json:"likes_count"`
	CreatedAt        time.Time          `                                    json:"created_at"`
	UpdatedAt        time.Time          `                                    json:"updated_at"`
	DeletedAt        gorm.DeletedAt     `gorm:"index"                        json:"-"`
}

// ─── Pledge ───────────────────────────────────────────────────────────────────

type Pledge struct {
	ID        uint      `gorm:"primaryKey"                  json:"id"`
	UserID    uint      `gorm:"not null"                    json:"user_id"`
	User      User      `gorm:"foreignKey:UserID"           json:"user,omitempty"`
	ProjectID uint      `gorm:"not null"                    json:"project_id"`
	Amount    float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	CreatedAt time.Time `                                   json:"created_at"`
}

// ─── Comment ──────────────────────────────────────────────────────────────────

type Comment struct {
	ID        uint      `gorm:"primaryKey"        json:"id"`
	UserID    uint      `gorm:"not null"          json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProjectID uint      `gorm:"not null"          json:"project_id"`
	Content   string    `gorm:"not null"          json:"content"`
	CreatedAt time.Time `                         json:"created_at"`
}

// ─── Like ─────────────────────────────────────────────────────────────────────

type Like struct {
	ID        uint      `gorm:"primaryKey"                                            json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_project_like"            json:"user_id"`
	ProjectID uint      `gorm:"not null;uniqueIndex:idx_user_project_like"            json:"project_id"`
	CreatedAt time.Time `                                                             json:"created_at"`
}

// ─── Notification ─────────────────────────────────────────────────────────────

type NotificationType string

const (
	NotifAIPassed NotificationType = "ai_passed"
	NotifAIFailed NotificationType = "ai_failed"
	NotifInvite   NotificationType = "invite"
)

type Notification struct {
	ID        uint             `gorm:"primaryKey"    json:"id"`
	UserID    uint             `gorm:"not null"      json:"user_id"`
	ProjectID *uint            `                     json:"project_id"`
	Type      NotificationType `gorm:"not null"      json:"type"`
	Title     string           `gorm:"not null"      json:"title"`
	Body      string           `                     json:"body"`
	IsRead    bool             `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time        `                     json:"created_at"`
}

// ─── Message ──────────────────────────────────────────────────────────────────

type Message struct {
	ID          uint      `gorm:"primaryKey"                  json:"id"`
	SenderID    uint      `gorm:"not null"                    json:"sender_id"`
	Sender      User      `gorm:"foreignKey:SenderID"         json:"sender,omitempty"`
	RecipientID uint      `gorm:"not null"                    json:"recipient_id"`
	Recipient   User      `gorm:"foreignKey:RecipientID"      json:"recipient,omitempty"`
	ProjectID   uint      `gorm:"not null"                    json:"project_id"`
	Project     Project   `gorm:"foreignKey:ProjectID"        json:"project,omitempty"`
	Title       string    `gorm:"not null"                    json:"title"`
	Content     string    `gorm:"not null"                    json:"content"`
	IsRead      bool      `gorm:"default:false"               json:"is_read"`
	CreatedAt   time.Time `                                   json:"created_at"`
}
