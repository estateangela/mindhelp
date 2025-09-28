package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CounselingCenter 諮商所模型
type CounselingCenter struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Address     string    `gorm:"type:text" json:"address"`
	Phone       string    `gorm:"size:50" json:"phone"`
	Email       string    `gorm:"size:255" json:"email"`
	Website     string    `gorm:"size:500" json:"website"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Services    string    `gorm:"type:text" json:"services"` // JSON 格式
	Hours       string    `gorm:"type:jsonb" json:"hours"` // JSON 格式
	PriceRange  string    `gorm:"size:100" json:"price_range"`
	Rating      float64   `gorm:"default:0" json:"rating"`
	ReviewCount int       `gorm:"default:0" json:"review_count"`
	IsVerified  bool      `gorm:"default:false" json:"is_verified"`
	IsActive    bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 關聯
	Reviews []Review `gorm:"foreignKey:ResourceID;references:ID" json:"reviews,omitempty"`
}
