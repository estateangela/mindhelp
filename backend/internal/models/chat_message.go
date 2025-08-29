package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ChatMessage 聊天訊息資料模型
type ChatMessage struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	Role      string         `json:"role" gorm:"size:10;not null"` // 'user' 或 'bot'
	Content   string         `json:"content" gorm:"type:text;not null"`
	Timestamp int64          `json:"timestamp" gorm:"not null"` // Unix milliseconds
	Model     string         `json:"model" gorm:"size:50"`       // AI 模型名稱
	Tokens    int            `json:"tokens" gorm:"default:0"`    // 使用的 token 數量
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 關聯
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定資料表名稱
func (ChatMessage) TableName() string {
	return "chat_messages"
}

// BeforeCreate 在創建前設定 UUID 和時間戳
func (cm *ChatMessage) BeforeCreate(tx *gorm.DB) error {
	if cm.ID == uuid.Nil {
		cm.ID = uuid.New()
	}
	if cm.Timestamp == 0 {
		cm.Timestamp = time.Now().UnixMilli()
	}
	return nil
}

// IsUser 檢查是否為使用者訊息
func (cm *ChatMessage) IsUser() bool {
	return cm.Role == "user"
}

// IsBot 檢查是否為機器人訊息
func (cm *ChatMessage) IsBot() bool {
	return cm.Role == "bot"
}
