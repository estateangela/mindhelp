package dto

import (
	"time"
)

// NotificationResponse 通知回應結構
type NotificationResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// NotificationListResponse 通知列表回應結構
type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	Total         int64                  `json:"total"`
	Page          int                    `json:"page"`
	PageSize      int                    `json:"page_size"`
	HasMore       bool                   `json:"has_more"`
}

// MarkNotificationReadRequest 標記通知為已讀請求
type MarkNotificationReadRequest struct {
	NotificationID string `json:"notification_id" binding:"required"`
}

// MarkAllNotificationsReadRequest 標記所有通知為已讀請求
type MarkAllNotificationsReadRequest struct {
	Type string `json:"type,omitempty"` // 可選：指定通知類型
}
