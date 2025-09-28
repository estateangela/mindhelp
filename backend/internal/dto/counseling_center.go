package dto

import "time"

// CounselingCenterRequest 諮商所請求結構
type CounselingCenterRequest struct {
	Name            string `json:"name" validate:"required"`
	Address         string `json:"address"`
	Phone           string `json:"phone"`
	OnlineCounseling bool  `json:"online_counseling"`
}

// CounselingCenterResponse 諮商所回應結構
type CounselingCenterResponse struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Address         string    `json:"address"`
	Phone           string    `json:"phone"`
	OnlineCounseling bool     `json:"online_counseling"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CounselingCenterListResponse 諮商所列表回應結構
type CounselingCenterListResponse struct {
	CounselingCenters []CounselingCenterResponse `json:"counseling_centers"`
	Total             int64                      `json:"total"`
	Page              int                        `json:"page"`
	PageSize          int                        `json:"page_size"`
}
