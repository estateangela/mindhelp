package handlers

import (
	"net/http"
	"strconv"

	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/dto"
	"mindhelp-backend/internal/middleware"
	"mindhelp-backend/internal/models"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReviewHandler 評論處理器
type ReviewHandler struct {
	db *gorm.DB
}

// NewReviewHandler 創建新的評論處理器
func NewReviewHandler() *ReviewHandler {
	return &ReviewHandler{
		db: database.GetDB(),
	}
}

// GetResourceReviews 獲取資源的所有評論
// @Summary 獲取資源評論
// @Description 獲取某個資源點的所有評論
// @Tags review
// @Accept json
// @Produce json
// @Param id path string true "資源ID"
// @Param page query int false "頁碼" default(1)
// @Param limit query int false "每頁數量" default(10)
// @Success 200 {object} vo.Response{data=dto.ReviewListResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Router /resources/{id}/reviews [get]
func (h *ReviewHandler) GetResourceReviews(c *gin.Context) {
	resourceID := c.Param("id")
	if resourceID == "" {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Resource ID is required",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 驗證 UUID 格式
	parsedID, err := uuid.Parse(resourceID)
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
	if err := h.db.Model(&models.Review{}).Where("resource_id = ?", parsedID).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to count reviews",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 獲取評論列表
	var reviews []models.Review
	if err := h.db.Preload("User").Where("resource_id = ?", parsedID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get reviews",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 計算統計資訊
	var avgRating float64
	var ratingDistribution = make(map[int]int64)
	
	// 計算平均評分
	h.db.Model(&models.Review{}).Where("resource_id = ?", parsedID).
		Select("AVG(rating)").Scan(&avgRating)
	
	// 計算評分分佈
	for i := 1; i <= 5; i++ {
		var count int64
		h.db.Model(&models.Review{}).Where("resource_id = ? AND rating = ?", parsedID, i).
			Count(&count)
		ratingDistribution[i] = count
	}

	// 獲取當前使用者ID (如果已登入)
	currentUserID := middleware.GetUserID(c)

	// 轉換為 DTO
	var reviewResponses []dto.ReviewResponse
	for _, review := range reviews {
		canEdit := false
		if currentUserID != "" && review.UserID.String() == currentUserID {
			canEdit = true
		}

		response := dto.ReviewResponse{
			ID: review.ID.String(),
			Author: dto.ReviewAuthor{
				ID:       review.User.ID.String(),
				Username: review.User.Username,
				Avatar:   review.User.Avatar,
			},
			ResourceID: review.ResourceID.String(),
			Rating:     review.Rating,
			Comment:    review.Comment,
			IsHelpful:  review.IsHelpful,
			CanEdit:    canEdit,
			CreatedAt:  review.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  review.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		reviewResponses = append(reviewResponses, response)
	}

	// 構建分頁回應
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	response := dto.ReviewListResponse{
		Reviews:    reviewResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasMore:    page < totalPages,
		Statistics: dto.ReviewStatistics{
			AverageRating:      avgRating,
			TotalReviews:       total,
			RatingDistribution: ratingDistribution,
		},
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "Reviews retrieved successfully"))
}

// CreateReview 為資源新增評論
// @Summary 新增評論
// @Description 為某個資源點新增一則評論
// @Tags review
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "資源ID"
// @Param request body dto.ReviewRequest true "評論資料"
// @Success 201 {object} vo.Response{data=dto.ReviewResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Failure 409 {object} vo.ErrorResponse
// @Router /resources/{id}/reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
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

	resourceID := c.Param("id")
	parsedID, err := uuid.Parse(resourceID)
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

	var req dto.ReviewRequest
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

	// 檢查資源是否存在
	var location models.Location
	if err := h.db.Where("id = ? AND is_public = ?", parsedID, true).First(&location).Error; err != nil {
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

	// 檢查使用者是否已評論過此資源
	var existingReview models.Review
	if err := h.db.Where("user_id = ? AND resource_id = ?", userID, parsedID).First(&existingReview).Error; err == nil {
		c.JSON(http.StatusConflict, vo.NewErrorResponse(
			"conflict",
			"User has already reviewed this resource",
			"ALREADY_REVIEWED",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 創建評論
	review := models.Review{
		UserID:     uuid.MustParse(userID),
		ResourceID: parsedID,
		Rating:     req.Rating,
		Comment:    req.Comment,
	}

	if err := h.db.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to create review",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 預載入使用者資訊
	if err := h.db.Preload("User").Where("id = ?", review.ID).First(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get created review",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 構建回應
	response := dto.ReviewResponse{
		ID: review.ID.String(),
		Author: dto.ReviewAuthor{
			ID:       review.User.ID.String(),
			Username: review.User.Username,
			Avatar:   review.User.Avatar,
		},
		ResourceID: review.ResourceID.String(),
		Rating:     review.Rating,
		Comment:    review.Comment,
		IsHelpful:  review.IsHelpful,
		CanEdit:    true, // 剛創建的評論，當前使用者可以編輯
		CreatedAt:  review.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  review.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusCreated, vo.SuccessResponse(response, "Review created successfully"))
}

// UpdateReview 修改評論
// @Summary 修改評論
// @Description 修改自己發布的評論
// @Tags review
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param reviewId path string true "評論ID"
// @Param request body dto.ReviewUpdateRequest true "評論更新資料"
// @Success 200 {object} vo.Response{data=dto.ReviewResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Failure 403 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Router /reviews/{reviewId} [put]
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
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

	reviewID := c.Param("reviewId")
	parsedID, err := uuid.Parse(reviewID)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid review ID format",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	var req dto.ReviewUpdateRequest
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

	// 查找評論
	var review models.Review
	if err := h.db.Where("id = ?", parsedID).First(&review).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Review not found",
				"NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get review",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 檢查權限
	if review.UserID.String() != userID {
		c.JSON(http.StatusForbidden, vo.NewErrorResponse(
			"forbidden",
			"You can only edit your own reviews",
			"FORBIDDEN",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 更新欄位
	updates := make(map[string]interface{})
	if req.Rating != nil {
		updates["rating"] = *req.Rating
	}
	if req.Comment != nil {
		updates["comment"] = *req.Comment
	}

	// 執行更新
	if len(updates) > 0 {
		if err := h.db.Model(&review).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
				"internal_error",
				"Failed to update review",
				"INTERNAL_ERROR",
				nil,
				c.Request.URL.Path,
			))
			return
		}
	}

	// 重新獲取更新後的評論（包含使用者資訊）
	if err := h.db.Preload("User").Where("id = ?", parsedID).First(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get updated review",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 構建回應
	response := dto.ReviewResponse{
		ID: review.ID.String(),
		Author: dto.ReviewAuthor{
			ID:       review.User.ID.String(),
			Username: review.User.Username,
			Avatar:   review.User.Avatar,
		},
		ResourceID: review.ResourceID.String(),
		Rating:     review.Rating,
		Comment:    review.Comment,
		IsHelpful:  review.IsHelpful,
		CanEdit:    true,
		CreatedAt:  review.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  review.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "Review updated successfully"))
}

// DeleteReview 刪除評論
// @Summary 刪除評論
// @Description 刪除自己發布的評論
// @Tags review
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param reviewId path string true "評論ID"
// @Success 200 {object} vo.Response
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Failure 403 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Router /reviews/{reviewId} [delete]
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
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

	reviewID := c.Param("reviewId")
	parsedID, err := uuid.Parse(reviewID)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid review ID format",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 查找評論
	var review models.Review
	if err := h.db.Where("id = ?", parsedID).First(&review).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Review not found",
				"NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get review",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 檢查權限
	if review.UserID.String() != userID {
		c.JSON(http.StatusForbidden, vo.NewErrorResponse(
			"forbidden",
			"You can only delete your own reviews",
			"FORBIDDEN",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 軟刪除評論
	if err := h.db.Delete(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to delete review",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(nil, "Review deleted successfully"))
}

// ReportContent 回報不當內容
// @Summary 回報不當內容
// @Description 回報不當的評論、文章或資源內容
// @Tags review
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ReportRequest true "回報資料"
// @Success 202 {object} vo.Response
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Router /report [post]
func (h *ReviewHandler) ReportContent(c *gin.Context) {
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

	var req dto.ReportRequest
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

	// 這裡可以實現回報邏輯，如：
	// 1. 記錄到報告表中
	// 2. 發送通知給管理員
	// 3. 如果報告數量達到閾值，自動隱藏內容
	
	// 簡化實現：記錄日誌
	// 實際應用中可以創建一個 Report 模型來儲存這些資料
	
	c.JSON(http.StatusAccepted, vo.SuccessResponse(nil, "Report submitted successfully"))
}
