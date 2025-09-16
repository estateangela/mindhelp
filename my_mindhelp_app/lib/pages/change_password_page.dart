// lib/pages/change_password_page.dart

import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/input_field.dart';
import '../widgets/primary_button.dart';
import '../services/auth_service.dart';

class ChangePasswordPage extends StatefulWidget {
  const ChangePasswordPage({super.key});

  @override
  State<ChangePasswordPage> createState() => _ChangePasswordPageState();
}

class _ChangePasswordPageState extends State<ChangePasswordPage> {
  final _oldPasswordController = TextEditingController();
  final _newPasswordController = TextEditingController();
  final _confirmNewPasswordController = TextEditingController();
  final _authService = AuthService();
  bool _isLoading = false;

  @override
  void dispose() {
    _oldPasswordController.dispose();
    _newPasswordController.dispose();
    _confirmNewPasswordController.dispose();
    super.dispose();
  }

  void _showSnackBar(String message) {
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(message)),
      );
    }
  }

  Future<void> _changePassword() async {
    final oldPassword = _oldPasswordController.text.trim();
    final newPassword = _newPasswordController.text.trim();
    final confirmNewPassword = _confirmNewPasswordController.text.trim();

    if (oldPassword.isEmpty ||
        newPassword.isEmpty ||
        confirmNewPassword.isEmpty) {
      _showSnackBar('所有欄位都是必填的。');
      return;
    }

    if (newPassword.length < 8 || newPassword.length > 16) {
      _showSnackBar('新密碼長度必須在 8 到 16 個字元之間。');
      return;
    }

    if (newPassword != confirmNewPassword) {
      _showSnackBar('兩次輸入的新密碼不一致。');
      return;
    }

    setState(() {
      _isLoading = true;
    });

    try {
      await _authService.changePassword(
        oldPassword: oldPassword,
        newPassword: newPassword,
      );
      _showSnackBar('密碼修改成功！');
      if (mounted) {
        Navigator.of(context).pop();
      }
    } catch (e) {
      _showSnackBar('修改失敗：${e.toString()}');
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: Text('修改密碼', style: Theme.of(context).textTheme.headlineLarge),
        centerTitle: true,
        backgroundColor: Colors.transparent,
        elevation: 0,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 32),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text('舊密碼', style: Theme.of(context).textTheme.bodyMedium),
            const SizedBox(height: 8),
            InputField(
              controller: _oldPasswordController,
              label: '請輸入舊密碼',
              prefixIcon: Icons.lock_outline,
              obscureText: true,
            ),
            const SizedBox(height: 24),
            Text('新密碼', style: Theme.of(context).textTheme.bodyMedium),
            const SizedBox(height: 8),
            InputField(
              controller: _newPasswordController,
              label: '8-16 碼',
              prefixIcon: Icons.lock_outline,
              obscureText: true,
            ),
            const SizedBox(height: 24),
            Text('確認新密碼', style: Theme.of(context).textTheme.bodyMedium),
            const SizedBox(height: 8),
            InputField(
              controller: _confirmNewPasswordController,
              label: '再次輸入新密碼',
              prefixIcon: Icons.lock_outline,
              obscureText: true,
            ),
            const SizedBox(height: 32),
            PrimaryButton(
              text: _isLoading ? '更新中...' : '確認更新',
              onPressed: _isLoading ? () {} : _changePassword,
            ),
          ],
        ),
      ),
    );
  }
}
