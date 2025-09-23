// lib/pages/chat_page.dart
import 'package:flutter/material.dart';

import '../core/theme.dart';
import '../models/chat_message.dart';
import '../services/ai_service.dart'; // 導入新的 AI 服務
import '../widgets/input_field.dart';

class ChatPage extends StatefulWidget {
  const ChatPage({super.key});

  @override
  ChatPageState createState() => ChatPageState();
}

class ChatPageState extends State<ChatPage> {
  final TextEditingController _ctrl = TextEditingController();
  final ScrollController _scrollController = ScrollController();
  final AiService _aiService = AiService(); // 實例化服務

  final List<ChatMessage> _msgs = [];
  bool _isBotTyping = false;
  bool _isInputEmpty = true; // 新增變數來追蹤輸入框狀態

  // 定義 AI 的系統提示（角色設定）
  final String _systemPrompt = "你是一個有幫助的心理健康助手，專注於提供支持與鼓勵。請用溫暖且富有同理心的語氣回答問題。";

  @override
  void initState() {
    super.initState();
    // 監聽輸入框的變化
    _ctrl.addListener(_updateInputState);
  }

  @override
  void dispose() {
    _ctrl.removeListener(_updateInputState);
    _ctrl.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  void _updateInputState() {
    setState(() {
      _isInputEmpty = _ctrl.text.trim().isEmpty;
    });
  }

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
      _isBotTyping = true;
      _ctrl.clear(); // 在此處清除文字
    });

    _scrollToBottom();

    try {
      final botReply = await _aiService.getOpenRouterCompletion(
        userMessage: userMsg.content,
        systemPrompt: _systemPrompt,
      );
      final botMsg = ChatMessage(
        role: 'bot',
        content: botReply,
        timestamp: DateTime.now().millisecondsSinceEpoch,
      );
      setState(() {
        _msgs.add(botMsg);
        _isBotTyping = false;
      });
    } catch (e) {
      String errorText = '無法連接到伺服器：$e';
      if (e.toString().contains('402')) {
        errorText = 'API 連線失敗：錯誤 402。請檢查您的 OpenRouter API Key 或帳戶餘額。';
      } else if (e.toString().contains('401')) {
        errorText = 'API 連線失敗：錯誤 401。API Key 無效或授權失敗。';
      }

      final errorMsg = ChatMessage(
        role: 'bot',
        content: errorText,
        timestamp: DateTime.now().millisecondsSinceEpoch,
      );
      setState(() {
        _msgs.add(errorMsg);
        _isBotTyping = false;
      });
    }
    _scrollToBottom();
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
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
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
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        if (m.content.isNotEmpty)
                          Text(
                            m.content,
                            style: TextStyle(
                              color: isUser ? Colors.white : AppColors.textBody,
                            ),
                          ),
                      ],
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
                  onPressed: _isBotTyping || _isInputEmpty ? null : _send,
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
        ],
      ),
    );
  }
}
