import 'package:flutter/material.dart';
import 'core/theme.dart';
import 'pages/login_page.dart';
import 'pages/sign_up_page.dart';
import 'pages/forgot_code_page.dart';
import 'pages/forgot_reset_page.dart';
import 'pages/home_page.dart';
import 'pages/chat_page.dart';
//import 'pages/counselor_list_page.dart';
//import 'pages/counselor_detail_page.dart';
//import 'pages/profile_page.dart';

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
        '/signup': (_) => SignUpPage(),
        '/forgot_code': (_) => ForgotCodePage(),
        '/forgot_reset': (_) => ForgotResetPage(),
        '/home': (_) => HomePage(),
        '/chat': (_) => ChatPage(),
        // '/counselors': (_) => CounselorListPage(),
        // '/detail': (_) => CounselorDetailPage(),
        //'/profile': (_) => ProfilePage(),
      },
    );
  }
}
