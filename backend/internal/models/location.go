package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Location 位置資料模型
type Location struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Description string         `json:"description" gorm:"type:text"`
	Address     string         `json:"address" gorm:"size:255"`
	Latitude    float64        `json:"latitude" gorm:"type:decimal(10,8);not null"`
	Longitude   float64        `json:"longitude" gorm:"type:decimal(11,8);not null"`
	Category    string         `json:"category" gorm:"size:50"` // 心理健康資源類別
	Phone       string         `json:"phone" gorm:"size:20"`
	Website     string         `json:"website" gorm:"size:255"`
	Rating      float64        `json:"rating" gorm:"type:decimal(3,2);default:0"`
	IsPublic    bool           `json:"is_public" gorm:"default:false"` // 是否為公開資源
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 關聯
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定資料表名稱
func (Location) TableName() string {
	return "locations"
}

// BeforeCreate 在創建前設定 UUID
func (l *Location) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}
