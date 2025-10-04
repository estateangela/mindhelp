package handlers

import (
	"net/http"
	"time"

	"mindhelp-backend/internal/scheduler"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
)

// AdminSchedulerHandler 管理員定時任務處理器
type AdminSchedulerHandler struct {
	scheduler *scheduler.Scheduler
}

// NewAdminSchedulerHandler 創建管理員定時任務處理器
func NewAdminSchedulerHandler(scheduler *scheduler.Scheduler) *AdminSchedulerHandler {
	return &AdminSchedulerHandler{
		scheduler: scheduler,
	}
}

// GetScheduledJobs 獲取定時任務狀態
// @Summary 獲取定時任務狀態
// @Description 獲取所有已排程的定時任務資訊
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} vo.Response{data=map[string]interface{}}
// @Failure 401 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /admin/scheduler/jobs [get]
func (h *AdminSchedulerHandler) GetScheduledJobs(c *gin.Context) {
	jobs := h.scheduler.GetScheduledJobs()

	c.JSON(http.StatusOK, vo.SuccessResponse(map[string]interface{}{
		"jobs":      jobs,
		"timestamp": time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}, "Scheduled jobs retrieved successfully"))
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
func (h *AdminSchedulerHandler) TriggerHourlyNotification(c *gin.Context) {
	if err := h.scheduler.TriggerHourlyNotification(); err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to trigger hourly notification",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(map[string]interface{}{
		"triggered_at": time.Now().Format("2006-01-02T15:04:05Z07:00"),
		"job_type":     "hourly_notification",
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
func (h *AdminSchedulerHandler) TriggerWeeklyNotification(c *gin.Context) {
	if err := h.scheduler.TriggerWeeklyNotification(); err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to trigger weekly notification",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(map[string]interface{}{
		"triggered_at": time.Now().Format("2006-01-02T15:04:05Z07:00"),
		"job_type":     "weekly_notification",
	}, "Weekly notification triggered successfully"))
}
