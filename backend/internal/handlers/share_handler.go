package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"mindhelp-backend/internal/config"
	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/dto"
	"mindhelp-backend/internal/middleware"
	"mindhelp-backend/internal/models"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ShareHandler 分享處理器
type ShareHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewShareHandler 創建新的分享處理器
func NewShareHandler(cfg *config.Config) *ShareHandler {
	return &ShareHandler{
		db:  database.GetDB(),
		cfg: cfg,
	}
}

// CreateShare 創建分享
// @Summary 創建分享連結
// @Description 為指定內容創建分享連結
// @Tags share
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ShareRequest true "分享請求"
// @Success 201 {object} vo.Response{data=dto.ShareResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Router /shares [post]
func (h *ShareHandler) CreateShare(c *gin.Context) {
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

	var req dto.ShareRequest
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

	contentID, err := uuid.Parse(req.ContentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid content ID format",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 驗證內容是否存在
	contentExists, err := h.checkContentExists(req.ContentType, contentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to check content",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	if !contentExists {
		c.JSON(http.StatusNotFound, vo.NewErrorResponse(
			"not_found",
			"Content not found",
			"CONTENT_NOT_FOUND",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 生成分享 URL
	shareURL := h.generateShareURL(req.ContentType, req.ContentID)
	shortURL := h.generateShortURL()

	// 創建分享記錄
	share := models.Share{
		UserID:      uuid.MustParse(userID),
		ContentType: req.ContentType,
		ContentID:   contentID,
		Platform:    req.Platform,
		ShareURL:    shareURL,
		ShortURL:    shortURL,
		Message:     req.Message,
		IsActive:    true,
	}

	if err := h.db.Create(&share).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to create share",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 生成 QR Code (簡化實現)
	qrCode := h.generateQRCode(shareURL)

	// 構建回應
	response := dto.ShareResponse{
		ShareID:     share.ID.String(),
		ContentType: share.ContentType,
		ContentID:   share.ContentID.String(),
		ShareURL:    share.ShareURL,
		ShortURL:    share.ShortURL,
		Platform:    share.Platform,
		QRCode:      qrCode,
		CreatedAt:   share.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, vo.SuccessResponse(response, "Share created successfully"))
}

// GetShare 獲取分享內容
// @Summary 獲取分享內容
// @Description 通過分享 ID 獲取分享內容詳情
// @Tags share
// @Accept json
// @Produce json
// @Param shareId path string true "分享ID"
// @Success 200 {object} vo.Response{data=dto.ShareContentResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Router /shares/{shareId} [get]
func (h *ShareHandler) GetShare(c *gin.Context) {
	shareID := c.Param("shareId")
	if shareID == "" {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Share ID is required",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	parsedShareID, err := uuid.Parse(shareID)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid share ID format",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 查找分享記錄
	var share models.Share
	if err := h.db.Where("id = ? AND is_active = ?", parsedShareID, true).First(&share).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Share not found",
				"SHARE_NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get share",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 檢查是否過期
	if share.IsExpired() {
		c.JSON(http.StatusNotFound, vo.NewErrorResponse(
			"not_found",
			"Share has expired",
			"SHARE_EXPIRED",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 記錄點擊
	h.recordClick(share.ID, c)

	// 增加觀看次數
	share.IncrementViewCount(h.db)

	// 獲取內容詳情
	contentDetails, err := h.getContentDetails(share.ContentType, share.ContentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get content details",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 獲取分享統計
	var shareCount int64
	h.db.Model(&models.Share{}).Where("content_type = ? AND content_id = ?",
		share.ContentType, share.ContentID).Count(&shareCount)

	// 構建回應
	response := dto.ShareContentResponse{
		ContentType: share.ContentType,
		ContentID:   share.ContentID.String(),
		Title:       contentDetails["title"].(string),
		Description: getStringValue(contentDetails, "description"),
		ImageURL:    getStringValue(contentDetails, "image_url"),
		URL:         h.generateContentURL(share.ContentType, share.ContentID.String()),
		Author:      getStringValue(contentDetails, "author"),
		IsValid:     true,
		ViewCount:   share.ViewCount,
		ShareCount:  shareCount,
		CreatedAt:   contentDetails["created_at"].(string),
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "Share content retrieved successfully"))
}

// GetUserShares 獲取使用者分享列表
// @Summary 獲取使用者分享列表
// @Description 獲取當前使用者的分享列表
// @Tags share
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param content_type query string false "內容類型"
// @Param platform query string false "分享平台"
// @Param page query int false "頁碼" default(1)
// @Param limit query int false "每頁數量" default(10)
// @Success 200 {object} vo.Response{data=dto.ShareListResponse}
// @Failure 401 {object} vo.ErrorResponse
// @Router /users/me/shares [get]
func (h *ShareHandler) GetUserShares(c *gin.Context) {
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
	contentType := strings.TrimSpace(c.Query("content_type"))
	platform := strings.TrimSpace(c.Query("platform"))
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

	// 構建查詢
	dbQuery := h.db.Model(&models.Share{}).Where("user_id = ? AND is_active = ?", userID, true)

	// 添加篩選條件
	if contentType != "" {
		dbQuery = dbQuery.Where("content_type = ?", contentType)
	}
	if platform != "" {
		dbQuery = dbQuery.Where("platform = ?", platform)
	}

	// 獲取總數
	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to count shares",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 獲取分享列表
	var shares []models.Share
	if err := dbQuery.Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&shares).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get shares",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 轉換為 DTO
	var shareResponses []dto.ShareResponse
	for _, share := range shares {
		response := dto.ShareResponse{
			ShareID:     share.ID.String(),
			ContentType: share.ContentType,
			ContentID:   share.ContentID.String(),
			ShareURL:    share.ShareURL,
			ShortURL:    share.ShortURL,
			Platform:    share.Platform,
			CreatedAt:   share.CreatedAt.Format(time.RFC3339),
		}
		if share.ExpiresAt != nil {
			response.ExpiresAt = share.ExpiresAt.Format(time.RFC3339)
		}
		shareResponses = append(shareResponses, response)
	}

	// 構建分頁回應
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	response := dto.ShareListResponse{
		Shares:     shareResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasMore:    page < totalPages,
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "User shares retrieved successfully"))
}

// GetShareStats 獲取分享統計
// @Summary 獲取分享統計
// @Description 獲取指定內容的分享統計資訊
// @Tags share
// @Accept json
// @Produce json
// @Param content_type query string true "內容類型"
// @Param content_id query string true "內容ID"
// @Success 200 {object} vo.Response{data=dto.ShareStatsResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Router /shares/stats [get]
func (h *ShareHandler) GetShareStats(c *gin.Context) {
	contentType := c.Query("content_type")
	contentIDStr := c.Query("content_id")

	if contentType == "" || contentIDStr == "" {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Content type and content ID are required",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	contentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid content ID format",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 獲取總分享數
	var totalShares int64
	h.db.Model(&models.Share{}).Where("content_type = ? AND content_id = ? AND is_active = ?",
		contentType, contentID, true).Count(&totalShares)

	// 獲取各平台分享統計
	var platformStats []struct {
		Platform string `json:"platform"`
		Count    int64  `json:"count"`
	}
	h.db.Model(&models.Share{}).
		Select("platform, COUNT(*) as count").
		Where("content_type = ? AND content_id = ? AND is_active = ?", contentType, contentID, true).
		Group("platform").
		Find(&platformStats)

	// 轉換為 map
	platforms := make(map[string]int64)
	for _, stat := range platformStats {
		platforms[stat.Platform] = stat.Count
	}

	// 獲取最近的分享記錄
	var recentShares []models.Share
	h.db.Where("content_type = ? AND content_id = ? AND is_active = ?", contentType, contentID, true).
		Order("created_at DESC").Limit(5).Find(&recentShares)

	var recentShareResponses []dto.ShareResponse
	for _, share := range recentShares {
		response := dto.ShareResponse{
			ShareID:     share.ID.String(),
			ContentType: share.ContentType,
			ContentID:   share.ContentID.String(),
			Platform:    share.Platform,
			CreatedAt:   share.CreatedAt.Format(time.RFC3339),
		}
		recentShareResponses = append(recentShareResponses, response)
	}

	// 構建回應
	response := dto.ShareStatsResponse{
		ContentID:    contentIDStr,
		ContentType:  contentType,
		TotalShares:  totalShares,
		Platforms:    platforms,
		RecentShares: recentShareResponses,
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "Share stats retrieved successfully"))
}

// 輔助方法

// checkContentExists 檢查內容是否存在
func (h *ShareHandler) checkContentExists(contentType string, contentID uuid.UUID) (bool, error) {
	switch contentType {
	case "article":
		var count int64
		err := h.db.Model(&models.Article{}).Where("id = ? AND is_published = ?", contentID, true).Count(&count).Error
		return count > 0, err
	case "location":
		var count int64
		err := h.db.Model(&models.Location{}).Where("id = ? AND is_public = ?", contentID, true).Count(&count).Error
		return count > 0, err
	case "quiz":
		var count int64
		err := h.db.Model(&models.Quiz{}).Where("id = ? AND is_active = ?", contentID, true).Count(&count).Error
		return count > 0, err
	default:
		return false, fmt.Errorf("unsupported content type: %s", contentType)
	}
}

// generateShareURL 生成分享 URL
func (h *ShareHandler) generateShareURL(contentType, contentID string) string {
	baseURL := "https://mindhelp.app" // 從配置中獲取
	return fmt.Sprintf("%s/share/%s/%s", baseURL, contentType, contentID)
}

// generateContentURL 生成內容 URL
func (h *ShareHandler) generateContentURL(contentType, contentID string) string {
	baseURL := "https://mindhelp.app"
	switch contentType {
	case "article":
		return fmt.Sprintf("%s/articles/%s", baseURL, contentID)
	case "location":
		return fmt.Sprintf("%s/locations/%s", baseURL, contentID)
	case "quiz":
		return fmt.Sprintf("%s/quizzes/%s", baseURL, contentID)
	default:
		return baseURL
	}
}

// generateShortURL 生成短連結
func (h *ShareHandler) generateShortURL() string {
	// 簡化實現 - 生成隨機短碼
	bytes := make([]byte, 6)
	rand.Read(bytes)
	shortCode := base64.URLEncoding.EncodeToString(bytes)[:8]
	return fmt.Sprintf("https://mh.ly/%s", shortCode)
}

// generateQRCode 生成 QR Code (簡化實現)
func (h *ShareHandler) generateQRCode(url string) string {
	// 實際實現中應該使用 QR Code 生成庫
	// 這裡返回一個佔位符 base64 字符串
	return "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/59BPwAHogJ/jTeLNAAAAAElFTkSuQmCC"
}

// getContentDetails 獲取內容詳情
func (h *ShareHandler) getContentDetails(contentType string, contentID uuid.UUID) (map[string]interface{}, error) {
	switch contentType {
	case "article":
		var article models.Article
		if err := h.db.Where("id = ?", contentID).First(&article).Error; err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"title":       article.Title,
			"description": article.Summary,
			"image_url":   article.ImageURL,
			"author":      article.Author,
			"created_at":  article.CreatedAt.Format(time.RFC3339),
		}, nil
	case "location":
		var location models.Location
		if err := h.db.Where("id = ?", contentID).First(&location).Error; err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"title":       location.Name,
			"description": location.Description,
			"image_url":   "",
			"author":      "",
			"created_at":  location.CreatedAt.Format(time.RFC3339),
		}, nil
	case "quiz":
		var quiz models.Quiz
		if err := h.db.Where("id = ?", contentID).First(&quiz).Error; err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"title":       quiz.Title,
			"description": quiz.Description,
			"image_url":   "",
			"author":      "",
			"created_at":  quiz.CreatedAt.Format(time.RFC3339),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}
}

// recordClick 記錄分享點擊
func (h *ShareHandler) recordClick(shareID uuid.UUID, c *gin.Context) {
	click := models.ShareClick{
		ShareID:   shareID,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Referer:   c.GetHeader("Referer"),
	}

	// 非阻塞記錄點擊
	go func() {
		h.db.Create(&click)
	}()
}

// getStringValue 安全地從 map 中獲取字符串值
func getStringValue(m map[string]interface{}, key string) string {
	if val, exists := m[key]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
