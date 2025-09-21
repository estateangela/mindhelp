import 'package:flutter/material.dart';
import 'pages/home_page.dart';
import 'pages/article_page.dart';
import 'pages/quiz_page.dart';
import 'pages/maps_page.dart';
import 'pages/chat_page.dart';
import 'pages/profile_page.dart';
import 'pages/login_page.dart';
import 'pages/sign_up_page.dart';
import 'pages/edit_nickname_page.dart';
import 'pages/change_password_page.dart';
import 'pages/forgot_code_page.dart';
import 'pages/forgot_reset_page.dart';
import 'pages/notify_page.dart';
import 'core/theme.dart';
import 'config/secrets.dart';

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
      initialRoute: '/login',
      routes: {
        '/notify': (context) => const NotifyPage(),
        '/login': (context) => const LoginPage(),
        '/signup': (context) => SignUpPage(),
        '/forgot_code': (context) => const ForgotCodePage(),
        '/forgot_reset': (context) => const ForgotResetPage(),
        '/home': (context) => const HomePage(),
        '/chat': (context) => const ChatPage(),
        '/maps': (context) => const MapsPage(),
        '/profile': (context) => const ProfilePage(),
        '/articles': (context) => ArticlePage(),
        '/quiz': (context) => const QuizPage(),
        '/edit-nickname': (context) => const EditNicknamePage(),
        '/change-password': (context) => const ChangePasswordPage(),
      },
    );
  }
}
