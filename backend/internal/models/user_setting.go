package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserSetting 使用?�設定�??�模??
type UserSetting struct {
	ID                uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID            uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	NotifyNewArticle  bool           `json:"notify_new_article" gorm:"default:true"`
	NotifyPromotions  bool           `json:"notify_promotions" gorm:"default:false"`
	NotifySystemUpdates bool         `json:"notify_system_updates" gorm:"default:true"`
	PushToken         string         `json:"push_token" gorm:"size:255"`
	Platform          string         `json:"platform" gorm:"size:10"` // ios, android
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`

	// 關聯
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定資料表名稱�?�?
func (UserSetting) TableName() string {
	return "user_settings"
}

// BeforeCreate 在創建前設定�?設�? UUID
func (us *UserSetting) BeforeCreate(tx *gorm.DB) error {
	if us.ID == uuid.Nil {
		us.ID = uuid.New()
	}
	return nil
}


