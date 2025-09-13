@echo off
REM MindHelp Backend Docker Build Script for Windows
REM 用於本地開發和測試的 Docker 構建腳本

echo.
echo 🧠 MindHelp Backend Docker 構建開始...
echo.

REM 檢查 Docker 是否安裝
docker --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker 未安裝或不在 PATH 中
    pause
    exit /b 1
)

REM 設定映像標籤
set IMAGE_NAME=mindhelp-backend
set TAG=%1
if "%TAG%"=="" set TAG=latest
set FULL_TAG=%IMAGE_NAME%:%TAG%

echo ℹ️  構建 Docker 映像: %FULL_TAG%
echo.

REM 構建 Docker 映像
docker build --tag "%FULL_TAG%" --tag "%IMAGE_NAME%:latest" --file "Dockerfile" .

if errorlevel 1 (
    echo ❌ Docker 映像構建失敗
    pause
    exit /b 1
)

echo.
echo ✅ Docker 映像構建成功: %FULL_TAG%
echo.

REM 詢問是否執行測試
set /p choice="是否要測試運行容器? (y/N): "
if /i "%choice%" neq "y" goto :skip_test

echo.
echo ℹ️  啟動測試容器...

REM 停止並移除可能存在的測試容器
docker stop mindhelp-test 2>nul
docker rm mindhelp-test 2>nul

REM 運行測試容器
docker run -d --name mindhelp-test --publish 8080:8080 --env GIN_MODE=debug --env SERVER_PORT=8080 --env DB_HOST=localhost --env DB_PORT=5432 --env DB_NAME=mindhelp --env DB_USER=postgres --env DB_PASSWORD=password --env DB_SSLMODE=disable --env JWT_SECRET=your-secret-key --env OPENROUTER_API_KEY=your-openrouter-key "%FULL_TAG%"

echo ℹ️  等待容器啟動...
timeout /t 5 /nobreak >nul

REM 檢查容器狀態
docker ps | findstr mindhelp-test >nul
if errorlevel 1 (
    echo ❌ 容器啟動失敗
    echo ℹ️  容器日誌:
    docker logs mindhelp-test
    docker rm mindhelp-test 2>nul
    pause
    exit /b 1
)

echo ✅ 容器啟動成功！
echo ℹ️  容器名稱: mindhelp-test
echo ℹ️  端口映射: http://localhost:8080
echo ℹ️  健康檢查: http://localhost:8080/health
echo ℹ️  API 文檔: http://localhost:8080/swagger/index.html
echo.

echo ℹ️  容器日誌 (前幾行):
docker logs mindhelp-test

echo.
echo ℹ️  測試容器管理指令:
echo   查看日誌: docker logs mindhelp-test
echo   停止容器: docker stop mindhelp-test
echo   移除容器: docker rm mindhelp-test
echo   進入容器: docker exec -it mindhelp-test sh

goto :end_test

:skip_test
echo ℹ️  跳過容器測試

:end_test
echo.
echo 🎉 Docker 構建完成！
echo.
echo ℹ️  映像標籤:
docker images | findstr mindhelp-backend

echo.
echo ℹ️  後續步驟:
echo 1. 設定適當的環境變數檔案 (.env)
echo 2. 啟動 PostgreSQL 資料庫
echo 3. 運行容器: docker run -p 8080:8080 --env-file .env %FULL_TAG%
echo 4. 訪問 API: http://localhost:8080
echo.

pause
