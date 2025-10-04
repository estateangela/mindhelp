# 資料種子部署指南

## 🎯 問題診斷

您的 API 回應顯示地址資料為空：

```json
{
  "data": {
    "addresses": null,
    "format": "google_maps_ready", 
    "total": 0
  },
  "message": "Google Maps 地址資訊已準備就緒",
  "success": true
}
```

這表示資料庫中沒有地址資料，需要運行資料種子程序。

## 🔧 解決方案

### 方案 1: 在 Render.com 上運行 Seeder

1. **SSH 到您的 Render 服務** (如果支援):
   ```bash
   # 在 Render 容器中執行
   cd /app
   go run cmd/seed/main.go
   ```

2. **或者通過 Render Dashboard**:
   - 進入您的服務設定
   - 添加一個 Build Command 或 Start Command
   - 暫時修改為運行 seeder

### 方案 2: 通過 Docker 本地測試

1. **本地建置並測試**:
   ```bash
   # 建置 Docker 映像
   docker build -t mindhelp-backend .
   
   # 運行 seeder (需要資料庫連接)
   docker run --rm -e DATABASE_URL="your-database-url" mindhelp-backend go run cmd/seed/main.go
   ```

### 方案 3: 修改主程序包含 Seeder

創建一個環境變數來控制是否運行 seeder：

```go
// 在 main.go 中添加
if os.Getenv("RUN_SEEDER") == "true" {
    log.Println("Running seeder...")
    // 運行 seeder 邏輯
}
```

然後在 Render 中設定環境變數 `RUN_SEEDER=true`。

## 📊 資料種子內容

我們的 seeder 會創建以下範例資料：

### 諮商師 (3 筆)
- 王心理師 - 台北市信義區信義路五段7號101大樓
- 李諮商師 - 台北市大安區復興南路一段390號  
- 陳心理師 - 台北市中山區南京東路二段125號

### 諮商所 (4 筆)
- 台北心理健康中心 - 台北市中正區中山南路1號2樓
- 信義諮商所 - 台北市信義區信義路四段1號8樓
- 大安心理診所 - 台北市大安區敦化南路二段216號3樓
- 松山諮商中心 - 台北市松山區八德路四段138號5樓

### 推薦醫師 (4 筆)
- 張精神科醫師 - 台大醫院
- 林心理師 - 榮總醫院
- 黃醫師 - 馬偕醫院
- 吳心理師 - 新光醫院

## 🚀 快速部署步驟

### 選項 A: 修改 Dockerfile 包含 Seeder

1. **更新 Dockerfile**:
   ```dockerfile
   # 在現有 Dockerfile 末尾添加
   COPY cmd/ /app/cmd/
   
   # 添加環境變數檢查
   RUN if [ "$RUN_SEEDER" = "true" ]; then go run cmd/seed/main.go; fi
   ```

2. **在 Render 中設定環境變數**:
   - `RUN_SEEDER=true`

### 選項 B: 創建 Init Container

1. **創建初始化腳本**:
   ```bash
   #!/bin/bash
   echo "Running database seeder..."
   go run cmd/seed/main.go
   echo "Starting main application..."
   exec ./main
   ```

2. **修改 Dockerfile**:
   ```dockerfile
   COPY init.sh /app/init.sh
   RUN chmod +x /app/init.sh
   CMD ["/app/init.sh"]
   ```

## 🔍 驗證步驟

運行 seeder 後，驗證資料是否成功插入：

1. **檢查個別 API**:
   ```bash
   curl https://mindhelp.onrender.com/api/v1/counselors
   curl https://mindhelp.onrender.com/api/v1/counseling-centers
   curl https://mindhelp.onrender.com/api/v1/recommended-doctors
   ```

2. **檢查地址 API**:
   ```bash
   curl https://mindhelp.onrender.com/api/v1/maps/addresses
   curl https://mindhelp.onrender.com/api/v1/maps/google-addresses
   ```

3. **運行診斷腳本**:
   ```powershell
   .\debug_database.ps1
   ```

## 📝 建議的立即行動

1. **立即解決方案** - 在 Render 中：
   - 進入 Service Settings
   - 添加環境變數 `RUN_SEEDER=true`
   - 重新部署服務

2. **長期解決方案** - 設定自動化：
   - 創建資料庫遷移腳本
   - 設定 CI/CD 管道包含 seeder
   - 實現管理員 API 來管理資料

## ⚠️ 注意事項

- Seeder 會檢查重複資料，不會覆蓋現有記錄
- 確保資料庫連接正常
- 生產環境建議先備份資料庫
- CSV 文件如果不存在不會影響範例資料的插入

## 🎯 預期結果

成功運行 seeder 後，您應該看到：

```json
{
  "data": {
    "addresses": [
      {
        "id": "...",
        "name": "王心理師",
        "address": "台北市信義區信義路五段7號101大樓",
        "type": "counselor"
      },
      // ... 更多地址
    ],
    "total": 11,
    "format": "google_maps_ready"
  },
  "success": true
}
```

這樣 Google Maps API 就有資料可以使用了！🎉
