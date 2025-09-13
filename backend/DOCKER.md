# 🐳 MindHelp Backend Docker 部署指南

本文檔說明如何使用 Docker 來建置、測試和部署 MindHelp 後端服務。

## 📋 前置需求

- [Docker](https://www.docker.com/get-started) 20.10+
- [Docker Compose](https://docs.docker.com/compose/) 2.0+

## 🚀 快速開始

### 1. 生產環境部署

```bash
# 構建並啟動生產環境
docker-compose up -d

# 查看服務狀態
docker-compose ps

# 查看日誌
docker-compose logs -f app
```

訪問服務：
- **API**: http://localhost:8080
- **健康檢查**: http://localhost:8080/health
- **API 文檔**: http://localhost:8080/swagger/index.html

### 2. 開發環境部署

```bash
# 使用開發環境配置
docker-compose -f docker-compose.dev.yml up -d

# 查看所有服務
docker-compose -f docker-compose.dev.yml ps
```

開發環境包含額外服務：
- **API**: http://localhost:8080
- **PostgreSQL**: localhost:5433
- **Redis**: localhost:6379  
- **pgAdmin**: http://localhost:5050 (admin@mindhelp.dev / admin123)

## 🛠️ 構建腳本

### Linux/macOS
```bash
# 使用 Shell 腳本構建
./docker-build.sh [tag]

# 範例：構建特定版本
./docker-build.sh v1.0.0
```

### Windows
```batch
# 使用批次檔構建
docker-build.bat [tag]

# 範例：構建並測試
docker-build.bat latest
```

## ⚙️ 環境變數配置

### 必要環境變數

| 變數名稱 | 說明 | 預設值 | 必要 |
|---------|------|-------|------|
| `SERVER_PORT` | 服務端口 | `8080` | ✅ |
| `DB_HOST` | 資料庫主機 | `localhost` | ✅ |
| `DB_PORT` | 資料庫端口 | `5432` | ✅ |
| `DB_USER` | 資料庫使用者 | - | ✅ |
| `DB_PASSWORD` | 資料庫密碼 | - | ✅ |
| `DB_NAME` | 資料庫名稱 | - | ✅ |
| `JWT_SECRET` | JWT 密鑰 | - | ✅ |

### 可選環境變數

| 變數名稱 | 說明 | 預設值 |
|---------|------|-------|
| `GIN_MODE` | Gin 執行模式 | `release` |
| `DB_SSLMODE` | 資料庫 SSL 模式 | `require` |
| `OPENROUTER_API_KEY` | OpenRouter API 密鑰 | - |
| `CORS_ALLOWED_ORIGINS` | CORS 允許來源 | - |
| `LOG_LEVEL` | 日誌等級 | `info` |

## 📊 健康檢查

服務提供多個健康檢查端點：

```bash
# 基本健康檢查
curl http://localhost:8080/health

# 詳細健康檢查
curl http://localhost:8080/health/detailed

# 就緒檢查
curl http://localhost:8080/health/ready

# 存活檢查
curl http://localhost:8080/health/live
```

## 🐛 故障排除

### 容器無法啟動

```bash
# 查看容器日誌
docker-compose logs app

# 檢查容器狀態
docker-compose ps

# 重新構建映像
docker-compose build --no-cache app
```

### 資料庫連接失敗

```bash
# 檢查資料庫容器狀態
docker-compose logs db

# 測試資料庫連接
docker-compose exec db pg_isready -U postgres -d mindhelp

# 重新啟動資料庫服務
docker-compose restart db
```

### 端口衝突

如果遇到端口衝突，可以修改 `docker-compose.yml` 中的端口映射：

```yaml
services:
  app:
    ports:
      - "8081:8080"  # 使用不同的本地端口
```

## 🔧 開發工作流程

### 1. 開發環境設置

```bash
# 1. 啟動開發環境
docker-compose -f docker-compose.dev.yml up -d

# 2. 查看 API 文檔
open http://localhost:8080/swagger/index.html

# 3. 使用 pgAdmin 管理資料庫
open http://localhost:5050
```

### 2. 代碼變更後重新部署

```bash
# 重新構建並啟動
docker-compose -f docker-compose.dev.yml up -d --build

# 或者只重新構建特定服務
docker-compose -f docker-compose.dev.yml build app
docker-compose -f docker-compose.dev.yml restart app
```

### 3. 調試和日誌

```bash
# 實時查看日誌
docker-compose -f docker-compose.dev.yml logs -f app

# 進入容器調試
docker-compose -f docker-compose.dev.yml exec app sh

# 執行資料庫命令
docker-compose -f docker-compose.dev.yml exec db psql -U mindhelp_dev -d mindhelp_dev
```

## 🚢 生產部署最佳實踐

### 1. 環境變數管理

生產環境不要將敏感資料寫在 docker-compose.yml 中：

```bash
# 創建 .env 檔案
cp .env.production.example .env

# 編輯敏感資料
vim .env

# 使用環境變數檔案啟動
docker-compose --env-file .env up -d
```

### 2. 資料持久化

確保資料庫資料持久化：

```bash
# 檢查 volume 狀態
docker volume ls | grep mindhelp

# 備份資料庫
docker-compose exec db pg_dump -U postgres mindhelp > backup.sql

# 恢復資料庫
cat backup.sql | docker-compose exec -T db psql -U postgres -d mindhelp
```

### 3. 監控和日誌

```bash
# 設置日誌輪轉
# 在 docker-compose.yml 中添加：
services:
  app:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## 📈 效能優化

### 1. 多階段構建

Dockerfile 已經使用多階段構建來減少映像大小：

```dockerfile
# 構建階段
FROM golang:1.21-alpine AS builder
# ... 構建代碼

# 運行階段
FROM alpine:latest
# ... 只複製必要文件
```

### 2. 資源限制

在生產環境中設置資源限制：

```yaml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '0.50'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
```

## 🔐 安全注意事項

1. **更改預設密碼**：生產環境務必更改所有預設密碼
2. **JWT 密鑰**：使用強密碼作為 JWT_SECRET
3. **資料庫連接**：生產環境啟用 SSL (`DB_SSLMODE=require`)
4. **容器使用者**：不要使用 root 使用者運行容器
5. **網路隔離**：使用自定義網路隔離服務

## 📞 支援

如果遇到問題，請：

1. 檢查本文檔的故障排除章節
2. 查看 [GitHub Issues](https://github.com/yourusername/mindhelp/issues)
3. 提交新的 Issue 並包含：
   - Docker 版本：`docker --version`
   - Docker Compose 版本：`docker-compose --version`
   - 錯誤日誌：`docker-compose logs`
   - 系統資訊：`docker system info`

---

**🧠 祝您部署順利！如有任何問題，歡迎隨時聯繫團隊。**
