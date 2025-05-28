import 'package:flutter/material.dart';

class PrimaryButton extends StatelessWidget {
  final String text;
  final VoidCallback onPressed;
  final double? width; // ← 新增

  const PrimaryButton({
    required this.text,
    required this.onPressed,
    this.width, // ← 新增
  });

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: width, // ← 使用
      child: ElevatedButton(
        onPressed: onPressed,
        child: Text(text),
      ),
    );
  }
}
