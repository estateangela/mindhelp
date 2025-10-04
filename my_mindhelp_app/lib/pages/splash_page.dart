import 'package:flutter/material.dart';
import 'dart:async';
import '../core/theme.dart';
import '../services/location_service.dart'; // 導入地圖服務

class SplashPage extends StatefulWidget {
  const SplashPage({super.key});

  @override
  State<SplashPage> createState() => _SplashPageState();
}

class _SplashPageState extends State<SplashPage> {
  @override
  void initState() {
    super.initState();
    _initializeApp();
  }

  Future<void> _initializeApp() async {
    // 模擬載入時間，確保使用者能看到 Logo
    await Future.delayed(const Duration(seconds: 3));

    // 移除載入 API 的邏輯，直接跳轉
    if (mounted) {
      Navigator.pushReplacementNamed(context, '/home');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Image.asset(
              'assets/images/logo.png', // 確保你的 Logo 檔案路徑正確
              width: 200,
              fit: BoxFit.contain,
            ),
            const SizedBox(height: 20),
            const Text(
              '用心，陪伴你我',
              style: TextStyle(
                fontSize: 20,
                color: AppColors.textHigh,
              ),
            ),
            const SizedBox(height: 40),
            CircularProgressIndicator(
              valueColor: AlwaysStoppedAnimation<Color>(AppColors.accent),
            ),
          ],
        ),
      ),
    );
  }
}
