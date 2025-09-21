import 'package:json_annotation/json_annotation.dart';

part 'chat_message.g.dart';

@JsonSerializable()
class ChatMessage {
  final String id;
  final String role; // 'user' or 'assistant'
  final String content;
  final String createdAt;
  final String sessionId;

  ChatMessage({
    required this.id,
    required this.role,
    required this.content,
    required this.createdAt,
    required this.sessionId,
  });

  factory ChatMessage.fromJson(Map<String, dynamic> json) => _$ChatMessageFromJson(json);
  Map<String, dynamic> toJson() => _$ChatMessageToJson(this);

  bool get isUser => role == 'user';
  bool get isAssistant => role == 'assistant';
}

@JsonSerializable()
class ChatSession {
  final String id;
  final String firstMessageSnippet;
  final String lastUpdatedAt;
  final int messageCount;

  ChatSession({
    required this.id,
    required this.firstMessageSnippet,
    required this.lastUpdatedAt,
    required this.messageCount,
  });

  factory ChatSession.fromJson(Map<String, dynamic> json) => _$ChatSessionFromJson(json);
  Map<String, dynamic> toJson() => _$ChatSessionToJson(this);
}
