package handlers

import (
	"encoding/json"
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

// NotificationHandler 通知處理器
type NotificationHandler struct {
	db *gorm.DB
}

// NewNotificationHandler 創建新的通知處理器
func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{
		db: database.GetDB(),
	}
}

// GetNotifications 獲取通知列表
// @Summary 獲取通知列表
// @Description 獲取使用者的通知列表
// @Tags notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "頁碼" default(1)
// @Param limit query int false "每頁數量" default(20)
// @Param unread_only query bool false "只顯示未讀通知" default(false)
// @Success 200 {object} vo.Response{data=dto.NotificationListResponse}
// @Failure 401 {object} vo.ErrorResponse
// @Router /notifications [get]
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
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
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	unreadOnly := c.Query("unread_only") == "true"

	// 參數驗證
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// 構建查詢
	dbQuery := h.db.Model(&models.Notification{}).Where("user_id = ?", userID)
	if unreadOnly {
		dbQuery = dbQuery.Where("is_read = ?", false)
	}

	// 獲取總數
	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to count notifications",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 獲取未讀數量
	var unreadCount int64
	h.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&unreadCount)

	// 獲取通知列表
	var notifications []models.Notification
	if err := dbQuery.Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get notifications",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 轉換為 DTO
	var notificationResponses []dto.NotificationResponse
	for _, notification := range notifications {
		// 解析 payload
		var payload map[string]interface{}
		if notification.Payload != "" {
			json.Unmarshal([]byte(notification.Payload), &payload)
		}

		response := dto.NotificationResponse{
			ID:        notification.ID.String(),
			Type:      notification.Type,
			Title:     notification.Title,
			Body:      notification.Body,
			IsRead:    notification.IsRead,
			Payload:   payload,
			CreatedAt: notification.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		notificationResponses = append(notificationResponses, response)
	}

	// 構建分頁回應
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	response := dto.NotificationListResponse{
		Notifications: notificationResponses,
		Total:         total,
		UnreadCount:   unreadCount,
		Page:          page,
		Limit:         limit,
		TotalPages:    totalPages,
		HasMore:       page < totalPages,
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "Notifications retrieved successfully"))
}

// MarkAsRead 標記通知為已讀
// @Summary 標記通知已讀
// @Description 將指定的通知標示為已讀
// @Tags notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.MarkAsReadRequest true "標記已讀請求"
// @Success 200 {object} vo.Response
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Router /notifications/mark-as-read [post]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
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

	var req dto.MarkAsReadRequest
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

	// 解析通知ID
	var notificationIDs []uuid.UUID
	for _, idStr := range req.NotificationIDs {
		if id, err := uuid.Parse(idStr); err == nil {
			notificationIDs = append(notificationIDs, id)
		}
	}

	if len(notificationIDs) == 0 {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"No valid notification IDs provided",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 標記為已讀（只能標記自己的通知）
	result := h.db.Model(&models.Notification{}).
		Where("user_id = ? AND id IN ?", userID, notificationIDs).
		Update("is_read", true)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to mark notifications as read",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(map[string]interface{}{
		"updated_count": result.RowsAffected,
	}, "Notifications marked as read successfully"))
}

// GetNotificationSettings 獲取通知設定
// @Summary 獲取通知設定
// @Description 獲取使用者的通知偏好設定
// @Tags notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} vo.Response{data=dto.NotificationSettingsResponse}
// @Failure 401 {object} vo.ErrorResponse
// @Router /users/me/notification-settings [get]
func (h *NotificationHandler) GetNotificationSettings(c *gin.Context) {
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

	// 查找使用者設定
	var settings models.UserSetting
	if err := h.db.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果沒有設定，返回預設值
			response := dto.NotificationSettingsResponse{
				NotifyNewArticle:    true,
				NotifyPromotions:    false,
				NotifySystemUpdates: true,
			}
			c.JSON(http.StatusOK, vo.SuccessResponse(response, "Notification settings retrieved successfully"))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get notification settings",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 構建回應
	response := dto.NotificationSettingsResponse{
		NotifyNewArticle:    settings.NotifyNewArticle,
		NotifyPromotions:    settings.NotifyPromotions,
		NotifySystemUpdates: settings.NotifySystemUpdates,
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "Notification settings retrieved successfully"))
}

// UpdateNotificationSettings 更新通知設定
// @Summary 更新通知設定
// @Description 更新使用者的通知偏好設定
// @Tags notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.NotificationSettingsRequest true "通知設定"
// @Success 200 {object} vo.Response{data=dto.NotificationSettingsResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Router /users/me/notification-settings [put]
func (h *NotificationHandler) UpdateNotificationSettings(c *gin.Context) {
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

	var req dto.NotificationSettingsRequest
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

	// 查找或創建使用者設定
	var settings models.UserSetting
	err := h.db.Where("user_id = ?", userID).First(&settings).Error
	if err == gorm.ErrRecordNotFound {
		// 創建新設定
		settings = models.UserSetting{
			UserID:              uuid.MustParse(userID),
			NotifyNewArticle:    true,
			NotifyPromotions:    false,
			NotifySystemUpdates: true,
		}
		if err := h.db.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
				"internal_error",
				"Failed to create notification settings",
				"INTERNAL_ERROR",
				nil,
				c.Request.URL.Path,
			))
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get notification settings",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 更新欄位
	updates := make(map[string]interface{})
	if req.NotifyNewArticle != nil {
		updates["notify_new_article"] = *req.NotifyNewArticle
	}
	if req.NotifyPromotions != nil {
		updates["notify_promotions"] = *req.NotifyPromotions
	}
	if req.NotifySystemUpdates != nil {
		updates["notify_system_updates"] = *req.NotifySystemUpdates
	}

	// 執行更新
	if len(updates) > 0 {
		if err := h.db.Model(&settings).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
				"internal_error",
				"Failed to update notification settings",
				"INTERNAL_ERROR",
				nil,
				c.Request.URL.Path,
			))
			return
		}
	}

	// 重新獲取更新後的設定
	if err := h.db.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get updated notification settings",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 構建回應
	response := dto.NotificationSettingsResponse{
		NotifyNewArticle:    settings.NotifyNewArticle,
		NotifyPromotions:    settings.NotifyPromotions,
		NotifySystemUpdates: settings.NotifySystemUpdates,
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "Notification settings updated successfully"))
}

// UpdatePushToken 註冊/更新推播 token
// @Summary 更新推播 Token
// @Description 註冊或更新裝置的推播 token
// @Tags notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.PushTokenRequest true "推播 Token"
// @Success 204 "No Content"
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Router /users/me/push-token [post]
func (h *NotificationHandler) UpdatePushToken(c *gin.Context) {
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

	var req dto.PushTokenRequest
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

	// 查找或創建使用者設定
	var settings models.UserSetting
	err := h.db.Where("user_id = ?", userID).First(&settings).Error
	if err == gorm.ErrRecordNotFound {
		// 創建新設定
		settings = models.UserSetting{
			UserID:              uuid.MustParse(userID),
			NotifyNewArticle:    true,
			NotifyPromotions:    false,
			NotifySystemUpdates: true,
			PushToken:          req.Token,
			Platform:           req.Platform,
		}
		if err := h.db.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
				"internal_error",
				"Failed to save push token",
				"INTERNAL_ERROR",
				nil,
				c.Request.URL.Path,
			))
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get user settings",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	} else {
		// 更新現有設定
		updates := map[string]interface{}{
			"push_token": req.Token,
			"platform":   req.Platform,
		}
		if err := h.db.Model(&settings).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
				"internal_error",
				"Failed to update push token",
				"INTERNAL_ERROR",
				nil,
				c.Request.URL.Path,
			))
			return
		}
	}

	c.Status(http.StatusNoContent)
}
