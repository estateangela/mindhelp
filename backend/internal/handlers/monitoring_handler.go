package handlers

import (
	"net/http"
	"runtime"
	"time"

	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// MonitoringHandler 監控處理器
type MonitoringHandler struct {
	startTime time.Time
}

// NewMonitoringHandler 創建監控處理器
func NewMonitoringHandler() *MonitoringHandler {
	return &MonitoringHandler{
		startTime: time.Now(),
	}
}

// HealthCheck 健康檢查
func (h *MonitoringHandler) HealthCheck(c *gin.Context) {
	// 檢查資料庫連線 (不影響整體健康狀態)
	dbStatus := "healthy"
	db := database.GetDB()
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "unhealthy"
		}
	} else {
		dbStatus = "connecting" // 改為 connecting 而不是 disconnected
	}

	// 系統資訊
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// 服務本身是健康的，即使資料庫還在連接中
	health := gin.H{
		"status":    "ok", // 始終返回 ok，除非有嚴重錯誤
		"service":   "mindhelp-backend",
		"version":   "1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(h.startTime).String(),
		"checks": gin.H{
			"database": dbStatus,
		},
		"system": gin.H{
			"goroutines":   runtime.NumGoroutine(),
			"memory_alloc": memStats.Alloc / 1024 / 1024,      // MB
			"memory_total": memStats.TotalAlloc / 1024 / 1024, // MB
			"memory_sys":   memStats.Sys / 1024 / 1024,        // MB
			"gc_runs":      memStats.NumGC,
		},
	}

	statusCode := http.StatusOK
	// 只有在資料庫完全無法連接時才返回 503，connecting 狀態仍返回 200
	if dbStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
		health["status"] = "degraded"
	} else if dbStatus == "connecting" {
		health["status"] = "starting"
	}

	c.JSON(statusCode, health)
}

// DetailedHealthCheck 詳細健康檢查
func (h *MonitoringHandler) DetailedHealthCheck(c *gin.Context) {
	// 獲取資料庫統計資訊
	dbStats, dbErr := database.GetConnectionStats()
	if dbErr != nil {
		dbStats = map[string]interface{}{"error": dbErr.Error()}
	}

	// 系統統計
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	c.JSON(http.StatusOK, gin.H{
		"service":   "mindhelp-backend",
		"version":   "1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(h.startTime).String(),
		"database":  dbStats,
		"system": gin.H{
			"goroutines":     runtime.NumGoroutine(),
			"memory_alloc":   memStats.Alloc,
			"memory_total":   memStats.TotalAlloc,
			"memory_sys":     memStats.Sys,
			"gc_runs":        memStats.NumGC,
			"last_gc":        time.Unix(0, int64(memStats.LastGC)).Format(time.RFC3339),
			"heap_objects":   memStats.HeapObjects,
			"stack_in_use":   memStats.StackInuse,
		},
	})
}

// ReadinessCheck 就緒檢查 - 檢查服務是否準備好處理請求
func (h *MonitoringHandler) ReadinessCheck(c *gin.Context) {
	// 檢查資料庫連接
	if err := database.CheckConnection(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"reason": "database not available",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

// LivenessCheck 存活檢查 - 檢查服務是否還活著
func (h *MonitoringHandler) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
		"uptime": time.Since(h.startTime).String(),
	})
}

// Metrics 獲取指標
func (h *MonitoringHandler) Metrics(c *gin.Context) {
	metrics := middleware.GetMetrics()
	c.JSON(http.StatusOK, metrics)
}
