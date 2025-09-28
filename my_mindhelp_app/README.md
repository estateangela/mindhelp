# MindHelp Flutter 應用程式

<div align="center">

![MindHelp Logo](assets/images/mindhelp.png)

**心理健康支援移動應用程式** 📱🧠

[![Flutter](https://img.shields.io/badge/Flutter-3.6.2+-blue.svg)](https://flutter.dev/)
[![Dart](https://img.shields.io/badge/Dart-3.6.2+-0175C2.svg)](https://dart.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

</div>

## 📱 專案概述

MindHelp Flutter 應用程式是心理健康支援平台的前端部分，提供現代化的移動端使用者體驗。應用程式採用 Flutter 框架開發，支援多平台部署（Android、iOS、Web、Windows、macOS、Linux）。

### 🌟 核心功能

- **📚 專家文章** - 心理健康專業知識分享
- **🧠 心理測驗** - 科學化的心理健康評估
- **🗺️ 資源地圖** - 整合全台心理健康資源
- **💬 AI 聊天** - 24/7 心理健康支援

## 🚀 快速開始

### 前置需求

#### 必要工具
- **Flutter SDK 3.6.2+**
- **Dart SDK**
- **Android Studio** (Android 開發)
- **VS Code** (推薦編輯器)
- **Git**

#### 平台特定需求
- **Android**: Android SDK, Android Studio
- **iOS**: Xcode (macOS), CocoaPods
- **Web**: Chrome 瀏覽器
- **Windows**: Visual Studio 2019+
- **macOS**: Xcode
- **Linux**: CMake, Ninja

### 安裝步驟

#### 1. 環境設置

```bash
# 檢查 Flutter 環境
flutter doctor

# 如果 Flutter 未安裝，請參考官方指南：
# https://docs.flutter.dev/get-started/install
```

#### 2. 獲取專案

```bash
# 克隆專案
git clone <repository-url>
cd mindhelp/my_mindhelp_app

# 安裝依賴項
flutter pub get
```

#### 3. 運行應用程式

```bash
# 運行應用程式
flutter run

# 在特定設備上運行
flutter run -d <device-id>

# 查看可用設備
flutter devices
```

## 📁 專案結構

```
lib/
├── core/                    # 核心配置
│   └── theme.dart          # 應用程式主題
├── models/                 # 資料模型
│   ├── article.dart        # 文章模型
│   ├── chat_message.dart   # 聊天訊息模型
│   ├── counseling_center.dart #心理資源模型
│   ├── map_item.dart       # 地圖項目模型
│   └── resource.dart       # 資源模型
├── pages/                  # 應用程式頁面
│   ├── splash_page.dart    # 跳轉頁面
│   ├── home_page.dart      # 首頁
│   ├── article_page.dart   # 文章頁面
│   ├── quiz_page.dart      # 測驗頁面
│   ├── maps_page.dart      # 地圖頁面
│   ├── chat_page.dart      # 聊天頁面
│   └── notify_page.dart    # 通知頁面
├── services/               # 業務邏輯服務
│   ├── ai_service.dart     # AI 服務
│   └── location_service.dart # 位置服務
├── widgets/                # 共用 UI 組件
│   ├── custom_app_bar.dart # 自訂應用欄
│   ├── input_field.dart    # 輸入欄位
│   └── primary_button.dart # 主要按鈕
└── main.dart               # 應用程式入口點
```

## 🛠️ 技術棧

### 核心技術
- **Flutter 3.6.2+** - 跨平台 UI 框架
- **Dart 3.6.2+** - 程式語言
- **Material Design 3** - UI 設計語言

### 主要依賴項
```yaml
dependencies:
  # UI 和導航
  flutter: sdk: flutter
  cupertino_icons: ^1.0.8
  
  # 地圖和位置服務
  google_maps_flutter: ^2.3.0
  geocoding: ^2.2.1
  geolocator: ^12.0.0
  
  # 網路和資料
  http: ^1.4.0
  path_provider: ^2.0.13
  path: ^1.8.3
  
  # Firebase 服務
  firebase_core: ^4.0.0
  firebase_auth: ^6.0.0
  cloud_firestore: ^6.0.0
  
  # UI 組件
  flutter_markdown: ^0.7.7+1
```

## 🔧 開發指南

### 開發命令

```bash
# 安裝依賴項
flutter pub get

# 運行應用程式
flutter run

# 熱重載 (開發中)
# 按 'r' 鍵進行熱重載
# 按 'R' 鍵進行熱重啟

# 建置應用程式
flutter build apk          # Android APK
flutter build appbundle    # Android App Bundle
flutter build ios          # iOS (需要 macOS)
flutter build web          # Web 應用程式
flutter build windows      # Windows 應用程式
flutter build macos        # macOS 應用程式
flutter build linux        # Linux 應用程式
```

### 程式碼品質

```bash
# 分析程式碼
flutter analyze

# 格式化程式碼
flutter format .

# 運行測試
flutter test

# 測試覆蓋率
flutter test --coverage
```

### 除錯和效能

```bash
# 除錯模式運行
flutter run --debug

# 發布模式運行
flutter run --release

# 效能分析
flutter run --profile
```

## 📱 功能特色

### 已實現功能 ✅
- **現代化 UI 設計** - Material Design 3 設計語言
- **響應式布局** - 支援各種螢幕尺寸
- **多平台支援** - Android、iOS、Web、Windows、macOS、Linux
- **導航結構** - 清晰的頁面導航
- **主題配置** - 統一的視覺風格
- **基本頁面框架** - 所有核心頁面已建立

### 待實現功能 🔄
- **後端 API 整合** - 連接 Go 後端服務
- **真實資料載入** - 動態資料獲取

## 🎨 設計系統

### 主題配置
應用程式使用統一的設計系統，定義在 `lib/core/theme.dart`：

```dart
// 主要顏色
- Primary: 心理健康主題色
- Secondary: 輔助色
- Background: 背景色
- Surface: 表面色
- Error: 錯誤色

// 字體
- 中文: 粉圓體
- 英文: 粉圓體
```

### 組件庫
- **CustomAppBar** - 自訂應用欄
- **InputField** - 統一的輸入欄位
- **PrimaryButton** - 主要操作按鈕

## 📊 建置和部署

### Android 部署

```bash
# 建置 APK
flutter build apk --release

# 建置 App Bundle (推薦)
flutter build appbundle --release

# 簽署 APK (生產環境)
flutter build apk --release --build-name=1.0.0 --build-number=1
```

### iOS 部署

```bash
# 建置 iOS 應用程式 (需要 macOS)
flutter build ios --release

# 開啟 Xcode 進行進一步配置
open ios/Runner.xcworkspace
```

### Web 部署

```bash
# 建置 Web 應用程式
flutter build web --release

# 部署到 Firebase Hosting
firebase deploy
```

### Windows 部署

```bash
# 建置 Windows 應用程式
flutter build windows --release

# 創建安裝程式 (需要額外工具)
```

## 🧪 測試

### 單元測試

```bash
# 運行所有測試
flutter test

# 運行特定測試
flutter test test/widget_test.dart

# 測試覆蓋率
flutter test --coverage
```

### 整合測試

```bash
# 運行整合測試
flutter drive --target=test_driver/app.dart
```

### 手動測試

- **UI 測試** - 檢查所有頁面的 UI 元素
- **導航測試** - 測試頁面間導航
- **響應式測試** - 測試不同螢幕尺寸
- **平台測試** - 測試各平台相容性

## 🐛 常見問題

### 環境問題

#### Flutter 命令無法識別
```bash
# 檢查 Flutter 安裝
flutter doctor

# 確認 PATH 設置
echo $PATH  # Linux/macOS
echo %PATH% # Windows
```

#### Android 建置失敗
```bash
# 接受 Android 許可證
flutter doctor --android-licenses

# 檢查 Android SDK
flutter doctor
```

#### iOS 建置問題 (macOS)
```bash
# 安裝 CocoaPods 依賴項
cd ios && pod install

# 檢查 Xcode 設置
flutter doctor
```

### 依賴項問題

#### 依賴項衝突
```bash
# 清理並重新安裝
flutter clean
flutter pub get

# 更新依賴項
flutter pub upgrade
```

#### 版本相容性
```bash
# 檢查過時依賴項
flutter pub outdated

# 更新特定依賴項
flutter pub upgrade package_name
```

## 📞 支援和貢獻

### 獲取幫助
1. 查看 [Flutter 官方文檔](https://docs.flutter.dev/)
2. 檢查專案的 [Issues](https://github.com/your-repo/mindhelp/issues)
3. 聯繫開發團隊

### 貢獻指南
1. Fork 專案
2. 創建功能分支
3. 提交變更
4. 開啟 Pull Request

### 開發規範
- 遵循 Dart 程式碼規範
- 撰寫單元測試
- 更新相關文檔
- 確保所有測試通過

## 📄 授權

本專案採用 MIT 授權條款 - 查看 [LICENSE](../LICENSE) 文件了解詳情。

---

<div align="center">

**MindHelp Flutter App** - 讓心理健康支援更貼近每個人 📱🧠💚

[![Made with Flutter](https://img.shields.io/badge/Made%20with-Flutter-blue.svg)](https://flutter.dev/)
[![Made with ❤️](https://img.shields.io/badge/Made%20with-❤️-red.svg)](https://github.com/your-repo/mindhelp)

</div>