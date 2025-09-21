import '../core/api_client.dart';
import '../core/api_config.dart';
import '../models/chat_message.dart';

class ChatService {
  final ApiClient _apiClient = ApiClient();

  /// 獲取聊天會話列表
  Future<List<ChatSession>> getChatSessions({
    int page = 1,
    int limit = 20,
  }) async {
    final response = await _apiClient.dio.get(
      '${ApiConfig.chat}/sessions',
      queryParameters: {
        'page': page,
        'limit': limit,
      },
    );

    if (response.statusCode == 200) {
      final List<dynamic> sessions = response.data['data'];
      return sessions.map((json) => ChatSession.fromJson(json)).toList();
    }
    throw Exception('獲取聊天會話失敗');
  }

  /// 建立新的聊天會話
  Future<ChatSession> createChatSession() async {
    final response = await _apiClient.dio.post('${ApiConfig.chat}/sessions');

    if (response.statusCode == 201) {
      return ChatSession.fromJson(response.data['data']);
    }
    throw Exception('建立聊天會話失敗');
  }

  /// 獲取會話中的訊息歷史
  Future<List<ChatMessage>> getChatMessages(String sessionId) async {
    final response = await _apiClient.dio.get(
      '${ApiConfig.chat}/sessions/$sessionId/messages',
    );

    if (response.statusCode == 200) {
      final List<dynamic> messages = response.data['data'];
      return messages.map((json) => ChatMessage.fromJson(json)).toList();
    }
    throw Exception('獲取聊天訊息失敗');
  }

  /// 發送訊息並獲取 AI 回覆
  Future<ChatMessage> sendMessage({
    required String sessionId,
    required String content,
  }) async {
    final response = await _apiClient.dio.post(
      '${ApiConfig.chat}/sessions/$sessionId/messages',
      data: {
        'content': content,
      },
    );

    if (response.statusCode == 200) {
      return ChatMessage.fromJson(response.data['data']);
    }
    throw Exception('發送訊息失敗');
  }
}
