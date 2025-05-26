// lib/pages/profile_page.dart

import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/custom_app_bar.dart';
import '../widgets/primary_button.dart';

class ProfilePage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: const CustomAppBar(
        titleWidget: Text(
          '我的資料',
          style: TextStyle(fontSize: 24, color: AppColors.textHigh),
        ),
      ),
      body: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
        child: Column(
          children: [
            // 使用者資訊卡片
            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(12),
              ),
              child: Row(
                children: [
                  CircleAvatar(radius: 28, backgroundColor: AppColors.accent),
                  const SizedBox(width: 16),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text('a123456789@gmail.com',
                          style: Theme.of(context).textTheme.bodyMedium),
                      const SizedBox(height: 4),
                      Text('王小美',
                          style: Theme.of(context).textTheme.bodyMedium),
                    ],
                  ),
                ],
              ),
            ),

            const SizedBox(height: 32),
            PrimaryButton(text: '修改信箱', onPressed: () {}),
            const SizedBox(height: 16),
            PrimaryButton(text: '修改暱稱', onPressed: () {}),
            const SizedBox(height: 16),
            PrimaryButton(text: '修改密碼', onPressed: () {}),
            const SizedBox(height: 16),
            PrimaryButton(text: '查詢預約紀錄', onPressed: () {}),

            const Spacer(),
          ],
        ),
      ),
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: 3, // Profile 在 index=3
        selectedItemColor: AppColors.accent,
        unselectedItemColor: AppColors.textBody,
        onTap: (i) {
          switch (i) {
            case 0:
              Navigator.pushReplacementNamed(context, '/home');
              break;
            case 1:
              Navigator.pushReplacementNamed(context, '/maps');
              break;
            case 2:
              Navigator.pushReplacementNamed(context, '/chat');
              break;
            case 3:
              // already here
              break;
          }
        },
        items: const [
          BottomNavigationBarItem(icon: Icon(Icons.home), label: 'Home'),
          BottomNavigationBarItem(icon: Icon(Icons.location_on), label: 'Maps'),
          BottomNavigationBarItem(icon: Icon(Icons.chat_bubble), label: 'Chat'),
          BottomNavigationBarItem(icon: Icon(Icons.person), label: 'Profile'),
        ],
      ),
    );
  }
}
