import 'package:flutter/material.dart';

class AppColors {
  static const primary = Color(0xFFFFBE57);
  static const accent = Color(0xFFE3AD57);
  static const background = Color(0xFFFFF6E7);
  static const textHigh = Color(0xFF636363);
  static const textBody = Color(0xFF9C9C9C);
}

class AppTheme {
  static final lightTheme = ThemeData(
    fontFamily: 'Huninn', // pubspec.yaml 註冊的 font family
    primaryColor: AppColors.primary,
    scaffoldBackgroundColor: AppColors.background,

    textTheme: TextTheme(
      // 新版 TextTheme 屬性
      headlineLarge: TextStyle(
        fontSize: 36,
        color: AppColors.textHigh,
      ),
      headlineMedium: TextStyle(
        fontSize: 24,
        color: AppColors.textHigh,
      ),
      bodyLarge: TextStyle(
        fontSize: 24,
        color: AppColors.textHigh,
      ),
      bodyMedium: TextStyle(
        fontSize: 15,
        color: AppColors.textBody,
      ),
    ),

    elevatedButtonTheme: ElevatedButtonThemeData(
  style: ElevatedButton.styleFrom(
    backgroundColor: Color(0xFFFFD590),   // 按钮背景
    foregroundColor: Colors.white,         // ← 文字／图标  白色
    shadowColor: Colors.black26,
    elevation: 4,
    textStyle: TextStyle(
      fontFamily: 'Huninn',
      fontSize: 24,
      fontWeight: FontWeight.bold,
      // 这里可以不用再写 color，因为 foregroundColor 已经生效
    ),
    padding: const EdgeInsets.symmetric(vertical: 30, horizontal: 10),
    shape: RoundedRectangleBorder(
      borderRadius: BorderRadius.circular(20),
      side: BorderSide(color: AppColors.accent, width: 1),
    ),
  ),
),

        padding: const EdgeInsets.symmetric(vertical: 30, horizontal: 10),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(20),
          side: BorderSide(color: AppColors.accent, width: 1),
        ),
      ),
    ),

    inputDecorationTheme: InputDecorationTheme(
      filled: true,
      fillColor: Colors.white,
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(8),
        borderSide: BorderSide(color: AppColors.accent),
      ),
    ),
  );
}

/// 讓舊版呼叫 .headline1 / .bodyText1 不會報錯
extension CustomTextTheme on TextTheme {
  TextStyle get headline1 => headlineLarge!;
  TextStyle get bodyText1 => bodyMedium!;
}
