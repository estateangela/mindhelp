package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Counselor 諮商師模型
type Counselor struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name            string    `gorm:"size:255;not null" json:"name"`
	Title           string    `gorm:"size:255" json:"title"`
	Specialization  string    `gorm:"size:500" json:"specialization"`
	LicenseNumber   string    `gorm:"size:100;index" json:"license_number"`
	ExperienceYears int       `gorm:"default:0" json:"experience_years"`
	Education       string    `gorm:"type:text" json:"education"`
	Certifications  string    `gorm:"type:text" json:"certifications"`
	WorkLocation    string    `gorm:"type:text" json:"work_location"`
	ContactPhone    string    `gorm:"size:50" json:"contact_phone"`
	ContactEmail    string    `gorm:"size:255" json:"contact_email"`
	Website         string    `gorm:"size:500" json:"website"`
	Description     string    `gorm:"type:text" json:"description"`
	Languages       string    `gorm:"type:text" json:"languages"` // JSON 格式
	Availability    string    `gorm:"type:jsonb" json:"availability"` // JSON 格式
	PriceRange      string    `gorm:"size:100" json:"price_range"`
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
