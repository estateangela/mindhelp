package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RecommendedDoctor 推薦醫師模型
type RecommendedDoctor struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name            string    `gorm:"size:255;not null" json:"name"`
	Specialty       string    `gorm:"size:255" json:"specialty"`
	Description     string    `gorm:"type:text" json:"description"`
	Hospital        string    `gorm:"size:255" json:"hospital"`
	Location        string    `gorm:"size:255" json:"location"`
	Contact         string    `gorm:"size:255" json:"contact"`
	ExperienceYears int       `gorm:"default:0" json:"experience_years"`
	Rating          float64   `gorm:"default:0" json:"rating"`
	ReviewCount     int       `gorm:"default:0" json:"review_count"`
	IsVerified      bool      `gorm:"default:false" json:"is_verified"`
	IsActive        bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// 關聯
	Reviews []Review `gorm:"foreignKey:ResourceID;references:ID" json:"reviews,omitempty"`
}
