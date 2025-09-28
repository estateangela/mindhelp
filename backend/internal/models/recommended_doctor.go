package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RecommendedDoctor 網友推薦醫師＆診所模型
type RecommendedDoctor struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name            string         `json:"name" gorm:"size:255;not null"`        // 名稱
	Description     string         `json:"description" gorm:"type:text"`         // 描述
	ExperienceCount int            `json:"experience_count" gorm:"default:0"`    // 經驗次數
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (RecommendedDoctor) TableName() string {
	return "recommended_doctors"
}
