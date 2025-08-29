# MindHelp Backend

基於 Go + Gin 的心理健康支援應用程式後端服務

## 專案概述

MindHelp Backend 是一個使用 Go 語言開發的 RESTful API 服務，提供心理健康支援應用程式的後端功能，包括使用者認證、AI 聊天支援、位置服務等。

## 技術架構

- **語言**: Go 1.21+
- **Web 框架**: Gin
- **ORM**: GORM
- **資料庫**: PostgreSQL (Supabase)
- **認證**: JWT
- **API 文檔**: Swagger
- **容器化**: Docker
- **部署**: Render

## 專案結構

```
backend/
├── internal/
│   ├── config/          # 配置管理
│   ├── database/        # 資料庫連接和遷移
│   ├── dto/            # 資料傳輸物件 (Request)
│   ├── handlers/       # HTTP 處理器
│   ├── middleware/     # 中間件
│   ├── models/         # 資料模型
│   ├── routes/         # 路由配置
│   └── vo/             # 視圖物件 (Response)
├── .env.example        # 環境變數範例
├── docker-compose.yml  # Docker Compose 配置
├── Dockerfile          # Docker 映像配置
├── go.mod              # Go 模組依賴
└── main.go             # 主程式入口
```

## 主要功能

### 🔐 認證系統
- 使用者註冊/登入
- JWT token 管理
- 密碼加密 (bcrypt)
- Token 刷新機制

### 💬 AI 聊天
- 與 OpenRouter API 整合
- 聊天記錄儲存
- 支援多種 AI 模型
- Token 使用統計

### 🗺️ 位置服務
- 心理健康資源位置管理
- 地理位置搜尋
- 距離計算
- 公開/私有位置控制

### 📊 資料管理
- 使用者資料管理
- 聊天歷史記錄
- 位置資訊 CRUD
- 軟刪除支援

## 快速開始

### 前置需求

- Go 1.21 或更高版本
- PostgreSQL 資料庫
- Docker 和 Docker Compose (可選)

### 本地開發

1. **複製專案**
   ```bash
   git clone <repository-url>
   cd backend
   ```

2. **安裝依賴**
   ```bash
   go mod download
   ```

3. **設定環境變數**
   ```bash
   cp env.example .env
   # 編輯 .env 文件，設定資料庫和 API 金鑰
   ```

4. **啟動資料庫**
   ```bash
   docker-compose up db -d
   ```

5. **執行應用程式**
   ```bash
   go run main.go
   ```

### 使用 Docker

1. **構建和啟動**
   ```bash
   docker-compose up --build
   ```

2. **僅啟動資料庫**
   ```bash
   docker-compose up db -d
   ```

## API 端點

### 認證端點
- `POST /api/v1/auth/register` - 使用者註冊
- `POST /api/v1/auth/login` - 使用者登入
- `POST /api/v1/auth/refresh` - 刷新 token

### 聊天端點
- `POST /api/v1/chat/send` - 發送聊天訊息
- `GET /api/v1/chat/history` - 獲取聊天歷史

### 位置端點
- `POST /api/v1/locations` - 創建位置
- `GET /api/v1/locations/search` - 搜尋位置
- `GET /api/v1/locations/:id` - 獲取位置詳情
- `PUT /api/v1/locations/:id` - 更新位置
- `DELETE /api/v1/locations/:id` - 刪除位置

### 健康檢查
- `GET /health` - 服務健康狀態
- `GET /swagger/*` - API 文檔

## 環境變數

| 變數名稱 | 說明 | 預設值 |
|---------|------|--------|
| `PORT` | 服務端口 | `8080` |
| `GIN_MODE` | Gin 模式 | `release` |
| `DB_HOST` | 資料庫主機 | `localhost` |
| `DB_PORT` | 資料庫端口 | `5432` |
| `DB_USER` | 資料庫使用者 | `postgres` |
| `DB_PASSWORD` | 資料庫密碼 | - |
| `DB_NAME` | 資料庫名稱 | `mindhelp` |
| `DB_SSL_MODE` | SSL 模式 | `disable` |
| `JWT_SECRET` | JWT 密鑰 | - |
| `JWT_EXPIRY` | JWT 過期時間 | `24h` |
| `OPENROUTER_API_KEY` | OpenRouter API 金鑰 | - |
| `ALLOWED_ORIGINS` | 允許的 CORS 來源 | `http://localhost:3000` |

## 部署到 Render

### 1. 準備 Supabase 資料庫
- 在 Supabase 創建新專案
- 獲取資料庫連接資訊
- 設定環境變數

### 2. 在 Render 部署
- 連接 GitHub 專案
- 選擇 Go 環境
- 設定環境變數
- 設定構建指令：`go build -o main .`
- 設定啟動指令：`./main`

### 3. 環境變數設定
```bash
PORT=10000
GIN_MODE=release
DB_HOST=your-project.supabase.co
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=postgres
DB_SSL_MODE=require
JWT_SECRET=your-super-secret-jwt-key
OPENROUTER_API_KEY=your-openrouter-api-key
ALLOWED_ORIGINS=https://yourdomain.com
```

## 開發指南

### 新增 API 端點

1. 在 `internal/dto/` 定義請求/回應結構
2. 在 `internal/handlers/` 實作處理邏輯
3. 在 `internal/routes/` 註冊路由
4. 添加 Swagger 註解

### 資料庫遷移

使用 GORM 自動遷移：
```go
db.AutoMigrate(&models.User{}, &models.ChatMessage{}, &models.Location{})
```

### 測試

```bash
# 執行所有測試
go test ./...

# 執行特定包的測試
go test ./internal/handlers

# 執行測試並顯示覆蓋率
go test -cover ./...
```

## 安全考量

- 使用 bcrypt 加密密碼
- JWT token 驗證
- CORS 配置
- 輸入驗證和清理
- SQL 注入防護 (GORM)
- 環境變數管理

## 監控和日誌

- 健康檢查端點
- 結構化日誌
- 錯誤處理和回應
- 請求/回應記錄

## 貢獻指南

1. Fork 專案
2. 創建功能分支
3. 實作功能
4. 添加測試
5. 提交 Pull Request

## 授權

此專案為私人專案，不適用於公開發行。
