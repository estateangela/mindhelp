package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Share ?�享資�?模�?
type Share struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	ContentType string         `json:"content_type" gorm:"size:20;not null"` // article, location, quiz
	ContentID   uuid.UUID      `json:"content_id" gorm:"type:uuid;not null;index"`
	Platform    string         `json:"platform" gorm:"size:20"`            // facebook, twitter, line, whatsapp, email, copy
	ShareURL    string         `json:"share_url" gorm:"size:500;not null"` // 完整?�享???
	ShortURL    string         `json:"short_url" gorm:"size:100"`          // ?��??
	Message     string         `json:"message" gorm:"size:500"`            // ?�享?��?訊息
	ViewCount   int64          `json:"view_count" gorm:"default:0"`        // 點�?次數
	ExpiresAt   *time.Time     `json:"expires_at"`                         // ?��??��? (?�選)
	IsActive    bool           `json:"is_active" gorm:"default:true"`      // ?�否?�用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 關聯
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// ShareClick ?�享點�?記�?模�?
type ShareClick struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ShareID   uuid.UUID      `json:"share_id" gorm:"type:uuid;not null;index"`
	IPAddress string         `json:"ip_address" gorm:"size:45"` // ?�援 IPv6
	UserAgent string         `json:"user_agent" gorm:"size:500"`
	Referer   string         `json:"referer" gorm:"size:500"`
	Country   string         `json:"country" gorm:"size:10"` // ?�家�?��
	City      string         `json:"city" gorm:"size:100"`   // ?��??�稱
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 關聯
	Share Share `json:"share,omitempty" gorm:"foreignKey:ShareID"`
}

// TableName 指定資料表名稱�?�?
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

// IsExpired 檢查?�享?�否已�???
func (s *Share) IsExpired() bool {
	if s.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*s.ExpiresAt)
}

// IncrementViewCount 增�?觀?�次??
func (s *Share) IncrementViewCount(db *gorm.DB) error {
	return db.Model(s).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}


