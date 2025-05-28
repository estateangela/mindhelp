import 'package:flutter/material.dart';

class AppColors {
  static const primary = Color(0xFFFFBE57);
  static const accent = Color(0xFFE3AD57);
  static const background = Color(0xFFFFF6E7);
  static const textHigh = Color(0xFF636363);
  static const textBody = Color(0xFF9C9C9C);
  static const error = Color(0xFFE53935);
}

class AppTheme {
  static final lightTheme = ThemeData(
    fontFamily: 'Huninn', // pubspec.yaml 注册的 font family
    primaryColor: AppColors.primary,
    scaffoldBackgroundColor: AppColors.background,

    textTheme: TextTheme(
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
        backgroundColor: Color(0xFFFFD590), // 按钮背景
        foregroundColor: Colors.white, // 文字／图标 白色
        shadowColor: Colors.black26,
        elevation: 4,
        textStyle: TextStyle(
          fontFamily: 'Huninn',
          fontSize: 24,
          fontWeight: FontWeight.bold,
        ),
        padding: const EdgeInsets.symmetric(vertical: 20, horizontal: 10),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(20),
          side: BorderSide(color: AppColors.accent, width: 1),
        ),
      ),
    ),

    inputDecorationTheme: InputDecorationTheme(
      filled: true,
      fillColor: const Color.fromARGB(255, 227, 227, 227),
      enabledBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(10),
        borderSide: BorderSide(color: AppColors.accent, width: 1),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(10),
        borderSide: BorderSide(color: AppColors.primary, width: 2),
      ),
    ),
    textSelectionTheme: TextSelectionThemeData(
      cursorColor: AppColors.primary,
      selectionColor: AppColors.accent.withOpacity(0.3),
      selectionHandleColor: AppColors.primary,
    ),
  );
}

/// 让旧版调用 .headline1 / .bodyText1 不会报错
extension CustomTextTheme on TextTheme {
  TextStyle get headline1 => headlineLarge!;
  TextStyle get bodyText1 => bodyLarge!;
}
