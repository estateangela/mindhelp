package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Share 分享模型
type Share struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ResourceID  uuid.UUID `gorm:"type:uuid;not null;index" json:"resource_id"`
	ResourceType string   `gorm:"size:50;not null;index" json:"resource_type"`
	ShareToken  string    `gorm:"size:100;unique;index" json:"share_token"`
	Title       string    `gorm:"size:255" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	ExpiresAt   *time.Time `json:"expires_at"`
	ViewCount   int       `gorm:"default:0" json:"view_count"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 關聯
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
