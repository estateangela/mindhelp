package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserSetting 用戶設定模型
type UserSetting struct {
	ID                  uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID              uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	NotifyNewArticle    bool           `json:"notify_new_article" gorm:"default:true"`
	NotifyPromotions    bool           `json:"notify_promotions" gorm:"default:false"`
	NotifySystemUpdates bool           `json:"notify_system_updates" gorm:"default:true"`
	PushToken           string         `json:"push_token" gorm:"size:500"`
	Platform            string         `json:"platform" gorm:"size:20"` // ios, android, web
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`

	// 關聯
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (UserSetting) TableName() string {
	return "user_settings"
}
