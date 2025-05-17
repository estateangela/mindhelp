import 'package:flutter/material.dart';
import 'core/theme.dart';
import 'pages/login_page.dart';
import 'pages/forgot_password_page.dart';
import 'pages/home_page.dart';
import 'pages/chat_page.dart';
import 'pages/counselor_list_page.dart';
import 'pages/counselor_detail_page.dart';

void main() {
  runApp(MindHelpApp());
}

class MindHelpApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'MindHelp',
      theme: AppTheme.lightTheme,
      initialRoute: '/login',
      routes: {
        '/login': (_) => LoginPage(),
        '/forgot': (_) => ForgotPasswordPage(),
        '/home': (_) => HomePage(),
        '/chat': (_) => ChatPage(),
        '/counselors': (_) => CounselorListPage(),
        '/detail': (_) => CounselorDetailPage(),
      },
    );
  }
}
