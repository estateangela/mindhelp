package dto

import (
	"github.com/go-playground/validator/v10"
)

// QuizRequest 測驗請求 (管理員用)
type QuizRequest struct {
	Title       string              `json:"title" binding:"required,max=200" validate:"required,max=200"`
	Description string              `json:"description" binding:"omitempty" validate:"omitempty"`
	Category    string              `json:"category" binding:"required,max=50" validate:"required,max=50"`
	Questions   []QuizQuestionRequest `json:"questions" binding:"required,min=1" validate:"required,min=1"`
}

// QuizQuestionRequest 測驗題目請求
type QuizQuestionRequest struct {
	Question string   `json:"question" binding:"required" validate:"required"`
	Options  []string `json:"options" binding:"required,min=2" validate:"required,min=2"`
	OrderNum int      `json:"order_num" binding:"required,min=1" validate:"required,min=1"`
}

// QuizSubmissionRequest 測驗提交請求
type QuizSubmissionRequest struct {
	QuizID  string             `json:"quiz_id" binding:"required,uuid" validate:"required,uuid"`
	Answers map[string]int     `json:"answers" binding:"required" validate:"required"` // question_id -> option_index
}

// QuizListRequest 測驗列表請求
type QuizListRequest struct {
	Category string `json:"category" form:"category" binding:"omitempty" validate:"omitempty"`
	Page     int    `json:"page" form:"page" binding:"omitempty,min=1" validate:"omitempty,min=1"`
	Limit    int    `json:"limit" form:"limit" binding:"omitempty,min=1,max=50" validate:"omitempty,min=1,max=50"`
}

// QuizResponse 測驗回應
type QuizResponse struct {
	ID          string              `json:"id"`
	Title       string              `json:"title"`
	Description string              `json:"description,omitempty"`
	Category    string              `json:"category"`
	Questions   []QuizQuestionResponse `json:"questions,omitempty"` // 列表中不包含，詳情中包含
	IsActive    bool                `json:"is_active"`
	CreatedAt   string              `json:"created_at"`
}

// QuizQuestionResponse 測驗題目回應
type QuizQuestionResponse struct {
	ID       string   `json:"id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
	OrderNum int      `json:"order_num"`
}

// QuizSubmissionResponse 測驗提交回應
type QuizSubmissionResponse struct {
	ID          string `json:"id"`
	QuizTitle   string `json:"quiz_title"`
	Score       int    `json:"score"`
	Result      string `json:"result"`
	CompletedAt string `json:"completed_at"`
	CreatedAt   string `json:"created_at"`
}

// QuizListResponse 測驗列表回應
type QuizListResponse struct {
	Quizzes    []QuizResponse `json:"quizzes"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
	HasMore    bool           `json:"has_more"`
}

// QuizHistoryResponse 測驗歷史回應
type QuizHistoryResponse struct {
	Submissions []QuizSubmissionResponse `json:"submissions"`
	Total       int64                    `json:"total"`
	Page        int                      `json:"page"`
	Limit       int                      `json:"limit"`
	TotalPages  int                      `json:"total_pages"`
	HasMore     bool                     `json:"has_more"`
}

// Validate 驗證請求資料
func (r *QuizRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *QuizSubmissionRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *QuizListRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
