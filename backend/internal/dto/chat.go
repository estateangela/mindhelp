package dto

import "github.com/go-playground/validator/v10"

// ChatMessageRequest 聊天訊息請求
type ChatMessageRequest struct {
	SessionID *string `json:"session_id" binding:"omitempty,uuid" validate:"omitempty,uuid"` // 可選，不提供則創建新 session
	Content   string  `json:"content" binding:"required,min=1,max=2000" validate:"required,min=1,max=2000"`
	Model     string  `json:"model" binding:"omitempty" validate:"omitempty"`
}

// ChatMessageResponse 聊天訊息回應
type ChatMessageResponse struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	SessionID *string `json:"session_id"`
	Role      string  `json:"role"`
	Content   string  `json:"content"`
	Timestamp int64   `json:"timestamp"`
	Model     string  `json:"model,omitempty"`
	Tokens    int     `json:"tokens,omitempty"`
	CreatedAt string  `json:"created_at"`
}

// ChatHistoryRequest 聊天歷史請求
type ChatHistoryRequest struct {
	Page     int `json:"page" binding:"omitempty,min=1" validate:"omitempty,min=1"`
	PageSize int `json:"page_size" binding:"omitempty,min=1,max=100" validate:"omitempty,min=1,max=100"`
}

// ChatHistoryResponse 聊天歷史回應
type ChatHistoryResponse struct {
	Messages  []ChatMessageResponse `json:"messages"`
	Total     int64                 `json:"total"`
	Page      int                   `json:"page"`
	PageSize  int                   `json:"page_size"`
	HasMore   bool                  `json:"has_more"`
}

// OpenRouterRequest OpenRouter API 請求
type OpenRouterRequest struct {
	Model       string    `json:"model" binding:"required" validate:"required"`
	Messages    []Message `json:"messages" binding:"required,min=1" validate:"required,min=1"`
	Temperature float64   `json:"temperature" binding:"omitempty,min=0,max=2" validate:"omitempty,min=0,max=2"`
	MaxTokens   int       `json:"max_tokens" binding:"omitempty,min=1,max=4000" validate:"omitempty,min=1,max=4000"`
}

// Message 聊天訊息結構
type Message struct {
	Role    string `json:"role" binding:"required,oneof=user assistant system" validate:"required,oneof=user assistant system"`
	Content string `json:"content" binding:"required" validate:"required"`
}

// OpenRouterResponse OpenRouter API 回應
type OpenRouterResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// ChatSessionRequest 聊天會話請求
type ChatSessionRequest struct {
	Title string `json:"title" binding:"omitempty,max=200" validate:"omitempty,max=200"`
}

// ChatSessionResponse 聊天會話回應
type ChatSessionResponse struct {
	ID                  string `json:"id"`
	UserID              string `json:"user_id"`
	Title               string `json:"title,omitempty"`
	FirstMessageSnippet string `json:"first_message_snippet,omitempty"`
	LastUpdatedAt       string `json:"last_updated_at"`
	MessageCount        int    `json:"message_count"`
	IsActive            bool   `json:"is_active"`
	CreatedAt           string `json:"created_at"`
}

// ChatSessionListResponse 聊天會話列表回應
type ChatSessionListResponse struct {
	Sessions   []ChatSessionResponse `json:"sessions"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int                   `json:"total_pages"`
	HasMore    bool                  `json:"has_more"`
}

// SessionMessagesResponse 會話訊息回應
type SessionMessagesResponse struct {
	SessionID string                `json:"session_id"`
	Messages  []ChatMessageResponse `json:"messages"`
	Total     int64                 `json:"total"`
	Page      int                   `json:"page"`
	Limit     int                   `json:"limit"`
	HasMore   bool                  `json:"has_more"`
}

// Validate 驗證請求資料
func (r *ChatMessageRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *ChatHistoryRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *OpenRouterRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *ChatSessionRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
