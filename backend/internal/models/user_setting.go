package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserSetting 使用者設定模型
type UserSetting struct {
	ID                    uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID                uuid.UUID `gorm:"type:uuid;not null;index;unique" json:"user_id"`
	Theme                 string    `gorm:"size:20;default:light" json:"theme"`
	Language              string    `gorm:"size:10;default:zh-TW" json:"language"`
	EmailNotifications    bool      `gorm:"default:true" json:"email_notifications"`
	PushNotifications     bool      `gorm:"default:true" json:"push_notifications"`
	LocationServices      bool      `gorm:"default:true" json:"location_services"`
	PrivacyMode           bool      `gorm:"default:false" json:"privacy_mode"`
	DataCollection        bool      `gorm:"default:true" json:"data_collection"`
	AutoSave              bool      `gorm:"default:true" json:"auto_save"`
	FontSize              string    `gorm:"size:10;default:medium" json:"font_size"`
	Timezone              string    `gorm:"size:50;default:Asia/Taipei" json:"timezone"`
	Currency              string    `gorm:"size:10;default:TWD" json:"currency"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`

	// 關聯
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
