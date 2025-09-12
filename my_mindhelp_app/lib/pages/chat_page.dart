// lib/pages/chat_page.dart
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import '../core/theme.dart';
import '../models/chat_message.dart';
import '../widgets/input_field.dart';

class ChatPage extends StatefulWidget {
  const ChatPage({Key? key}) : super(key: key);

  @override
  _ChatPageState createState() => _ChatPageState();
}

class _ChatPageState extends State<ChatPage> {
  final TextEditingController _ctrl = TextEditingController();
  final ScrollController _scrollController = ScrollController();

  // 使用一個臨時列表來儲存聊天記錄，不使用任何資料庫
  List<ChatMessage> _msgs = [];
  bool _isBotTyping = false;

  void _scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scrollController.hasClients) {
        _scrollController.animateTo(
          _scrollController.position.maxScrollExtent,
          duration: const Duration(milliseconds: 300),
          curve: Curves.easeOut,
        );
      }
    });
  }

  Future<void> _send() async {
    final txt = _ctrl.text.trim();
    if (txt.isEmpty) return;

    final userMsg = ChatMessage(
      role: 'user',
      content: txt,
      timestamp: DateTime.now().millisecondsSinceEpoch,
    );

    setState(() {
      _msgs.add(userMsg);
      _isBotTyping = true; // 顯示機器人正在打字
    });

    _ctrl.clear();
    _scrollToBottom();

    try {
      final botReply = await _getBotResponse(userMsg.content);
      final botMsg = ChatMessage(
        role: 'bot',
        content: botReply,
        timestamp: DateTime.now().millisecondsSinceEpoch,
      );
      setState(() {
        _msgs.add(botMsg);
        _isBotTyping = false; // 機器人回覆完成
      });
    } catch (e) {
      final errorMsg = ChatMessage(
        role: 'bot',
        content: '無法連接到伺服器：$e',
        timestamp: DateTime.now().millisecondsSinceEpoch,
      );
      setState(() {
        _msgs.add(errorMsg);
        _isBotTyping = false; // 機器人回覆完成
      });
    }
    _scrollToBottom();
  }

  Future<String> _getBotResponse(String userMessage) async {
    // TODO: 將 'https://your-backend.com/api/chat' 替換為您的後端 API 端點
    const String backendUrl = 'https://your-backend.com/api/chat';

    final response = await http.post(
      Uri.parse(backendUrl),
      headers: {
        'Content-Type': 'application/json',
      },
      body: jsonEncode({
        'message': userMessage,
      }),
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      // TODO: 根據您的後端回傳格式調整解析邏輯
      return data['reply'];
    } else {
      throw Exception(
          'Failed to load response from backend: ${response.statusCode}');
    }
  }

  void _onNavTap(int idx) {
    switch (idx) {
      case 0:
        Navigator.pushReplacementNamed(context, '/home');
        break;
      case 1:
        Navigator.pushReplacementNamed(context, '/maps');
        break;
      case 2:
        // already on Chat
        break;
      case 3:
        Navigator.pushReplacementNamed(context, '/profile');
        break;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.settings, color: AppColors.textHigh),
          onPressed: () => Navigator.pushNamed(context, '/settings'),
        ),
        title: Text('Chat', style: Theme.of(context).textTheme.headlineLarge),
        centerTitle: true,
        backgroundColor: Colors.transparent,
        elevation: 0,
        actions: [
          IconButton(
            icon: const Icon(Icons.notifications, color: AppColors.textHigh),
            onPressed: () => Navigator.pushNamed(context, '/notify'),
          ),
        ],
      ),
      body: Column(
        children: [
          Expanded(
            child: ListView.builder(
              controller: _scrollController,
              padding: const EdgeInsets.all(16),
              itemCount: _msgs.length + (_isBotTyping ? 1 : 0),
              itemBuilder: (_, i) {
                if (_isBotTyping && i == _msgs.length) {
                  return Align(
                    alignment: Alignment.centerLeft,
                    child: Container(
                      margin: const EdgeInsets.symmetric(vertical: 6),
                      padding: const EdgeInsets.all(12),
                      decoration: BoxDecoration(
                        color: Colors.white,
                        border: Border.all(color: AppColors.accent),
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: const Text(
                        '機器人正在回覆...',
                        style: TextStyle(
                            fontStyle: FontStyle.italic,
                            color: AppColors.textBody),
                      ),
                    ),
                  );
                }
                final m = _msgs[i];
                final isUser = m.role == 'user';
                return Align(
                  alignment:
                      isUser ? Alignment.centerRight : Alignment.centerLeft,
                  child: Container(
                    margin: const EdgeInsets.symmetric(vertical: 6),
                    padding: const EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: isUser ? AppColors.accent : Colors.white,
                      border:
                          isUser ? null : Border.all(color: AppColors.accent),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      m.content,
                      style: TextStyle(
                        color: isUser ? Colors.white : AppColors.textBody,
                      ),
                    ),
                  ),
                );
              },
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            child: Row(
              children: [
                Expanded(
                  child: InputField(
                    controller: _ctrl,
                    label: '',
                    prefixIcon: null,
                  ),
                ),
                const SizedBox(width: 8),
                IconButton(
                  icon: const Icon(Icons.send, color: AppColors.accent),
                  onPressed: _isBotTyping ? null : _send,
                ),
              ],
            ),
          ),
        ],
      ),
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: 2,
        selectedItemColor: AppColors.accent,
        unselectedItemColor: AppColors.textBody,
        onTap: _onNavTap,
        items: const [
          BottomNavigationBarItem(icon: Icon(Icons.home), label: 'Home'),
          BottomNavigationBarItem(icon: Icon(Icons.location_on), label: 'Maps'),
          BottomNavigationBarItem(icon: Icon(Icons.chat_bubble), label: 'Chat'),
          BottomNavigationBarItem(icon: Icon(Icons.person), label: 'Profile'),
        ],
      ),
    );
  }
}
