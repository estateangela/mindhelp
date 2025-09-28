package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ChatSession 聊天會話資料模型
type ChatSession struct {
	ID                     uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID                 uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	Title                  string         `json:"title" gorm:"size:200"` // 會話標題，可從第一則訊息自動生成
	FirstMessageSnippet    string         `json:"first_message_snippet" gorm:"size:100"`
	LastUpdatedAt          time.Time      `json:"last_updated_at" gorm:"not null"`
	MessageCount           int            `json:"message_count" gorm:"default:0"`
	IsActive               bool           `json:"is_active" gorm:"default:true"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `json:"-" gorm:"index"`

	// 關聯
	User     User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Messages []ChatMessage `json:"messages,omitempty" gorm:"foreignKey:SessionID"`
}

// TableName 指定資料表名稱名稱
func (ChatSession) TableName() string {
	return "chat_sessions"
}

// BeforeCreate 在創建前設定前設定 UUID
func (cs *ChatSession) BeforeCreate(tx *gorm.DB) error {
	if cs.ID == uuid.Nil {
		cs.ID = uuid.New()
	}
	if cs.LastUpdatedAt.IsZero() {
		cs.LastUpdatedAt = time.Now()
	}
	return nil
}

// UpdateLastActivity 更新最後活動時間
func (cs *ChatSession) UpdateLastActivity() {
	cs.LastUpdatedAt = time.Now()
	cs.MessageCount++
}

