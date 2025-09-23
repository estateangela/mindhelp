package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// GoogleMapsUsage Google Maps API 使用統計
type GoogleMapsUsage struct {
	Endpoint    string    `json:"endpoint"`
	Method      string    `json:"method"`
	StatusCode  int       `json:"status_code"`
	Duration    int64     `json:"duration_ms"`
	RequestSize int64     `json:"request_size_bytes"`
	ResponseSize int64    `json:"response_size_bytes"`
	Timestamp   time.Time `json:"timestamp"`
	UserID      string    `json:"user_id,omitempty"`
	ClientIP    string    `json:"client_ip"`
	UserAgent   string    `json:"user_agent"`
	Error       string    `json:"error,omitempty"`
}

// GoogleMapsMetricsMiddleware Google Maps API 指標中間件
func GoogleMapsMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只對 Google Maps API 路由生效
		if !isGoogleMapsAPI(c.Request.URL.Path) {
			c.Next()
			return
		}

		startTime := time.Now()
		
		// 記錄請求體大小
		var requestSize int64
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestSize = int64(len(bodyBytes))
			// 重新設置請求體以供後續使用
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 創建自定義 ResponseWriter 來捕獲回應
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 處理請求
		c.Next()

		// 計算處理時間
		duration := time.Since(startTime)
		
		// 獲取使用者資訊（如果有認證）
		var userID string
		if uid, exists := c.Get("user_id"); exists {
			if uidStr, ok := uid.(string); ok {
				userID = uidStr
			}
		}

		// 記錄使用統計
		usage := GoogleMapsUsage{
			Endpoint:     c.Request.URL.Path,
			Method:       c.Request.Method,
			StatusCode:   c.Writer.Status(),
			Duration:     duration.Milliseconds(),
			RequestSize:  requestSize,
			ResponseSize: int64(blw.body.Len()),
			Timestamp:    startTime,
			UserID:       userID,
			ClientIP:     c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
		}

		// 如果有錯誤，記錄錯誤資訊
		if len(c.Errors) > 0 {
			usage.Error = c.Errors.String()
		}

		// 記錄到日誌
		logGoogleMapsUsage(usage)
	}
}

// bodyLogWriter 自定義 ResponseWriter 來捕獲回應內容
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// isGoogleMapsAPI 檢查是否為 Google Maps API 路由
func isGoogleMapsAPI(path string) bool {
	return len(path) > 12 && path[:13] == "/api/v1/google-maps"
}

// logGoogleMapsUsage 記錄 Google Maps API 使用情況
func logGoogleMapsUsage(usage GoogleMapsUsage) {
	// 將使用統計轉換為 JSON
	usageJSON, err := json.Marshal(usage)
	if err != nil {
		log.Printf("Failed to marshal Google Maps usage: %v", err)
		return
	}

	// 記錄到日誌（可以擴展為發送到監控系統）
	log.Printf("GOOGLE_MAPS_API_USAGE: %s", string(usageJSON))

	// 如果回應時間過長，記錄警告
	if usage.Duration > 5000 { // 超過 5 秒
		log.Printf("SLOW_GOOGLE_MAPS_API: %s took %dms", usage.Endpoint, usage.Duration)
	}

	// 如果有錯誤，記錄錯誤
	if usage.Error != "" {
		log.Printf("GOOGLE_MAPS_API_ERROR: %s - %s", usage.Endpoint, usage.Error)
	}
}

// GoogleMapsRateLimitMiddleware Google Maps API 速率限制中間件
func GoogleMapsRateLimitMiddleware() gin.HandlerFunc {
	// 簡單的記憶體速率限制器
	requestCounts := make(map[string][]time.Time)
	
	return func(c *gin.Context) {
		// 只對 Google Maps API 路由生效
		if !isGoogleMapsAPI(c.Request.URL.Path) {
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		now := time.Now()
		
		// 清理過期的請求記錄（1 分鐘前的）
		if requests, exists := requestCounts[clientIP]; exists {
			var validRequests []time.Time
			for _, reqTime := range requests {
				if now.Sub(reqTime) < time.Minute {
					validRequests = append(validRequests, reqTime)
				}
			}
			requestCounts[clientIP] = validRequests
		}

		// 檢查是否超過速率限制（每分鐘最多 60 個請求）
		if len(requestCounts[clientIP]) >= 60 {
			c.JSON(429, gin.H{
				"error": "Too Many Requests",
				"message": "API rate limit exceeded. Maximum 60 requests per minute.",
				"retry_after": 60,
			})
			c.Abort()
			return
		}

		// 記錄此次請求
		requestCounts[clientIP] = append(requestCounts[clientIP], now)
		
		c.Next()
	}
}

// GoogleMapsAPIKeyValidationMiddleware API Key 驗證中間件
func GoogleMapsAPIKeyValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只對 Google Maps API 路由生效
		if !isGoogleMapsAPI(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 檢查是否有 Google Maps API Key 配置
		// 這裡可以從配置或環境變數中檢查
		// 實際實現時應該從 config 中獲取
		
		c.Next()
	}
}
