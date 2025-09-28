package dto

import (
	"github.com/go-playground/validator/v10"
)

// ShareRequest 分享請求
type ShareRequest struct {
	ContentType string `json:"content_type" binding:"required,oneof=article location quiz" validate:"required,oneof=article location quiz"`
	ContentID   string `json:"content_id" binding:"required,uuid" validate:"required,uuid"`
	Platform    string `json:"platform" binding:"omitempty,oneof=facebook twitter line whatsapp email copy" validate:"omitempty,oneof=facebook twitter line whatsapp email copy"`
	Message     string `json:"message" binding:"omitempty,max=500" validate:"omitempty,max=500"`
}

// ShareResponse 分享回應
type ShareResponse struct {
	ShareID     string                 `json:"share_id"`
	ContentType string                 `json:"content_type"`
	ContentID   string                 `json:"content_id"`
	ShareURL    string                 `json:"share_url"`
	ShortURL    string                 `json:"short_url,omitempty"`
	Platform    string                 `json:"platform,omitempty"`
	QRCode      string                 `json:"qr_code,omitempty"` // Base64 encoded QR code image
	ExpiresAt   string                 `json:"expires_at,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   string                 `json:"created_at"`
}

// ShareContentResponse 分享內容詳情
type ShareContentResponse struct {
	ContentType string                 `json:"content_type"`
	ContentID   string                 `json:"content_id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	ImageURL    string                 `json:"image_url,omitempty"`
	URL         string                 `json:"url"`
	Author      string                 `json:"author,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	IsValid     bool                   `json:"is_valid"`
	ViewCount   int64                  `json:"view_count"`
	ShareCount  int64                  `json:"share_count"`
	CreatedAt   string                 `json:"created_at"`
}

// ShareStatsResponse 分享統計回應
type ShareStatsResponse struct {
	ContentID    string           `json:"content_id"`
	ContentType  string           `json:"content_type"`
	TotalShares  int64            `json:"total_shares"`
	Platforms    map[string]int64 `json:"platforms"` // platform -> count
	RecentShares []ShareResponse  `json:"recent_shares,omitempty"`
}

// SocialShareMetadata 社交媒體分享元數據
type SocialShareMetadata struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ImageURL    string   `json:"image_url"`
	URL         string   `json:"url"`
	SiteName    string   `json:"site_name"`
	Author      string   `json:"author,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// ShareListRequest 分享列表請求
type ShareListRequest struct {
	ContentType string `json:"content_type" form:"content_type" binding:"omitempty,oneof=article location quiz" validate:"omitempty,oneof=article location quiz"`
	Platform    string `json:"platform" form:"platform" binding:"omitempty" validate:"omitempty"`
	Page        int    `json:"page" form:"page" binding:"omitempty,min=1" validate:"omitempty,min=1"`
	Limit       int    `json:"limit" form:"limit" binding:"omitempty,min=1,max=50" validate:"omitempty,min=1,max=50"`
}

// ShareListResponse 分享列表回應
type ShareListResponse struct {
	Shares     []ShareResponse `json:"shares"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
	HasMore    bool            `json:"has_more"`
}

// Validate 驗證請求資料
func (r *ShareRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *ShareListRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
