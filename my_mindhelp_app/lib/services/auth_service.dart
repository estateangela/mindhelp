import 'package:http/http.dart' as http;
import 'dart:convert';

class AuthService {
  final String _baseUrl = 'https://mindhelp.onrender.com/v1';
  final String _authHeader =
      'Bearer <YOUR_JWT_TOKEN>'; // TODO: 請替換為您實際的 JWT Token

  // 新增：使用者登入的 API 呼叫
  Future<Map<String, dynamic>> login({
    required String email,
    required String password,
  }) async {
    final url = Uri.parse('$_baseUrl/auth/login');

    final response = await http.post(
      url,
      headers: {
        'Content-Type': 'application/json',
      },
      body: jsonEncode({
        'email': email,
        'password': password,
      }),
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return data['data']; // 回傳包含 access_token 和 user 資訊的 Map
    } else {
      final errorData = jsonDecode(response.body);
      throw Exception(errorData['error']['message'] ?? '登入失敗，請檢查信箱或密碼。');
    }
  }

  Future<void> register({
    required String email,
    required String password,
    required String nickname,
  }) async {
    final url = Uri.parse('$_baseUrl/auth/register');

    final response = await http.post(
      url,
      headers: {
        'Content-Type': 'application/json',
      },
      body: jsonEncode({
        'email': email,
        'password': password,
        'nickname': nickname,
      }),
    );

    if (response.statusCode == 201) {
      // 註冊成功，不做任何事，讓前端處理導航
    } else {
      final errorData = jsonDecode(response.body);
      throw Exception(errorData['error']['message'] ?? '註冊失敗，請重試。');
    }
  }

  // 新增：更新暱稱的 API 呼叫
  Future<void> updateNickname({
    required String nickname,
  }) async {
    final url = Uri.parse('$_baseUrl/users/me');

    final response = await http.put(
      url,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': _authHeader,
      },
      body: jsonEncode({
        'nickname': nickname,
      }),
    );

    if (response.statusCode == 200) {
      // 更新成功
    } else {
      final errorData = jsonDecode(response.body);
      throw Exception(errorData['error']['message'] ?? '暱稱更新失敗，請重試。');
    }
  }

  // 新增：修改密碼的 API 呼叫
  Future<void> changePassword({
    required String oldPassword,
    required String newPassword,
  }) async {
    final url = Uri.parse('$_baseUrl/users/me/password');

    final response = await http.put(
      url,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': _authHeader,
      },
      body: jsonEncode({
        'oldPassword': oldPassword,
        'newPassword': newPassword,
      }),
    );

    if (response.statusCode == 200) {
      // 密碼更新成功
    } else {
      final errorData = jsonDecode(response.body);
      throw Exception(errorData['error']['message'] ?? '密碼修改失敗，請重試。');
    }
  }
}
