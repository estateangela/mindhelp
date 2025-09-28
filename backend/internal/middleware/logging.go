package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// LogEntry 結構化日誌條目
type LogEntry struct {
	Timestamp    string `json:"timestamp"`
	Level        string `json:"level"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	StatusCode   int    `json:"status_code"`
	ResponseTime int64  `json:"response_time_ms"`
	IP           string `json:"client_ip"`
	UserAgent    string `json:"user_agent"`
	RequestID    string `json:"request_id"`
	RequestSize  int    `json:"request_size"`
	ResponseSize int    `json:"response_size"`
	Error        string `json:"error,omitempty"`
}

// StructuredLogger 結構化日誌中間件
func StructuredLogger(logFile string) gin.HandlerFunc {
	// 設定日誌輸出
	var writer io.Writer = os.Stdout
	
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("Failed to open log file %s: %v", logFile, err)
		} else {
			writer = io.MultiWriter(os.Stdout, file)
		}
	}

	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			logEntry := LogEntry{
				Timestamp:    param.TimeStamp.Format(time.RFC3339),
				Level:        "INFO",
				Method:       param.Method,
				Path:         param.Path,
				StatusCode:   param.StatusCode,
				ResponseTime: param.Latency.Milliseconds(),
				IP:           param.ClientIP,
				UserAgent:    param.Request.UserAgent(),
				RequestID:    param.Request.Header.Get("X-Request-ID"),
				RequestSize:  int(param.Request.ContentLength),
				ResponseSize: param.BodySize,
			}

			// 如果有錯誤，記錄錯誤級別
			if param.ErrorMessage != "" {
				logEntry.Level = "ERROR"
				logEntry.Error = param.ErrorMessage
			} else if param.StatusCode >= 400 {
				logEntry.Level = "WARN"
			}

			logJSON, _ := json.Marshal(logEntry)
			return string(logJSON) + "\n"
		},
		Output: writer,
	})
}

// MetricsMiddleware 效能指標中間件
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// 處理請求
		c.Next()
		
		// 計算響應時間
		latency := time.Since(start)
		
		// 記錄到全局指標
		RecordMetrics(MetricsData{
			Method:       c.Request.Method,
			Path:         c.FullPath(),
			StatusCode:   c.Writer.Status(),
			ResponseTime: latency,
			RequestSize:  c.Request.ContentLength,
			ResponseSize: int64(c.Writer.Size()),
		})
	}
}

// MetricsData 指標資料
type MetricsData struct {
	Method       string
	Path         string
	StatusCode   int
	ResponseTime time.Duration
	RequestSize  int64
	ResponseSize int64
}

// 全局指標儲存
var (
	requestCount    = make(map[string]int64)
	responseTime    = make(map[string][]time.Duration)
	errorCount      = make(map[string]int64)
	totalRequests   int64
	totalErrors     int64
)

// RecordMetrics 記錄指標
func RecordMetrics(data MetricsData) {
	key := fmt.Sprintf("%s %s", data.Method, data.Path)
	
	// 增加請求計數
	requestCount[key]++
	totalRequests++
	
	// 記錄響應時間
	responseTime[key] = append(responseTime[key], data.ResponseTime)
	
	// 記錄錯誤
	if data.StatusCode >= 400 {
		errorCount[key]++
		totalErrors++
	}
}

// GetMetrics 獲取指標
func GetMetrics() map[string]interface{} {
	return map[string]interface{}{
		"total_requests": totalRequests,
		"total_errors":   totalErrors,
		"request_count":  requestCount,
		"error_count":    errorCount,
		"timestamp":      time.Now().Format(time.RFC3339),
	}
}
