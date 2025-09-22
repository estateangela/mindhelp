// lib/pages/forgot_code_page.dart

import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/input_field.dart';
import '../widgets/primary_button.dart';
import '../services/auth_service.dart';

class ForgotCodePage extends StatefulWidget {
  const ForgotCodePage({super.key});

  @override
  _ForgotCodePageState createState() => _ForgotCodePageState();
}

class _ForgotCodePageState extends State<ForgotCodePage> {
  final _emailController = TextEditingController();
  final _newPwdController = TextEditingController();
  final _confirmController = TextEditingController();
  final _authService = AuthService();
  bool _isLoading = false;

  @override
  void dispose() {
    _emailController.dispose();
    _newPwdController.dispose();
    _confirmController.dispose();
    super.dispose();
  }

  void _showSnackBar(String message) {
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(message)),
      );
    }
  }

  Future<void> _resetPassword() async {
    final email = _emailController.text.trim();
    final newPassword = _newPwdController.text.trim();
    final confirmPassword = _confirmController.text.trim();

    if (email.isEmpty || newPassword.isEmpty || confirmPassword.isEmpty) {
      _showSnackBar('所有欄位都是必填的。');
      return;
    }
    if (newPassword != confirmPassword) {
      _showSnackBar('兩次輸入的新密碼不一致。');
      return;
    }
    if (newPassword.length < 8 || newPassword.length > 16) {
      _showSnackBar('密碼長度必須在 8 到 16 個字元之間。');
      return;
    }

    setState(() {
      _isLoading = true;
    });

    try {
      // TODO: 呼叫重設密碼 API
      // await _authService.resetPassword(email: email, newPassword: newPassword);
      _showSnackBar('密碼重設成功！');
      if (mounted) {
        Navigator.popUntil(context, ModalRoute.withName('/login'));
      }
    } catch (e) {
      _showSnackBar('密碼重設失敗：${e.toString()}');
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  Widget _buildCircle(double size, Alignment alignment) {
    return Align(
      alignment: alignment,
      child: Container(
        width: size,
        height: size,
        decoration: BoxDecoration(
          color: AppColors.accent.withOpacity(0.3),
          shape: BoxShape.circle,
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        title: Text('重設密碼', style: Theme.of(context).textTheme.headlineLarge),
        centerTitle: true,
      ),
      body: Stack(
        children: [
          _buildCircle(140, Alignment.bottomLeft),
          _buildCircle(100, Alignment.topRight),
          SingleChildScrollView(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 24),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  Text(
                    '請輸入你的電子郵件',
                    style: Theme.of(context).textTheme.bodyMedium,
                  ),
                  const SizedBox(height: 8),
                  InputField(
                    controller: _emailController,
                    label: '',
                    prefixIcon: Icons.email_outlined,
                  ),
                  const SizedBox(height: 24),
                  Text(
                    '請輸入新密碼',
                    style: Theme.of(context).textTheme.bodyMedium,
                  ),
                  const SizedBox(height: 8),
                  InputField(
                    controller: _newPwdController,
                    label: '',
                    prefixIcon: Icons.lock_outline,
                    obscureText: true,
                  ),
                  const SizedBox(height: 24),
                  Text(
                    '確認密碼',
                    style: Theme.of(context).textTheme.bodyMedium,
                  ),
                  const SizedBox(height: 8),
                  InputField(
                    controller: _confirmController,
                    label: '',
                    prefixIcon: Icons.lock_outline,
                    obscureText: true,
                  ),
                  const SizedBox(height: 24),
                  PrimaryButton(
                    text: _isLoading ? '重設中...' : '確認',
                    onPressed: _isLoading ? () {} : _resetPassword,
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}
