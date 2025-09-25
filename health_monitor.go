package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// HealthCheckResponse 健康檢查回應結構
type HealthCheckResponse struct {
	Status    string                 `json:"status"`
	Service   string                 `json:"service"`
	Version   string                 `json:"version"`
	Timestamp string                 `json:"timestamp"`
	Uptime    string                 `json:"uptime"`
	Checks    map[string]string      `json:"checks"`
	System    map[string]interface{} `json:"system"`
}

// MonitorConfig 監控配置
type MonitorConfig struct {
	URL          string        `json:"url"`
	Interval     time.Duration `json:"interval"`
	Timeout      time.Duration `json:"timeout"`
	LogFile      string        `json:"log_file"`
	AlertOnError bool          `json:"alert_on_error"`
	MaxRetries   int           `json:"max_retries"`
}

// HealthMonitor 健康監控器
type HealthMonitor struct {
	config     MonitorConfig
	httpClient *http.Client
	logger     *log.Logger
	logFile    *os.File
}

// NewHealthMonitor 創建新的健康監控器
func NewHealthMonitor(config MonitorConfig) (*HealthMonitor, error) {
	// 開啟日誌檔案
	var logFile *os.File
	var err error
	if config.LogFile != "" {
		logFile, err = os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("無法開啟日誌檔案: %w", err)
		}
	}

	// 建立 logger
	logger := log.New(logFile, "[HEALTH_MONITOR] ", log.LstdFlags|log.Lshortfile)

	// 建立 HTTP 客戶端
	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	return &HealthMonitor{
		config:     config,
		httpClient: httpClient,
		logger:     logger,
		logFile:    logFile,
	}, nil
}

// Close 關閉監控器
func (hm *HealthMonitor) Close() {
	if hm.logFile != nil {
		hm.logFile.Close()
	}
}

// CheckHealth 執行健康檢查
func (hm *HealthMonitor) CheckHealth() (*HealthCheckResponse, error) {
	resp, err := hm.httpClient.Get(hm.config.URL)
	if err != nil {
		return nil, fmt.Errorf("健康檢查請求失敗: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("讀取回應失敗: %w", err)
	}

	var healthResp HealthCheckResponse
	if err := json.Unmarshal(body, &healthResp); err != nil {
		return nil, fmt.Errorf("解析 JSON 失敗: %w", err)
	}

	return &healthResp, nil
}

// LogHealthStatus 記錄健康狀態
func (hm *HealthMonitor) LogHealthStatus(health *HealthCheckResponse, err error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	if err != nil {
		hm.logger.Printf("[ERROR] %s - 健康檢查失敗: %v", timestamp, err)
		if hm.config.AlertOnError {
			fmt.Printf("🚨 [%s] 健康檢查失敗: %v\n", timestamp, err)
		}
		return
	}

	// 記錄基本狀態
	status := "✅"
	if health.Status != "ok" {
		status = "⚠️"
	}

	hm.logger.Printf("[%s] %s 服務狀態: %s, 運行時間: %s",
		timestamp, status, health.Status, health.Uptime)

	// 記錄詳細檢查結果
	for checkName, checkStatus := range health.Checks {
		checkIcon := "✅"
		if checkStatus != "healthy" {
			checkIcon = "❌"
		}
		hm.logger.Printf("[%s] %s %s: %s", timestamp, checkIcon, checkName, checkStatus)
	}

	// 記錄系統資源
	if system, ok := health.System["goroutines"].(float64); ok {
		hm.logger.Printf("[%s] 📊 Goroutines: %.0f", timestamp, system)
	}
	if memory, ok := health.System["memory_alloc"].(float64); ok {
		hm.logger.Printf("[%s] 📊 記憶體使用: %.0f MB", timestamp, memory)
	}

	// 控制台輸出
	fmt.Printf("%s [%s] %s 狀態: %s | 運行: %s\n",
		status, timestamp, health.Service, health.Status, health.Uptime)
}

// Start 開始監控
func (hm *HealthMonitor) Start() {
	hm.logger.Printf("開始健康監控 - URL: %s, 間隔: %v", hm.config.URL, hm.config.Interval)
	fmt.Printf("🏥 開始健康監控服務\n")
	fmt.Printf("📍 監控端點: %s\n", hm.config.URL)
	fmt.Printf("⏰ 檢查間隔: %v\n", hm.config.Interval)
	fmt.Printf("📝 日誌檔案: %s\n", hm.config.LogFile)
	fmt.Printf("%s\n", strings.Repeat("=", 50))

	ticker := time.NewTicker(hm.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			health, err := hm.CheckHealth()
			hm.LogHealthStatus(health, err)
		}
	}
}

func main() {
	// 預設配置
	config := MonitorConfig{
		URL:          "https://mind-map-api.estateangela.dpdns.org/health",
		Interval:     30 * time.Second,
		Timeout:      10 * time.Second,
		LogFile:      "health_monitor.log",
		AlertOnError: true,
		MaxRetries:   3,
	}

	// 從環境變數讀取配置
	if url := os.Getenv("HEALTH_CHECK_URL"); url != "" {
		config.URL = url
	}
	if interval := os.Getenv("HEALTH_CHECK_INTERVAL"); interval != "" {
		if duration, err := time.ParseDuration(interval); err == nil {
			config.Interval = duration
		}
	}
	if logFile := os.Getenv("HEALTH_LOG_FILE"); logFile != "" {
		config.LogFile = logFile
	}

	// 建立監控器
	monitor, err := NewHealthMonitor(config)
	if err != nil {
		log.Fatalf("建立健康監控器失敗: %v", err)
	}
	defer monitor.Close()

	// 開始監控
	monitor.Start()
}
