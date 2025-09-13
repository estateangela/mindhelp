package dto

import (
	"github.com/go-playground/validator/v10"
)

// UpdateUserRequest 更新使用者資料請求
type UpdateUserRequest struct {
	Username *string `json:"username" binding:"omitempty,min=3,max=50" validate:"omitempty,min=3,max=50"`
	FullName *string `json:"full_name" binding:"omitempty,min=2,max=100" validate:"omitempty,min=2,max=100"`
	Phone    *string `json:"phone" binding:"omitempty,len=10" validate:"omitempty,len=10"`
	Avatar   *string `json:"avatar" binding:"omitempty,url" validate:"omitempty,url"`
}

// ChangePasswordRequest 變更密碼請求 (從 auth.go 移動到這裡)
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required" validate:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8" validate:"required,min=8"`
}

// UserProfileResponse 使用者檔案回應
type UserProfileResponse struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Username  string  `json:"username"`
	FullName  string  `json:"full_name"`
	Phone     string  `json:"phone,omitempty"`
	Avatar    string  `json:"avatar,omitempty"`
	IsActive  bool    `json:"is_active"`
	LastLogin *string `json:"last_login,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// UserStatsResponse 使用者統計回應
type UserStatsResponse struct {
	TotalChatSessions   int64 `json:"total_chat_sessions"`
	TotalQuizzes        int64 `json:"total_quizzes"`
	TotalBookmarks      int64 `json:"total_bookmarks"`
	TotalReviews        int64 `json:"total_reviews"`
	UnreadNotifications int64 `json:"unread_notifications"`
}

// DeleteAccountRequest 刪除帳號請求
type DeleteAccountRequest struct {
	Password string `json:"password" binding:"required" validate:"required"`
	Reason   string `json:"reason" binding:"omitempty,max=500" validate:"omitempty,max=500"`
}

// Validate 驗證請求資料
func (r *UpdateUserRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *ChangePasswordRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *DeleteAccountRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
