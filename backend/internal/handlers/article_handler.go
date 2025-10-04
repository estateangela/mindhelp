package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/dto"
	"mindhelp-backend/internal/middleware"
	"mindhelp-backend/internal/models"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ArticleHandler 文章處理器
type ArticleHandler struct {
}

// NewArticleHandler 創建新的文章處理器
func NewArticleHandler() *ArticleHandler {
	return &ArticleHandler{}
}

// GetArticles 獲取文章列表
// @Summary 獲取文章列表
// @Description 搜尋和獲取專家文章列表
// @Tags article
// @Accept json
// @Produce json
// @Param q query string false "搜尋關鍵字"
// @Param tag query string false "依標籤篩選"
// @Param sort_by query string false "排序依據" Enums(publish_date,popularity) default(publish_date)
// @Param page query int false "頁碼" default(1)
// @Param limit query int false "每頁數量" default(10)
// @Success 200 {object} vo.Response{data=dto.ArticleListResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Router /articles [get]
func (h *ArticleHandler) GetArticles(c *gin.Context) {
	// 解析查詢參數
	query := strings.TrimSpace(c.Query("q"))
	tag := strings.TrimSpace(c.Query("tag"))
	sortBy := c.DefaultQuery("sort_by", "publish_date")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// 參數驗證
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	offset := (page - 1) * limit

	// 獲取資料庫連接
	db, err := database.GetDBSafely()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, vo.NewErrorResponse(
			"database_unavailable",
			"Database service is currently unavailable",
			"SERVICE_UNAVAILABLE",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 構建查詢
	dbQuery := db.Model(&models.Article{}).Where("is_published = ?", true)

	// 添加搜尋條件
	if query != "" {
		dbQuery = dbQuery.Where("title ILIKE ? OR content ILIKE ? OR author ILIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	// 添加標籤篩選
	if tag != "" {
		dbQuery = dbQuery.Where("tags ILIKE ?", "%\""+tag+"\"%")
	}

	// 添加排序
	switch sortBy {
	case "popularity":
		dbQuery = dbQuery.Order("view_count DESC, publish_date DESC")
	default:
		dbQuery = dbQuery.Order("publish_date DESC")
	}

	// 獲取總數
	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to count articles",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 獲取文章列表
	var articles []models.Article
	if err := dbQuery.Offset(offset).Limit(limit).Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get articles",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 獲取當前使用者ID (如果已登入)
	userID := middleware.GetUserID(c)

	// 轉換為 DTO
	var articleResponses []dto.ArticleResponse
	for _, article := range articles {
		// 解析標籤
		var tags []string
		if article.Tags != "" {
			json.Unmarshal([]byte(article.Tags), &tags)
		}

		// 檢查是否已收藏
		isBookmarked := false
		if userID != "" {
			var count int64
			db.Model(&models.Bookmark{}).Where(
				"user_id = ? AND resource_type = ? AND article_id = ?",
				userID, "article", article.ID).Count(&count)
			isBookmarked = count > 0
		}

		response := dto.ArticleResponse{
			ID:           article.ID.String(),
			Title:        article.Title,
			Author:       article.Author,
			AuthorTitle:  article.AuthorTitle,
			PublishDate:  article.PublishDate.Format("2006-01-02T15:04:05Z07:00"),
			Summary:      article.Summary,
			Tags:         tags,
			ImageURL:     article.ImageURL,
			IsBookmarked: isBookmarked,
			ViewCount:    article.ViewCount,
			CreatedAt:    article.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		articleResponses = append(articleResponses, response)
	}

	// 構建分頁回應
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	response := dto.ArticleListResponse{
		Articles:   articleResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasMore:    page < totalPages,
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "Articles retrieved successfully"))
}

// GetArticle 獲取單篇文章詳情
// @Summary 獲取文章詳情
// @Description 獲取特定文章的詳細內容
// @Tags article
// @Accept json
// @Produce json
// @Param id path string true "文章ID"
// @Success 200 {object} vo.Response{data=dto.ArticleResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Router /articles/{id} [get]
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	articleID := c.Param("id")
	if articleID == "" {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Article ID is required",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 驗證 UUID 格式
	parsedID, err := uuid.Parse(articleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid article ID format",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	var article models.Article
	if err := db.Where("id = ? AND is_published = ?", parsedID, true).First(&article).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Article not found",
				"NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get article",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 增加瀏覽次數
	db.Model(&article).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))

	// 解析標籤
	var tags []string
	if article.Tags != "" {
		json.Unmarshal([]byte(article.Tags), &tags)
	}

	// 檢查是否已收藏
	userID := middleware.GetUserID(c)
	isBookmarked := false
	if userID != "" {
		var count int64
		db.Model(&models.Bookmark{}).Where(
			"user_id = ? AND resource_type = ? AND article_id = ?",
			userID, "article", article.ID).Count(&count)
		isBookmarked = count > 0
	}

	// 構建回應
	response := dto.ArticleResponse{
		ID:           article.ID.String(),
		Title:        article.Title,
		Author:       article.Author,
		AuthorTitle:  article.AuthorTitle,
		PublishDate:  article.PublishDate.Format("2006-01-02T15:04:05Z07:00"),
		Summary:      article.Summary,
		Content:      article.Content, // 詳情包含完整內容
		Tags:         tags,
		ImageURL:     article.ImageURL,
		IsBookmarked: isBookmarked,
		ViewCount:    article.ViewCount + 1, // 包含本次瀏覽
		CreatedAt:    article.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "Article retrieved successfully"))
}

// BookmarkArticle 收藏文章
// @Summary 收藏文章
// @Description 將文章加入使用者收藏列表
// @Tags article
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "文章ID"
// @Success 204 "No Content"
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Failure 409 {object} vo.ErrorResponse
// @Router /articles/{id}/bookmark [post]
func (h *ArticleHandler) BookmarkArticle(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, vo.NewErrorResponse(
			"unauthorized",
			"User not authenticated",
			"UNAUTHORIZED",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	articleID := c.Param("id")
	parsedID, err := uuid.Parse(articleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid article ID format",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 檢查文章是否存在
	var article models.Article
	if err := db.Where("id = ? AND is_published = ?", parsedID, true).First(&article).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Article not found",
				"NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to check article",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 檢查是否已收藏
	var existingBookmark models.Bookmark
	if err := db.Where("user_id = ? AND resource_type = ? AND article_id = ?",
		userID, "article", parsedID).First(&existingBookmark).Error; err == nil {
		c.JSON(http.StatusConflict, vo.NewErrorResponse(
			"conflict",
			"Article already bookmarked",
			"ALREADY_BOOKMARKED",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 創建收藏
	bookmark := models.Bookmark{
		UserID:       uuid.MustParse(userID),
		ResourceType: "article",
		ArticleID:    &parsedID,
	}

	if err := db.Create(&bookmark).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to bookmark article",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.Status(http.StatusNoContent)
}

// UnbookmarkArticle 取消收藏文章
// @Summary 取消收藏文章
// @Description 從使用者收藏列表中移除文章
// @Tags article
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "文章ID"
// @Success 204 "No Content"
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Router /articles/{id}/bookmark [delete]
func (h *ArticleHandler) UnbookmarkArticle(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, vo.NewErrorResponse(
			"unauthorized",
			"User not authenticated",
			"UNAUTHORIZED",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	articleID := c.Param("id")
	parsedID, err := uuid.Parse(articleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid article ID format",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 查找並刪除收藏
	result := db.Where("user_id = ? AND resource_type = ? AND article_id = ?",
		userID, "article", parsedID).Delete(&models.Bookmark{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to unbookmark article",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, vo.NewErrorResponse(
			"not_found",
			"Bookmark not found",
			"NOT_FOUND",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.Status(http.StatusNoContent)
}
