import 'package:flutter/material.dart';
import '../widgets/input_field.dart';
import '../widgets/primary_button.dart';

class ChatPage extends StatefulWidget {
  @override
  _ChatPageState createState() => _ChatPageState();
}

class _ChatPageState extends State<ChatPage> {
  final List<Map<String, String>> _messages = [];
  final _inputController = TextEditingController();

  void _sendMessage() {
    final text = _inputController.text;
    if (text.isEmpty) return;
    setState(() {
      _messages.add({'role': 'user', 'text': text});
      // TODO: 呼叫 API 拿回 AI 回應
      _messages.add({'role': 'bot', 'text': '這是一個回覆示例。'});
    });
    _inputController.clear();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('AI 諮詢聊天')),
      body: Column(
        children: [
          Expanded(
            child: ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: _messages.length,
              itemBuilder: (_, i) {
                final msg = _messages[i];
                final isUser = msg['role'] == 'user';
                return Align(
                  alignment:
                      isUser ? Alignment.centerRight : Alignment.centerLeft,
                  child: Container(
                    margin: EdgeInsets.symmetric(vertical: 4),
                    padding: EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: isUser ? Colors.blue[100] : Colors.grey[200],
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Text(msg['text']!, style: TextStyle(fontSize: 16)),
                  ),
                );
              },
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(8),
            child: Row(
              children: [
                Expanded(
                    child: InputField(
                        controller: _inputController, label: '輸入訊息')),
                SizedBox(width: 8),
                PrimaryButton(text: '送出', onPressed: _sendMessage),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
