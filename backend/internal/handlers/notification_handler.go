package handlers

import (
	"net/http"
	"strconv"

	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/dto"
	"mindhelp-backend/internal/models"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// NotificationHandler 通知處理器
type NotificationHandler struct{}

// NewNotificationHandler 創建通知處理器
func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

// GetNotifications 獲取用戶通知列表
// @Summary 獲取用戶通知列表
// @Description 獲取當前用戶的所有通知，支援分頁和篩選
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "頁碼" default(1)
// @Param page_size query int false "每頁數量" default(10)
// @Param type query string false "通知類型篩選"
// @Param unread_only query bool false "只顯示未讀通知" default(false)
// @Success 200 {object} vo.NotificationListVO
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /notifications [get]
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	// 從上下文獲取用戶ID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "User ID not found in context",
		})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Invalid user ID format",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "INVALID_USER_ID",
			Message: "Invalid user ID format",
		})
		return
	}

	// 解析查詢參數
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	notificationType := c.Query("type")
	unreadOnly := c.Query("unread_only") == "true"

	// 確保分頁參數合理
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 構建查詢
	query := database.GetDB().Model(&models.Notification{}).Where("user_id = ?", userID)

	// 添加篩選條件
	if notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}
	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}

	// 獲取總數
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to count notifications",
			Error:   err.Error(),
		})
		return
	}

	// 獲取未讀通知數量
	var unreadCount int64
	unreadQuery := database.GetDB().Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false)
	if notificationType != "" {
		unreadQuery = unreadQuery.Where("type = ?", notificationType)
	}
	unreadQuery.Count(&unreadCount)

	// 獲取通知列表
	var notifications []models.Notification
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to fetch notifications",
			Error:   err.Error(),
		})
		return
	}

	// 轉換為回應格式
	var notificationVOs []vo.NotificationVO
	for _, notification := range notifications {
		notificationVOs = append(notificationVOs, vo.NotificationVO{
			ID:        notification.ID,
			Title:     notification.Title,
			Content:   notification.Content,
			Type:      notification.Type,
			IsRead:    notification.IsRead,
			CreatedAt: notification.CreatedAt,
		})
	}

	// 計算是否有更多資料
	hasMore := int64(page*pageSize) < total

	c.JSON(http.StatusOK, vo.NotificationListVO{
		Notifications: notificationVOs,
		Total:         total,
		Page:          page,
		PageSize:      pageSize,
		HasMore:       hasMore,
		UnreadCount:   unreadCount,
	})
}

// MarkNotificationRead 標記通知為已讀
// @Summary 標記通知為已讀
// @Description 將指定的通知標記為已讀
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "通知ID"
// @Success 200 {object} vo.SuccessResponse
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /notifications/{id}/read [put]
func (h *NotificationHandler) MarkNotificationRead(c *gin.Context) {
	// 從上下文獲取用戶ID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "User ID not found in context",
		})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Invalid user ID format",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "INVALID_USER_ID",
			Message: "Invalid user ID format",
		})
		return
	}

	// 獲取通知ID
	notificationIDStr := c.Param("id")
	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "INVALID_NOTIFICATION_ID",
			Message: "Invalid notification ID format",
		})
		return
	}

	// 查找通知並驗證所有權
	var notification models.Notification
	if err := database.GetDB().Where("id = ? AND user_id = ?", notificationID, userID).First(&notification).Error; err != nil {
		c.JSON(http.StatusNotFound, vo.ErrorResponse{
			Code:    "NOTIFICATION_NOT_FOUND",
			Message: "Notification not found or access denied",
		})
		return
	}

	// 標記為已讀
	if err := database.GetDB().Model(&notification).Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to mark notification as read",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(nil, "Notification marked as read"))
}

// MarkAllNotificationsRead 標記所有通知為已讀
// @Summary 標記所有通知為已讀
// @Description 將當前用戶的所有通知標記為已讀，可選擇特定類型
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.MarkAllNotificationsReadRequest true "標記所有通知為已讀請求"
// @Success 200 {object} vo.SuccessResponse
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /notifications/read-all [put]
func (h *NotificationHandler) MarkAllNotificationsRead(c *gin.Context) {
	// 從上下文獲取用戶ID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "User ID not found in context",
		})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Invalid user ID format",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "INVALID_USER_ID",
			Message: "Invalid user ID format",
		})
		return
	}

	// 解析請求
	var req dto.MarkAllNotificationsReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	// 構建更新查詢
	query := database.GetDB().Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false)

	// 如果指定了類型，添加類型篩選
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	// 執行更新
	result := query.Update("is_read", true)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to mark notifications as read",
			Error:   result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(map[string]interface{}{
		"updated_count": result.RowsAffected,
	}, "All notifications marked as read"))
}

// DeleteNotification 刪除通知
// @Summary 刪除通知
// @Description 刪除指定的通知（軟刪除）
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "通知ID"
// @Success 200 {object} vo.SuccessResponse
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /notifications/{id} [delete]
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	// 從上下文獲取用戶ID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "User ID not found in context",
		})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Invalid user ID format",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "INVALID_USER_ID",
			Message: "Invalid user ID format",
		})
		return
	}

	// 獲取通知ID
	notificationIDStr := c.Param("id")
	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "INVALID_NOTIFICATION_ID",
			Message: "Invalid notification ID format",
		})
		return
	}

	// 查找通知並驗證所有權
	var notification models.Notification
	if err := database.GetDB().Where("id = ? AND user_id = ?", notificationID, userID).First(&notification).Error; err != nil {
		c.JSON(http.StatusNotFound, vo.ErrorResponse{
			Code:    "NOTIFICATION_NOT_FOUND",
			Message: "Notification not found or access denied",
		})
		return
	}

	// 軟刪除通知
	if err := database.GetDB().Delete(&notification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to delete notification",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(nil, "Notification deleted successfully"))
}
