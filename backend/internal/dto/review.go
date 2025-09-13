package dto

import (
	"github.com/go-playground/validator/v10"
)

// ReviewRequest 評論請求
type ReviewRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5" validate:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"omitempty,max=1000" validate:"omitempty,max=1000"`
}

// ReviewUpdateRequest 評論更新請求
type ReviewUpdateRequest struct {
	Rating  *int    `json:"rating" binding:"omitempty,min=1,max=5" validate:"omitempty,min=1,max=5"`
	Comment *string `json:"comment" binding:"omitempty,max=1000" validate:"omitempty,max=1000"`
}

// ReviewResponse 評論回應
type ReviewResponse struct {
	ID         string       `json:"id"`
	Author     ReviewAuthor `json:"author"`
	ResourceID string       `json:"resource_id"`
	Rating     int          `json:"rating"`
	Comment    string       `json:"comment,omitempty"`
	IsHelpful  int          `json:"is_helpful"`
	CanEdit    bool         `json:"can_edit"` // 當前使用者是否可編輯/刪除
	CreatedAt  string       `json:"created_at"`
	UpdatedAt  string       `json:"updated_at"`
}

// ReviewAuthor 評論作者資訊 (簡化的使用者資訊)
type ReviewAuthor struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar,omitempty"`
}

// ReviewListResponse 評論列表回應
type ReviewListResponse struct {
	Reviews    []ReviewResponse `json:"reviews"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
	HasMore    bool             `json:"has_more"`
	Statistics ReviewStatistics `json:"statistics"`
}

// ReviewStatistics 評論統計資訊
type ReviewStatistics struct {
	AverageRating float64          `json:"average_rating"`
	TotalReviews  int64            `json:"total_reviews"`
	RatingDistribution map[int]int64 `json:"rating_distribution"` // rating -> count
}

// ReportRequest 回報不當內容請求
type ReportRequest struct {
	ContentType string `json:"content_type" binding:"required,oneof=review article resource" validate:"required,oneof=review article resource"`
	ContentID   string `json:"content_id" binding:"required,uuid" validate:"required,uuid"`
	Reason      string `json:"reason" binding:"required,oneof=spam inappropriate incorrect_info" validate:"required,oneof=spam inappropriate incorrect_info"`
	Details     string `json:"details" binding:"omitempty,max=1000" validate:"omitempty,max=1000"`
}

// Validate 驗證請求資料
func (r *ReviewRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *ReviewUpdateRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *ReportRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
