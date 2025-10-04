import 'package:http/http.dart' as http;
import 'dart:convert';
import '../config/secrets.dart';

class AiService {
  Future<String> getOpenRouterCompletion({
    required String userMessage,
    required String systemPrompt,
  }) async {
    if (Secrets.openRouterApiKey.isEmpty) {
      throw Exception('OpenRouter API Key is not set in secrets.dart.');
    }

    const String apiUrl = 'https://openrouter.ai/api/v1/chat/completions';
    const String model = 'google/gemini-2.5-flash-lite';

    final response = await http.post(
      Uri.parse(apiUrl),
      headers: {
        'Authorization': 'Bearer ${Secrets.openRouterApiKey}',
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
