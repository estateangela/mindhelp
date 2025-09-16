// lib/pages/edit_nickname_page.dart

import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/input_field.dart';
import '../widgets/primary_button.dart';
import '../services/auth_service.dart';

class EditNicknamePage extends StatefulWidget {
  const EditNicknamePage({super.key});

  @override
  State<EditNicknamePage> createState() => _EditNicknamePageState();
}

class _EditNicknamePageState extends State<EditNicknamePage> {
  final _nicknameController = TextEditingController();
  final _authService = AuthService();
  bool _isLoading = false;

  @override
  void dispose() {
    _nicknameController.dispose();
    super.dispose();
  }

  void _showSnackBar(String message) {
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(message)),
      );
    }
  }

  Future<void> _updateNickname() async {
    final newNickname = _nicknameController.text.trim();
    if (newNickname.isEmpty) {
      _showSnackBar('暱稱不能為空。');
      return;
    }

    setState(() {
      _isLoading = true;
    });

    try {
      await _authService.updateNickname(nickname: newNickname);
      _showSnackBar('暱稱更新成功！');
      // 檢查 widget 是否仍然存在於 widget 樹中
      if (mounted) {
        Navigator.of(context).pop();
      }
    } catch (e) {
      _showSnackBar('更新失敗：${e.toString()}');
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
        title: Text('修改暱稱', style: Theme.of(context).textTheme.headlineLarge),
        centerTitle: true,
        backgroundColor: Colors.transparent,
        elevation: 0,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 32),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text('新暱稱', style: Theme.of(context).textTheme.bodyMedium),
            const SizedBox(height: 8),
            InputField(
              controller: _nicknameController,
              label: '請輸入新的暱稱',
              prefixIcon: Icons.face_outlined,
            ),
            const SizedBox(height: 32),
            PrimaryButton(
              text: _isLoading ? '更新中...' : '確認更新',
              onPressed: _isLoading ? () {} : _updateNickname,
            ),
          ],
        ),
      ),
    );
  }
}
