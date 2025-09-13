// lib/main.dart

import 'package:flutter/material.dart';
import 'package:firebase_core/firebase_core.dart';

import 'core/theme.dart';
import 'pages/login_page.dart';
import 'pages/sign_up_page.dart';
import 'pages/forgot_code_page.dart';
import 'pages/forgot_reset_page.dart';
import 'pages/home_page.dart';
import 'pages/chat_page.dart';
import 'pages/maps_page.dart';
import 'pages/profile_page.dart';
import 'pages/notify_page.dart';
import 'pages/settings_page.dart';
import 'pages/article_page.dart';
import 'pages/quiz_page.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  runApp(const MindHelpApp());
}

class MindHelpApp extends StatelessWidget {
  const MindHelpApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'MindHelp',
      theme: AppTheme.lightTheme,
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
        '/articles': (_) => ArticlePage(),
        '/quiz': (_) => QuizPage(),
      },
    );
  }
}
