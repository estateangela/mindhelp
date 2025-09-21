import '../core/api_client.dart';
import '../core/api_config.dart';
import '../models/user.dart';

class AuthService {
  final ApiClient _apiClient = ApiClient();

  /// 註冊新使用者
  Future<AuthResponse> register({
    required String email,
    required String password,
    String? nickname,
  }) async {
    final request = RegisterRequest(
      email: email,
      password: password,
      nickname: nickname,
    );

    final response = await _apiClient.dio.post(
      '${ApiConfig.auth}/register',
      data: request.toJson(),
    );

    if (response.statusCode == 201) {
      final authResponse = AuthResponse.fromJson(response.data['data']);
      await _apiClient.setAuthToken(authResponse.token);
      return authResponse;
    }
    throw Exception('註冊失敗');
  }

  /// 登入
  Future<AuthResponse> login({
    required String email,
    required String password,
  }) async {
    final request = LoginRequest(
      email: email,
      password: password,
    );

    final response = await _apiClient.dio.post(
      '${ApiConfig.auth}/login',
      data: request.toJson(),
    );

    if (response.statusCode == 200) {
      final authResponse = AuthResponse.fromJson(response.data['data']);
      await _apiClient.setAuthToken(authResponse.token);
      return authResponse;
    }
    throw Exception('登入失敗');
  }

  /// 登出
  Future<void> logout() async {
    await _apiClient.dio.post('${ApiConfig.auth}/logout');
    await _apiClient.clearAuthToken();
  }

  /// 獲取當前使用者資訊
  Future<User> getCurrentUser() async {
    final response = await _apiClient.dio.get(ApiConfig.users + '/me');
    
    if (response.statusCode == 200) {
      return User.fromJson(response.data['data']);
    }
    throw Exception('獲取使用者資訊失敗');
  }

  /// 更新暱稱
  Future<User> updateNickname({
    required String nickname,
  }) async {
    final request = UpdateUserRequest(nickname: nickname);

    final response = await _apiClient.dio.put(
      ApiConfig.users + '/me',
      data: request.toJson(),
    );

    if (response.statusCode == 200) {
      return User.fromJson(response.data['data']);
    }
    throw Exception('暱稱更新失敗');
  }

  /// 修改密碼
  Future<void> changePassword({
    required String oldPassword,
    required String newPassword,
  }) async {
    final request = ChangePasswordRequest(
      oldPassword: oldPassword,
      newPassword: newPassword,
    );

    final response = await _apiClient.dio.put(
      '${ApiConfig.users}/me/password',
      data: request.toJson(),
    );

    if (response.statusCode != 200) {
      throw Exception('密碼修改失敗');
    }
  }

  /// 刪除帳號
  Future<void> deleteAccount() async {
    final response = await _apiClient.dio.delete(ApiConfig.users + '/me');
    
    if (response.statusCode == 200) {
      await _apiClient.clearAuthToken();
    } else {
      throw Exception('帳號刪除失敗');
    }
  }

  /// 檢查是否已登入
  bool get isLoggedIn => _apiClient.isAuthenticated;
}
