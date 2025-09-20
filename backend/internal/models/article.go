package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Article 專家文章資料模型
type Article struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title       string         `json:"title" gorm:"size:200;not null"`
	Author      string         `json:"author" gorm:"size:100;not null"`
	AuthorTitle string         `json:"author_title" gorm:"size:100"`
	PublishDate time.Time      `json:"publish_date" gorm:"not null"`
	Summary     string         `json:"summary" gorm:"size:500"`
	Content     string         `json:"content" gorm:"type:text;not null"`
	Tags        string         `json:"tags" gorm:"type:text"` // JSON array stored as text
	ImageURL    string         `json:"image_url" gorm:"size:255"`
	IsPublished bool           `json:"is_published" gorm:"default:true"`
	ViewCount   int64          `json:"view_count" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 關聯
	Bookmarks []Bookmark `json:"bookmarks,omitempty" gorm:"foreignKey:ArticleID"`
}

// TableName 指定資料表名稱名稱
func (Article) TableName() string {
	return "articles"
}

// BeforeCreate 在創建前設定前設定 UUID
func (a *Article) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

