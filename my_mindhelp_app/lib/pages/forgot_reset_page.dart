import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/input_field.dart';
import '../widgets/primary_button.dart';
import '../services/auth_service.dart';

class ForgotResetPage extends StatelessWidget {
  final _newPwdController = TextEditingController();
  final _confirmController = TextEditingController();
  final _authService = AuthService();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(backgroundColor: Colors.transparent, elevation: 0),
      body: Stack(
        children: [
          Positioned(top: -60, right: -60, child: _buildCircle(180)),
          SingleChildScrollView(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 60),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Center(
                      child: Image.asset(
                    'assets/images/mindhelp.png',
                    width: 180,
                  )),
                  const SizedBox(height: 48),
                  Text('請輸入新密碼', style: Theme.of(context).textTheme.bodyText1),
                  const SizedBox(height: 4),
                  InputField(
                    controller: _newPwdController,
                    label: '',
                    prefixIcon: Icons.lock_outline,
                    obscureText: true,
                  ),
                  const SizedBox(height: 24),
                  Text('確認密碼', style: Theme.of(context).textTheme.bodyText1),
                  const SizedBox(height: 4),
                  InputField(
                    controller: _confirmController,
                    label: '',
                    prefixIcon: Icons.lock_outline,
                    obscureText: true,
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
                      const SizedBox(width: 16),
                      Expanded(
                        child: PrimaryButton(
                          text: '確認',
                          onPressed: () async {
                            final newPassword = _newPwdController.text.trim();
                            final confirmPassword =
                                _confirmController.text.trim();

                            if (newPassword.isEmpty ||
                                confirmPassword.isEmpty) {
                              _showSnackBar(context, '所有欄位都是必填的。');
                              return;
                            }
                            if (newPassword != confirmPassword) {
                              _showSnackBar(context, '兩次輸入的新密碼不一致。');
                              return;
                            }
                            if (newPassword.length < 8 ||
                                newPassword.length > 16) {
                              _showSnackBar(context, '密碼長度必須在 8 到 16 個字元之間。');
                              return;
                            }

                            try {
                              // TODO: 呼叫重設密碼 API
                              // await _authService.resetPassword(password: newPassword);
                              _showSnackBar(context, '密碼重設成功！');
                              Navigator.popUntil(
                                  context, ModalRoute.withName('/login'));
                            } catch (e) {
                              _showSnackBar(context, '密碼重設失敗：${e.toString()}');
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

  void _showSnackBar(BuildContext context, String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message)),
    );
  }
}
