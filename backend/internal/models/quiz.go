package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Quiz 測驗模型
type Quiz struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string    `gorm:"size:500;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Questions   string    `gorm:"type:jsonb" json:"questions"` // JSON 格式存儲問題
	Category    string    `gorm:"size:100;index" json:"category"`
	Difficulty  string    `gorm:"size:20;default:easy" json:"difficulty"`
	TimeLimit   int       `gorm:"default:0" json:"time_limit"` // 0 表示無時間限制（分鐘）
	MaxScore    int       `gorm:"default:100" json:"max_score"`
	PassScore   int       `gorm:"default:60" json:"pass_score"`
	IsActive    bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 關聯
	Results []QuizResult `gorm:"foreignKey:QuizID" json:"results,omitempty"`
}
