import 'package:flutter/material.dart';
import '../widgets/input_field.dart';
import '../widgets/primary_button.dart';

class ForgotPasswordPage extends StatefulWidget {
  @override
  _ForgotPasswordPageState createState() => _ForgotPasswordPageState();
}

class _ForgotPasswordPageState extends State<ForgotPasswordPage> {
  final _emailController = TextEditingController();
  final _codeController = TextEditingController();
  final _newPwdController = TextEditingController();
  final _confirmController = TextEditingController();
  bool _codeSent = false;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('重設密碼')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            if (!_codeSent) ...[
              InputField(
                  controller: _emailController,
                  label: '信箱',
                  prefixIcon: Icons.email),
              SizedBox(height: 24),
              PrimaryButton(
                  text: '寄送驗證碼',
                  onPressed: () {
                    // TODO: 呼叫後端寄送驗證碼
                    setState(() => _codeSent = true);
                  }),
            ] else ...[
              InputField(
                  controller: _codeController,
                  label: '驗證碼',
                  prefixIcon: Icons.verified_user),
              SizedBox(height: 12),
              InputField(
                  controller: _newPwdController,
                  label: '新密碼',
                  prefixIcon: Icons.lock,
                  obscureText: true),
              SizedBox(height: 12),
              InputField(
                  controller: _confirmController,
                  label: '確認密碼',
                  prefixIcon: Icons.lock,
                  obscureText: true),
              SizedBox(height: 24),
              PrimaryButton(
                  text: '確認重設',
                  onPressed: () {
                    // TODO: 呼叫後端驗證並更新密碼
                    Navigator.pop(context);
                  }),
            ]
          ],
        ),
      ),
    );
  }
}
