import 'package:http/http.dart' as http;
import 'dart:convert';

class AiService {
  static const String openRouterApiKey =
      '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'; // 換成你自己的金鑰

  Future<String> getOpenRouterCompletion({
    required String userMessage,
    required String systemPrompt,
  }) async {
    if (openRouterApiKey.isEmpty) {
      throw Exception('OpenRouter API Key is not set.');
    }

    const String apiUrl = 'https://openrouter.ai/api/v1/chat/completions';
    const String model = 'openai/gpt-oss-120b:free'; // 用一個確定存在的模型

    final response = await http.post(
      Uri.parse(apiUrl),
      headers: {
        'Authorization': 'Bearer $openRouterApiKey',
        'Content-Type': 'application/json',
      },
      body: jsonEncode({
        'model': model,
        'messages': [
          {'role': 'system', 'content': systemPrompt},
          {'role': 'user', 'content': userMessage},
        ],
      }),
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return data['choices'][0]['message']['content'];
    } else {
      throw Exception(
          'Failed to load response: ${response.statusCode}, body: ${response.body}');
    }
  }
}
