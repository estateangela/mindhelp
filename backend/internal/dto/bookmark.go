package dto

import (
	"github.com/go-playground/validator/v10"
)

// BookmarkRequest 收藏請求
type BookmarkRequest struct {
	ResourceType string `json:"resource_type" binding:"required,oneof=article location" validate:"required,oneof=article location"`
	ResourceID   string `json:"resource_id" binding:"required,uuid" validate:"required,uuid"`
}

// BookmarkResponse 收藏回應
type BookmarkResponse struct {
	ID           string                  `json:"id"`
	ResourceType string                  `json:"resource_type"`
	Resource     interface{}             `json:"resource"` // ArticleResponse 或 LocationResponse
	CreatedAt    string                  `json:"created_at"`
}

// BookmarkListResponse 收藏列表回應
type BookmarkListResponse struct {
	Bookmarks  []BookmarkResponse `json:"bookmarks"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
	HasMore    bool               `json:"has_more"`
}

// ArticleBookmarkResponse 文章收藏回應
type ArticleBookmarkResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Summary   string `json:"summary,omitempty"`
	ImageURL  string `json:"image_url,omitempty"`
	CreatedAt string `json:"created_at"`
}

// LocationBookmarkResponse 地點收藏回應
type LocationBookmarkResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	Category    string  `json:"category"`
	Rating      float64 `json:"rating"`
	Phone       string  `json:"phone,omitempty"`
	Website     string  `json:"website,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

// Validate 驗證請求資料
func (r *BookmarkRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
