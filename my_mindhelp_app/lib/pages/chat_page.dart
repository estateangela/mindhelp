// lib/pages/chat_page.dart

import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../models/chat_message.dart';
import '../utils/db_helper.dart';
import '../widgets/input_field.dart';

class ChatPage extends StatefulWidget {
  const ChatPage({Key? key}) : super(key: key);

  @override
  _ChatPageState createState() => _ChatPageState();
}

class _ChatPageState extends State<ChatPage> {
  final TextEditingController _ctrl = TextEditingController();
  List<ChatMessage> _msgs = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadHistory();
  }

  Future<void> _loadHistory() async {
    final hist = await DBHelper().getAllMessages();
    setState(() {
      _msgs = hist;
      _isLoading = false;
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
    await DBHelper().insertMessage(userMsg);
    setState(() => _msgs.add(userMsg));
    _ctrl.clear();

    // 模擬機器人回覆
    Future.delayed(const Duration(milliseconds: 500), () async {
      final botMsg = ChatMessage(
        role: 'bot',
        content: '這是一個示範回覆。',
        timestamp: DateTime.now().millisecondsSinceEpoch,
      );
      await DBHelper().insertMessage(botMsg);
      setState(() => _msgs.add(botMsg));
    });
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
          icon: Icon(Icons.settings, color: AppColors.textHigh),
          onPressed: () => Navigator.pushNamed(context, '/settings'),
        ),
        title: Text('Chat', style: Theme.of(context).textTheme.headlineLarge),
        centerTitle: true,
        backgroundColor: Colors.transparent,
        elevation: 0,
        actions: [
          IconButton(
            icon: Icon(Icons.notifications, color: AppColors.textHigh),
            onPressed: () => Navigator.pushNamed(context, '/notify'),
          ),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : Column(
              children: [
                Expanded(
                  child: ListView.builder(
                    padding: const EdgeInsets.all(16),
                    itemCount: _msgs.length,
                    itemBuilder: (_, i) {
                      final m = _msgs[i];
                      final isUser = m.role == 'user';
                      return Align(
                        alignment: isUser
                            ? Alignment.centerRight
                            : Alignment.centerLeft,
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
                  padding:
                      const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
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
                        icon: Icon(Icons.send, color: AppColors.accent),
                        onPressed: _send,
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
