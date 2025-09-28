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
	// 獲取詳細的系統資訊
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// 資料庫連線檢查
	dbChecks := gin.H{}
	db := database.GetDB()
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			stats := sqlDB.Stats()
			dbChecks["connection_pool"] = gin.H{
				"open_connections": stats.OpenConnections,
				"in_use":           stats.InUse,
				"idle":             stats.Idle,
				"wait_count":       stats.WaitCount,
				"wait_duration":    stats.WaitDuration.String(),
			}

			// 測試查詢
			var result int
			err = db.Raw("SELECT 1").Scan(&result).Error
			dbChecks["query_test"] = gin.H{
				"status": map[bool]string{true: "pass", false: "fail"}[err == nil],
				"error":  map[bool]interface{}{true: nil, false: err.Error()}[err == nil],
			}
		}
	}

	health := gin.H{
		"service":   "mindhelp-backend",
		"version":   "1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(h.startTime).String(),
		"database":  dbChecks,
		"system": gin.H{
			"cpu_count":  runtime.NumCPU(),
			"goroutines": runtime.NumGoroutine(),
			"memory": gin.H{
				"alloc":       memStats.Alloc,
				"total_alloc": memStats.TotalAlloc,
				"sys":         memStats.Sys,
				"heap_alloc":  memStats.HeapAlloc,
				"heap_sys":    memStats.HeapSys,
			},
			"gc": gin.H{
				"num_gc":      memStats.NumGC,
				"pause_total": memStats.PauseTotalNs,
				"last_gc":     time.Unix(0, int64(memStats.LastGC)).Format(time.RFC3339),
			},
		},
	}

	c.JSON(http.StatusOK, health)
}

// Metrics 取得效能指標
func (h *MonitoringHandler) Metrics(c *gin.Context) {
	metrics := middleware.GetMetrics()

	// 添加系統指標
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics["system"] = gin.H{
		"goroutines": runtime.NumGoroutine(),
		"memory_mb":  memStats.Alloc / 1024 / 1024,
		"uptime":     time.Since(h.startTime).String(),
	}

	c.JSON(http.StatusOK, metrics)
}

// ReadinessCheck 就緒檢查
func (h *MonitoringHandler) ReadinessCheck(c *gin.Context) {
	// 檢查所有必要服務是否就緒
	ready := true
	checks := gin.H{}

	// 檢查資料庫
	db := database.GetDB()
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			ready = false
			checks["database"] = "not ready"
		} else {
			checks["database"] = "ready"
		}
	} else {
		ready = false
		checks["database"] = "not connected"
	}

	response := gin.H{
		"ready":  ready,
		"checks": checks,
	}

	statusCode := http.StatusOK
	if !ready {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// LivenessCheck 存活檢查
func (h *MonitoringHandler) LivenessCheck(c *gin.Context) {
	// 簡單的存活檢查
	c.JSON(http.StatusOK, gin.H{
		"alive":     true,
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(h.startTime).String(),
	})
}
