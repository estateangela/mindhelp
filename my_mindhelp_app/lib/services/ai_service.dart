import 'package:http/http.dart' as http;
import 'dart:convert';
import 'dart:typed_data';

class AiService {
  static const String openRouterApiKey =
      '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';

  Future<String> getOpenRouterCompletion({
    required String userMessage,
    required String systemPrompt,
    Uint8List? imageBytes,
  }) async {
    if (openRouterApiKey.isEmpty) {
      throw Exception('OpenRouter API Key is not set.');
    }

    const String apiUrl = 'https://openrouter.ai/api/v1/chat/completions';
    const String model = 'google/gemini-2.0-flash-exp:free';

    // 構建訊息列表，根據是否有圖片動態添加
    List<Map<String, dynamic>> messages = [
      {'role': 'system', 'content': systemPrompt},
    ];

    List<Map<String, dynamic>> userContent = [];

    // 添加文字內容
    userContent.add({
      'type': 'text',
      'text': userMessage,
    });

    // 如果有圖片，添加圖片內容並進行 Base64 編碼
    if (imageBytes != null) {
      String base64Image = base64Encode(imageBytes);
      userContent.add({
        'type': 'image_url',
        'image_url': {
          'url': 'data:image/png;base64,$base64Image',
        },
      });
    }

    messages.add({
      'role': 'user',
      'content': userContent,
    });

    final response = await http.post(
      Uri.parse(apiUrl),
      headers: {
        'Authorization': 'Bearer $openRouterApiKey',
        'Content-Type': 'application/json',
      },
      body: jsonEncode({
        'model': model,
        'messages': messages,
      }),
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return data['choices'][0]['message']['content'];
    } else {
      throw Exception(
          'Failed to load response from OpenRouter: ${response.statusCode}');
    }
  }
}
