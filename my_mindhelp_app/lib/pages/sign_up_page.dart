import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/input_field.dart';
import '../widgets/primary_button.dart';

class SignUpPage extends StatelessWidget {
  final _accountController = TextEditingController();
  final _passwordController = TextEditingController();
  final _emailController = TextEditingController();
  final _nicknameController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      body: Stack(
        children: [
          Positioned(top: -50, left: -50, child: _buildCircle(160)),
          Positioned(bottom: -40, right: -40, child: _buildCircle(120)),
          SingleChildScrollView(
            child: Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 32, vertical: 120),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('註冊', style: Theme.of(context).textTheme.headline1),
                  SizedBox(height: 48),
                  Text('帳號', style: Theme.of(context).textTheme.bodyText1),
                  SizedBox(height: 4),
                  InputField(
                      controller: _accountController,
                      label: '',
                      prefixIcon: Icons.person_outline),
                  SizedBox(height: 24),
                  Text('密碼', style: Theme.of(context).textTheme.bodyText1),
                  SizedBox(height: 4),
                  InputField(
                    controller: _passwordController,
                    label: '',
                    prefixIcon: Icons.lock_outline,
                    obscureText: true,
                  ),
                  SizedBox(height: 24),
                  Text('電子郵件', style: Theme.of(context).textTheme.bodyText1),
                  SizedBox(height: 4),
                  InputField(
                      controller: _emailController,
                      label: '',
                      prefixIcon: Icons.email_outlined),
                  SizedBox(height: 24),
                  Text('暱稱', style: Theme.of(context).textTheme.bodyText1),
                  SizedBox(height: 4),
                  InputField(
                      controller: _nicknameController,
                      label: '',
                      prefixIcon: Icons.face_outlined),
                  SizedBox(height: 48),
                  Row(
                    children: [
                      Expanded(
                        child: PrimaryButton(
                          text: '返回',
                          onPressed: () => Navigator.pop(context),
                        ),
                      ),
                      SizedBox(width: 16),
                      Expanded(
                        child: PrimaryButton(
                          text: '確認',
                          onPressed: () {
                            // TODO: 呼叫註冊 API
                            Navigator.pop(context);
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
}
