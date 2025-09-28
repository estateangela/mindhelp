package vo

import "time"

// Response 通用 API 回應結構
type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Path      string      `json:"path,omitempty"`
}

// PaginationResponse 分頁回應結構
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
	HasMore    bool        `json:"has_more"`
	HasPrev    bool        `json:"has_prev"`
}

// ErrorResponse 錯誤回應結構
type ErrorResponse struct {
	Success   bool      `json:"success"`
	Error     string    `json:"error"`
	Message   string    `json:"message"`
	Code      string    `json:"code,omitempty"`
	Details   []string  `json:"details,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path,omitempty"`
}

// SuccessResponse 成功回應
func SuccessResponse(data interface{}, message string) Response {
	return Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// ErrorResponse 錯誤回應
func NewErrorResponse(error, message, code string, details []string, path string) ErrorResponse {
	return ErrorResponse{
		Success:   false,
		Error:     error,
		Message:   message,
		Code:      code,
		Details:   details,
		Timestamp: time.Now(),
		Path:      path,
	}
}

// NewResponse 創建成功回應
func NewResponse(message string, data interface{}) Response {
	return Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// PaginationResponse 分頁回應
func NewPaginationResponse(data interface{}, total int64, page, pageSize int) PaginationResponse {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	
	return PaginationResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		HasMore:    page < totalPages,
		HasPrev:    page > 1,
	}
}
