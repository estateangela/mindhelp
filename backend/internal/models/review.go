package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Review 評�?資�?模�?
type Review struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	ResourceID uuid.UUID      `json:"resource_id" gorm:"type:uuid;not null;index"` // Location ID
	Rating     int            `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Comment    string         `json:"comment" gorm:"size:1000"`
	IsHelpful  int            `json:"is_helpful" gorm:"default:0"` // ?�用票數
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// 關聯
	User     User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Location Location `json:"location,omitempty" gorm:"foreignKey:ResourceID"`
}

// TableName 指定資料表名稱�?�?
func (Review) TableName() string {
	return "reviews"
}

// BeforeCreate 在創建前設定�?設�? UUID
func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}


