/// MindHelp Backend API 配置
class ApiConfig {
  // 生產環境 API 基礎 URL
  static const String baseUrl = 'https://mindhelp.onrender.com/api/v1';
  
  // API 端點
  static const String auth = '/auth';
  static const String users = '/users';
  static const String articles = '/articles';
  static const String resources = '/resources';
  static const String counselors = '/counselors';
  static const String counselingCenters = '/counseling-centers';
  static const String recommendedDoctors = '/recommended-doctors';
  static const String maps = '/maps';
  static const String quizzes = '/quizzes';
  static const String chat = '/chat';
  static const String notifications = '/notifications';
  static const String config = '/config';
  static const String bookmarks = '/bookmarks';
  static const String reviews = '/reviews';
  
  // 請求超時設定
  static const int connectTimeout = 30000; // 30 秒
  static const int receiveTimeout = 30000; // 30 秒
  
  // 分頁設定
  static const int defaultPageSize = 10;
  static const int maxPageSize = 100;
}
