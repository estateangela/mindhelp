# Google Maps API 設定指南

## 🚀 快速開始

### 1. 環境變數設定

在 `.env` 文件中添加以下配置：

```env
# Google Maps API Configuration
GOOGLE_MAPS_API_KEY=your-google-maps-api-key-here
GOOGLE_MAPS_BASE_URL=https://maps.googleapis.com/maps/api
GOOGLE_MAPS_GEOCODING_URL=https://maps.googleapis.com/maps/api/geocode/json
GOOGLE_MAPS_PLACES_URL=https://maps.googleapis.com/maps/api/place
GOOGLE_MAPS_DIRECTIONS_URL=https://maps.googleapis.com/maps/api/directions/json
GOOGLE_MAPS_DISTANCE_MATRIX_URL=https://maps.googleapis.com/maps/api/distancematrix/json
```

### 2. 獲取 Google Maps API Key

1. 前往 [Google Cloud Console](https://console.cloud.google.com/)
2. 創建或選擇一個專案
3. 啟用以下 API：
   - Geocoding API
   - Places API
   - Directions API
   - Distance Matrix API
4. 創建 API 金鑰並設定適當的限制

### 3. Docker 建置修復

如果遇到 `golang.org/x/time/rate` 依賴問題，請確保：

1. `go.mod` 文件包含：
```go
golang.org/x/time v0.8.0
```

2. `go.sum` 文件包含：
```
golang.org/x/time v0.8.0 h1:9i3RxcPv3PZnitoVGMPDKZSq1xW1gK1Xy3ArNOGZfEg=
golang.org/x/time v0.8.0/go.mod h1:3BpzKBy/shNhVucY/MWOyx10tF3SFh9QdLuxbVysPQM=
```

### 4. 本地測試

運行測試腳本：

```powershell
# Windows PowerShell
.\test_google_maps_api.ps1

# 或者手動測試編譯
go run build_test.go
```

### 5. Docker 建置

```bash
# 標準建置
docker build -t mindhelp-backend .

# 如果仍有依賴問題，可以嘗試清理建置
docker build --no-cache -t mindhelp-backend .
```

## 🔧 故障排除

### 問題 1: 缺少 golang.org/x/time/rate 依賴

**解決方案：**
1. 確保 `go.mod` 和 `go.sum` 文件已正確更新
2. 如果有 Go 環境，運行 `go mod tidy`
3. 重新建置 Docker 映像

### 問題 2: API Key 未設定

**症狀：** API 回應 "Google Maps API Key 未設定"

**解決方案：**
1. 檢查 `.env` 文件中的 `GOOGLE_MAPS_API_KEY`
2. 確保環境變數正確載入

### 問題 3: API 配額超限

**症狀：** API 回應 "OVER_QUERY_LIMIT"

**解決方案：**
1. 檢查 Google Cloud Console 中的配額使用情況
2. 考慮啟用計費以獲得更高配額
3. 實現更積極的快取策略

### 問題 4: 速率限制

**症狀：** API 回應 "Too Many Requests"

**解決方案：**
1. 調整 `services/google_maps_service.go` 中的速率限制設定
2. 增加快取時間以減少 API 呼叫

## 📊 監控和優化

### API 使用統計

訪問 `/api/v1/google-maps/usage-stats` 來查看：
- 快取命中率
- API 配置狀態
- 可用端點列表

### 快取管理

- 地理編碼結果快取 1 小時
- 反向地理編碼結果快取 24 小時
- 附近搜尋結果快取 30 分鐘
- 路線規劃結果快取 15 分鐘

手動清除快取：`POST /api/v1/google-maps/clear-cache`

### 效能調優

1. **批次處理：** 使用 `/batch-geocode` 端點處理多個地址
2. **並發控制：** 預設最多 5 個並發請求
3. **超時設定：** 30 秒請求超時
4. **速率限制：** 每秒最多 10 個請求，突發最多 20 個

## 🔐 安全考量

1. **API Key 保護：** 永不將 API Key 提交到版本控制
2. **域名限制：** 在 Google Cloud Console 中設定 API Key 的使用限制
3. **配額監控：** 定期檢查 API 使用量避免意外費用
4. **錯誤處理：** 適當處理 API 錯誤避免洩露敏感資訊

## 🚀 部署檢查清單

- [ ] 設定 Google Maps API Key
- [ ] 啟用必要的 Google Maps API
- [ ] 配置環境變數
- [ ] 測試 API 連接
- [ ] 驗證快取功能
- [ ] 檢查速率限制設定
- [ ] 監控 API 使用量
- [ ] 設定錯誤警報

## 📞 支援

如果遇到問題，請檢查：
1. Google Cloud Console 中的 API 狀態
2. 應用程式日誌中的錯誤訊息
3. 網路連接和防火牆設定
4. API Key 的權限和限制設定
