// lib/api/openrouter_api.dart

import 'dart:convert';
import 'package:http/http.dart' as http;

class OpenRouterApi {
  static const _apiKey =
      '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';
  // 继续使用 instruction + prompt 模式的模型 ID
  static const _modelId = 'google/gemma-3n-e4b-it:free';
  static const _baseUrl = 'https://openrouter.ai/api/v1/chat/completions';

  /// 向 Gemma-3N 模型发送 prompt，默认让它用繁体中文回答
  static Future<String> sendPrompt({
    required String prompt,
    double temperature = 0.7,
    int maxTokens = 512,
  }) async {
    final uri = Uri.parse(_baseUrl);
    final headers = {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $_apiKey',
    };

    final body = jsonEncode({
      'model': _modelId,
      'instruction': '你是一個心理諮詢助理，請溫和地用繁體中文回覆所有訊息。',
      'prompt': prompt,
      'temperature': temperature,
      'max_tokens': maxTokens,
    });

    final response = await http
        .post(uri, headers: headers, body: body)
        .timeout(const Duration(seconds: 30), onTimeout: () {
      throw Exception('網路請求超時');
    });

    if (response.statusCode == 200) {
      // 用 utf8.decode 保证中文不会乱码
      final decoded = utf8.decode(response.bodyBytes);
      final data = jsonDecode(decoded);

      // 注意：使用 instruction + prompt 返回时，回答在 choices[i].text
      final choices = data['choices'] as List<dynamic>?;
      if (choices == null || choices.isEmpty) {
        throw Exception('OpenRouter 回傳結構異常：choices 為空');
      }

      final firstChoice = choices[0] as Map<String, dynamic>;
      final content = firstChoice['text'] as String?;
      if (content == null) {
        throw Exception('OpenRouter 回傳結構異常：choices[0]["text"] 為 null');
      }
      return content.trim();
    } else {
      // 错误时也用 utf8.decode 看到完整 JSON
      final errText = utf8.decode(response.bodyBytes);
      throw Exception('OpenRouter API 錯誤: ${response.statusCode} $errText');
    }
  }
}
