package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Article 文章模型
type Article struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string    `gorm:"size:500;not null" json:"title"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	Summary     string    `gorm:"size:1000" json:"summary"`
	Author      string    `gorm:"size:255" json:"author"`
	Category    string    `gorm:"size:100;index" json:"category"`
	Tags        string    `gorm:"type:text" json:"tags"` // JSON 格式存儲標籤
	Views       int       `gorm:"default:0" json:"views"`
	PublishDate time.Time `gorm:"index" json:"publish_date"`
	IsPublished bool      `gorm:"default:false;index" json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 關聯
	Bookmarks []Bookmark `gorm:"foreignKey:ResourceID;references:ID" json:"bookmarks,omitempty"`
}
