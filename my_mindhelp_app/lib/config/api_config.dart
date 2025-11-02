class ApiConfig {
  // 允許以 --dart-define=API_BASE_URL=... 覆寫
  static const String _overrideBaseUrl = String.fromEnvironment('API_BASE_URL');

  static String get baseUrl {
    if (_overrideBaseUrl.isNotEmpty) {
      return _overrideBaseUrl;
    }
    // 預設使用雲端後端
    return 'https://mindhelp.onrender.com/api/v1';
  }
}


