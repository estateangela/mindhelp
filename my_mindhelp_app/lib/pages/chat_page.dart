// lib/pages/chat_page.dart

import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/input_field.dart';

class ChatPage extends StatefulWidget {
  @override
  _ChatPageState createState() => _ChatPageState();
}

class _ChatPageState extends State<ChatPage> {
  final List<Map<String, String>> _msgs = [];
  final _ctrl = TextEditingController();

  void _send() {
    final txt = _ctrl.text.trim();
    if (txt.isEmpty) return;
    setState(() => _msgs.add({'role': 'user', 'text': txt}));
    _ctrl.clear();

    // TODO: call your LLM API
    Future.delayed(const Duration(milliseconds: 500), () {
      setState(() {
        _msgs.add({'role': 'bot', 'text': '這是一個示範回覆。'});
      });
    });
  }

  void _onNavTap(int index) {
    switch (index) {
      case 0:
        Navigator.pushReplacementNamed(context, '/home');
        break;
      case 1:
        Navigator.pushReplacementNamed(context, '/counselors');
        break;
      case 2:
        // already on chat
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
        title: Text('Chat', style: Theme.of(context).textTheme.headlineLarge),
        backgroundColor: Colors.transparent,
        elevation: 0,
        centerTitle: true,
      ),
      body: Column(
        children: [
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
                    child: Text(
                      m['text']!,
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
