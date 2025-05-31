// lib/pages/chat_page.dart

import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter_markdown/flutter_markdown.dart';
import '../core/theme.dart';
import '../widgets/input_field.dart';
import '../widgets/custom_app_bar.dart';
import '../api/openrouter_api.dart';

class ChatPage extends StatefulWidget {
  @override
  _ChatPageState createState() => _ChatPageState();
}

class _ChatPageState extends State<ChatPage> {
  final List<Map<String, String>> _msgs = [];
  final TextEditingController _ctrl = TextEditingController();
  bool _isWaiting = false;

  void _send() async {
    final txt = _ctrl.text.trim();
    if (txt.isEmpty) return;

    // 加入使用者訊息
    setState(() {
      _msgs.add({'role': 'user', 'text': txt});
      _isWaiting = true;
    });
    _ctrl.clear();

    try {
      // 呼叫 OpenRouter API，取得 AI 回覆
      final answer = await OpenRouterApi.sendPrompt(prompt: txt);
      setState(() {
        _msgs.add({'role': 'assistant', 'text': answer});
      });
    } catch (e) {
      setState(() {
        _msgs.add({'role': 'assistant', 'text': '呼叫 OpenRouter 失敗：$e'});
      });
    } finally {
      setState(() {
        _isWaiting = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: const CustomAppBar(
        showBackButton: true,
        titleWidget: Text(
          'AI 諮詢',
          style: TextStyle(fontSize: 24, color: AppColors.textHigh),
        ),
      ),
      body: Column(
        children: [
          // 訊息列表區
          Expanded(
            child: ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: _msgs.length,
              itemBuilder: (_, i) {
                final m = _msgs[i];
                final isUser = m['role'] == 'user';
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
                    child: MarkdownBody(
                      data: m['text']!,
                      styleSheet: MarkdownStyleSheet(
                        p: TextStyle(
                          color: isUser ? Colors.white : AppColors.textBody,
                          fontSize: 15,
                          height: 1.4,
                        ),
                        h1: TextStyle(
                          color: isUser ? Colors.white : AppColors.textHigh,
                          fontSize: 24,
                          fontWeight: FontWeight.bold,
                        ),
                        strong: TextStyle(
                          color: isUser ? Colors.white : AppColors.textHigh,
                          fontWeight: FontWeight.bold,
                        ),
                        em: TextStyle(
                          fontStyle: FontStyle.italic,
                          color: isUser ? Colors.white : AppColors.textHigh,
                        ),
                        listBullet: TextStyle(
                          color: isUser ? Colors.white : AppColors.textBody,
                        ),
                        code: TextStyle(
                          color: isUser ? Colors.white : AppColors.textBody,
                          fontFamily: 'monospace',
                          backgroundColor: isUser
                              ? AppColors.accent.withOpacity(0.3)
                              : Colors.grey[200],
                        ),
                      ),
                      selectable: false,
                    ),
                  ),
                );
              },
            ),
          ),

          // 輪詢中顯示進度指示
          if (_isWaiting)
            const Padding(
              padding: EdgeInsets.only(bottom: 8.0),
              child: CircularProgressIndicator(),
            ),

          // 使用者輸入區
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            child: Row(
              children: [
                Expanded(
                  child: InputField(
                    controller: _ctrl,
                    label: '請輸入詢問內容',
                    prefixIcon: Icons.message_outlined,
                  ),
                ),
                const SizedBox(width: 8),
                IconButton(
                  icon: Icon(Icons.send, color: AppColors.accent),
                  onPressed: _isWaiting ? null : _send,
                )
              ],
            ),
          ),
        ],
      ),

      // 下方導覽列
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: 2, // Chat 索引
        selectedItemColor: AppColors.accent,
        unselectedItemColor: AppColors.textBody,
        onTap: (i) {
          switch (i) {
            case 0:
              Navigator.pushReplacementNamed(context, '/home');
              break;
            case 1:
              Navigator.pushReplacementNamed(context, '/maps');
              break;
            case 2:
              // 已經在 Chat 頁面
              break;
            case 3:
              Navigator.pushReplacementNamed(context, '/profile');
              break;
          }
        },
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
