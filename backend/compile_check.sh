#!/bin/bash

echo "🔍 檢查 Go 程式碼編譯..."

# 檢查語法錯誤
echo "檢查語法錯誤..."
go vet ./...

# 檢查未使用的 imports
echo "檢查未使用的 imports..."
go mod tidy

# 嘗試編譯
echo "嘗試編譯..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

if [ $? -eq 0 ]; then
    echo "✅ 編譯成功！"
    rm -f main  # 清理編譯產出
else
    echo "❌ 編譯失敗！"
    exit 1
fi

echo "🎉 所有檢查通過！"
