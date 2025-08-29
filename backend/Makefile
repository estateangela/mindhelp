.PHONY: help build run test clean docker-build docker-run docker-stop swagger

# 預設目標
.DEFAULT_GOAL := help

# 變數
APP_NAME := mindhelp-backend
BINARY_NAME := main
DOCKER_IMAGE := mindhelp-backend
DOCKER_TAG := latest

# 幫助信息
help: ## 顯示幫助信息
	@echo "MindHelp Backend 開發工具"
	@echo ""
	@echo "可用命令:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# 開發相關
build: ## 構建應用程式
	@echo "構建應用程式..."
	@go build -o $(BINARY_NAME) .

run: ## 運行應用程式
	@echo "運行應用程式..."
	@go run main.go

test: ## 執行測試
	@echo "執行測試..."
	@go test -v ./...

test-coverage: ## 執行測試並顯示覆蓋率
	@echo "執行測試並顯示覆蓋率..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

clean: ## 清理構建文件
	@echo "清理構建文件..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out
	@go clean

# 依賴管理
deps: ## 下載依賴
	@echo "下載依賴..."
	@go mod download

deps-update: ## 更新依賴
	@echo "更新依賴..."
	@go get -u ./...
	@go mod tidy

# Docker 相關
docker-build: ## 構建 Docker 映像
	@echo "構建 Docker 映像..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## 使用 Docker 運行應用程式
	@echo "使用 Docker 運行應用程式..."
	@docker-compose up --build

docker-stop: ## 停止 Docker 容器
	@echo "停止 Docker 容器..."
	@docker-compose down

docker-logs: ## 查看 Docker 日誌
	@echo "查看 Docker 日誌..."
	@docker-compose logs -f

# 資料庫相關
db-up: ## 啟動資料庫
	@echo "啟動資料庫..."
	@docker-compose up db -d

db-down: ## 停止資料庫
	@echo "停止資料庫..."
	@docker-compose stop db

db-reset: ## 重置資料庫
	@echo "重置資料庫..."
	@docker-compose down -v
	@docker-compose up db -d

# Swagger 文檔
swagger: ## 生成 Swagger 文檔
	@echo "生成 Swagger 文檔..."
	@swag init -g main.go

# 開發環境設定
dev-setup: ## 設定開發環境
	@echo "設定開發環境..."
	@cp env.example .env
	@echo "請編輯 .env 文件設定環境變數"
	@echo "然後執行: make db-up"
	@echo "最後執行: make run"

# 部署相關
deploy-check: ## 檢查部署準備
	@echo "檢查部署準備..."
	@go vet ./...
	@go test ./...
	@go build -o $(BINARY_NAME) .
	@echo "部署檢查完成"

# 程式碼品質
lint: ## 執行程式碼檢查
	@echo "執行程式碼檢查..."
	@go vet ./...
	@golangci-lint run

format: ## 格式化程式碼
	@echo "格式化程式碼..."
	@go fmt ./...

# 監控和日誌
health: ## 檢查服務健康狀態
	@echo "檢查服務健康狀態..."
	@curl -f http://localhost:8080/health || echo "服務未運行"

# 快速開發循環
dev: ## 開發模式 (自動重啟)
	@echo "啟動開發模式..."
	@go install github.com/cosmtrek/air@latest
	@air

# 安裝開發工具
install-tools: ## 安裝開發工具
	@echo "安裝開發工具..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
