# 🚀 MindHelp Backend 部署指南

## 📋 概述
本指南提供 MindHelp Backend 的完整部署流程，包含安全配置、監控設定和最佳實務。

## 🔧 環境需求

### 基本需求
- **Go**: 1.23.0 或更高版本
- **SQL Server**: 2019 或更高版本  
- **作業系統**: Windows Server 2019+ / Linux / macOS
- **記憶體**: 最少 2GB RAM
- **磁碟**: 最少 10GB 可用空間

### 網路需求
- 對 SQL Server (140.131.114.241:1433) 的連接
- HTTPS/TLS 支援（生產環境）
- 防火牆規則允許應用程式埠（預設 8080）

## 🔐 安全配置

### 1. 環境變數設定

```bash
# 複製並編輯環境配置
cp env.production.example .env

# 設定資料庫連線（生產環境）
export DB_HOST="140.131.114.241"
export DB_PORT="1433"
export DB_USER="MindHelp114"
export DB_PASSWORD="your-secure-password"
export DB_NAME="114-MindHelp"

# 啟用 SSL/TLS 加密
export DB_ENCRYPT="true"
export DB_TRUST_CERT="false"  # 生產環境建議設為 false

# JWT 安全配置
export JWT_SECRET="your-256-bit-production-secret-key"
export JWT_EXPIRY="24h"
export JWT_REFRESH_EXPIRY="168h"
```

### 2. SSL/TLS 配置

```bash
# 應用程式層 TLS
export ENABLE_TLS="true"
export TLS_CERT_FILE="/path/to/certificate.crt"
export TLS_KEY_FILE="/path/to/private.key"

# 資料庫層加密
export DB_ENCRYPT="true"
export DB_CONNECTION_TIMEOUT="30"
export DB_MAX_POOL_SIZE="100"
```

### 3. 憑證產生（開發環境）

```bash
# 產生自簽憑證（僅開發用）
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```

## 📊 日誌和監控配置

### 1. 日誌設定

```bash
# 結構化日誌
export LOG_LEVEL="info"
export LOG_FORMAT="json"
export LOG_FILE="/var/log/mindhelp/app.log"
export ENABLE_STRUCTURED_LOGGING="true"
```

### 2. 監控端點

應用程式提供以下監控端點：

- **健康檢查**: `GET /health`
- **詳細健康檢查**: `GET /health/detailed`
- **就緒檢查**: `GET /health/ready`
- **存活檢查**: `GET /health/live`  
- **效能指標**: `GET /metrics`

### 3. 監控配置

```bash
# Prometheus 指標
export ENABLE_METRICS="true"
export METRICS_PORT="9090"
export HEALTH_CHECK_INTERVAL="30s"
```

## 🏭 部署步驟

### 1. 編譯應用程式

```bash
# 建置生產版本
go build -ldflags="-w -s" -o mindhelp-backend .

# 或使用 Makefile
make build-prod
```

### 2. 設定系統服務（Linux）

創建 systemd 服務檔案 `/etc/systemd/system/mindhelp-backend.service`：

```ini
[Unit]
Description=MindHelp Backend Service
After=network.target

[Service]
Type=simple
User=mindhelp
WorkingDirectory=/opt/mindhelp
ExecStart=/opt/mindhelp/mindhelp-backend
Restart=always
RestartSec=5
EnvironmentFile=/opt/mindhelp/.env

# 安全設定
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ReadWritePaths=/var/log/mindhelp

[Install]
WantedBy=multi-user.target
```

啟動服務：
```bash
sudo systemctl daemon-reload
sudo systemctl enable mindhelp-backend
sudo systemctl start mindhelp-backend
```

### 3. 設定 Windows 服務

使用 NSSM 或內建服務管理：

```cmd
# 使用 NSSM
nssm install MindHelpBackend "C:\path\to\mindhelp-backend.exe"
nssm set MindHelpBackend AppDirectory "C:\path\to\app"
nssm start MindHelpBackend
```

### 4. Docker 部署

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o mindhelp-backend .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/mindhelp-backend .
COPY --from=builder /app/env.production.example .env
CMD ["./mindhelp-backend"]
```

## 🔍 健康檢查和監控

### 1. 健康檢查腳本

```bash
#!/bin/bash
# health_check.sh

HEALTH_URL="http://localhost:8080/health"
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" $HEALTH_URL)

if [ $RESPONSE -eq 200 ]; then
    echo "✅ Service is healthy"
    exit 0
else
    echo "❌ Service is unhealthy (HTTP $RESPONSE)"
    exit 1
fi
```

### 2. Prometheus 配置

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'mindhelp-backend'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 30s
```

### 3. 日誌輪替

```bash
# /etc/logrotate.d/mindhelp-backend
/var/log/mindhelp/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 0644 mindhelp mindhelp
    postrotate
        systemctl reload mindhelp-backend
    endscript
}
```

## 🛡️ 安全最佳實務

### 1. 防火牆設定

```bash
# Ubuntu/Debian
sudo ufw allow 8080/tcp
sudo ufw allow from 140.131.114.241 to any port 1433

# CentOS/RHEL
firewall-cmd --add-port=8080/tcp --permanent
firewall-cmd --add-rich-rule="rule family='ipv4' source address='140.131.114.241' port port='1433' protocol='tcp' accept" --permanent
firewall-cmd --reload
```

### 2. 資料庫安全

```sql
-- SQL Server 安全設定
-- 限制連線 IP
USE master;
CREATE LOGIN [MindHelp114] WITH PASSWORD = 'SecurePassword123!';
USE [114-MindHelp];
CREATE USER [MindHelp114] FOR LOGIN [MindHelp114];
ALTER ROLE db_datareader ADD MEMBER [MindHelp114];
ALTER ROLE db_datawriter ADD MEMBER [MindHelp114];
```

### 3. SSL 憑證管理

```bash
# 使用 Let's Encrypt（生產環境）
certbot certonly --standalone -d yourdomain.com
export TLS_CERT_FILE="/etc/letsencrypt/live/yourdomain.com/fullchain.pem"
export TLS_KEY_FILE="/etc/letsencrypt/live/yourdomain.com/privkey.pem"
```

## 🚨 故障排除

### 1. 常見問題

**資料庫連線失敗**
```bash
# 檢查連線
telnet 140.131.114.241 1433
# 檢查 DNS 解析
nslookup 140.131.114.241
```

**SSL 憑證問題**
```bash
# 檢查憑證有效性
openssl x509 -in certificate.crt -text -noout
```

**記憶體洩漏**
```bash
# 監控記憶體使用
curl -s http://localhost:8080/health/detailed | jq '.system.memory'
```

### 2. 日誌分析

```bash
# 查看錯誤日誌
tail -f /var/log/mindhelp/app.log | grep '"level":"ERROR"'

# 分析效能指標
curl -s http://localhost:8080/metrics | jq .
```

## 📈 效能優化

### 1. 資料庫連線池

```go
// 建議設定
DB_MAX_POOL_SIZE=100
DB_CONNECTION_TIMEOUT=30
```

### 2. 記憶體優化

```bash
# 設定 Go 記憶體限制
export GOGC=100
export GOMEMLIMIT=2GiB
```

### 3. 負載均衡

```nginx
# Nginx 配置
upstream mindhelp_backend {
    server localhost:8080 max_fails=3 fail_timeout=30s;
    server localhost:8081 max_fails=3 fail_timeout=30s;
}

server {
    listen 443 ssl;
    server_name yourdomain.com;
    
    location / {
        proxy_pass http://mindhelp_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    location /health {
        proxy_pass http://mindhelp_backend/health;
        access_log off;
    }
}
```

## 🔄 備份和災難恢復

### 1. 資料庫備份

```sql
-- SQL Server 備份
BACKUP DATABASE [114-MindHelp] 
TO DISK = 'C:\Backups\mindhelp_backup.bak'
WITH FORMAT, INIT, COMPRESSION;
```

### 2. 應用程式備份

```bash
# 每日備份腳本
#!/bin/bash
DATE=$(date +%Y%m%d)
tar -czf "/backup/mindhelp-app-$DATE.tar.gz" /opt/mindhelp/
```

---

**📞 支援聯絡**
如遇部署問題，請聯絡開發團隊或查閱專案 Wiki。
