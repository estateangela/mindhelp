package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Share 分享資料模型
type Share struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uniqueidentifier;primary_key"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uniqueidentifier;not null;index"`
	ContentType string         `json:"content_type" gorm:"size:20;not null"` // article, location, quiz
	ContentID   uuid.UUID      `json:"content_id" gorm:"type:uniqueidentifier;not null;index"`
	Platform    string         `json:"platform" gorm:"size:20"`            // facebook, twitter, line, whatsapp, email, copy
	ShareURL    string         `json:"share_url" gorm:"size:500;not null"` // 完整分享連結
	ShortURL    string         `json:"short_url" gorm:"size:100"`          // 短連結
	Message     string         `json:"message" gorm:"size:500"`            // 分享時的訊息
	ViewCount   int64          `json:"view_count" gorm:"default:0"`        // 點擊次數
	ExpiresAt   *time.Time     `json:"expires_at"`                         // 過期時間 (可選)
	IsActive    bool           `json:"is_active" gorm:"default:true"`      // 是否啟用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 關聯
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// ShareClick 分享點擊記錄模型
type ShareClick struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uniqueidentifier;primary_key"`
	ShareID   uuid.UUID      `json:"share_id" gorm:"type:uniqueidentifier;not null;index"`
	IPAddress string         `json:"ip_address" gorm:"size:45"` // 支援 IPv6
	UserAgent string         `json:"user_agent" gorm:"size:500"`
	Referer   string         `json:"referer" gorm:"size:500"`
	Country   string         `json:"country" gorm:"size:10"` // 國家代碼
	City      string         `json:"city" gorm:"size:100"`   // 城市名稱
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 關聯
	Share Share `json:"share,omitempty" gorm:"foreignKey:ShareID"`
}

// TableName 指定資料表名稱
func (Share) TableName() string {
	return "shares"
}

func (ShareClick) TableName() string {
	return "share_clicks"
}

// BeforeCreate hooks
func (s *Share) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func (sc *ShareClick) BeforeCreate(tx *gorm.DB) error {
	if sc.ID == uuid.Nil {
		sc.ID = uuid.New()
	}
	return nil
}

// IsExpired 檢查分享是否已過期
func (s *Share) IsExpired() bool {
	if s.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*s.ExpiresAt)
}

// IncrementViewCount 增加觀看次數
func (s *Share) IncrementViewCount(db *gorm.DB) error {
	return db.Model(s).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}
