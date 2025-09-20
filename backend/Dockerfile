# 使用官方 Go 映像作為構建階段
FROM golang:1.24-alpine AS builder

# 設定工作目錄
WORKDIR /app

# 安裝必要的系統依賴
RUN apk add --no-cache git ca-certificates tzdata

# 複製 go mod 文件
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製源代碼
COPY . .

# 構建應用程式
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 使用輕量級的 alpine 映像作為運行階段
FROM alpine:latest

# 安裝 ca-certificates 和 tzdata
RUN apk --no-cache add ca-certificates tzdata

# 創建非 root 使用者
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 設定工作目錄
WORKDIR /root/

# 從構建階段複製二進制文件
COPY --from=builder /app/main .

# 複製環境變數文件
COPY --from=builder /app/env.example .env

# 複製 CSV 資料文件
COPY --from=builder /app/document/ ./document/

# 複製 Swagger 文檔
COPY --from=builder /app/docs/ ./docs/

# 變更檔案擁有者
RUN chown -R appuser:appgroup /root/

# 切換到非 root 使用者
USER appuser

# 暴露端口
EXPOSE 8080

# 健康檢查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 運行應用程式
CMD ["./main"]
