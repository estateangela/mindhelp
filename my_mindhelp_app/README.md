# MindHelp App

[English](#english) | [繁體中文](#繁體中文)

---

## English

### Overview

MindHelp is a comprehensive mental health support Flutter application designed to provide users with accessible mental health resources, AI-powered chat support, and location-based mental health services.

### Features

- **🔐 User Authentication**: Secure login/signup with Firebase Auth
- **💬 AI Chat Support**: Intelligent conversation with OpenRouter AI API
- **🗺️ Location Services**: Google Maps integration for finding nearby mental health resources
- **📱 Cross-Platform**: Available on iOS, Android, Web, Windows, macOS, and Linux
- **🎨 Modern UI**: Clean, intuitive interface with custom Huninn font
- **💾 Local Storage**: SQLite database for offline chat history
- **☁️ Cloud Sync**: Firebase Firestore for data synchronization
- **🔔 Notifications**: Built-in notification system
- **⚙️ Settings**: Customizable app preferences

### Tech Stack

- **Frontend**: Flutter 3.6.2+
- **Backend**: Firebase (Auth, Firestore)
- **AI Integration**: OpenRouter API
- **Maps**: Google Maps Flutter
- **Database**: SQLite (local), Firestore (cloud)
- **State Management**: Flutter built-in state management
- **UI Components**: Material Design with custom theming

### Project Structure

```
lib/
├── api/           # API integrations (OpenRouter)
├── core/          # App theme and core configurations
├── models/        # Data models (ChatMessage)
├── pages/         # App screens and navigation
├── utils/         # Helper functions and services
├── widgets/       # Reusable UI components
└── main.dart      # App entry point
```

### Getting Started

#### Prerequisites

- Flutter SDK 3.6.2 or higher
- Dart SDK
- Android Studio / Xcode (for mobile development)
- Firebase project setup

#### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd my_mindhelp_app
   ```

2. **Install dependencies**
   ```bash
   flutter pub get
   ```

3. **Configure Firebase**
   - Create a Firebase project
   - Add your `google-services.json` (Android) and `GoogleService-Info.plist` (iOS)
   - Enable Authentication and Firestore

4. **Configure API Keys**
   - Set up OpenRouter API key for AI chat functionality
   - Configure Google Maps API key

5. **Run the app**
   ```bash
   flutter run
   ```

### Available Pages

- **Login/Signup**: User authentication
- **Home**: Main dashboard and navigation
- **Chat**: AI-powered mental health conversations
- **Maps**: Location-based mental health resources
- **Profile**: User profile management
- **Settings**: App configuration
- **Notifications**: User alerts and updates

### Building for Production

```bash
# Android APK
flutter build apk --release

# iOS
flutter build ios --release

# Web
flutter build web

# Windows
flutter build windows

# macOS
flutter build macos

# Linux
flutter build linux
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

### License

This project is private and not intended for public distribution.

---

## 繁體中文

### 專案概述

MindHelp 是一個全面的心理健康支援 Flutter 應用程式，旨在為使用者提供可及的心理健康資源、AI 驅動的聊天支援，以及基於位置的心理健康服務。

### 主要功能

- **🔐 使用者認證**: 使用 Firebase Auth 的安全登入/註冊
- **💬 AI 聊天支援**: 與 OpenRouter AI API 的智能對話
- **🗺️ 位置服務**: Google Maps 整合，尋找附近的心理健康資源
- **📱 跨平台支援**: 支援 iOS、Android、Web、Windows、macOS 和 Linux
- **🎨 現代化介面**: 乾淨、直觀的介面，使用自訂 Huninn 字體
- **💾 本地儲存**: SQLite 資料庫用於離線聊天記錄
- **☁️ 雲端同步**: Firebase Firestore 用於資料同步
- **🔔 通知系統**: 內建通知功能
- **⚙️ 設定**: 可自訂的應用程式偏好設定

### 技術架構

- **前端**: Flutter 3.6.2+
- **後端**: Firebase (認證、Firestore)
- **AI 整合**: OpenRouter API
- **地圖服務**: Google Maps Flutter
- **資料庫**: SQLite (本地)、Firestore (雲端)
- **狀態管理**: Flutter 內建狀態管理
- **UI 元件**: Material Design 與自訂主題

### 專案結構

```
lib/
├── api/           # API 整合 (OpenRouter)
├── core/          # 應用程式主題和核心配置
├── models/        # 資料模型 (ChatMessage)
├── pages/         # 應用程式畫面與導航
├── utils/         # 輔助函數和服務
├── widgets/       # 可重複使用的 UI 元件
└── main.dart      # 應用程式進入點
```

### 開始使用

#### 前置需求

- Flutter SDK 3.6.2 或更高版本
- Dart SDK
- Android Studio / Xcode (用於行動開發)
- Firebase 專案設定

#### 安裝步驟

1. **複製專案**
   ```bash
   git clone <repository-url>
   cd my_mindhelp_app
   ```

2. **安裝依賴套件**
   ```bash
   flutter pub get
   ```

3. **設定 Firebase**
   - 建立 Firebase 專案
   - 新增 `google-services.json` (Android) 和 `GoogleService-Info.plist` (iOS)
   - 啟用認證和 Firestore

4. **設定 API 金鑰**
   - 設定 OpenRouter API 金鑰用於 AI 聊天功能
   - 配置 Google Maps API 金鑰

5. **執行應用程式**
   ```bash
   flutter run
   ```

### 可用頁面

- **登入/註冊**: 使用者認證
- **首頁**: 主要儀表板和導航
- **聊天**: AI 驅動的心理健康對話
- **地圖**: 基於位置的心理健康資源
- **個人資料**: 使用者資料管理
- **設定**: 應用程式配置
- **通知**: 使用者提醒和更新

### 建置生產版本

```bash
# Android APK
flutter build apk --release

# iOS
flutter build ios --release

# Web
flutter build web

# Windows
flutter build windows

# macOS
flutter build macos

# Linux
flutter build linux
```

### 貢獻指南

1. Fork 專案
2. 建立功能分支
3. 進行您的修改
4. 如適用，新增測試
5. 提交 Pull Request

### 授權

此專案為私人專案，不適用於公開發行。
