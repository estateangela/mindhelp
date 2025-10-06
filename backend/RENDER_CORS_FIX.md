# Render.com CORS 修復指南

## 問題診斷

您的 API 請求失敗主要原因：

### 1. URL 拼寫錯誤
- ❌ 錯誤：`mindhelp.onrenderr.com` (雙r)
- ✅ 正確：`mindhelp.onrender.com` (單r)

### 2. CORS 配置問題
原始配置只允許 `localhost` 來源，生產環境的 Swagger UI 被阻擋。

## 修復步驟

### 第一步：更新 Render.com 環境變數

在您的 Render.com Dashboard 中設定以下環境變數：

```bash
# CORS 配置 - 允許多個來源
CORS_ALLOWED_ORIGINS=https://mindhelp.onrender.com,https://www.mindhelp.onrender.com,http://localhost:3000,http://localhost:3001

# 確保生產模式
GIN_MODE=release

# 其他必要變數
JWT_SECRET=your-production-jwt-secret-key
OPENROUTER_API_KEY=your-openrouter-api-key
DATABASE_URL=your-supabase-connection-string
```

### 第二步：驗證修復效果

使用正確的 URL 測試：

```bash
curl -X 'GET' \
  'https://mindhelp.onrender.com/api/v1/counselors?page=1&page_size=10' \
  -H 'accept: application/json'
```

### 第三步：測試 Swagger UI

訪問：https://mindhelp.onrender.com/swagger/index.html

現在應該能夠成功執行 API 測試。

## 程式碼修復詳情

### 1. 修復了 CORS 配置解析
- 支援逗號分隔的多個來源
- 自動添加生產環境來源
- 增強的 CORS 標頭支援

### 2. 增強的 CORS 中間件
- 添加更多允許的 HTTP 方法
- 增加標頭支援
- 設定預快取時間

## 常見問題解決

### Q: 仍然看到 CORS 錯誤？
A: 確認環境變數已正確設定並重新部署應用。

### Q: Swagger UI 仍無法工作？
A: 檢查瀏覽器開發者工具的 Network 標籤，確認請求 URL 正確。

### Q: 資料庫連接問題？
A: 確認 `DATABASE_URL` 環境變數正確設定。

## 部署檢查清單

- [ ] 更新 `CORS_ALLOWED_ORIGINS` 環境變數
- [ ] 設定 `GIN_MODE=release`
- [ ] 驗證 `DATABASE_URL` 正確
- [ ] 重新部署應用
- [ ] 測試 Swagger UI
- [ ] 測試 API 端點

---

**注意**：修復後需要重新部署應用才會生效。
