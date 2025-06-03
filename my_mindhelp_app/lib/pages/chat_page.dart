// lib/pages/chat_page.dart

import 'package:flutter/material.dart';
import 'package:flutter_markdown/flutter_markdown.dart';
import '../core/theme.dart';
import '../widgets/input_field.dart';
import '../widgets/custom_app_bar.dart';
import '../api/openrouter_api.dart';
import '../utils/db_helper.dart';

class ChatPage extends StatefulWidget {
  const ChatPage({Key? key}) : super(key: key);

  @override
  _ChatPageState createState() => _ChatPageState();
}

class _ChatPageState extends State<ChatPage> {
  final List<Map<String, String>> _msgs = [];
  final TextEditingController _ctrl = TextEditingController();
  final ScrollController _scrollCtrl = ScrollController();
  bool _isWaiting = false;
  bool _hasNewMessages = false;

  @override
  void initState() {
    super.initState();
    _loadHistory();
  }

  @override
  void dispose() {
    _ctrl.dispose();
    _scrollCtrl.dispose();
    super.dispose();
  }

  Future<void> _loadHistory() async {
    final hist = await DBHelper().getAllMessages();
    setState(() {
      _msgs.clear();
      for (final msg in hist) {
        _msgs.add({'role': msg.role, 'text': msg.content});
      }
    });
    _scrollToBottom();
  }

  void _scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scrollCtrl.hasClients) {
        _scrollCtrl.jumpTo(_scrollCtrl.position.maxScrollExtent);
      }
    });
  }

  Future<void> _saveAllPending() async {
    // 目前每次執行 _send() 都會立刻 insertMessage，
    // 所以這裡直接重置旗標即可
    _hasNewMessages = false;
  }

  Future<bool> _onWillPop() async {
    if (!_hasNewMessages) return true;

    final result = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('離開前是否要儲存對話？'),
        content: const Text('尚有新訊息，請選擇：'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(ctx).pop(false),
            child: const Text('不儲存'),
          ),
          TextButton(
            onPressed: () => Navigator.of(ctx).pop(true),
            child: const Text('儲存並退出'),
          ),
        ],
      ),
    );

    if (result == true) {
      await _saveAllPending();
      return true;
    } else if (result == false) {
      return true;
    } else {
      return false; // 點 dialog 外面或按取消就不離開
    }
  }

  void _send() async {
    final txt = _ctrl.text.trim();
    if (txt.isEmpty) return;

    // 1. 先把使用者訊息存入 SQLite (mobile) 或 in‐memory (web)
    final userMsg = ChatMessage(
      role: 'user',
      content: txt,
      timestamp: DateTime.now().millisecondsSinceEpoch,
    );
    await DBHelper().insertMessage(userMsg);

    setState(() {
      _msgs.add({'role': 'user', 'text': txt});
      _isWaiting = true;
      _hasNewMessages = true;
    });
    _ctrl.clear();
    _scrollToBottom();

    try {
      // 2. 呼叫 OpenRouter AI API 取得回覆
      final answer = await OpenRouterApi.sendPrompt(prompt: txt);

      // 3. 把 AI 回覆也存入資料庫
      final botMsg = ChatMessage(
        role: 'assistant',
        content: answer,
        timestamp: DateTime.now().millisecondsSinceEpoch,
      );
      await DBHelper().insertMessage(botMsg);

      setState(() {
        _msgs.add({'role': 'assistant', 'text': answer});
      });
    } catch (e) {
      final err = '呼叫 AI 失敗：$e';
      final errorMsg = ChatMessage(
        role: 'assistant',
        content: err,
        timestamp: DateTime.now().millisecondsSinceEpoch,
      );
      await DBHelper().insertMessage(errorMsg);

      setState(() {
        _msgs.add({'role': 'assistant', 'text': err});
      });
    } finally {
      setState(() {
        _isWaiting = false;
      });
      _scrollToBottom();
    }
  }

  @override
  Widget build(BuildContext context) {
    return WillPopScope(
      onWillPop: _onWillPop,
      child: Scaffold(
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
            Expanded(
              child: ListView.builder(
                controller: _scrollCtrl,
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
                        border: isUser
                            ? null
                            : Border.all(color: AppColors.accent),
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: MarkdownBody(
                        data: m['text']!,
                        styleSheet: MarkdownStyleSheet(
                          p: TextStyle(
                            color:
                                isUser ? Colors.white : AppColors.textBody,
                            fontSize: 15,
                            height: 1.4,
                          ),
                        ),
                      ),
                    ),
                  );
                },
              ),
            ),

            if (_isWaiting)
              const Padding(
                padding: EdgeInsets.only(bottom: 8.0),
                child: CircularProgressIndicator(),
              ),

            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
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
          onTap: (i) {
            switch (i) {
              case 0:
                Navigator.pushReplacementNamed(context, '/home');
                break;
              case 1:
                Navigator.pushReplacementNamed(context, '/maps');
                break;
              case 2:
                // 已在 Chat
                break;
              case 3:
                Navigator.pushReplacementNamed(context, '/profile');
                break;
            }
          },
          items: const [
            BottomNavigationBarItem(icon: Icon(Icons.home), label: 'Home'),
            BottomNavigationBarItem(
                icon: Icon(Icons.location_on), label: 'Maps'),
            BottomNavigationBarItem(
                icon: Icon(Icons.chat_bubble), label: 'Chat'),
            BottomNavigationBarItem(
                icon: Icon(Icons.person), label: 'Profile'),
          ],
        ),
      ),
    );
  }
}
