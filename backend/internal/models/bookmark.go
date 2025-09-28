package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Bookmark 書籤資料模型
type Bookmark struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	ResourceType string         `json:"resource_type" gorm:"size:20;not null"` // article, location
	ArticleID    *uuid.UUID     `json:"article_id" gorm:"type:uuid;index"`
	LocationID   *uuid.UUID     `json:"location_id" gorm:"type:uuid;index"`
	Title        string         `json:"title" gorm:"size:200;not null"`
	Description  string         `json:"description" gorm:"size:500"`
	URL          string         `json:"url" gorm:"size:255"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// 關聯
	User     User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Article  *Article  `json:"article,omitempty" gorm:"foreignKey:ArticleID"`
	Location *Location `json:"location,omitempty" gorm:"foreignKey:LocationID"`
}

// TableName 指定資料表名稱名稱
func (Bookmark) TableName() string {
	return "bookmarks"
}

// BeforeCreate 在創建前設定前設定 UUID
func (b *Bookmark) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

