import 'package:dio/dio.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'api_config.dart';

/// 統一的 API 客戶端，處理認證和錯誤
class ApiClient {
  static final ApiClient _instance = ApiClient._internal();
  factory ApiClient() => _instance;
  ApiClient._internal();

  late Dio _dio;
  String? _authToken;

  /// 初始化 API 客戶端
  Future<void> initialize() async {
    _dio = Dio(BaseOptions(
      baseUrl: ApiConfig.baseUrl,
      connectTimeout: Duration(milliseconds: ApiConfig.connectTimeout),
      receiveTimeout: Duration(milliseconds: ApiConfig.receiveTimeout),
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
    ));

    // 添加請求攔截器
    _dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) async {
        // 自動添加認證標頭
        if (_authToken != null) {
          options.headers['Authorization'] = 'Bearer $_authToken';
        }
        handler.next(options);
      },
      onError: (error, handler) {
        // 統一錯誤處理
        _handleError(error);
        handler.next(error);
      },
    ));

    // 從本地儲存載入 token
    await _loadAuthToken();
  }

  /// 從本地儲存載入認證 token
  Future<void> _loadAuthToken() async {
    final prefs = await SharedPreferences.getInstance();
    _authToken = prefs.getString('auth_token');
  }

  /// 設定認證 token
  Future<void> setAuthToken(String token) async {
    _authToken = token;
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('auth_token', token);
  }

  /// 清除認證 token
  Future<void> clearAuthToken() async {
    _authToken = null;
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('auth_token');
  }

  /// 獲取 Dio 實例
  Dio get dio => _dio;

  /// 檢查是否已認證
  bool get isAuthenticated => _authToken != null;

  /// 統一錯誤處理
  void _handleError(DioException error) {
    switch (error.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.sendTimeout:
      case DioExceptionType.receiveTimeout:
        throw Exception('網路連線超時，請檢查網路狀態');
      case DioExceptionType.badResponse:
        final statusCode = error.response?.statusCode;
        final message = _getErrorMessage(error.response?.data);
        
        if (statusCode == 401) {
          // 認證失敗，清除本地 token
          clearAuthToken();
          throw Exception('認證已過期，請重新登入');
        } else if (statusCode == 403) {
          throw Exception('權限不足');
        } else if (statusCode == 404) {
          throw Exception('請求的資源不存在');
        } else if (statusCode == 500) {
          throw Exception('伺服器內部錯誤，請稍後再試');
        } else {
          throw Exception(message ?? '請求失敗');
        }
      case DioExceptionType.cancel:
        throw Exception('請求已取消');
      case DioExceptionType.unknown:
        throw Exception('網路連線失敗，請檢查網路狀態');
      default:
        throw Exception('未知錯誤');
    }
  }

  /// 從回應中提取錯誤訊息
  String? _getErrorMessage(dynamic responseData) {
    if (responseData is Map<String, dynamic>) {
      if (responseData.containsKey('error')) {
        final error = responseData['error'];
        if (error is Map<String, dynamic> && error.containsKey('message')) {
          return error['message'] as String?;
        }
      }
      if (responseData.containsKey('message')) {
        return responseData['message'] as String?;
      }
    }
    return null;
  }
}
