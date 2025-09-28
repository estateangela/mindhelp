package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CounselingCenter 台北諮商所模型
type CounselingCenter struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name            string         `json:"name" gorm:"size:255;not null"`        // 機構名稱
	Address         string         `json:"address" gorm:"size:500"`              // 地址
	Phone           string         `json:"phone" gorm:"size:50"`                 // 電話
	OnlineCounseling bool          `json:"online_counseling" gorm:"default:false"` // 通訊心理諮商
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (CounselingCenter) TableName() string {
	return "counseling_centers"
}
