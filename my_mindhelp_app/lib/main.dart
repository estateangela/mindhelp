// lib/main.dart

import 'package:flutter/material.dart';
import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_auth/firebase_auth.dart';

import 'core/theme.dart';
import 'pages/login_page.dart';
import 'pages/sign_up_page.dart';
import 'pages/forgot_code_page.dart';
import 'pages/forgot_reset_page.dart';
import 'pages/home_page.dart';
import 'pages/chat_page.dart';
import 'pages/mapspage.dart';
import 'pages/profile_page.dart';
import 'pages/notify_page.dart';
import 'pages/settings_page.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await Firebase.initializeApp();
  runApp(const MindHelpApp());
}

class MindHelpApp extends StatelessWidget {
  const MindHelpApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'MindHelp',
      theme: AppTheme.lightTheme,
      initialRoute: '/login',
      routes: {
        '/notify': (_) => const NotifyPage(),
        '/settings': (_) => const SettingsPage(),
        '/login': (_) => const LoginPage(),
        '/signup': (_) => const SignUpPage(),
        '/forgot_code': (_) => const ForgotCodePage(),
        '/forgot_reset': (_) => const ForgotResetPage(),
        '/home': (_) => const HomePage(),
        '/chat': (_) => const ChatPage(),
        '/maps': (_) => const MapsPage(),
        '/profile': (_) => const ProfilePage(),
      },
    );
  }
}
