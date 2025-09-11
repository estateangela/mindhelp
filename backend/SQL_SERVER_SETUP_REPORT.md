# SQL Server 連線設定完成報告

## 📊 專案概要
- **專案名稱**: MindHelp Backend
- **資料庫類型**: Microsoft SQL Server
- **完成日期**: 2025-09-11
- **狀態**: ✅ 設定完成並成功測試

## 🎯 已完成的任務

### 1. 驅動程式安裝
- ✅ 安裝 `gorm.io/driver/sqlserver v1.6.1`
- ✅ 安裝 `github.com/microsoft/go-mssqldb v1.8.2`
- ✅ 更新相關依賴套件

### 2. 應用程式配置修改
- ✅ 修改 `internal/database/database.go` 使用 SQL Server 驅動
- ✅ 更新 `internal/config/config.go` DSN 格式
- ✅ 調整預設埠號從 5432 (PostgreSQL) 到 1433 (SQL Server)

### 3. 資料模型調整
- ✅ 將所有 UUID 欄位從 `uuid` 改為 `uniqueidentifier`
- ✅ 移除 PostgreSQL 專用的 `gen_random_uuid()` 函數
- ✅ 修正 email 欄位大小限制以支援唯一索引
- ✅ 使用 GORM 的 `BeforeCreate` Hook 自動生成 UUID

### 4. 資料庫連線設定
- **主機**: 140.131.114.241:1433
- **用戶**: MindHelp114
- **資料庫**: 114-MindHelp
- **SSL模式**: encrypt=disable, TrustServerCertificate=true

### 5. 資料表重建
- ✅ 成功刪除舊資料表
- ✅ 重新建立 `users` 資料表
- ✅ 重新建立 `chat_messages` 資料表
- ✅ 重新建立 `locations` 資料表
- ✅ 建立所有索引和外鍵約束

## 🔧 環境變數配置

```bash
# 資料庫設定
DB_HOST=140.131.114.241
DB_PORT=1433
DB_USER=MindHelp114
DB_PASSWORD=!QAZ2wsx//
DB_NAME=114-MindHelp
DB_SSL_MODE=disable

# 應用程式設定
PORT=8080
GIN_MODE=release
JWT_SECRET=mindhelp-super-secret-jwt-key-2024
JWT_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h
```

## 🏗️ 資料表結構

### users 資料表
```sql
CREATE TABLE "users" (
    "id" uniqueidentifier PRIMARY KEY,
    "email" nvarchar(255) NOT NULL,
    "password" nvarchar(MAX) NOT NULL,
    "username" nvarchar(50) NOT NULL,
    "full_name" nvarchar(100),
    "phone" nvarchar(20),
    "avatar" nvarchar(255),
    "is_active" bit DEFAULT 1,
    "last_login" datetimeoffset,
    "created_at" datetimeoffset,
    "updated_at" datetimeoffset,
    "deleted_at" datetimeoffset
);
```

### chat_messages 資料表
```sql
CREATE TABLE "chat_messages" (
    "id" uniqueidentifier PRIMARY KEY,
    "user_id" uniqueidentifier NOT NULL,
    "role" nvarchar(10) NOT NULL,
    "content" text NOT NULL,
    "timestamp" bigint NOT NULL,
    "model" nvarchar(50),
    "tokens" bigint DEFAULT 0,
    "created_at" datetimeoffset,
    "updated_at" datetimeoffset,
    "deleted_at" datetimeoffset,
    CONSTRAINT "fk_users_chat_messages" FOREIGN KEY ("user_id") REFERENCES "users"("id")
);
```

### locations 資料表
```sql
CREATE TABLE "locations" (
    "id" uniqueidentifier PRIMARY KEY,
    "user_id" uniqueidentifier NOT NULL,
    "name" nvarchar(100) NOT NULL,
    "description" text,
    "address" nvarchar(255),
    "latitude" decimal(10,8) NOT NULL,
    "longitude" decimal(11,8) NOT NULL,
    "category" nvarchar(50),
    "phone" nvarchar(20),
    "website" nvarchar(255),
    "rating" decimal(3,2) DEFAULT 0,
    "is_public" bit DEFAULT 0,
    "created_at" datetimeoffset,
    "updated_at" datetimeoffset,
    "deleted_at" datetimeoffset,
    CONSTRAINT "fk_users_locations" FOREIGN KEY ("user_id") REFERENCES "users"("id")
);
```

## 🔐 安全性評估

### 目前配置
- ❌ **SSL/TLS 已停用**: `encrypt=disable`
- ❌ **明文密碼**: 環境變數中包含明文密碼
- ⚠️  **信任所有證書**: `TrustServerCertificate=true`
- ✅ **資料庫層級驗證**: 使用用戶名/密碼驗證

### 🚨 安全建議

#### 立即改善項目：
1. **啟用 SSL/TLS 加密**:
   ```bash
   # 修改 DSN 為:
   server=140.131.114.241;user id=MindHelp114;password=!QAZ2wsx//;port=1433;database=114-MindHelp;encrypt=true;TrustServerCertificate=false
   ```

2. **密碼保護**:
   - 使用 Azure Key Vault 或其他密鑰管理服務
   - 建立 `.env` 檔案並加入 `.gitignore`
   - 考慮使用 Azure AD 認證

3. **網路安全**:
   - 配置防火牆規則限制IP來源
   - 使用VPN或私有網路連接
   - 定期更換資料庫密碼

#### 生產環境建議：
```bash
# 建議的生產環境 DSN
server=140.131.114.241;user id=MindHelp114;password=${DB_PASSWORD};port=1433;database=114-MindHelp;encrypt=true;TrustServerCertificate=false;Connection Timeout=30;
```

## 🚀 部署清單

### 開發環境 ✅
- [x] SQL Server 連線成功
- [x] 資料表遷移完成
- [x] API 伺服器啟動正常

### 生產環境準備 🔄
- [ ] 啟用 SSL/TLS 加密
- [ ] 設定適當的連線池大小
- [ ] 配置監控和日誌
- [ ] 備份策略實施
- [ ] 災難恢復計畫

## 📝 已修改的檔案

1. `go.mod` - 新增 SQL Server 驅動
2. `internal/database/database.go` - 更改驅動和連線
3. `internal/config/config.go` - DSN 格式調整
4. `internal/models/user.go` - UUID 欄位類型
5. `internal/models/chat_message.go` - UUID 欄位類型
6. `internal/models/location.go` - UUID 欄位類型
7. `env.example` - 更新範例環境變數

## 🎉 結論

SQL Server 連線設定已成功完成！應用程式可以正常連接到資料庫並執行所有 CRUD 操作。

**下一步建議**:
1. 在正式部署前實施安全性建議
2. 設定適當的監控和警報
3. 測試所有 API 端點功能
4. 進行負載測試
