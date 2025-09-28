package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Review 評論模型
type Review struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ResourceID  uuid.UUID `gorm:"type:uuid;not null;index" json:"resource_id"`
	ResourceType string   `gorm:"size:50;not null;index" json:"resource_type"`
	Rating      int       `gorm:"check:rating >= 1 AND rating <= 5" json:"rating"`
	Title       string    `gorm:"size:255" json:"title"`
	Content     string    `gorm:"type:text" json:"content"`
	IsVerified  bool      `gorm:"default:false" json:"is_verified"`
	IsHelpful   int       `gorm:"default:0" json:"is_helpful"` // 有幫助的投票數
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 關聯
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
