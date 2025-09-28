package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User 使用者資料模型
type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email     string         `json:"email" gorm:"size:255;uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"` // 密碼不會在 JSON 中返回
	Username  string         `json:"username" gorm:"size:50;not null"`
	FullName  string         `json:"full_name" gorm:"size:100"`
	Phone     string         `json:"phone" gorm:"size:20"`
	Avatar    string         `json:"avatar" gorm:"size:255"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	LastLogin *time.Time     `json:"last_login"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 關聯
	ChatMessages []ChatMessage `json:"chat_messages,omitempty" gorm:"foreignKey:UserID"`
	Locations    []Location    `json:"locations,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定資料表名稱
func (User) TableName() string {
	return "users"
}

// BeforeCreate 在創建前設定 UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
