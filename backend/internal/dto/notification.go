package dto

import (
	"github.com/go-playground/validator/v10"
)

// NotificationResponse 通知回應
type NotificationResponse struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Body      string                 `json:"body"`
	IsRead    bool                   `json:"is_read"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
	CreatedAt string                 `json:"created_at"`
}

// NotificationListResponse 通知列表回應
type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	Total         int64                  `json:"total"`
	UnreadCount   int64                  `json:"unread_count"`
	Page          int                    `json:"page"`
	Limit         int                    `json:"limit"`
	TotalPages    int                    `json:"total_pages"`
	HasMore       bool                   `json:"has_more"`
}

// MarkAsReadRequest 標記已讀請求
type MarkAsReadRequest struct {
	NotificationIDs []string `json:"notification_ids" binding:"required,min=1" validate:"required,min=1"`
}

// NotificationSettingsRequest 通知設定請求
type NotificationSettingsRequest struct {
	NotifyNewArticle    *bool `json:"notify_new_article" binding:"omitempty" validate:"omitempty"`
	NotifyPromotions    *bool `json:"notify_promotions" binding:"omitempty" validate:"omitempty"`
	NotifySystemUpdates *bool `json:"notify_system_updates" binding:"omitempty" validate:"omitempty"`
}

// NotificationSettingsResponse 通知設定回應
type NotificationSettingsResponse struct {
	NotifyNewArticle    bool `json:"notify_new_article"`
	NotifyPromotions    bool `json:"notify_promotions"`
	NotifySystemUpdates bool `json:"notify_system_updates"`
}

// PushTokenRequest 推播 Token 請求
type PushTokenRequest struct {
	Token    string `json:"token" binding:"required" validate:"required"`
	Platform string `json:"platform" binding:"required,oneof=ios android" validate:"required,oneof=ios android"`
}

// Validate 驗證請求資料
func (r *MarkAsReadRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *NotificationSettingsRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *PushTokenRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
