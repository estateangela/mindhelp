package dto

import (
	"github.com/go-playground/validator/v10"
)

// ArticleRequest 文章請求 (管理員用)
type ArticleRequest struct {
	Title       string    `json:"title" binding:"required,max=200" validate:"required,max=200"`
	Author      string    `json:"author" binding:"required,max=100" validate:"required,max=100"`
	AuthorTitle string    `json:"author_title" binding:"omitempty,max=100" validate:"omitempty,max=100"`
	Summary     string    `json:"summary" binding:"omitempty,max=500" validate:"omitempty,max=500"`
	Content     string    `json:"content" binding:"required" validate:"required"`
	Tags        []string  `json:"tags" binding:"omitempty" validate:"omitempty"`
	ImageURL    string    `json:"image_url" binding:"omitempty,url" validate:"omitempty,url"`
	IsPublished bool      `json:"is_published"`
}

// ArticleSearchRequest 文章搜尋請求
type ArticleSearchRequest struct {
	Query   string `json:"query" form:"q" binding:"omitempty" validate:"omitempty"`
	Tag     string `json:"tag" form:"tag" binding:"omitempty" validate:"omitempty"`
	SortBy  string `json:"sort_by" form:"sort_by" binding:"omitempty,oneof=publish_date popularity" validate:"omitempty,oneof=publish_date popularity"`
	Page    int    `json:"page" form:"page" binding:"omitempty,min=1" validate:"omitempty,min=1"`
	Limit   int    `json:"limit" form:"limit" binding:"omitempty,min=1,max=50" validate:"omitempty,min=1,max=50"`
}

// ArticleResponse 文章回應
type ArticleResponse struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Author       string    `json:"author"`
	AuthorTitle  string    `json:"author_title,omitempty"`
	PublishDate  string    `json:"publish_date"`
	Summary      string    `json:"summary,omitempty"`
	Content      string    `json:"content,omitempty"` // 列表中不包含，詳情中包含
	Tags         []string  `json:"tags,omitempty"`
	ImageURL     string    `json:"image_url,omitempty"`
	IsBookmarked bool      `json:"is_bookmarked"`
	ViewCount    int64     `json:"view_count"`
	CreatedAt    string    `json:"created_at"`
}

// ArticleListResponse 文章列表回應
type ArticleListResponse struct {
	Articles   []ArticleResponse `json:"articles"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
	HasMore    bool              `json:"has_more"`
}

// ArticleAuthor 文章作者資訊
type ArticleAuthor struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

// Validate 驗證請求資料
func (r *ArticleRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *ArticleSearchRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
