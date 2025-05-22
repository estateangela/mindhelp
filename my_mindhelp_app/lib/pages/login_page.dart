import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/input_field.dart';
import '../widgets/primary_button.dart';

class LoginPage extends StatelessWidget {
  final _accountController = TextEditingController();
  final _passwordController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      body: Stack(
        children: [
          // 背景圓點
          Positioned(
            top: -50,
            left: -60,
            child: _buildCircle(180),
          ),
          Positioned(
            bottom: -50,
            left: -50,
            child: _buildCircle(140),
          ),
          Positioned(
            top: 100,
            right: -30,
            child: _buildCircle(100),
          ),

          SingleChildScrollView(
            child: Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 30, vertical: 100),
              child: Column(
                children: [
                  Image.asset('assets/images/logo.png', width: 180),
                  SizedBox(height: 10),
                  Align(
                    alignment: Alignment.centerLeft,
                    child: Text('帳號',
                        style: Theme.of(context).textTheme.bodyText1),
                  ),
                  SizedBox(height: 3),
                  InputField(
                      controller: _accountController,
                      label: '',
                      prefixIcon: Icons.person_outline),
                  SizedBox(height: 10),
                  Align(
                    alignment: Alignment.centerLeft,
                    child: Text('密碼',
                        style: Theme.of(context).textTheme.bodyText1),
                  ),
                  SizedBox(height: 3),
                  InputField(
                    controller: _passwordController,
                    label: '',
                    prefixIcon: Icons.lock_outline,
                    obscureText: true,
                  ),
                  SizedBox(height: 20),
                  Row(
                    children: [
                      Expanded(
                        child: PrimaryButton(
                          text: '登入',
                          onPressed: () {
                            Navigator.pushReplacementNamed(context, '/home');
                          },
                        ),
                      ),
                      SizedBox(width: 20),
                      Expanded(
                        child: PrimaryButton(
                          text: '註冊',
                          onPressed: () {
                            Navigator.pushNamed(context, '/signup');
                          },
                        ),
                      ),
                    ],
                  ),
                  SizedBox(height: 20),
                  GestureDetector(
                    onTap: () => Navigator.pushNamed(context, '/forgot_code'),
                    child: Text(
                      '忘記密碼？',
                      style: TextStyle(
                        color: AppColors.textBody,
                        decoration: TextDecoration.underline,
                      ),
                    ),
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
