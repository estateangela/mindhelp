#!/bin/bash

# MindHelp Flutter 應用程式設置腳本

echo "🚀 開始設置 MindHelp Flutter 應用程式..."

# 檢查 Flutter 是否已安裝
if ! command -v flutter &> /dev/null; then
    echo "❌ Flutter 未安裝，請先安裝 Flutter SDK 3.6.2+"
    echo "📖 安裝指南: https://docs.flutter.dev/get-started/install"
    exit 1
fi

echo "✅ Flutter 已安裝: $(flutter --version | head -n 1)"

# 進入專案目錄
cd "$(dirname "$0")"

# 清理舊的建置快取
echo "🧹 清理舊的建置快取..."
flutter clean

# 獲取依賴項
echo "📦 安裝依賴項..."
flutter pub get

# 生成 JSON 序列化代碼
echo "🔧 生成 JSON 序列化代碼..."
flutter packages pub run build_runner build --delete-conflicting-outputs

# 分析程式碼
echo "🔍 分析程式碼..."
flutter analyze

# 格式化程式碼
echo "✨ 格式化程式碼..."
flutter format .

echo "🎉 MindHelp Flutter 應用程式設置完成！"
echo ""
echo "📱 運行應用程式:"
echo "   flutter run"
echo ""
echo "🔧 其他有用命令:"
echo "   flutter test          # 運行測試"
echo "   flutter build apk     # 建置 Android APK"
echo "   flutter build web     # 建置 Web 版本"
echo ""
echo "📖 API 文檔: https://mindhelp.onrender.com/swagger/index.html"
