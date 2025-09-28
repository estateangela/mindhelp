package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RequestMetrics struct {
	TotalRequests    int64                  `json:"total_requests"`
	ActiveRequests   int64                  `json:"active_requests"`
	RequestsByMethod map[string]int64       `json:"requests_by_method"`
	RequestsByStatus map[string]int64       `json:"requests_by_status"`
	AverageLatency   float64                `json:"average_latency_ms"`
	StartTime        time.Time              `json:"start_time"`
	Uptime           string                 `json:"uptime"`
}

var (
	metrics = &RequestMetrics{
		RequestsByMethod: make(map[string]int64),
		RequestsByStatus: make(map[string]int64),
		StartTime:        time.Now(),
	}
	metricsMutex sync.RWMutex
	totalLatency time.Duration
)

// MetricsMiddleware 記錄請求指標
func MetricsMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		metricsMutex.Lock()
		defer metricsMutex.Unlock()

		// 更新指標
		metrics.TotalRequests++
		metrics.RequestsByMethod[param.Method]++
		metrics.RequestsByStatus[string(rune(param.StatusCode))]++

		// 更新平均延遲
		totalLatency += param.Latency
		if metrics.TotalRequests > 0 {
			metrics.AverageLatency = float64(totalLatency.Nanoseconds()) / float64(metrics.TotalRequests) / 1e6 // 轉換為毫秒
		}

		// 更新運行時間
		metrics.Uptime = time.Since(metrics.StartTime).String()

		// 返回空字符串以避免重複日誌
		return ""
	})
}

// GetMetrics 獲取當前指標
func GetMetrics() *RequestMetrics {
	metricsMutex.RLock()
	defer metricsMutex.RUnlock()

	// 創建副本避免競態條件
	result := &RequestMetrics{
		TotalRequests:    metrics.TotalRequests,
		ActiveRequests:   metrics.ActiveRequests,
		RequestsByMethod: make(map[string]int64),
		RequestsByStatus: make(map[string]int64),
		AverageLatency:   metrics.AverageLatency,
		StartTime:        metrics.StartTime,
		Uptime:           time.Since(metrics.StartTime).String(),
	}

	for k, v := range metrics.RequestsByMethod {
		result.RequestsByMethod[k] = v
	}

	for k, v := range metrics.RequestsByStatus {
		result.RequestsByStatus[k] = v
	}

	return result
}
