// lib/main.dart

import 'package:flutter/material.dart';
import 'core/theme.dart';

// 仅从自己的页面文件里导入需要的类，避免重复导入
import 'pages/login_page.dart';
import 'pages/sign_up_page.dart';
import 'pages/forgot_code_page.dart';
import 'pages/forgot_reset_page.dart';
import 'pages/home_page.dart';
import 'pages/chat_page.dart';
import 'pages/map_page.dart';
import 'pages/profile_page.dart';
import 'pages/notify_page.dart';
import 'pages/settings_page.dart';

void main() {
  runApp(MindHelpApp());
}

class MindHelpApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'MindHelp',
      theme: AppTheme.lightTheme,
      debugShowCheckedModeBanner: false,
      // 初始路由设为登录页
      initialRoute: '/login',
      routes: {
        '/notify': (_) => NotifyPage(),
        '/settings': (_) => SettingsPage(),
        '/login': (_) => LoginPage(),
        '/signup': (_) => SignUpPage(),
        '/forgot_code': (_) => ForgotCodePage(),
        '/forgot_reset': (_) => ForgotResetPage(),
        '/home': (_) => HomePage(),
        '/chat': (_) => ChatPage(),
        '/maps': (_) => MapsPage(),
        '/profile': (_) => ProfilePage(),
      },
    );
  }
}
