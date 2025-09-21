import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/input_field.dart';
import '../widgets/primary_button.dart';
import '../services/auth_service.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  late final TextEditingController _accountController;
  late final TextEditingController _passwordController;
  bool _isInputEmpty = true;
  bool _showValidation = false; // 新增變數來控制提示顯示

  @override
  void initState() {
    super.initState();
    _accountController = TextEditingController();
    _passwordController = TextEditingController();

    _accountController.addListener(_updateLoginButtonState);
    _passwordController.addListener(_updateLoginButtonState);
  }

  @override
  void dispose() {
    _accountController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  void _updateLoginButtonState() {
    setState(() {
      _isInputEmpty =
          _accountController.text.isEmpty || _passwordController.text.isEmpty;
    });
  }

  void _showSnackBar(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message)),
    );
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
                      label: _showValidation && _accountController.text.isEmpty
                          ? '必填'
                          : '',
                      prefixIcon: Icons.person_outline,
                    ),
                    const SizedBox(height: 20),
                    Text('密碼', style: Theme.of(context).textTheme.bodyMedium),
                    const SizedBox(height: 8),
                    InputField(
                      controller: _passwordController,
                      label: _showValidation && _passwordController.text.isEmpty
                          ? '必填'
                          : '',
                      prefixIcon: Icons.lock_outline,
                      obscureText: true,
                    ),
                    const SizedBox(height: 32),
                    Row(
                      children: [
                        Expanded(
                          child: PrimaryButton(
                            text: '登入',
                            onPressed: () {
                              setState(() {
                                _showValidation = true;
                              });
                              if (!_isInputEmpty) {
                                _showSnackBar('登入成功!');
                                Navigator.pushReplacementNamed(
                                    context, '/home');
                              }
                            },
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
                      onTap: () => Navigator.pushNamed(context, '/forgot_code'),
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
