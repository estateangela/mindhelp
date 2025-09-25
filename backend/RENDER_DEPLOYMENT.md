# Render 部署指南 - Supabase 資料庫修正

## 🚨 問題診斷

當前部署出現的問題：
- 資料庫連接被拒絕
- Supabase 可能處於休眠狀態
- 連接字串配置可能需要優化

## 🛠️ 解決方案

### 步驟 1: 更新 Render 環境變數

在 Render Dashboard 中更新以下環境變數：

#### 資料庫設定（選擇其中一種）

**選項 A: 使用完整連接字串（推薦）**
```
DATABASE_URL=postgresql://postgres.haunuvdhisdygfradaya:MIND_HELP_2025@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres?sslmode=require&connect_timeout=30&pool_mode=transaction
```

**選項 B: 使用個別參數**
```
DB_HOST=aws-1-ap-southeast-1.pooler.supabase.com
DB_PORT=6543
DB_USER=postgres.haunuvdhisdygfradaya
DB_PASSWORD=MIND_HELP_2025
DB_NAME=postgres
DB_SSL_MODE=require
```

### 步驟 2: 重新部署

1. 在 Render Dashboard 中：
   - 進入你的服務
   - 點擊 **Environment** 標籤
   - 更新環境變數
   - 點擊 **Save changes**
   - 點擊 **Deploy latest commit**

2. 或者從命令列：
   ```bash
   # 觸發重新部署
   git commit --allow-empty -m "Fix Supabase connection"
   git push origin main
   ```

### 步驟 3: 驗證連接

部署完成後，檢查日誌：
- 成功連接：`資料庫連接成功!`
- 失敗：查看詳細錯誤訊息

## 🔧 進階診斷

如果問題持續，嘗試以下方法：

### 方法 A: 檢查 Supabase 狀態

1. 登入 Supabase Dashboard
2. 檢查資料庫是否處於活動狀態
3. 如果資料庫已休眠，點擊 **Resume** 喚醒

### 方法 B: 測試不同連接配置

```bash
# 嘗試不同的 SSL 模式
DATABASE_URL=postgresql://postgres.haunuvdhisdygfradaya:MIND_HELP_2025@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres?sslmode=disable

# 或者使用不同的主機
DATABASE_URL=postgresql://postgres.haunuvdhisdygfradaya:db.supabase.co:5432/postgres?sslmode=require
```

### 方法 C: 檢查網路連通性

在 Render 服務日誌中查找：
- `connection refused` - 資料庫無法連接到
- `timeout` - 網路超時
- `authentication failed` - 認證失敗

## 📊 監控與日誌

### 檢查 Render 日誌
1. Render Dashboard → 你的服務
2. 點擊 **Logs** 標籤
3. 查看實時日誌

### 健康檢查端點
部署後測試：
```bash
curl https://your-app.render.com/health
```

期望回應：
```json
{
  "status": "ok",
  "checks": {
    "database": "healthy"
  }
}
```

## 🚨 緊急故障排除

如果所有方法都失敗：

1. **檢查 Supabase 專案狀態**
   - 確認專案沒有被暫停
   - 檢查是否有未付費用戶

2. **聯繫 Render 支援**
   - 檢查是否有 IP 限制
   - 確認網路連通性

3. **臨時解決方案**
   - 考慮使用本地資料庫進行測試
   - 設定更長的連接超時時間

## 📝 預期結果

修正後應該看到：
- ✅ 資料庫連接成功
- ✅ API 端點正常回應
- ✅ 健康檢查顯示資料庫狀態正常
- ✅ 應用程式可以正常運行

## 🔄 下一步

1. 更新 Render 環境變數
2. 重新部署
3. 監控日誌
4. 驗證功能正常
