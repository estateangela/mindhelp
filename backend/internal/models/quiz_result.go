package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// QuizResult 測驗結果模型
type QuizResult struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	QuizID    uuid.UUID `gorm:"type:uuid;not null;index" json:"quiz_id"`
	Score     int       `gorm:"not null" json:"score"`
	MaxScore  int       `gorm:"not null" json:"max_score"`
	Answers   string    `gorm:"type:jsonb" json:"answers"` // JSON 格式存儲答案
	Completed bool      `gorm:"default:false" json:"completed"`
	TimeSpent int       `gorm:"default:0" json:"time_spent"` // 花費時間（秒）
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 關聯
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Quiz Quiz `gorm:"foreignKey:QuizID" json:"quiz,omitempty"`
}
