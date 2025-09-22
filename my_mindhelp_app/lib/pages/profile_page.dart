// lib/pages/profile_page.dart

import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/custom_app_bar.dart';
import '../widgets/primary_button.dart';

class ProfilePage extends StatelessWidget {
  const ProfilePage({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,

      // 使用自定義的 AppBar，中央顯示「我的資料」，右側顯示通知鈴鐺
      appBar: CustomAppBar(
        showBackButton: false,
        titleWidget: const Text(
          '我的資料',
          style: TextStyle(fontSize: 24, color: AppColors.textHigh),
        ),
        rightIcon: IconButton(
          icon: const Icon(Icons.notifications, color: AppColors.textHigh),
          // TODO: 已經完成 - 導向通知頁
          onPressed: () {
            Navigator.pushNamed(context, '/notify');
          },
        ),
      ),

      // 讓整個頁面在垂直空間不足時可以滾動
      body: SingleChildScrollView(
        // 加一點底部 padding，避免內容貼著 BottomNavigationBar
        padding: const EdgeInsets.only(bottom: 16),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 24),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              // 1. 用戶資訊卡（Email + 暱稱 + 圓形頭像）
              Container(
                padding: const EdgeInsets.all(16),
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Row(
                  children: [
                    // 圓形頭像
                    CircleAvatar(
                      radius: 28,
                      backgroundColor: AppColors.accent,
                      child: const Icon(
                        Icons.person,
                        color: Colors.white,
                        size: 32,
                      ),
                    ),

                    const SizedBox(width: 16),

                    // 這裡把文字區包到 Expanded 裡，讓 email 過長時自動換行
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          // Email
                          Text(
                            'a123456789@gmail.com',
                            style: Theme.of(context).textTheme.bodyLarge,
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                          ),
                          const SizedBox(height: 4),
                          // 暱稱
                          Text(
                            '王小美',
                            style: Theme.of(context).textTheme.bodyMedium,
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),

              const SizedBox(height: 32),

              // 2. 各項操作按鈕：修改暱稱、修改密碼、查詢預約紀錄、登出
              PrimaryButton(
                text: '修改暱稱',
                onPressed: () {
                  // TODO: 已經完成 - 導向「修改暱稱」頁面
                  Navigator.pushNamed(context, '/edit-nickname');
                },
              ),
              const SizedBox(height: 16),

              PrimaryButton(
                text: '修改密碼',
                onPressed: () {
                  // TODO: 已經完成 - 導向「修改密碼」頁面
                  Navigator.pushNamed(context, '/change-password');
                },
              ),
              const SizedBox(height: 16),
              PrimaryButton(
                text: '登出',
                onPressed: () {
                  // TODO: 已經完成 - 實作登出邏輯
                  // 假設登入頁面路由為 '/' 或 '/login'
                  Navigator.pushNamedAndRemoveUntil(
                    context,
                    '/login',
                    (Route<dynamic> route) => false,
                  );
                },
              ),
            ],
          ),
        ),
      ),

      // 3. 底部導覽列
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: 3, // Profile 頁的索引
        selectedItemColor: AppColors.accent,
        unselectedItemColor: AppColors.textBody,
        onTap: (idx) {
          switch (idx) {
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
              // 已經在 Profile，不做動作
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
