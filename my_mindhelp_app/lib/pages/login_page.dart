import 'package:flutter/material.dart';
import '../widgets/input_field.dart';
import '../widgets/primary_button.dart';

class LoginPage extends StatelessWidget {
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  final _nicknameController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('註冊 / 登入')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            InputField(
                controller: _emailController,
                label: '信箱 (Google 為主)',
                prefixIcon: Icons.email),
            SizedBox(height: 12),
            InputField(
                controller: _passwordController,
                label: '密碼',
                prefixIcon: Icons.lock,
                obscureText: true),
            SizedBox(height: 12),
            InputField(
                controller: _nicknameController,
                label: '使用者暱稱',
                prefixIcon: Icons.person),
            SizedBox(height: 24),
            PrimaryButton(
              text: '註冊 / 登入',
              onPressed: () {
                Navigator.pushReplacementNamed(context, '/home');
              },
            ),
            TextButton(
              onPressed: () => Navigator.pushNamed(context, '/forgot'),
              child: Text('忘記密碼？'),
            ),
          ],
        ),
      ),
    );
  }
}
