import 'package:flutter/material.dart';

class AppColors {
  static const primary = Color(0xFF4A90E2);
  static const background = Color(0xFFF5F5F5);
  static const accent = Color(0xFF50E3C2);
  static const textPrimary = Color(0xFF333333);
}

class AppTextStyles {
  static const headline = TextStyle(
      fontSize: 24, fontWeight: FontWeight.bold, color: AppColors.textPrimary);
  static const body = TextStyle(fontSize: 16, color: AppColors.textPrimary);
  static const button =
      TextStyle(fontSize: 16, fontWeight: FontWeight.w600, color: Colors.white);
}

class AppTheme {
  static final lightTheme = ThemeData(
    primaryColor: AppColors.primary,
    scaffoldBackgroundColor: AppColors.background,
    fontFamily: 'Roboto',
    textTheme: TextTheme(
      headline1: AppTextStyles.headline,
      bodyText1: AppTextStyles.body,
    ),
    elevatedButtonTheme: ElevatedButtonThemeData(
      style: ElevatedButton.styleFrom(
        primary: AppColors.primary,
        textStyle: AppTextStyles.button,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
      ),
    ),
    inputDecorationTheme: InputDecorationTheme(
      border: OutlineInputBorder(borderRadius: BorderRadius.circular(8)),
      filled: true,
      fillColor: Colors.white,
    ),
  );
}
