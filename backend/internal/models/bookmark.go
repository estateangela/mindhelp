package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Bookmark 收藏資料模型
type Bookmark struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uniqueidentifier;primary_key"`
	UserID       uuid.UUID      `json:"user_id" gorm:"type:uniqueidentifier;not null;index"`
	ResourceType string         `json:"resource_type" gorm:"size:20;not null"` // article, location
	ArticleID    *uuid.UUID     `json:"article_id" gorm:"type:uniqueidentifier;index"`
	LocationID   *uuid.UUID     `json:"location_id" gorm:"type:uniqueidentifier;index"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// 關聯
	User     User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Article  *Article  `json:"article,omitempty" gorm:"foreignKey:ArticleID"`
	Location *Location `json:"location,omitempty" gorm:"foreignKey:LocationID"`
}

// TableName 指定資料表名稱
func (Bookmark) TableName() string {
	return "bookmarks"
}

// BeforeCreate 在創建前設定 UUID
func (b *Bookmark) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// IsArticleBookmark 檢查是否為文章收藏
func (b *Bookmark) IsArticleBookmark() bool {
	return b.ResourceType == "article" && b.ArticleID != nil
}

// IsLocationBookmark 檢查是否為地點收藏
func (b *Bookmark) IsLocationBookmark() bool {
	return b.ResourceType == "location" && b.LocationID != nil
}
