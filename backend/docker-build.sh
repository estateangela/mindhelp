#!/bin/bash

# MindHelp Backend Docker Build Script
# 用於本地開發和測試的 Docker 構建腳本

set -e

# 顏色輸出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 函數：打印彩色訊息
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# 檢查是否有 Docker
if ! command -v docker &> /dev/null; then
    print_error "Docker 未安裝或不在 PATH 中"
    exit 1
fi

# 獲取專案根目錄
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

print_info "MindHelp Backend Docker 構建開始..."
print_info "專案根目錄: $PROJECT_ROOT"

# 設定映像標籤
IMAGE_NAME="mindhelp-backend"
TAG=${1:-latest}
FULL_TAG="$IMAGE_NAME:$TAG"

print_info "構建 Docker 映像: $FULL_TAG"

# 構建 Docker 映像
docker build \
    --tag "$FULL_TAG" \
    --tag "$IMAGE_NAME:latest" \
    --file "$SCRIPT_DIR/Dockerfile" \
    "$SCRIPT_DIR"

if [ $? -eq 0 ]; then
    print_success "Docker 映像構建成功: $FULL_TAG"
else
    print_error "Docker 映像構建失敗"
    exit 1
fi

# 詢問是否執行測試
echo ""
read -p "是否要測試運行容器? (y/N): " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]; then
    print_info "啟動測試容器..."
    
    # 停止並移除可能存在的測試容器
    docker stop mindhelp-test 2>/dev/null || true
    docker rm mindhelp-test 2>/dev/null || true
    
    # 運行測試容器
    docker run -d \
        --name mindhelp-test \
        --publish 8080:8080 \
        --env GIN_MODE=debug \
        --env SERVER_PORT=8080 \
        --env DB_HOST=aws-1-ap-southeast-1.pooler.supabase.com \
        --env DB_PORT=6543 \
        --env DB_NAME=postgres \
        --env DB_USER=postgres.haunuvdhisdygfradaya \
        --env DB_PASSWORD=MIND_HELP_2025 \
        --env DB_SSLMODE=require \
        --env JWT_SECRET=your-super-secret-jwt-key-change-in-production \
        --env OPENROUTER_API_KEY=your-openrouter-api-key \
        --env OPENROUTER_BASE_URL=https://openrouter.ai/api/v1 \
        --env CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001 \
        --env LOG_LEVEL=info \
        --env LOG_FORMAT=json \
        "$FULL_TAG"
    
    print_info "等待容器啟動..."
    sleep 5
    
    # 檢查容器狀態
    if docker ps | grep mindhelp-test > /dev/null; then
        print_success "容器啟動成功！"
        print_info "容器名稱: mindhelp-test"
        print_info "端口映射: http://localhost:8080"
        print_info "健康檢查: http://localhost:8080/health"
        print_info "API 文檔: http://localhost:8080/swagger/index.html"
        
        echo ""
        print_info "容器日誌 (前 20 行):"
        docker logs mindhelp-test | head -20
        
        echo ""
        print_info "測試健康檢查端點..."
        sleep 2
        if curl -s http://localhost:8080/health > /dev/null; then
            print_success "健康檢查端點正常回應"
        else
            print_warning "健康檢查端點無回應 (可能是資料庫連線問題)"
        fi
        
        echo ""
        print_info "測試容器管理指令:"
        echo "  查看日誌: docker logs mindhelp-test"
        echo "  停止容器: docker stop mindhelp-test"
        echo "  移除容器: docker rm mindhelp-test"
        echo "  進入容器: docker exec -it mindhelp-test sh"
        
    else
        print_error "容器啟動失敗"
        print_info "容器日誌:"
        docker logs mindhelp-test
        docker rm mindhelp-test 2>/dev/null || true
        exit 1
    fi
else
    print_info "跳過容器測試"
fi

echo ""
print_success "🎉 Docker 構建完成！"
print_info "映像標籤:"
docker images | grep mindhelp-backend | head -5

echo ""
print_info "後續步驟:"
echo "1. 設定適當的環境變數檔案 (.env)"
echo "2. 啟動 PostgreSQL 資料庫"
echo "3. 運行容器: docker run -p 8080:8080 --env-file .env $FULL_TAG"
echo "4. 訪問 API: http://localhost:8080"
