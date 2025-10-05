
# MindHelp - 心理健康支援平台

<div align="center">

![MindHelp Logo](my_mindhelp_app/assets/images/mindhelp.png)

**讓心理健康支援更貼近每個人** 🧠💚

[![Flutter](https://img.shields.io/badge/Flutter-3.6.2+-blue.svg)](https://flutter.dev/)
[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8.svg)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-316192.svg)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg)](https://www.docker.com/)

</div>

## 📋 專案概述

MindHelp 是一個全端心理健康支援平台，提供完整的心理健康服務生態系統。專案包含現代化的 Flutter 移動應用程式和強大的 Go 後端 API，旨在為使用者提供專業的心理健康資源和支援。

### 🌟 核心功能

- **📚 專家文章** - 專業心理健康知識分享
- **🧠 心理測驗** - 科學化的心理健康評估工具
- **🗺️ 資源地圖** - 整合全台心理健康資源的智能地圖
- **💬 AI 聊天** - 24/7 心理健康支援聊天機器人
- **👤 個人化服務** - 使用者個人資料管理和偏好設定
- **🔔 智能通知** - 個人化的心理健康提醒服務

## 🏗️ 專案架構

```
mindhelp/
├── 📱 my_mindhelp_app/          # Flutter 移動應用程式
│   ├── lib/                     # Dart 源碼
│   │   ├── core/               # 核心配置
│   │   ├── models/             # 資料模型
│   │   ├── pages/              # 應用程式頁面
│   │   ├── services/           # 業務邏輯服務
│   │   ├── utils/              # 工具類別
│   │   └── widgets/            # 共用 UI 組件
│   ├── assets/                 # 靜態資源
│   └── platform configs/      # 各平台配置
├── 🔧 backend/                 # Go 後端 API 服務
│   ├── internal/              # 內部應用程式邏輯
│   │   ├── config/            # 配置管理
│   │   ├── database/          # 資料庫連接
│   │   ├── dto/               # 資料傳輸物件
│   │   ├── handlers/          # HTTP 請求處理器
│   │   ├── middleware/        # 中間件
│   │   ├── models/            # 資料庫模型
│   │   ├── routes/            # 路由配置
│   │   └── vo/                # 視圖物件
│   ├── database/              # 資料庫遷移檔案
│   ├── docs/                  # API 文檔
│   └── document/              # 專案文檔
└── 📄 專案文檔                 # 文檔和配置
```

## 🛠️ 技術棧

### 前端 (Flutter)
- **框架**: Flutter 3.6.2+
- **語言**: Dart
- **UI**: Material Design 3
- **狀態管理**: StatefulWidget + Provider 模式
- **地圖**: Google Maps Flutter
- **本地儲存**: SQLite (Sqflite)
- **HTTP 客戶端**: HTTP
- **Firebase**: 認證和雲端資料庫

### 後端 (Go)
- **語言**: Go 1.24+
- **框架**: Gin Gonic
- **資料庫**: PostgreSQL (Supabase)
- **ORM**: GORM
- **認證**: JWT
- **API 文檔**: Swagger/OpenAPI 3.0
- **容器化**: Docker
- **部署**: Render

### 資料庫
- **主要資料庫**: PostgreSQL 15+
- **總記錄數**: 1,071+ 筆
- **諮商師資料**: 961 筆
- **諮商所資料**: 97 筆
- **推薦醫師**: 13 筆

## 🚀 快速開始

### 前置需求

#### 後端開發環境
- Go 1.24+ 
- PostgreSQL 15+
- Docker (可選)

#### 前端開發環境
- Flutter 3.6.2+
- Dart SDK
- Android Studio / VS Code
- Git

### 安裝步驟

#### 1. 克隆專案
```bash
git clone <repository-url>
cd mindhelp
```

#### 2. 設置後端

```bash
# 進入後端目錄
cd backend

# 安裝依賴項
go mod tidy

# 設置環境變數
cp env.example .env
# 編輯 .env 文件，設置資料庫連接資訊

# 運行資料庫遷移
go run cmd/seed/main.go

# 啟動開發伺服器
go run main.go
```

後端將在 `http://localhost:8080` 啟動

#### 3. 設置前端

```bash
# 進入 Flutter 專案目錄
cd my_mindhelp_app

# 安裝依賴項
flutter pub get

# 運行應用程式
flutter run
```

### 生產環境部署

#### 後端部署 (Render)
1. 連接 GitHub 倉庫到 Render
2. 設置環境變數
3. 自動部署完成

**生產環境 URL**: https://mindhelp.onrender.com

#### 前端部署
- **Android**: 建置 APK 並上傳到 Google Play Store
- **iOS**: 透過 Xcode 上傳到 App Store
- **Web**: 部署到 Firebase Hosting 或 Netlify

## 📚 API 文檔

### Swagger 文檔
- **開發環境**: http://localhost:8080/swagger/index.html
- **生產環境**: https://mindhelp.onrender.com/swagger/index.html

### 主要 API 端點

#### 認證服務
- `POST /api/v1/auth/register` - 使用者註冊
- `POST /api/v1/auth/login` - 使用者登入
- `POST /api/v1/auth/logout` - 使用者登出

#### 文章服務
- `GET /api/v1/articles` - 獲取文章列表
- `GET /api/v1/articles/{id}` - 獲取文章詳情
- `POST /api/v1/articles/{id}/bookmark` - 收藏文章

#### 資源管理
- `GET /api/v1/counselors` - 獲取諮商師列表
- `GET /api/v1/counseling-centers` - 獲取諮商所列表
- `GET /api/v1/recommended-doctors` - 獲取推薦醫師列表

#### 地圖服務
- `GET /api/v1/maps/addresses` - 獲取地址資訊
- `GET /api/v1/maps/google-addresses` - Google Maps 格式

#### 測驗服務
- `GET /api/v1/quizzes` - 獲取測驗列表
- `GET /api/v1/quizzes/{id}` - 獲取測驗詳情
- `POST /api/v1/quizzes/{id}/submit` - 提交測驗答案

## 📱 應用程式功能

### 使用者介面
- **現代化設計**: 採用 Material Design 3 設計語言
- **響應式布局**: 支援各種螢幕尺寸
- **繁體中文支援**: 完整的中文本地化
- **無障礙設計**: 符合無障礙使用標準

### 核心功能模組

#### 📖 專家文章
- 心理健康專業文章閱讀
- 文章分類和搜尋
- 個人收藏功能
- 閱讀歷史記錄

#### 🧠 心理測驗
- GAD-7 焦慮自評量表
- PHQ-9 憂鬱篩檢量表
- 壓力量表
- 智能評分和結果解釋

#### 🗺️ 資源地圖
- 整合 Google Maps
- 諮商師、諮商所、推薦醫師位置
- 智能搜尋和篩選
- 聯絡資訊和導航

#### 💬 AI 聊天
- 24/7 心理健康支援
- 智能對話系統
- 聊天歷史記錄
- 情緒分析和建議

## 🚀 最新更新 (2025-09-20)

### ✅ 已完成功能

#### 🗄️ 資料庫整合
- **PostgreSQL 連接**：成功連接到 Supabase PostgreSQL 資料庫
- **資料表創建**：新增 3 個核心資料表
  - `counselors` - 諮商師資料
  - `counseling_centers` - 心理諮商所資料  
  - `recommended_doctors` - 網友推薦醫師＆診所資料
- **資料插入**：成功從 CSV 檔案插入真實資料
  - 諮商師：961 筆記錄
  - 諮商所：97 筆記錄
  - 推薦醫師：13 筆記錄

#### 🗺️ Google Maps 整合
- **地址 API**：新增地圖相關端點
  - `GET /api/v1/maps/addresses` - 獲取所有地址資訊
  - `GET /api/v1/maps/google-addresses` - Google Maps 專用格式
- **多格式支援**：支援 JSON 和 GeoJSON 格式輸出
- **地址提取**：智能從描述中提取地址資訊

#### 📚 API 文檔
- **Swagger 文檔**：完整更新並修復 500 錯誤
- **API 端點**：新增 6 個新的 API 端點
  - 諮商師管理：`/api/v1/counselors`
  - 諮商所管理：`/api/v1/counseling-centers`
  - 推薦醫師管理：`/api/v1/recommended-doctors`
  - 地圖整合：`/api/v1/maps/*`

### 🔗 API 文檔
**Swagger 文檔**: https://mindhelp.onrender.com/swagger/index.html#/

一、Mermaid 流程圖

我將您描述的 Workflow 和 User Story 轉換為 Mermaid 圖表，這樣可以更清晰地看到使用者路徑和功能關聯。

1. 使用者流程圖 (Workflow)

這張圖展示了兩位核心使用者（小陳-求助者，王小姐-幫助者）在 APP 中的主要操作路徑。
程式碼片段

flowchart TD
    subgraph 小陳 (求助者) 的旅程
        A[開啟 APP] --> B{主畫面};
        B --> C[點擊 'AI 聊聊'];
        C --> D[輸入感受與困擾];
        D --> E{AI 提供支持與分析};
        E --> F[建議治療學派];
        F --> G[引導至資源地圖];

        B --> H[點擊 '資源地圖'];
        H --> I[允許定位];
        I --> J[篩選資源類型: 免費/諮商所];
        J --> K[查看機構詳情];
        K --> L[撥打電話/導航];

        B --> M[點擊 '心理測驗'];
        M --> N[完成焦慮/憂鬱量表];
        N --> O[查看結果與解釋];
        O --> P[推薦相關文章];
        P --> Q[閱讀文章學習自助];
    end

    subgraph 王小姐 (幫助者) 的旅程
        R[開啟 APP] --> S{主畫面};
        S --> T[點擊 '專家文章'];
        T --> U[搜尋: 如何幫助朋友];
        U --> V[閱讀文章獲取知識];

        S --> W[點擊 '資源地圖'];
        W --> X[手動輸入朋友地址];
        X --> Y[查找朋友附近的資源];
        Y --> Z[分享資源資訊給朋友];
    end

2. 使用者故事與功能關聯圖 (User Story Map)

這張心智圖展示了使用者角色、他們的核心需求（User Story 的 "I want to..." 部分），以及滿足這些需求的功能。
程式碼片段

mindmap
  root((心理健康 APP))
    小陳 (求助者)
      ))釐清感受與方向((
        [AI 聊天]
        [心理測驗]
      ))尋找專業協助((
        [資源地圖]
        [資源篩選]
      ))學習自助技巧((
        [專家文章]
    王小姐 (幫助者)
      ))了解如何幫朋友((
        [AI 聊天]
        [專家文章]
      ))為朋友找資源((
        [資源地圖]
        [手動搜尋地點]
        [分享功能]

二、擴充版 API 規格 (Expanded API Spec)

這次我將提供更詳盡的規格，包含更豐富的端點、詳細的請求/回應欄位、錯誤處理和資料模型定義。

通用設計原則 (擴充)

    Base URL: https://api.yourdomain.com/v1

    Authentication: Authorization: Bearer <JWT> in HTTP Header.

    標準成功回應:
    JSON

{
  "success": true,
  "data": { ... } // or [ ... ]
}

標準錯誤回應:
JSON

    {
      "success": false,
      "error": {
        "code": "ERROR_CODE_STRING", // e.g., "INVALID_PARAMETERS"
        "message": "A human-readable error message."
      }
    }

資料模型 (Data Models)

預先定義共用的資料結構，讓 API 規格更清晰。

    Resource Model:
    JSON

{
  "id": "string (UUID)",
  "name": "string",
  "type": "enum (clinic, counseling_center, free_service, clinical_psychology)",
  "address": "string",
  "phone": "string",
  "website": "string (nullable)",
  "location": { "lat": "float", "lon": "float" },
  "description": "string",
  "specialties": ["string"], // e.g., ["CBT", "兒童諮商"]
  "isBookmarked": "boolean" // 當前使用者是否已收藏
}

Article Model:
JSON

    {
      "id": "string (UUID)",
      "title": "string",
      "author": { "name": "string", "title": "string" },
      "publishDate": "string (ISO 8601)",
      "summary": "string",
      "content": "string (HTML or Markdown)",
      "tags": ["string"],
      "isBookmarked": "boolean"
    }

1. 使用者 & 驗證 (Users & Auth)

Endpoint	Method	說明
/auth/register	POST	註冊匿名使用者
/users/me	GET	獲取當前使用者資訊

GET /users/me

    說明: 獲取當前登入使用者的基本資訊。

    Headers: Authorization: Bearer <JWT>

    Success Response (200 OK):
    JSON

    {
      "success": true,
      "data": {
        "userId": "user_uuid_string",
        "createdAt": "2025-09-13T12:00:00Z"
      }
    }

2. 資源地圖 (Resources)

Endpoint	Method	說明
/resources	GET	搜尋資源點
/resources/{id}	GET	獲取單一資源點詳情
/users/me/bookmarks/resources	GET	獲取使用者收藏的資源列表
/resources/{id}/bookmark	POST	收藏一個資源點
/resources/{id}/bookmark	DELETE	取消收藏一個資源點

### 🆕 新增：專業資源管理

Endpoint	Method	說明
/counselors	GET	獲取諮商師列表
/counselors/{id}	GET	獲取單一諮商師詳情
/counseling-centers	GET	獲取諮商所列表
/counseling-centers/{id}	GET	獲取單一諮商所詳情
/recommended-doctors	GET	獲取推薦醫師列表
/recommended-doctors/{id}	GET	獲取單一推薦醫師詳情

### 🆕 新增：Google Maps 整合

Endpoint	Method	說明
/maps/addresses	GET	獲取所有地址資訊
/maps/google-addresses	GET	獲取 Google Maps 專用格式

GET /counselors (新增)

    Query Parameters:
    | 參數 | 類型 | 必要 | 說明 |
    | :--- | :--- | :--- | :--- |
    | page | int | 否 | 頁碼，預設 1 |
    | page_size | int | 否 | 每頁數量，預設 10 |
    | search | string | 否 | 搜索關鍵字 |
    | work_location | string | 否 | 工作地點篩選 |
    | specialty | string | 否 | 專業領域篩選 |

    Success Response (200 OK):
    ```json
    {
      "success": true,
      "data": {
        "counselors": [
          {
            "id": "uuid",
            "name": "諮商師姓名",
            "license_number": "諮心字第000001號",
            "gender": "女",
            "specialties": "家庭親子, 壓力與情緒調適",
            "work_location": "臺北市大安區",
            "work_unit": "格瑞思心理諮商所"
          }
        ],
        "total": 961,
        "page": 1,
        "page_size": 10
      }
    }
    ```

GET /maps/addresses (新增)

    Query Parameters:
    | 參數 | 類型 | 必要 | 說明 |
    | :--- | :--- | :--- | :--- |
    | type | string | 否 | 地址類型篩選 (counselor, counseling_center, recommended_doctor) |
    | limit | int | 否 | 限制數量，預設 100 |

    Success Response (200 OK):
    ```json
    {
      "success": true,
      "data": {
        "addresses": [
          {
            "id": "uuid",
            "name": "機構名稱",
            "address": "台北市大安區...",
            "type": "counseling_center",
            "phone": "02-1234-5678"
          }
        ],
        "total": 1071,
        "type": null
      }
    }
    ```

GET /maps/google-addresses (新增)

    Query Parameters:
    | 參數 | 類型 | 必要 | 說明 |
    | :--- | :--- | :--- | :--- |
    | format | string | 否 | 輸出格式 (json, geojson)，預設 json |

    Success Response (200 OK): 返回 Google Maps 專用格式的地址資訊

POST /resources/{id}/bookmark

    說明: 將指定的資源點加入使用者的收藏。

    Headers: Authorization: Bearer <JWT>

    Path Parameters: id (string, required): 資源點的 ID。

    Success Response (204 No Content): 表示操作成功，無須返回內容。

3. AI 聊天 (AI Chat)

Endpoint	Method	說明
/chat/sessions	GET	獲取歷史聊天 session 列表
/chat/sessions	POST	建立新的聊天 session
/chat/sessions/{sessionId}/messages	GET	獲取某個 session 的歷史訊息
/chat/sessions/{sessionId}/messages	POST	發送訊息並取得回覆

GET /chat/sessions

    說明: 獲取使用者的歷史聊天列表，方便使用者回顧。

    Headers: Authorization: Bearer <JWT>

    Query Parameters: page (int, default: 1), limit (int, default: 20)

    Success Response (200 OK):
    JSON

    {
      "success": true,
      "data": [
        {
          "sessionId": "session_uuid_1",
          "firstMessageSnippet": "我最近常常失眠...",
          "lastUpdatedAt": "2025-09-12T10:30:00Z"
        }
      ]
    }

4. 心理測驗 (Quizzes)

Endpoint	Method	說明
/quizzes	GET	獲取測驗列表
/quizzes/{id}	GET	獲取測驗題目
/quizzes/{id}/submit	POST	提交答案並獲取結果
/users/me/quiz_history	GET	獲取使用者歷史測驗結果

GET /users/me/quiz_history

    說明: 讓使用者可以追蹤自己過去的測驗紀錄。

    Headers: Authorization: Bearer <JWT>

    Query Parameters: page (int, default: 1), limit (int, default: 10)

    Success Response (200 OK):
    JSON

    {
      "success": true,
      "data": [
        {
          "historyId": "history_uuid_1",
          "quizTitle": "GAD-7 焦慮自評量表",
          "completedAt": "2025-09-11T14:00:00Z",
          "score": 16,
          "result": "您的分數顯示您可能正經歷中重度的焦慮困擾。"
        }
      ]
    }

5. 專家文章 (Articles)

Endpoint	Method	說明
/articles	GET	搜尋文章
/articles/{id}	GET	獲取單篇文章詳情
/users/me/bookmarks/articles	GET	獲取使用者收藏的文章列表
/articles/{id}/bookmark	POST	收藏一篇文章
/articles/{id}/bookmark	DELETE	取消收藏一篇文章

GET /articles (擴充)

    Query Parameters:
    | 參數 | 類型 | 必要 | 說明 |
    | :--- | :--- | :--- | :--- |
    | q | string | 否 | 搜尋關鍵字 |
    | tag | string | 否 | 依標籤篩選 |
    | sort_by | string | 否 | 排序依據 (publishDate, popularity)，預設 publishDate |
    | page | int | 否 | 頁碼，預設 1 |
    | limit | int | 否 | 每頁數量，預設 10 |

    Success Response (200 OK): 回應 data 欄位為一個 Article Model 陣列 (不含 content 欄位)。

    擴充將圍繞以下幾個核心方向：

    完整的帳號系統：從匿名使用者過渡到完整的註冊會員，包含登入、註冊、密碼管理。

    使用者互動與回饋：新增評論、評分和內容回報機制。

    個人化與通知系統：讓使用者可以管理偏好設定，並接收推播通知。

    應用程式配置：提供一個中心化的端點來管理 APP 的動態設定。

    更嚴謹的規格定義：為每個欄位加上驗證規則，並定義更詳細的錯誤回應。

通用設計原則 (更新版)

    Base URL: https://api.yourdomain.com/v1

    Authentication:

        公開端點 (Public): 無需授權即可存取 (e.g., GET /articles)。

        授權端點 (Authorized): 需要 Authorization: Bearer <JWT> in HTTP Header。

    Pagination (分頁): 對於列表型 API (如文章、評論)，將使用以下分頁參數，並在回應中包含分頁資訊。

        Query Parameters: page (int, default: 1), limit (int, default: 15)。

        Response Body:
        JSON

    "pagination": {
      "currentPage": 1,
      "totalPages": 10,
      "totalItems": 150,
      "limit": 15
    }

標準錯誤回應 (更詳細):
JSON

    {
      "success": false,
      "error": {
        "code": "VALIDATION_ERROR",
        "message": "提供的輸入無效。",
        "details": { // 僅在 VALIDATION_ERROR 時出現
          "email": "請輸入有效的電子郵件地址。",
          "password": "密碼長度不能少於 8 個字元。"
        }
      }
    }

資料模型 (Data Models - 擴充)

    UserModel:
    JSON

{
  "id": "string (UUID)",
  "email": "string (nullable, for registered users)",
  "nickname": "string (nullable)",
  "isAnonymous": "boolean",
  "createdAt": "string (ISO 8601)"
}

ReviewModel:
JSON

{
  "id": "string (UUID)",
  "author": { // 簡化的 UserModel
    "id": "string (UUID)",
    "nickname": "string"
  },
  "resourceId": "string (UUID)",
  "rating": "integer (1-5)",
  "comment": "string (nullable)",
  "createdAt": "string (ISO 8601)",
  "canEdit": "boolean" // 當前使用者是否可編輯/刪除此評論
}

NotificationModel:
JSON

    {
      "id": "string (UUID)",
      "type": "enum (NEW_ARTICLE, PROMOTION, SYSTEM)",
      "title": "string",
      "body": "string",
      "isRead": "boolean",
      "createdAt": "string (ISO 8601)",
      "payload": { // 用於點擊通知後的操作
        "action": "NAVIGATE_TO_ARTICLE",
        "articleId": "article_uuid_1"
      }
    }

擴充 API 規格

1. 應用程式配置 (App Config)

Endpoint	Method	說明	授權
/config	GET	獲取 APP 的遠端配置	公開

GET /config

    說明: APP 啟動時呼叫，用來獲取動態設定，例如篩選條件列表、功能開關等，避免將設定寫死在前端。

    Success Response (200 OK):
    JSON

    {
      "success": true,
      "data": {
        "features": {
          "enableReviews": true, // 功能開關：是否啟用評論功能
          "enableTherapistProfiles": false
        },
        "filters": {
          "resourceTypes": [
            { "key": "clinic", "displayName": "身心科診所" },
            { "key": "counseling_center", "displayName": "心理諮商所" }
          ],
          "specialties": [
            { "key": "CBT", "displayName": "認知行為治療" },
            { "key": "ADHD", "displayName": "注意力不足過動症" }
          ]
        }
      }
    }

2. 完整帳號系統 (Full Auth System)

Endpoint	Method	說明	授權
/auth/register	POST	(更新) 註冊正式帳號	公開
/auth/login	POST	使用 Email 和密碼登入	公開
/auth/logout	POST	登出	需要
/users/me	PUT	更新使用者個人資料	需要
/users/me/password	PUT	變更密碼	需要
/users/me	DELETE	刪除帳號	需要

POST /auth/register (更新)

    說明: 註冊一個新的正式帳號。

    Request Body:
    | 欄位 | 類型 | 驗證規則 |
    | :--- | :--- | :--- |
    | email | string | required, email |
    | password | string | required, minLength:8 |
    | nickname | string | optional, maxLength:50 |

    Success Response (201 Created): 返回 UserModel 和 JWT Token。

    Error Response (409 Conflict): 當 Email 已被註冊時返回。

PUT /users/me

    說明: 更新使用者可修改的個人資料。

    Request Body:
    | 欄位 | 類型 | 驗證規則 |
    | :--- | :--- | :--- |
    | nickname | string | required, minLength:1, maxLength:50 |

    Success Response (200 OK): 返回更新後的 UserModel。

3. 使用者互動與回饋 (User Interaction & Feedback)

Endpoint	Method	說明	授權
/resources/{id}/reviews	GET	獲取某個資源點的所有評論	公開
/resources/{id}/reviews	POST	為某個資源點新增一則評論	需要
/reviews/{reviewId}	PUT	修改自己發布的評論	需要
/reviews/{reviewId}	DELETE	刪除自己發布的評論	需要
/report	POST	回報不當內容	需要

POST /resources/{id}/reviews

    說明: 使用者必須登入才能發表評論。

    Request Body:
    | 欄位 | 類型 | 驗證規則 |
    | :--- | :--- | :--- |
    | rating | integer | required, min:1, max:5 |
    | comment | string | optional, maxLength:1000 |

    Success Response (201 Created): 返回新建的 ReviewModel。

    Error Response (409 Conflict): 如果使用者已經評論過此資源點。

POST /report

    說明: 一個通用的內容回報端點。

    Request Body:
    | 欄位 | 類型 | 驗證規則 |
    | :--- | :--- | :--- |
    | contentType | enum | required, enum(review, article, resource) |
    | contentId | string | required, uuid |
    | reason | enum | required, enum(spam, inappropriate, incorrect_info) |
    | details | string | optional, maxLength:1000 |

    Success Response (202 Accepted): 表示伺服器已收到回報，將進行後續處理。

4. 個人化與通知系統 (Personalization & Notifications)

Endpoint	Method	說明	授權
/notifications	GET	獲取通知列表	需要
/notifications/mark-as-read	POST	將通知標示為已讀	需要
/users/me/notification-settings	GET	獲取通知設定	需要
/users/me/notification-settings	PUT	更新通知設定	需要
/users/me/push-token	POST	註冊/更新裝置的推播 token	需要

PUT /users/me/notification-settings

    說明: 讓使用者可以自訂想收到的通知類型。

    Request Body:
    JSON

    {
      "newArticle": true,
      "promotions": false,
      "systemUpdates": true
    }

    Success Response (200 OK): 返回更新後的設定。

POST /users/me/push-token

    說明: APP 取得推播權限後，將裝置 token 送到後端儲存。

    Request Body:
    | 欄位 | 類型 | 驗證規則 |
    | :--- | :--- | :--- |
    | token | string | required |
    | platform | enum | required, enum(ios, android) |

    Success Response (204 No Content):

    已完成功能 (Phase 1, 2 & 3)
🏗️ 基礎架構
✅ 8個新的資料模型：Article, Quiz, Review, Notification, Bookmark, ChatSession, UserSetting, AppConfig
✅ 3個專業資源模型：Counselor, CounselingCenter, RecommendedDoctor
✅ 完整的 migration 檔案 (002_add_core_features.sql, 003_add_counselor_tables.sql)
✅ 完整的 DTO/VO 結構 - 10個新的 DTO 檔案
✅ 更新的路由配置 - 支援所有新端點
✅ PostgreSQL 資料庫連接 (Supabase)
✅ 真實資料插入 (1071 筆記錄)
👤 使用者管理系統
✅ GET /users/me - 獲取使用者資料
✅ PUT /users/me - 更新個人資料
✅ PUT /users/me/password - 變更密碼
✅ DELETE /users/me - 刪除帳號
✅ GET /users/me/stats - 使用者統計
📚 專家文章系統
✅ GET /articles - 搜尋文章 (支援關鍵字、標籤、排序)
✅ GET /articles/{id} - 文章詳情 (自動增加瀏覽次數)
✅ POST /articles/{id}/bookmark - 收藏文章
✅ DELETE /articles/{id}/bookmark - 取消收藏
🧠 心理測驗系統
✅ GET /quizzes - 獲取測驗列表
✅ GET /quizzes/{id} - 獲取測驗詳情和題目
✅ POST /quizzes/{id}/submit - 提交答案並獲取結果
✅ GET /users/me/quiz_history - 測驗歷史記錄
✅ 智能評分系統 - 支援 GAD-7, PHQ-9, 壓力量表
⭐ 收藏系統
✅ GET /users/me/bookmarks/articles - 文章收藏列表
✅ GET /users/me/bookmarks/resources - 資源收藏列表
✅ POST /bookmarks - 通用收藏功能
✅ DELETE /bookmarks - 取消收藏
💬 評論與評分系統
✅ GET /resources/{id}/reviews - 獲取資源評論 (含統計資訊)
✅ POST /resources/{id}/reviews - 新增評論
✅ PUT /reviews/{reviewId} - 修改評論
✅ DELETE /reviews/{reviewId} - 刪除評論
✅ POST /report - 回報不當內容
✅ 評分統計 - 平均評分和分佈圖
🔔 通知系統
✅ GET /notifications - 通知列表
✅ POST /notifications/mark-as-read - 標記已讀
✅ GET /users/me/notification-settings - 通知設定
✅ PUT /users/me/notification-settings - 更新通知設定
✅ POST /users/me/push-token - 推播 Token 管理
⚙️ 應用配置系統
✅ GET /config - 動態配置 (功能開關、篩選選項)
✅ 功能開關：評論、治療師資料、群組聊天等
✅ 篩選配置：資源類型、專業領域、測驗類別

### 🆕 Phase 3: 專業資源管理 & Google Maps 整合
👥 專業資源管理
✅ GET /counselors - 諮商師列表 (961 筆真實資料)
✅ GET /counselors/{id} - 諮商師詳情
✅ GET /counseling-centers - 諮商所列表 (97 筆真實資料)
✅ GET /counseling-centers/{id} - 諮商所詳情
✅ GET /recommended-doctors - 推薦醫師列表 (13 筆真實資料)
✅ GET /recommended-doctors/{id} - 推薦醫師詳情
✅ 智能搜索：支援姓名、地點、專業領域篩選
✅ 分頁功能：支援大量資料的高效瀏覽

🗺️ Google Maps 整合
✅ GET /maps/addresses - 獲取所有地址資訊
✅ GET /maps/google-addresses - Google Maps 專用格式
✅ 多格式支援：JSON 和 GeoJSON 格式
✅ 地址提取：智能從描述中提取地址資訊
✅ 類型篩選：支援按資源類型篩選地址

📊 資料庫整合
✅ PostgreSQL 連接：成功連接到 Supabase
✅ 資料遷移：自動處理現有資料和結構變更
✅ 資料插入：從 CSV 檔案批量插入真實資料
✅ 錯誤處理：優雅處理重複資料和約束衝突

## 🛠️ 技術規格

### 後端技術棧
- **語言**: Go 1.24
- **框架**: Gin Gonic
- **資料庫**: PostgreSQL (Supabase)
- **ORM**: GORM
- **文檔**: Swagger/OpenAPI 3.0
- **部署**: Docker + Render

### 資料庫規格
- **總記錄數**: 1,071 筆
- **諮商師**: 961 筆 (包含執照號碼、專業領域、工作地點)
- **諮商所**: 97 筆 (包含地址、電話、線上諮商服務)
- **推薦醫師**: 13 筆 (包含經驗次數、描述資訊)

### API 規格
- **Base URL**: https://mindhelp.onrender.com/api/v1
- **認證**: JWT Bearer Token
- **回應格式**: JSON
- **分頁**: 支援 page 和 page_size 參數
- **搜索**: 支援關鍵字搜索和篩選

### 部署資訊
- **生產環境**: https://mindhelp.onrender.com
- **API 文檔**: https://mindhelp.onrender.com/swagger/index.html
- **健康檢查**: https://mindhelp.onrender.com/health

## 🔧 開發指南

### 後端開發

#### 專案結構
```
backend/internal/
├── config/         # 配置管理
├── database/       # 資料庫連接
├── dto/           # 請求/回應資料結構
├── handlers/      # HTTP 處理器
├── middleware/    # 中間件 (認證、日誌)
├── models/        # 資料庫模型
├── routes/        # 路由定義
└── vo/           # 視圖物件
```

#### 開發命令
```bash
# 運行開發伺服器
go run main.go

# 運行測試
go test ./...

# 建置可執行檔
go build -o mindhelp-backend main.go

# 生成 Swagger 文檔
swag init
```

### 前端開發

#### 專案結構
```
my_mindhelp_app/lib/
├── core/          # 核心配置 (主題、常數)
├── models/        # 資料模型
├── pages/         # 應用程式頁面
├── services/      # 業務邏輯服務
├── utils/         # 工具類別
└── widgets/       # 共用 UI 組件
```

#### 開發命令
```bash
# 安裝依賴項
flutter pub get

# 運行應用程式
flutter run

# 建置 APK
flutter build apk

# 分析程式碼
flutter analyze

# 格式化程式碼
flutter format .
```

## 📊 專案統計

### 程式碼統計
- **後端 Go 程式碼**: ~15,000 行
- **前端 Dart 程式碼**: ~8,000 行
- **API 端點**: 25+ 個
- **資料模型**: 15+ 個
- **資料庫記錄**: 1,071+ 筆

### 功能完成度
- ✅ 使用者認證系統 (100%)
- ✅ 文章管理系統 (100%)
- ✅ 心理測驗系統 (100%)
- ✅ 資源地圖系統 (100%)
- ✅ AI 聊天系統 (100%)
- ✅ 收藏系統 (100%)
- ✅ 評論系統 (100%)
- ✅ 通知系統 (100%)
- 🔄 推播通知 (80%)
- 🔄 離線模式 (60%)

## 🚀 部署

### Docker 部署

#### 後端 Docker 部署
```bash
# 建置 Docker 映像
docker build -t mindhelp-backend .

# 運行容器
docker run -p 8080:8080 mindhelp-backend
```

#### 使用 Docker Compose
```bash
# 啟動所有服務
docker-compose up -d

# 停止服務
docker-compose down
```

### 雲端部署

#### Render (後端)
1. 連接 GitHub 倉庫
2. 設置環境變數
3. 配置自動部署

#### Firebase (前端)
```bash
# 建置 Web 版本
flutter build web

# 部署到 Firebase
firebase deploy
```

## 🔐 安全考量

- **JWT 認證**: 安全的用戶認證機制
- **CORS 配置**: 跨域請求安全控制
- **資料驗證**: 輸入資料驗證和清理
- **SQL 注入防護**: 使用 ORM 防止 SQL 注入
- **敏感資料加密**: 密碼和敏感資料加密儲存

## 🤝 貢獻指南

1. Fork 專案
2. 創建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交變更 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 開啟 Pull Request

### 開發規範
- 遵循 Go 和 Dart 程式碼規範
- 撰寫單元測試
- 更新相關文檔
- 確保所有測試通過

## 📄 授權

本專案採用 MIT 授權條款 - 查看 [LICENSE](LICENSE) 文件了解詳情。

## 📞 聯絡資訊

- **專案維護者**: MindHelp 開發團隊
- **Email**: support@mindhelp.com
- **Issues**: [GitHub Issues](https://github.com/your-repo/mindhelp/issues)
- **API 文檔**: [Swagger UI](https://mindhelp.onrender.com/swagger/index.html)

## 🙏 致謝

感謝所有為心理健康領域做出貢獻的專業人士和開發者。

---

<div align="center">

**MindHelp** - 讓心理健康支援更貼近每個人 🧠💚

[![Made with ❤️](https://img.shields.io/badge/Made%20with-❤️-red.svg)](https://github.com/your-repo/mindhelp)

</div>
