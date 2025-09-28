package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	//"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/dto"
	"mindhelp-backend/internal/middleware"
	"mindhelp-backend/internal/models"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BookmarkHandler 收藏處理器
type BookmarkHandler struct {
}

// NewBookmarkHandler 創建新的收藏處理器
func NewBookmarkHandler() *BookmarkHandler {
	return &BookmarkHandler{}
}

// GetArticleBookmarks 獲取文章收藏列表
// @Summary 獲取文章收藏
// @Description 獲取使用者收藏的文章列表
// @Tags bookmark
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "頁碼" default(1)
// @Param limit query int false "每頁數量" default(10)
// @Success 200 {object} vo.Response{data=dto.BookmarkListResponse}
// @Failure 401 {object} vo.ErrorResponse
// @Router /users/me/bookmarks/articles [get]
func (h *BookmarkHandler) GetArticleBookmarks(c *gin.Context) {
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

	// 解析查詢參數
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

	// 獲取總數
	var total int64
	if err := db.Model(&models.Bookmark{}).
		Where("user_id = ? AND resource_type = ?", userID, "article").
		Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to count bookmarks",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 獲取收藏列表
	var bookmarks []models.Bookmark
	if err := db.Preload("Article").
		Where("user_id = ? AND resource_type = ?", userID, "article").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&bookmarks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get bookmarks",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 轉換為 DTO
	var bookmarkResponses []dto.BookmarkResponse
	for _, bookmark := range bookmarks {
		if bookmark.Article != nil {
			// 解析文章標籤
			var tags []string
			if bookmark.Article.Tags != "" {
				json.Unmarshal([]byte(bookmark.Article.Tags), &tags)
			}

			articleResponse := dto.ArticleResponse{
				ID:           bookmark.Article.ID.String(),
				Title:        bookmark.Article.Title,
				Author:       bookmark.Article.Author,
				AuthorTitle:  bookmark.Article.AuthorTitle,
				PublishDate:  bookmark.Article.PublishDate.Format("2006-01-02T15:04:05Z07:00"),
				Summary:      bookmark.Article.Summary,
				Tags:         tags,
				ImageURL:     bookmark.Article.ImageURL,
				IsBookmarked: true, // 當然是已收藏
				ViewCount:    bookmark.Article.ViewCount,
				CreatedAt:    bookmark.Article.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}

			response := dto.BookmarkResponse{
				ID:           bookmark.ID.String(),
				ResourceType: bookmark.ResourceType,
				Resource:     articleResponse,
				CreatedAt:    bookmark.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
			bookmarkResponses = append(bookmarkResponses, response)
		}
	}

	// 構建分頁回應
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	response := dto.BookmarkListResponse{
		Bookmarks:  bookmarkResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasMore:    page < totalPages,
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "Article bookmarks retrieved successfully"))
}

// GetLocationBookmarks 獲取地點收藏列表
// @Summary 獲取地點收藏
// @Description 獲取使用者收藏的地點列表
// @Tags bookmark
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "頁碼" default(1)
// @Param limit query int false "每頁數量" default(10)
// @Success 200 {object} vo.Response{data=dto.BookmarkListResponse}
// @Failure 401 {object} vo.ErrorResponse
// @Router /users/me/bookmarks/resources [get]
func (h *BookmarkHandler) GetLocationBookmarks(c *gin.Context) {
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

	// 解析查詢參數
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

	// 獲取總數
	var total int64
	if err := db.Model(&models.Bookmark{}).
		Where("user_id = ? AND resource_type = ?", userID, "location").
		Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to count bookmarks",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 獲取收藏列表
	var bookmarks []models.Bookmark
	if err := db.Preload("Location").
		Where("user_id = ? AND resource_type = ?", userID, "location").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&bookmarks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get bookmarks",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 轉換為 DTO
	var bookmarkResponses []dto.BookmarkResponse
	for _, bookmark := range bookmarks {
		if bookmark.Location != nil {
			locationResponse := dto.LocationResponse{
				ID:          bookmark.Location.ID.String(),
				UserID:      bookmark.Location.UserID.String(),
				Name:        bookmark.Location.Name,
				Description: bookmark.Location.Description,
				Address:     bookmark.Location.Address,
				Latitude:    bookmark.Location.Latitude,
				Longitude:   bookmark.Location.Longitude,
				Category:    bookmark.Location.Category,
				Phone:       bookmark.Location.Phone,
				Website:     bookmark.Location.Website,
				Rating:      bookmark.Location.Rating,
				IsPublic:    bookmark.Location.IsPublic,
				CreatedAt:   bookmark.Location.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt:   bookmark.Location.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}

			response := dto.BookmarkResponse{
				ID:           bookmark.ID.String(),
				ResourceType: bookmark.ResourceType,
				Resource:     locationResponse,
				CreatedAt:    bookmark.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
			bookmarkResponses = append(bookmarkResponses, response)
		}
	}

	// 構建分頁回應
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	response := dto.BookmarkListResponse{
		Bookmarks:  bookmarkResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasMore:    page < totalPages,
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "Location bookmarks retrieved successfully"))
}

// BookmarkResource 收藏資源 (通用端點)
// @Summary 收藏資源
// @Description 收藏文章或地點資源
// @Tags bookmark
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.BookmarkRequest true "收藏請求"
// @Success 204 "No Content"
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Failure 409 {object} vo.ErrorResponse
// @Router /bookmarks [post]
func (h *BookmarkHandler) BookmarkResource(c *gin.Context) {
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

	var req dto.BookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid request data",
			"VALIDATION_ERROR",
			[]string{err.Error()},
			c.Request.URL.Path,
		))
		return
	}

	// 驗證請求資料
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Validation failed",
			"VALIDATION_ERROR",
			[]string{err.Error()},
			c.Request.URL.Path,
		))
		return
	}

	resourceID, err := uuid.Parse(req.ResourceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid resource ID format",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 檢查資源是否存在
	if err := h.checkResourceExists(req.ResourceType, resourceID); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Resource not found",
				"NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to check resource",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 檢查是否已收藏
	var existingBookmark models.Bookmark
	query := db.Where("user_id = ? AND resource_type = ?", userID, req.ResourceType)
	if req.ResourceType == "article" {
		query = query.Where("article_id = ?", resourceID)
	} else {
		query = query.Where("location_id = ?", resourceID)
	}

	if err := query.First(&existingBookmark).Error; err == nil {
		c.JSON(http.StatusConflict, vo.NewErrorResponse(
			"conflict",
			"Resource already bookmarked",
			"ALREADY_BOOKMARKED",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 創建收藏
	bookmark := models.Bookmark{
		UserID:       uuid.MustParse(userID),
		ResourceType: req.ResourceType,
	}

	if req.ResourceType == "article" {
		bookmark.ArticleID = &resourceID
	} else {
		bookmark.LocationID = &resourceID
	}

	if err := db.Create(&bookmark).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to bookmark resource",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.Status(http.StatusNoContent)
}

// UnbookmarkResource 取消收藏資源
// @Summary 取消收藏資源
// @Description 取消收藏文章或地點資源
// @Tags bookmark
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resource_type query string true "資源類型" Enums(article,location)
// @Param resource_id query string true "資源ID"
// @Success 204 "No Content"
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Router /bookmarks [delete]
func (h *BookmarkHandler) UnbookmarkResource(c *gin.Context) {
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

	resourceType := c.Query("resource_type")
	resourceIDStr := c.Query("resource_id")

	if resourceType == "" || resourceIDStr == "" {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Resource type and resource ID are required",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	if resourceType != "article" && resourceType != "location" {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid resource type",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	resourceID, err := uuid.Parse(resourceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid resource ID format",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 查找並刪除收藏
	query := db.Where("user_id = ? AND resource_type = ?", userID, resourceType)
	if resourceType == "article" {
		query = query.Where("article_id = ?", resourceID)
	} else {
		query = query.Where("location_id = ?", resourceID)
	}

	result := query.Delete(&models.Bookmark{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to unbookmark resource",
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

// checkResourceExists 檢查資源是否存在
func (h *BookmarkHandler) checkResourceExists(resourceType string, resourceID uuid.UUID) error {
	switch resourceType {
	case "article":
		var article models.Article
		return db.Where("id = ? AND is_published = ?", resourceID, true).First(&article).Error
	case "location":
		var location models.Location
		return db.Where("id = ? AND is_public = ?", resourceID, true).First(&location).Error
	default:
		return gorm.ErrRecordNotFound
	}
}
