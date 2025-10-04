package dto

import (
	"fmt"
)

// NotificationResponse 通知回應結構
type NotificationResponse struct {
	ID        string                 `json:"id"`
	Title     string                 `json:"title"`
	Content   string                 `json:"content"`
	Type      string                 `json:"type"`
	IsRead    bool                   `json:"is_read"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
	CreatedAt string                 `json:"created_at"`
}

// NotificationListResponse 通知列表回應結構
type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	Total         int64                  `json:"total"`
	Page          int                    `json:"page"`
	Limit         int                    `json:"limit"`
	TotalPages    int                    `json:"total_pages"`
	HasMore       bool                   `json:"has_more"`
	UnreadCount   int64                  `json:"unread_count"`
}

// MarkAsReadRequest 標記通知為已讀請求
type MarkAsReadRequest struct {
	NotificationIDs []string `json:"notification_ids" binding:"required"`
}

// Validate 驗證請求
func (r *MarkAsReadRequest) Validate() error {
	if len(r.NotificationIDs) == 0 {
		return fmt.Errorf("notification_ids cannot be empty")
	}
	return nil
}

// NotificationSettingsRequest 通知設定請求
type NotificationSettingsRequest struct {
	NotifyNewArticle    *bool `json:"notify_new_article,omitempty"`
	NotifyPromotions    *bool `json:"notify_promotions,omitempty"`
	NotifySystemUpdates *bool `json:"notify_system_updates,omitempty"`
}

// Validate 驗證請求
func (r *NotificationSettingsRequest) Validate() error {
	return nil
}

// NotificationSettingsResponse 通知設定回應
type NotificationSettingsResponse struct {
	NotifyNewArticle    bool `json:"notify_new_article"`
	NotifyPromotions    bool `json:"notify_promotions"`
	NotifySystemUpdates bool `json:"notify_system_updates"`
}

// PushTokenRequest 推播 Token 請求
type PushTokenRequest struct {
	Token    string `json:"token" binding:"required"`
	Platform string `json:"platform" binding:"required,oneof=ios android web"`
}

// Validate 驗證請求
func (r *PushTokenRequest) Validate() error {
	if r.Token == "" {
		return fmt.Errorf("token cannot be empty")
	}
	if r.Platform != "ios" && r.Platform != "android" && r.Platform != "web" {
		return fmt.Errorf("platform must be ios, android, or web")
	}
	return nil
}
