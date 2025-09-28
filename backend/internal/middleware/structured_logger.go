package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LogEntry struct {
	Timestamp      string `json:"timestamp"`
	Level          string `json:"level"`
	Method         string `json:"method"`
	Path           string `json:"path"`
	StatusCode     int    `json:"status_code"`
	ResponseTimeMs int64  `json:"response_time_ms"`
	ClientIP       string `json:"client_ip"`
	UserAgent      string `json:"user_agent"`
	RequestID      string `json:"request_id"`
	RequestSize    int64  `json:"request_size"`
	ResponseSize   int64  `json:"response_size"`
}

// StructuredLogger 創建結構化日誌中間件
func StructuredLogger(logFile string) gin.HandlerFunc {
	var writer io.Writer = os.Stdout

	// 如果指定了日誌文件，同時寫入文件和標準輸出
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			writer = io.MultiWriter(os.Stdout, file)
		}
	}

	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 生成請求 ID
		requestID := uuid.New().String()[:8]

		// 確定日誌級別
		level := "INFO"
		if param.StatusCode >= 400 && param.StatusCode < 500 {
			level = "WARN"
		} else if param.StatusCode >= 500 {
			level = "ERROR"
		}

		entry := LogEntry{
			Timestamp:      param.TimeStamp.Format(time.RFC3339),
			Level:          level,
			Method:         param.Method,
			Path:           param.Path,
			StatusCode:     param.StatusCode,
			ResponseTimeMs: param.Latency.Milliseconds(),
			ClientIP:       param.ClientIP,
			UserAgent:      param.Request.UserAgent(),
			RequestID:      requestID,
			RequestSize:    param.Request.ContentLength,
			ResponseSize:   int64(param.BodySize),
		}

		jsonBytes, err := json.Marshal(entry)
		if err != nil {
			return fmt.Sprintf("Failed to marshal log entry: %v\n", err)
		}

		fmt.Fprint(writer, string(jsonBytes)+"\n")
		return ""
	})
}
