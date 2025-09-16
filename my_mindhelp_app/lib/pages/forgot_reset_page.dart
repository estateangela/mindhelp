import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/input_field.dart';
import '../widgets/primary_button.dart';

class ForgotResetPage extends StatelessWidget {
  final _newPwdController = TextEditingController();
  final _confirmController = TextEditingController();

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
                          onPressed: () {
                            // TODO: 呼叫重設密碼 API
                            Navigator.popUntil(
                                context, ModalRoute.withName('/login'));
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
