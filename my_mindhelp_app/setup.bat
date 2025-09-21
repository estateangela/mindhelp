@echo off
REM MindHelp Flutter 應用程式設置腳本 (Windows)

echo 🚀 開始設置 MindHelp Flutter 應用程式...

REM 檢查 Flutter 是否已安裝
flutter --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Flutter 未安裝，請先安裝 Flutter SDK 3.6.2+
    echo 📖 安裝指南: https://docs.flutter.dev/get-started/install
    pause
    exit /b 1
)

echo ✅ Flutter 已安裝
flutter --version | findstr /R "Flutter"

REM 進入專案目錄
cd /d "%~dp0"

REM 清理舊的建置快取
echo 🧹 清理舊的建置快取...
flutter clean

REM 獲取依賴項
echo 📦 安裝依賴項...
flutter pub get

REM 生成 JSON 序列化代碼
echo 🔧 生成 JSON 序列化代碼...
flutter packages pub run build_runner build --delete-conflicting-outputs

REM 分析程式碼
echo 🔍 分析程式碼...
flutter analyze

REM 格式化程式碼
echo ✨ 格式化程式碼...
flutter format .

echo.
echo 🎉 MindHelp Flutter 應用程式設置完成！
echo.
echo 📱 運行應用程式:
echo    flutter run
echo.
echo 🔧 其他有用命令:
echo    flutter test          # 運行測試
echo    flutter build apk     # 建置 Android APK
echo    flutter build web     # 建置 Web 版本
echo.
echo 📖 API 文檔: https://mindhelp.onrender.com/swagger/index.html
echo.
pause
