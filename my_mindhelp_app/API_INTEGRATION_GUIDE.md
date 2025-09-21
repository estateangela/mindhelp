# MindHelp Flutter API 整合指南

## 📋 概述

本指南說明如何在 MindHelp Flutter 應用程式中使用已整合的 API 服務。

## 🏗️ 架構設計

### 服務層架構
```
services/
├── auth_service.dart      # 使用者認證
├── article_service.dart   # 文章管理
├── resource_service.dart  # 資源管理
├── quiz_service.dart      # 心理測驗
└── chat_service.dart      # AI 聊天
```

### 核心組件
- **ApiClient**: 統一的 HTTP 客戶端，處理認證和錯誤
- **ApiConfig**: API 配置和端點定義
- **Models**: 對應後端 API 的資料模型

## 🔧 使用範例

### 1. 認證服務

#### 登入
```dart
import '../services/auth_service.dart';

final authService = AuthService();

try {
  final response = await authService.login(
    email: 'user@example.com',
    password: 'password123',
  );
  // 登入成功，token 會自動儲存
  print('歡迎回來，${response.user.nickname}！');
} catch (e) {
  print('登入失敗: $e');
}
```

#### 註冊
```dart
try {
  final response = await authService.register(
    email: 'newuser@example.com',
    password: 'password123',
    nickname: '新使用者',
  );
  print('註冊成功: ${response.user.id}');
} catch (e) {
  print('註冊失敗: $e');
}
```

#### 檢查登入狀態
```dart
if (authService.isLoggedIn) {
  print('使用者已登入');
} else {
  print('使用者未登入');
}
```

### 2. 文章服務

#### 獲取文章列表
```dart
import '../services/article_service.dart';

final articleService = ArticleService();

try {
  final response = await articleService.getArticles(
    search: '壓力管理',
    tag: '心理健康',
    sortBy: 'publishDate',
    page: 1,
    limit: 10,
  );
  
  for (final article in response.articles) {
    print('標題: ${article.title}');
    print('作者: ${article.author.name}');
    print('摘要: ${article.summary}');
  }
} catch (e) {
  print('獲取文章失敗: $e');
}
```

#### 收藏文章
```dart
try {
  await articleService.bookmarkArticle('article_id');
  print('文章已收藏');
} catch (e) {
  print('收藏失敗: $e');
}
```

### 3. 資源服務

#### 獲取諮商師列表
```dart
import '../services/resource_service.dart';

final resourceService = ResourceService();

try {
  final response = await resourceService.getCounselors(
    search: '台北',
    workLocation: '台北市',
    specialty: '焦慮症',
    page: 1,
    pageSize: 10,
  );
  
  for (final counselor in response.counselors) {
    print('姓名: ${counselor.name}');
    print('執照: ${counselor.licenseNumber}');
    print('專業: ${counselor.specialties}');
    print('地點: ${counselor.workLocation}');
  }
} catch (e) {
  print('獲取諮商師失敗: $e');
}
```

#### 獲取地圖地址
```dart
try {
  final addresses = await resourceService.getMapAddresses(
    type: 'counseling_center',
    limit: 50,
  );
  
  print('找到 ${addresses['total']} 個諮商所地址');
} catch (e) {
  print('獲取地址失敗: $e');
}
```

### 4. 測驗服務

#### 獲取測驗列表
```dart
import '../services/quiz_service.dart';

final quizService = QuizService();

try {
  final response = await quizService.getQuizzes(
    category: '焦慮',
    page: 1,
    limit: 10,
  );
  
  for (final quiz in response.quizzes) {
    print('測驗: ${quiz.title}');
    print('描述: ${quiz.description}');
    print('題數: ${quiz.questions.length}');
  }
} catch (e) {
  print('獲取測驗失敗: $e');
}
```

#### 提交測驗答案
```dart
try {
  final result = await quizService.submitQuiz(
    quizId: 'quiz_id',
    answers: [
      QuizAnswer(questionId: 'q1', optionId: 'opt1'),
      QuizAnswer(questionId: 'q2', optionId: 'opt2'),
    ],
  );
  
  print('測驗完成');
  print('分數: ${result.score}');
  print('結果: ${result.result}');
} catch (e) {
  print('提交測驗失敗: $e');
}
```

### 5. 聊天服務

#### 獲取聊天會話
```dart
import '../services/chat_service.dart';

final chatService = ChatService();

try {
  final sessions = await chatService.getChatSessions(
    page: 1,
    limit: 20,
  );
  
  for (final session in sessions) {
    print('會話: ${session.id}');
    print('預覽: ${session.firstMessageSnippet}');
    print('更新時間: ${session.lastUpdatedAt}');
  }
} catch (e) {
  print('獲取會話失敗: $e');
}
```

#### 發送訊息
```dart
try {
  final message = await chatService.sendMessage(
    sessionId: 'session_id',
    content: '我最近感到壓力很大...',
  );
  
  print('AI 回覆: ${message.content}');
} catch (e) {
  print('發送訊息失敗: $e');
}
```

## 🎨 頁面整合範例

### 登入頁面整合
```dart
class LoginPage extends StatefulWidget {
  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final AuthService _authService = AuthService();
  bool _isLoading = false;

  Future<void> _handleLogin() async {
    setState(() => _isLoading = true);
    
    try {
      await _authService.login(
        email: _emailController.text.trim(),
        password: _passwordController.text,
      );
      
      if (mounted) {
        Navigator.pushReplacementNamed(context, '/home');
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('登入失敗: ${e.toString()}')),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isLoading = false);
      }
    }
  }
}
```

### 文章列表頁面整合
```dart
class ArticlePage extends StatefulWidget {
  @override
  State<ArticlePage> createState() => _ArticlePageState();
}

class _ArticlePageState extends State<ArticlePage> {
  final ArticleService _articleService = ArticleService();
  List<Article> _articles = [];
  bool _isLoading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadArticles();
  }

  Future<void> _loadArticles() async {
    try {
      final response = await _articleService.getArticles();
      setState(() {
        _articles = response.articles;
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    if (_error != null) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text('載入失敗: $_error'),
            ElevatedButton(
              onPressed: _loadArticles,
              child: const Text('重試'),
            ),
          ],
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: _loadArticles,
      child: ListView.builder(
        itemCount: _articles.length,
        itemBuilder: (context, index) {
          final article = _articles[index];
          return _buildArticleCard(article);
        },
      ),
    );
  }
}
```

## 🛡️ 錯誤處理

### 統一錯誤處理模式
```dart
Future<void> _handleApiCall() async {
  setState(() {
    _isLoading = true;
    _error = null;
  });

  try {
    final data = await someService.getData();
    setState(() {
      _data = data;
      _isLoading = false;
    });
  } catch (e) {
    setState(() {
      _error = e.toString();
      _isLoading = false;
    });
    
    // 顯示錯誤訊息
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('操作失敗: ${e.toString()}')),
      );
    }
  }
}
```

### 常見錯誤類型
- **網路錯誤**: 檢查網路連線
- **認證錯誤**: 重新登入
- **權限錯誤**: 檢查使用者權限
- **資料錯誤**: 檢查請求參數

## 🔄 狀態管理最佳實踐

### 1. 載入狀態
```dart
bool _isLoading = false;

// 開始載入
setState(() => _isLoading = true);

// 結束載入
setState(() => _isLoading = false);
```

### 2. 錯誤狀態
```dart
String? _error;

// 設置錯誤
setState(() => _error = e.toString());

// 清除錯誤
setState(() => _error = null);
```

### 3. 資料狀態
```dart
List<Article> _articles = [];

// 更新資料
setState(() {
  _articles = newArticles;
});
```

## 📱 響應式設計

### 處理不同螢幕尺寸
```dart
@override
Widget build(BuildContext context) {
  final screenWidth = MediaQuery.of(context).size.width;
  
  return screenWidth > 600
      ? _buildTabletLayout()
      : _buildMobileLayout();
}
```

### 處理方向變化
```dart
@override
Widget build(BuildContext context) {
  final orientation = MediaQuery.of(context).orientation;
  
  return orientation == Orientation.landscape
      ? _buildLandscapeLayout()
      : _buildPortraitLayout();
}
```

## 🚀 效能優化

### 1. 分頁載入
```dart
int _currentPage = 1;
bool _hasMore = true;

Future<void> _loadMoreArticles() async {
  if (!_hasMore) return;
  
  final response = await _articleService.getArticles(
    page: _currentPage + 1,
    limit: 10,
  );
  
  setState(() {
    _articles.addAll(response.articles);
    _currentPage++;
    _hasMore = response.articles.length == 10;
  });
}
```

### 2. 快取策略
```dart
// 使用 SharedPreferences 快取資料
final prefs = await SharedPreferences.getInstance();
await prefs.setString('cached_articles', jsonEncode(articles.toJson()));
```

### 3. 圖片優化
```dart
// 使用 CachedNetworkImage 快取網路圖片
CachedNetworkImage(
  imageUrl: article.imageUrl,
  placeholder: (context, url) => CircularProgressIndicator(),
  errorWidget: (context, url, error) => Icon(Icons.error),
)
```

## 📋 檢查清單

### 開發前檢查
- [ ] API 端點是否正確
- [ ] 模型類別是否完整
- [ ] 錯誤處理是否完善
- [ ] 載入狀態是否適當

### 測試檢查
- [ ] 網路連線測試
- [ ] 認證流程測試
- [ ] 錯誤情況測試
- [ ] 效能測試

### 部署前檢查
- [ ] API 基礎 URL 設定
- [ ] 環境變數配置
- [ ] 建置配置檢查
- [ ] 版本號更新

## 🔗 相關連結

- [後端 API 文檔](https://mindhelp.onrender.com/swagger/index.html)
- [Flutter 官方文檔](https://docs.flutter.dev/)
- [Dio HTTP 客戶端](https://pub.dev/packages/dio)
- [JSON 序列化](https://pub.dev/packages/json_annotation)
