package dto

import "time"

// RecommendedDoctorRequest 推薦醫師請求結構
type RecommendedDoctorRequest struct {
	Name            string `json:"name" validate:"required"`
	Description     string `json:"description"`
	ExperienceCount int    `json:"experience_count"`
}

// RecommendedDoctorResponse 推薦醫師回應結構
type RecommendedDoctorResponse struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ExperienceCount int       `json:"experience_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// RecommendedDoctorListResponse 推薦醫師列表回應結構
type RecommendedDoctorListResponse struct {
	RecommendedDoctors []RecommendedDoctorResponse `json:"recommended_doctors"`
	Total              int64                       `json:"total"`
	Page               int                         `json:"page"`
	PageSize           int                         `json:"page_size"`
}
