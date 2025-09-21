// lib/main.dart

import 'package:flutter/material.dart';

import 'core/theme.dart';
import 'core/api_client.dart';
import 'pages/login_page.dart';
import 'pages/sign_up_page.dart';
import 'pages/forgot_code_page.dart';
import 'pages/forgot_reset_page.dart';
import 'pages/home_page.dart';
import 'pages/chat_page.dart';
import 'pages/maps_page.dart';
import 'pages/profile_page.dart';
import 'pages/notify_page.dart';
import 'pages/article_page.dart';
import 'pages/quiz_page.dart';
import 'pages/edit_nickname_page.dart';
import 'pages/change_password_page.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // 初始化 API 客戶端
  await ApiClient().initialize();
  
  runApp(const MindHelpApp());
}

class MindHelpApp extends StatelessWidget {
  const MindHelpApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'MindHelp',
      theme: AppTheme.lightTheme,
      debugShowCheckedModeBanner: false,
      initialRoute: '/login',
      routes: {
        '/notify': (_) => NotifyPage(),
        '/login': (_) => LoginPage(),
        '/signup': (_) => SignUpPage(),
        '/forgot_code': (_) => ForgotCodePage(),
        '/forgot_reset': (_) => ForgotResetPage(),
        '/home': (_) => HomePage(),
        '/chat': (_) => ChatPage(),
        '/maps': (_) => MapsPage(),
        '/profile': (_) => ProfilePage(),
        '/articles': (_) => ArticlePage(),
        '/quiz': (_) => QuizPage(),
        '/edit_nickname': (_) => EditNicknamePage(),
        '/change_password': (_) => ChangePasswordPage(),
      },
    );
  }
}
