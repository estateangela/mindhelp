
# MindHelp - 心理健康支援平台

<div align="center">

![MindHelp Logo](my_mindhelp_app/assets/images/logo.png)

**讓心理健康支援更貼近每個人** 🧠💚

[![Flutter](https://img.shields.io/badge/Flutter-3.6.2+-blue.svg)](https://flutter.dev/)
[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8.svg)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-316192.svg)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg)](https://www.docker.com/)

</div>

## 📋 專案概述

MindHelp 是一個前後端一體的心理健康支援平台：
- 前端使用 Flutter 打造跨平台行動體驗。
- 後端使用 Go (Gin) 提供高效的 REST API，整合 PostgreSQL（Supabase）與 Google Maps 服務。

### 🌟 功能亮點
- **📚 專家文章**：心理健康知識內容與收藏
- **🧠 心理測驗**：GAD-7量表與結果解釋
- **🗺️ 資源地圖**：諮商師、諮商所、推薦醫師位置檢索
- **💬 AI 聊天**：情緒支持與歷史會話
- **🔔 通知**：訊息與偏好設定管理

## 🏗️ 專案結構

```
mindhelp/
├── backend/                    # Go 後端 API 服務
│   ├── internal/
│   │   ├── config/             # 設定載入與安全
│   │   ├── database/           # 資料庫連線與遷移
│   │   ├── dto/                # 請求/回應傳輸物件
│   │   ├── handlers/           # 業務處理器 (22 個檔案)
│   │   ├── middleware/         # 認證、日誌、CORS
│   │   ├── models/             # GORM 資料模型
│   │   ├── routes/             # 路由註冊 (Gin)
│   │   ├── scheduler/          # 定時任務 (通知等)
│   │   └── vo/                 # 視圖物件
│   ├── database/migrations/    # SQL 遷移
│   ├── docs/                   # Swagger/OpenAPI 文檔
│   └── document/               # 專案說明與資料
├── my_mindhelp_app/            # Flutter 應用程式
│   ├── lib/
│   │   ├── core/               # 主題/常數
│   │   ├── models/             # 資料模型
│   │   ├── pages/              # App 頁面
│   │   ├── services/           # 業務服務
│   │   └── widgets/            # 共用元件
│   └── README.md               # 前端細節使用說明
└── Docs/                       # 額外文檔與 Swagger 匯出
```

## 🚀 快速開始

### 後端 (Go)
1) 準備環境與依賴
```bash
cd backend
go mod tidy
cp env.example .env   # 請依需求編輯 .env
```

2) 啟動服務（自動在背景嘗試連線資料庫並執行遷移）
```bash
go run main.go
# 或使用 Makefile
make run
```

- 預設啟動於: http://localhost:8080
- 健康檢查: `/health`, `/health/ready`, `/health/live`, `/health/detailed`, `/metrics`
- Swagger: `/swagger/index.html`

3) 使用 Docker Compose 本機一鍵啟動
```bash
cd backend
docker-compose up -d         # 啟動 app + PostgreSQL
docker-compose logs -f app   # 觀察服務日誌
```

備註：`docker-compose.yml` 會自動掛載 `database/migrations` 至資料庫容器初始化，應用也提供 GORM AutoMigrate 與額外修正（例如 `recommended_doctors` 空名稱修復）。

### 前端 (Flutter)
```bash
cd my_mindhelp_app
flutter pub get
flutter run
```

更多平台建置/疑難排解請見 `my_mindhelp_app/README.md`。

## 🔧 常用開發腳本 (Makefile)
於 `backend/` 目錄：

```bash
make dev-setup     # 複製 env 並提示下一步
make db-up         # 啟動 Postgres (Docker)
make run           # 啟動後端
make test          # 執行測試
make swagger       # 產生 Swagger 文檔
make docker-run    # 以 docker-compose 啟動
make docker-stop   # 停止 docker-compose 服務
```

## 🔐 環境變數
請參考 `backend/env.example`，常見設定：
- `PORT`：服務監聽埠（Render 會注入 `PORT`）
- `GIN_MODE`：`release` / `debug`
- `DATABASE_URL`：完整連線字串（建議，Supabase 友善）
  - 或使用 `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSL_MODE`
- `JWT_SECRET`, `JWT_EXPIRY`
- `OPENROUTER_API_KEY`, `OPENROUTER_BASE_URL`
- `GOOGLE_MAPS_API_KEY` 及相關 `GOOGLE_MAPS_*`
- `ALLOWED_ORIGINS` / `CORS_ALLOWED_ORIGINS`
- `LOG_LEVEL`, `LOG_FORMAT`

## 📚 API 概覽
Base Path: `/api/v1`

### 公開端點
- 文章：`GET /articles`, `GET /articles/:id`
- 測驗：`GET /quizzes`, `GET /quizzes/:id`
- 應用配置：`GET /config`
- 評論（查詢）：`GET /resources/:id/reviews`
- 專業資源：
  - `GET /counselors`, `GET /counselors/:id`
  - `GET /counseling-centers`, `GET /counseling-centers/:id`
  - `GET /recommended-doctors`, `GET /recommended-doctors/:id`
- 地圖地址：`GET /maps/addresses`
- Google Maps 代理：`/google-maps/*`（如 `POST /google-maps/geocode`, `POST /google-maps/search-places` 等）

### 認證與受保護端點
- 認證：`POST /auth/register`, `POST /auth/login`, `POST /auth/refresh`
- 使用者：`GET/PUT /users/me`, `PUT /users/me/password`, `DELETE /users/me`, `GET /users/me/stats`
- 聊天：`GET/POST /chat/sessions`, `GET/POST /chat/sessions/:sessionId/messages`（亦保留舊版 `/chat/send`, `/chat/history`）
- 測驗提交：`POST /quizzes/:id/submit`; 歷史：`GET /users/me/quiz_history`
- 收藏：`GET /users/me/bookmarks/articles`, `GET /users/me/bookmarks/resources`, `POST /bookmarks`, `DELETE /bookmarks`
- 評論：`POST /resources/:id/reviews`, `PUT /reviews/:reviewId`, `DELETE /reviews/:reviewId`, `POST /report`
- 通知：`GET /notifications`, `POST /notifications/mark-as-read`, `GET/PUT /users/me/notification-settings`, `POST /users/me/push-token`
- 分享：`POST /shares`, `GET /users/me/shares`; 公開查閱：`GET /shares/:shareId`, `GET /shares/stats`

健康檢查與文檔：
- `GET /`（根）與 `/health*` 系列
- Swagger UI：`/swagger/index.html`

## 🧱 資料庫與遷移
- SQL 遷移：`backend/database/migrations/*.sql`（docker-compose 初始化掛載）
- 自動遷移：啟動時執行 GORM AutoMigrate；並含少量資料修正（例如 `recommended_doctors` 空名稱補值）

## ☁️ 部署建議
- Docker：
  ```bash
  cd backend
  docker build -t mindhelp-backend .
  docker run -p 8080:8080 --env-file .env mindhelp-backend
  ```
- Docker Compose：`docker-compose up -d`
- Render：請參考 `backend/RENDER_DEPLOYMENT.md` 與 `backend/RENDER_CORS_FIX.md`（若有）並確保設定 `PORT` 與資料庫連線。

## 🤝 貢獻
1. Fork 專案
2. 建立分支：`git checkout -b feature/your-feature`
3. 提交變更：`git commit -m "feat: add your feature"`
4. 推送分支：`git push origin feature/your-feature`
5. 送出 Pull Request

問題與建議請至 Issues 回報。

---

<div align="center">

**MindHelp** - 讓心理健康支援更貼近每個人 🧠💚

</div>
