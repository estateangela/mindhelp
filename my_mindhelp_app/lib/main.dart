import 'package:flutter/material.dart';
import 'pages/home_page.dart';
import 'pages/article_page.dart';
import 'pages/maps_page.dart';
import 'pages/chat_page.dart';
import 'pages/notify_page.dart';
import 'pages/splash_page.dart';
import 'core/theme.dart';
import 'pages/quiz_landing_page.dart.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  // 移除 Supabase 初始化
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
      initialRoute: '/splash',
      routes: {
        '/splash': (context) => const SplashPage(),
        '/notify': (context) => NotifyPage(),
        '/home': (context) => HomePage(),
        '/chat': (context) => ChatPage(),
        '/maps': (context) => MapsPage(),
        '/articles': (context) => ArticlePage(),
        '/quiz_landing': (context) => QuizLandingPage(),
      },
    );
  }
}
