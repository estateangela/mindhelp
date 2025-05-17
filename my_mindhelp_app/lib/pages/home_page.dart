import 'package:flutter/material.dart';
import '../widgets/primary_button.dart';

class HomePage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('MindHelp Home')),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            PrimaryButton(
                text: 'AI 諮詢聊天',
                onPressed: () => Navigator.pushNamed(context, '/chat')),
            SizedBox(height: 16),
            PrimaryButton(
                text: '尋找心理師',
                onPressed: () => Navigator.pushNamed(context, '/counselors')),
          ],
        ),
      ),
    );
  }
}
