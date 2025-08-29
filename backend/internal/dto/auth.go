package dto

import "github.com/go-playground/validator/v10"

// LoginRequest 登入請求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required,min=6" validate:"required,min=6"`
}

// RegisterRequest 註冊請求
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email" validate:"required,email"`
	Password  string `json:"password" binding:"required,min=6" validate:"required,min=6"`
	Username  string `json:"username" binding:"required,min=3,max=50" validate:"required,min=3,max=50"`
	FullName  string `json:"full_name" binding:"required,min=2,max=100" validate:"required,min=2,max=100"`
	Phone     string `json:"phone" binding:"omitempty,len=10" validate:"omitempty,len=10"`
	Avatar    string `json:"avatar" binding:"omitempty,url" validate:"omitempty,url"`
}

// ChangePasswordRequest 變更密碼請求
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required" validate:"required"`
	NewPassword    string `json:"new_password" binding:"required,min=6" validate:"required,min=6"`
}

// ForgotPasswordRequest 忘記密碼請求
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" validate:"required,email"`
}

// ResetPasswordRequest 重設密碼請求
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required" validate:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6" validate:"required,min=6"`
}

// AuthResponse 認證回應
type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
}

// UserResponse 使用者回應
type UserResponse struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Username  string  `json:"username"`
	FullName  string  `json:"full_name"`
	Phone     string  `json:"phone,omitempty"`
	Avatar    string  `json:"avatar,omitempty"`
	IsActive  bool    `json:"is_active"`
	LastLogin *string `json:"last_login,omitempty"`
	CreatedAt string  `json:"created_at"`
}

// Validate 驗證請求資料
func (r *LoginRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate 驗證請求資料
func (r *RegisterRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate 驗證請求資料
func (r *ChangePasswordRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate 驗證請求資料
func (r *ForgotPasswordRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate 驗證請求資料
func (r *ResetPasswordRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
