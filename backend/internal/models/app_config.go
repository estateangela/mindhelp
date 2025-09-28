package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AppConfig 應用程式配置模型
type AppConfig struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Key         string    `gorm:"size:100;not null;unique;index" json:"key"`
	Value       string    `gorm:"type:text;not null" json:"value"`
	Type        string    `gorm:"size:20;default:string" json:"type"` // string, int, bool, json
	Description string    `gorm:"size:500" json:"description"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	Category    string    `gorm:"size:50;index" json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
