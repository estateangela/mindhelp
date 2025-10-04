package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification 通知模型
type Notification struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	Title     string         `json:"title" gorm:"size:255;not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Type      string         `json:"type" gorm:"size:50;not null"` // hourly_reminder, weekly_bulletin, system, etc.
	IsRead    bool           `json:"is_read" gorm:"default:false"`
	Payload   string         `json:"payload" gorm:"type:text"` // JSON 格式的額外資料
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 關聯
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (Notification) TableName() string {
	return "notifications"
}
