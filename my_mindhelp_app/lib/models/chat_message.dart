// lib/models/chat_message.dart

class ChatMessage {
  final int? id;
  final String role; // 'user' or 'bot'
  final String content;
  final int timestamp; // Unix milliseconds

  ChatMessage({
    this.id,
    required this.role,
    required this.content,
    required this.timestamp,
  });

  Map<String, dynamic> toMap() => {
        'id': id,
        'role': role,
        'content': content,
        'timestamp': timestamp,
      };

  factory ChatMessage.fromMap(Map<String, dynamic> map) => ChatMessage(
        id: map['id'] as int?,
        role: map['role'] as String,
        content: map['content'] as String,
        timestamp: map['timestamp'] as int,
      );

  bool get isUser => role == 'user';
}
