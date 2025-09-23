import 'package:flutter/material.dart';
import 'pages/home_page.dart';
import 'pages/article_page.dart';
import 'pages/quiz_page.dart';
import 'pages/maps_page.dart';
import 'pages/chat_page.dart';
import 'pages/notify_page.dart';
import 'core/theme.dart';
import 'config/secrets.dart';
import 'pages/splash_page.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  runApp(const MindHelpApp());
}

class MindHelpApp extends StatelessWidget {
  const MindHelpApp({super.key});

  @override
  Widget build(BuildContext context) {
    // 根據平台設定 Google Maps API 金鑰
    // 注意：這部分邏輯不適用於 Web 端，Web 端的 API 金鑰需在 index.html 中設定
    if (Theme.of(context).platform == TargetPlatform.android) {
      // 在 Android 上，API 金鑰應在 AndroidManifest.xml 中設定
    } else if (Theme.of(context).platform == TargetPlatform.iOS) {
      // 在 iOS 上，API 金鑰應在 AppDelegate.swift 中設定
    }

    return MaterialApp(
      title: 'MindHelp',
      theme: AppTheme.lightTheme,
      debugShowCheckedModeBanner: false,
      initialRoute: '/splash',
      routes: {
        '/notify': (context) => NotifyPage(),
        '/home': (context) => HomePage(),
        '/chat': (context) => ChatPage(),
        '/maps': (context) => MapsPage(),
        '/articles': (context) => ArticlePage(),
        '/quiz': (context) => QuizPage(),
        '/splash': (context) => SplashPage(),
      },
    );
  }
}
