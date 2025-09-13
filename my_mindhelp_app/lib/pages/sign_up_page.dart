// lib/pages/sign_up_page.dart

import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/input_field.dart';
import '../widgets/primary_button.dart';
import '../services/auth_service.dart';

class SignUpPage extends StatelessWidget {
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  final _nicknameController = TextEditingController();

  final _authService = AuthService();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      body: Stack(
        children: [
          Positioned(top: -50, left: -50, child: _buildCircle(160)),
          Positioned(bottom: -40, right: -40, child: _buildCircle(120)),
          SafeArea(
            child: SingleChildScrollView(
              padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 60),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('註冊', style: Theme.of(context).textTheme.headlineLarge),
                  const SizedBox(height: 30),
                  Text('電子郵件', style: Theme.of(context).textTheme.bodyMedium),
                  const SizedBox(height: 4),
                  InputField(
                    controller: _emailController,
                    label: '請輸入有效的電子郵件',
                    prefixIcon: Icons.email_outlined,
                  ),
                  const SizedBox(height: 24),
                  Text('密碼', style: Theme.of(context).textTheme.bodyMedium),
                  const SizedBox(height: 4),
                  InputField(
                    controller: _passwordController,
                    label: '8-16 碼',
                    prefixIcon: Icons.lock_outline,
                    obscureText: true,
                  ),
                  const SizedBox(height: 24),
                  Text('暱稱', style: Theme.of(context).textTheme.bodyMedium),
                  const SizedBox(height: 4),
                  InputField(
                    controller: _nicknameController,
                    label: '請輸入暱稱',
                    prefixIcon: Icons.face_outlined,
                  ),
                  const SizedBox(height: 48),
                  Row(
                    children: [
                      Expanded(
                        child: PrimaryButton(
                          text: '返回',
                          onPressed: () =>
                              Navigator.pushReplacementNamed(context, '/login'),
                        ),
                      ),
                      const SizedBox(width: 24),
                      Expanded(
                        child: PrimaryButton(
                          text: '確認',
                          onPressed: () async {
                            final email = _emailController.text.trim();
                            final password = _passwordController.text.trim();
                            final nickname = _nicknameController.text.trim();

                            // 驗證邏輯
                            if (email.isEmpty ||
                                password.isEmpty ||
                                nickname.isEmpty) {
                              _showSnackBar(context, '所有欄位都是必填的。');
                              return;
                            }
                            if (!_isValidEmail(email)) {
                              _showSnackBar(context, '請輸入有效的電子郵件地址。');
                              return;
                            }
                            if (password.length < 8 || password.length > 16) {
                              _showSnackBar(context, '密碼長度必須在 8 到 16 個字元之間。');
                              return;
                            }

                            try {
                              await _authService.register(
                                email: email,
                                password: password,
                                nickname: nickname,
                              );
                              // 註冊成功，導向登入頁面
                              Navigator.pushReplacementNamed(context, '/login');
                            } catch (e) {
                              _showSnackBar(context, '註冊失敗：${e.toString()}');
                            }
                          },
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCircle(double size) => Container(
        width: size,
        height: size,
        decoration: BoxDecoration(
          color: AppColors.accent.withOpacity(0.3),
          shape: BoxShape.circle,
        ),
      );

  bool _isValidEmail(String email) {
    return RegExp(r'^[^@]+@[^@]+\.[^@]+').hasMatch(email);
  }

  void _showSnackBar(BuildContext context, String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message)),
    );
  }
}
