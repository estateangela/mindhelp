package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AppConfig 應用程式配置資料模型
type AppConfig struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Key         string         `json:"key" gorm:"size:50;uniqueIndex;not null"`
	Value       string         `json:"value" gorm:"type:text;not null"`
	Description string         `json:"description" gorm:"size:200"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定資料表名稱名稱
func (AppConfig) TableName() string {
	return "app_configs"
}

// BeforeCreate 在創建前設定前設定 UUID
func (ac *AppConfig) BeforeCreate(tx *gorm.DB) error {
	if ac.ID == uuid.Nil {
		ac.ID = uuid.New()
	}
	return nil
}

