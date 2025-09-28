package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification 通知模型
type Notification struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Message     string    `gorm:"type:text;not null" json:"message"`
	Type        string    `gorm:"size:50;not null;index" json:"type"`
	IsRead      bool      `gorm:"default:false;index" json:"is_read"`
	ActionURL   string    `gorm:"size:500" json:"action_url"`
	Data        string    `gorm:"type:jsonb" json:"data"` // 附加資料 JSON
	Priority    string    `gorm:"size:20;default:normal" json:"priority"`
	ExpiresAt   *time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 關聯
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
