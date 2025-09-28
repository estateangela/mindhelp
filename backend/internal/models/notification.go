package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification ?�知資�?模�?
type Notification struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	Type      string         `json:"type" gorm:"size:20;not null"` // NEW_ARTICLE, PROMOTION, SYSTEM
	Title     string         `json:"title" gorm:"size:200;not null"`
	Body      string         `json:"body" gorm:"size:500"`
	IsRead    bool           `json:"is_read" gorm:"default:false"`
	Payload   string         `json:"payload" gorm:"type:text"` // JSON stored as text
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 關聯
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定資料表名稱�?�?
func (Notification) TableName() string {
	return "notifications"
}

// BeforeCreate 在創建前設定�?設�? UUID
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

// MarkAsRead 標�??�已讀
func (n *Notification) MarkAsRead() {
	n.IsRead = true
}


