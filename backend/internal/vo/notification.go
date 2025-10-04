package vo

import (
	"time"

	"github.com/google/uuid"
)

// NotificationVO 通知視圖物件
type NotificationVO struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// NotificationListVO 通知列表視圖物件
type NotificationListVO struct {
	Notifications []NotificationVO `json:"notifications"`
	Total         int64           `json:"total"`
	Page          int             `json:"page"`
	PageSize      int             `json:"page_size"`
	HasMore       bool            `json:"has_more"`
	UnreadCount   int64           `json:"unread_count"`
}
