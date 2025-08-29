package dto

import "github.com/go-playground/validator/v10"

// LocationRequest 位置請求
type LocationRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=100" validate:"required,min=1,max=100"`
	Description string  `json:"description" binding:"omitempty,max=1000" validate:"omitempty,max=1000"`
	Address     string  `json:"address" binding:"omitempty,max=255" validate:"omitempty,max=255"`
	Latitude    float64 `json:"latitude" binding:"required,min=-90,max=90" validate:"required,min=-90,max=90"`
	Longitude   float64 `json:"longitude" binding:"required,min=-180,max=180" validate:"required,min=-180,max=180"`
	Category    string  `json:"category" binding:"omitempty,max=50" validate:"omitempty,max=50"`
	Phone       string  `json:"phone" binding:"omitempty,len=10" validate:"omitempty,len=10"`
	Website     string  `json:"website" binding:"omitempty,url" validate:"omitempty,url"`
	Rating      float64 `json:"rating" binding:"omitempty,min=0,max=5" validate:"omitempty,min=0,max=5"`
	IsPublic    bool    `json:"is_public"`
}

// LocationResponse 位置回應
type LocationResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Address     string  `json:"address,omitempty"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Category    string  `json:"category,omitempty"`
	Phone       string  `json:"phone,omitempty"`
	Website     string  `json:"website,omitempty"`
	Rating      float64 `json:"rating"`
	IsPublic    bool    `json:"is_public"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// LocationSearchRequest 位置搜尋請求
type LocationSearchRequest struct {
	Query     string  `json:"query" binding:"omitempty" validate:"omitempty"`
	Latitude  float64 `json:"latitude" binding:"omitempty,min=-90,max=90" validate:"omitempty,min=-90,max=90"`
	Longitude float64 `json:"longitude" binding:"omitempty,min=-180,max=180" validate:"omitempty,min=-180,max=180"`
	Radius    float64 `json:"radius" binding:"omitempty,min=0.1,max=100" validate:"omitempty,min=0.1,max=100"` // 公里
	Category  string  `json:"category" binding:"omitempty" validate:"omitempty"`
	Page      int     `json:"page" binding:"omitempty,min=1" validate:"omitempty,min=1"`
	PageSize  int     `json:"page_size" binding:"omitempty,min=1,max=100" validate:"omitempty,min=1,max=100"`
}

// LocationSearchResponse 位置搜尋回應
type LocationSearchResponse struct {
	Locations []LocationResponse `json:"locations"`
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
	HasMore   bool               `json:"has_more"`
}

// LocationUpdateRequest 位置更新請求
type LocationUpdateRequest struct {
	Name        *string  `json:"name,omitempty" binding:"omitempty,min=1,max=100" validate:"omitempty,min=1,max=100"`
	Description *string  `json:"description,omitempty" binding:"omitempty,max=1000" validate:"omitempty,max=1000"`
	Address     *string  `json:"address,omitempty" binding:"omitempty,max=255" validate:"omitempty,max=255"`
	Latitude    *float64 `json:"latitude,omitempty" binding:"omitempty,min=-90,max=90" validate:"omitempty,min=-90,max=90"`
	Longitude   *float64 `json:"longitude,omitempty" binding:"omitempty,min=-180,max=180" validate:"omitempty,min=-180,max=180"`
	Category    *string  `json:"category,omitempty" binding:"omitempty,max=50" validate:"omitempty,max=50"`
	Phone       *string  `json:"phone,omitempty" binding:"omitempty,len=10" validate:"omitempty,len=10"`
	Website     *string  `json:"website,omitempty" binding:"omitempty,url" validate:"omitempty,url"`
	Rating      *float64 `json:"rating,omitempty" binding:"omitempty,min=0,max=5" validate:"omitempty,min=0,max=5"`
	IsPublic   *bool    `json:"is_public,omitempty"`
}

// Validate 驗證請求資料
func (r *LocationRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate 驗證請求資料
func (r *LocationSearchRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate 驗證請求資料
func (r *LocationUpdateRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
