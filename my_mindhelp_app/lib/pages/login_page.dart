import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/input_field.dart';
import '../widgets/primary_button.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key}); // 注意這裡加上 const

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  // 將 TextEditingController 移到 State 類別中
  late final TextEditingController _accountController;
  late final TextEditingController _passwordController;

  @override
  void initState() {
    super.initState();
    _accountController = TextEditingController();
    _passwordController = TextEditingController();
  }

  @override
  void dispose() {
    _accountController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      resizeToAvoidBottomInset: true,
      body: Stack(
        children: [
          Positioned(top: -50, left: -60, child: _buildCircle(180)),
          Positioned(bottom: -50, left: -50, child: _buildCircle(140)),
          Positioned(top: 100, right: -30, child: _buildCircle(100)),
          Center(
            child: SingleChildScrollView(
              padding: const EdgeInsets.symmetric(horizontal: 30, vertical: 30),
              child: ConstrainedBox(
                constraints: const BoxConstraints(maxWidth: 400),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    Image.asset(
                      'assets/images/logo.png',
                      width: 180,
                    ),
                    const SizedBox(height: 24),
                    Text('帳號', style: Theme.of(context).textTheme.bodyMedium),
                    const SizedBox(height: 8),
                    InputField(
                      controller: _accountController,
                      label: '',
                      prefixIcon: Icons.person_outline,
                    ),
                    const SizedBox(height: 20),
                    Text('密碼', style: Theme.of(context).textTheme.bodyMedium),
                    const SizedBox(height: 8),
                    InputField(
                      controller: _passwordController,
                      label: '',
                      prefixIcon: Icons.lock_outline,
                      obscureText: true,
                    ),
                    const SizedBox(height: 32),
                    Row(
                      children: [
                        Expanded(
                          child: PrimaryButton(
                            text: '登入',
                            onPressed: () => Navigator.pushReplacementNamed(
                                context, '/home'),
                          ),
                        ),
                        const SizedBox(width: 16),
                        Expanded(
                          child: PrimaryButton(
                            text: '註冊',
                            onPressed: () =>
                                Navigator.pushNamed(context, '/signup'),
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 24),
                    GestureDetector(
                      onTap: () =>
                          Navigator.pushNamed(context, '/forgot-reset'),
                      child: Text(
                        '忘記密碼？',
                        textAlign: TextAlign.center,
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
