// lib/pages/settings_page.dart

import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/custom_app_bar.dart';

class SettingsPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,

      // 自定义 AppBar：左返回／右通知
      appBar: const CustomAppBar(
        showBackButton: true,
        titleWidget: Text(
          '設定',
          style: TextStyle(fontSize: 24, color: AppColors.textHigh),
        ),
      ),

      body: SingleChildScrollView(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 24),
        child: Column(
          children: [
            _buildOptionCard(
              context,
              label: '推播通知設定',
              onTap: () => Navigator.pushNamed(context, '/notify'),
            ),
            const SizedBox(height: 16),
            _buildOptionCard(
              context,
              label: '常見問題',
              onTap: () => Navigator.pushNamed(context, '/faq'),
            ),
            const SizedBox(height: 16),
            _buildOptionCard(
              context,
              label: '顯示設定',
              onTap: () => Navigator.pushNamed(context, '/display_settings'),
            ),
            const SizedBox(height: 16),
            _buildOptionCard(
              context,
              label: '關於我們',
              onTap: () => Navigator.pushNamed(context, '/about'),
            ),
            const SizedBox(height: 16),

            // 新增「登出」按鈕
            _buildOptionCard(
              context,
              label: '登出',
              onTap: () {
                // TODO: 在這裡執行登出邏輯（如清除 Token），然後返回登入頁
                Navigator.pushReplacementNamed(context, '/login');
              },
            ),
          ],
        ),
      ),

      bottomNavigationBar: BottomNavigationBar(
        currentIndex: 3, // Profile 的位置
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

  Widget _buildOptionCard(BuildContext ctx,
      {required String label, required VoidCallback onTap}) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        height: 50,
        width: double.infinity,
        alignment: Alignment.center,
        decoration: BoxDecoration(
          color: Colors.white,
          border: Border.all(color: AppColors.accent),
          borderRadius: BorderRadius.circular(8),
        ),
        child: Text(
          label,
          style: Theme.of(ctx).textTheme.bodyMedium,
        ),
      ),
    );
  }
}
