package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Bookmark 收藏模型
type Bookmark struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ResourceID uuid.UUID `gorm:"type:uuid;not null;index" json:"resource_id"`
	ResourceType string  `gorm:"size:50;not null;index" json:"resource_type"` // "article", "location", etc.
	Notes      string    `gorm:"type:text" json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// 關聯
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
