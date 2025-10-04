package handlers

import (
	"net/http"
	"time"

	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
)

// SchedulerTriggerHandler 定時任務觸發處理器
type SchedulerTriggerHandler struct{}

// NewSchedulerTriggerHandler 創建定時任務觸發處理器
func NewSchedulerTriggerHandler() *SchedulerTriggerHandler {
	return &SchedulerTriggerHandler{}
}

// TriggerHourlyNotification 手動觸發每小時通知
// @Summary 手動觸發每小時通知
// @Description 立即執行每小時通知任務，忽略時間和排程限制
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} vo.Response
// @Failure 401 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /admin/scheduler/trigger/hourly [post]
func (h *SchedulerTriggerHandler) TriggerHourlyNotification(c *gin.Context) {
	// 這裡暫時返回成功，實際的觸發邏輯會在後續實現
	c.JSON(http.StatusOK, vo.SuccessResponse(map[string]interface{}{
		"triggered_at": time.Now().Format("2006-01-02T15:04:05Z07:00"),
		"job_type":     "hourly_notification",
		"message":      "Hourly notification triggered successfully (placeholder)",
	}, "Hourly notification triggered successfully"))
}

// TriggerWeeklyNotification 手動觸發每週通知
// @Summary 手動觸發每週通知
// @Description 立即執行每週通知任務，忽略時間和排程限制
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} vo.Response
// @Failure 401 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /admin/scheduler/trigger/weekly [post]
func (h *SchedulerTriggerHandler) TriggerWeeklyNotification(c *gin.Context) {
	// 這裡暫時返回成功，實際的觸發邏輯會在後續實現
	c.JSON(http.StatusOK, vo.SuccessResponse(map[string]interface{}{
		"triggered_at": time.Now().Format("2006-01-02T15:04:05Z07:00"),
		"job_type":     "weekly_notification",
		"message":      "Weekly notification triggered successfully (placeholder)",
	}, "Weekly notification triggered successfully"))
}

// GetSchedulerStatus 獲取定時任務狀態
// @Summary 獲取定時任務狀態
// @Description 獲取所有已排程的定時任務資訊
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} vo.Response{data=map[string]interface{}}
// @Failure 401 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /admin/scheduler/status [get]
func (h *SchedulerTriggerHandler) GetSchedulerStatus(c *gin.Context) {
	c.JSON(http.StatusOK, vo.SuccessResponse(map[string]interface{}{
		"status":    "running",
		"timestamp": time.Now().Format("2006-01-02T15:04:05Z07:00"),
		"jobs": []map[string]interface{}{
			{
				"id":       "hourly_notification",
				"schedule": "0 * * * *",
				"next_run": "Next hour at minute 0",
			},
			{
				"id":       "weekly_notification",
				"schedule": "0 20 * * 1",
				"next_run": "Next Monday at 8:00 PM (Taipei time)",
			},
		},
	}, "Scheduler status retrieved successfully"))
}
