# 健康監控腳本

## 使用方法

### 1. 基本使用
```bash
go run health_monitor.go
```

### 2. 自訂配置
```bash
# 設定監控端點
export HEALTH_CHECK_URL="https://your-api.com/health"

# 設定檢查間隔（預設30秒）
export HEALTH_CHECK_INTERVAL="30s"

# 設定日誌檔案
export HEALTH_LOG_FILE="health_monitor.log"

go run health_monitor.go
```

### 3. 背景執行
```bash
# Windows PowerShell
Start-Process -FilePath "go" -ArgumentList "run", "health_monitor.go" -WindowStyle Hidden

# Linux/Mac
nohup go run health_monitor.go > health_monitor.out 2>&1 &
```

## 功能特色

- ✅ 每30秒自動檢查健康狀態
- 📝 詳細日誌記錄
- 🚨 異常警報
- 📊 系統資源監控
- ⚙️ 環境變數配置
- 🔄 自動重試機制

## 監控項目

- 服務狀態 (ok/degraded)
- 資料庫連接狀態
- 系統資源使用量
- Goroutines 數量
- 記憶體使用量
- 服務運行時間

## 日誌格式

```
[HEALTH_MONITOR] 2025-01-25 13:30:00 [INFO] ✅ 服務狀態: ok, 運行時間: 2h30m15s
[HEALTH_MONITOR] 2025-01-25 13:30:00 [INFO] ✅ database: healthy
[HEALTH_MONITOR] 2025-01-25 13:30:00 [INFO] 📊 Goroutines: 25
[HEALTH_MONITOR] 2025-01-25 13:30:00 [INFO] 📊 記憶體使用: 45 MB
```

## 異常處理

當健康檢查失敗時：
- 記錄錯誤日誌
- 顯示警報訊息
- 繼續下次檢查
- 支援重試機制
