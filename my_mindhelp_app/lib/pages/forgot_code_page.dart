// lib/pages/forgot_code_page.dart

import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/input_field.dart';
import '../widgets/primary_button.dart';

class ForgotCodePage extends StatefulWidget {
  @override
  _ForgotCodePageState createState() => _ForgotCodePageState();
}

class _ForgotCodePageState extends State<ForgotCodePage> {
  final _emailController = TextEditingController();
  final _codeController = TextEditingController();
  bool _codeSent = false;

  void _sendCode() {
    final email = _emailController.text.trim();
    if (email.isEmpty) return;
    // TODO: 呼叫後端 API 寄送驗證碼
    setState(() => _codeSent = true);
  }

  void _verifyCode() {
    final code = _codeController.text.trim();
    if (code.isEmpty) return;
    // TODO: 呼叫後端 API 驗證 code，若成功則跳到重設密碼頁
    Navigator.pushReplacementNamed(context, '/forgot_reset');
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
        title: Text('忘記密碼', style: Theme.of(context).textTheme.headlineLarge),
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
                  if (!_codeSent) ...[
                    Text(
                      '請輸入你註冊的信箱',
                      style: Theme.of(context).textTheme.bodyMedium,
                    ),
                    const SizedBox(height: 8),
                    InputField(
                      controller: _emailController,
                      label: '',
                      prefixIcon: Icons.email_outlined,
                    ),
                    const SizedBox(height: 24),
                    PrimaryButton(
                      text: '寄送驗證碼',
                      onPressed: _sendCode,
                    ),
                  ] else ...[
                    Text(
                      '請輸入收取的驗證碼',
                      style: Theme.of(context).textTheme.bodyMedium,
                    ),
                    const SizedBox(height: 8),
                    InputField(
                      controller: _codeController,
                      label: '',
                      prefixIcon: Icons.verified_user,
                    ),
                    const SizedBox(height: 8),
                    TextButton(
                      onPressed: _sendCode,
                      child: Text('重新寄送驗證碼'),
                    ),
                    const SizedBox(height: 24),
                    PrimaryButton(
                      text: '確認驗證',
                      onPressed: _verifyCode,
                    ),
                  ],
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}
