import 'package:http/http.dart' as http;
import 'dart:convert';

class AuthService {
  // 將基底 URL 替換為你的 Cloudflare API 網址
  final String _baseUrl = 'https://api.estateangela.dpdns.org/v1';

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
}
