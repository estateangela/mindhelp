package handlers

import (
	"net/http"

	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
)

// AdminHandler 管理員相關處理器
type AdminHandler struct{}

// NewAdminHandler 創建新的管理員處理器
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// SeedDatabase 資料庫種子資料
// @Summary 插入種子資料
// @Description 插入範例的諮商師、諮商所和推薦醫師資料
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} vo.ErrorResponse
// @Router /admin/seed-database [post]
func (h *AdminHandler) SeedDatabase(c *gin.Context) {
	// 獲取資料庫連接
	db, err := database.GetDBSafely()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, vo.NewErrorResponse(
			"database_unavailable",
			"資料庫暫時無法使用，請稍後再試",
			"DATABASE_UNAVAILABLE",
			err.Error(),
			c.Request.URL.Path,
		))
		return
	}

	// 基本健康檢查
	sqlDB, err := db.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"database_error",
			"無法獲取資料庫實例",
			"DATABASE_ERROR",
			err.Error(),
			c.Request.URL.Path,
		))
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"database_error",
			"資料庫連接測試失敗",
			"DATABASE_ERROR",
			err.Error(),
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusOK, vo.NewResponse("success", gin.H{
		"message": "資料庫種子資料功能已準備就緒",
		"note":    "具體的種子資料插入功能需要進一步實現",
	}))
}

// GetDatabaseStats 獲取資料庫統計資訊
// @Summary 獲取資料庫統計
// @Description 獲取資料庫連接和表統計資訊
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} vo.ErrorResponse
// @Router /admin/database-stats [get]
func (h *AdminHandler) GetDatabaseStats(c *gin.Context) {
	// 獲取資料庫連接統計
	stats, err := database.GetConnectionStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"database_error",
			"無法獲取資料庫統計資訊",
			"DATABASE_ERROR",
			err.Error(),
			c.Request.URL.Path,
		))
		return
	}

	// 檢查資料庫健康狀態
	healthy := database.IsHealthy()

	c.JSON(http.StatusOK, vo.NewResponse("success", gin.H{
		"database_healthy":      healthy,
		"connection_pool_stats": stats,
		"timestamp":             "now", // 可以改為實際時間戳
	}))
}
