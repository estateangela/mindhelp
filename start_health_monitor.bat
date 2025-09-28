@echo off
echo 啟動 MindHelp 健康監控服務...
echo.

REM 設定環境變數
set HEALTH_CHECK_URL=https://mind-map-api.estateangela.dpdns.org/health
set HEALTH_CHECK_INTERVAL=30s
set HEALTH_LOG_FILE=health_monitor.log

REM 檢查 Go 是否安裝
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo 錯誤: 未安裝 Go 語言環境
    echo 請先安裝 Go: https://golang.org/dl/
    pause
    exit /b 1
)

REM 啟動監控
echo 監控端點: %HEALTH_CHECK_URL%
echo 檢查間隔: %HEALTH_CHECK_INTERVAL%
echo 日誌檔案: %HEALTH_LOG_FILE%
echo.
echo 按 Ctrl+C 停止監控
echo ========================================

go run health_monitor.go

pause
