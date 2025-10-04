#!/bin/bash

echo "啟動 MindHelp 健康監控服務..."
echo

# 設定環境變數
export HEALTH_CHECK_URL="https://mind-map-api.estateangela.dpdns.org/health"
export HEALTH_CHECK_INTERVAL="30s"
export HEALTH_LOG_FILE="health_monitor.log"

# 檢查 Go 是否安裝
if ! command -v go &> /dev/null; then
    echo "錯誤: 未安裝 Go 語言環境"
    echo "請先安裝 Go: https://golang.org/dl/"
    exit 1
fi

# 啟動監控
echo "監控端點: $HEALTH_CHECK_URL"
echo "檢查間隔: $HEALTH_CHECK_INTERVAL"
echo "日誌檔案: $HEALTH_LOG_FILE"
echo
echo "按 Ctrl+C 停止監控"
echo "========================================"

go run health_monitor.go
