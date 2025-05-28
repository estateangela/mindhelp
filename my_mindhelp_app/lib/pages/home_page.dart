// lib/pages/home_page.dart

import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/custom_app_bar.dart';
import '../widgets/primary_button.dart';

class HomePage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,

      // 顶部改成共用的带 logo 的 AppBar
      appBar: const CustomAppBar(
        titleWidget: Image(
          image: AssetImage('assets/images/mindhelp.png'),
          width: 200,
          fit: BoxFit.contain,
        ),
      ),

      body: SafeArea(
        child: Column(
          children: [
            const SizedBox(height: 24),

            // 四格按鈕
            Expanded(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 10),
                child: GridView.count(
                  crossAxisCount: 2,
                  mainAxisSpacing: 10,
                  crossAxisSpacing: 10,
                  children: [
                    _buildTile(context, Icons.map, '尋找附近\n醫療資源', () {
                      Navigator.pushNamed(context, '/maps');
                    }),
                    _buildTile(context, Icons.menu_book, '心理師專欄', () {
                      Navigator.pushNamed(context, '/counselors');
                    }),
                    _buildTile(context, Icons.chat_bubble_outline, 'AI諮詢', () {
                      Navigator.pushNamed(context, '/chat');
                    }),
                    _buildTile(context, Icons.favorite_border, '心理測驗', () {
                      Navigator.pushNamed(context, '/quiz');
                    }),
                  ],
                ),
              ),
            ),

            // 底下两行功能卡片
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 2),
              child: Column(
                children: [
                  _buildFunctionCard(
                    context,
                    icon: Icons.info_outline,
                    label: '最新心理健康文章',
                    onTap: () => Navigator.pushNamed(context, '/articles'),
                  ),
                  const SizedBox(height: 8),
                  _buildFunctionCard(
                    context,
                    icon: Icons.event_note_outlined,
                    label: '預約紀錄查詢',
                    onTap: () => Navigator.pushNamed(context, '/appointments'),
                  ),
                ],
              ),
            ),

            const SizedBox(height: 20),
          ],
        ),
      ),

      bottomNavigationBar: _buildBottomNav(context, 0),
    );
  }

  Widget _buildTile(
      BuildContext ctx, IconData icon, String label, VoidCallback onTap) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        decoration: BoxDecoration(
          color: AppColors.accent,
          borderRadius: BorderRadius.circular(12),
        ),
        child: Center(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(icon, size: 36, color: Colors.white),
              const SizedBox(height: 8),
              Text(
                label,
                textAlign: TextAlign.center,
                style: const TextStyle(
                    color: Colors.white, fontWeight: FontWeight.bold),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildFunctionCard(BuildContext ctx,
      {required IconData icon,
      required String label,
      required VoidCallback onTap}) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        height: 50,
        decoration: BoxDecoration(
          color: Colors.white,
          border: Border.all(color: AppColors.accent),
          borderRadius: BorderRadius.circular(8),
        ),
        padding: const EdgeInsets.symmetric(horizontal: 16),
        child: Row(
          children: [
            Icon(icon, color: AppColors.accent),
            const SizedBox(width: 12),
            Text(label, style: Theme.of(ctx).textTheme.bodyMedium),
          ],
        ),
      ),
    );
  }

  Widget _buildBottomNav(BuildContext ctx, int idx) {
    return BottomNavigationBar(
      currentIndex: idx,
      selectedItemColor: AppColors.accent,
      unselectedItemColor: AppColors.textBody,
      onTap: (i) {
        switch (i) {
          case 0:
            Navigator.pushReplacementNamed(ctx, '/home');
            break;
          case 1:
            Navigator.pushReplacementNamed(ctx, '/maps');
            break;
          case 2:
            Navigator.pushReplacementNamed(ctx, '/chat');
            break;
          case 3:
            Navigator.pushReplacementNamed(ctx, '/profile');
            break;
        }
      },
      items: const [
        BottomNavigationBarItem(icon: Icon(Icons.home), label: 'Home'),
        BottomNavigationBarItem(icon: Icon(Icons.location_on), label: 'Maps'),
        BottomNavigationBarItem(icon: Icon(Icons.chat_bubble), label: 'Chat'),
        BottomNavigationBarItem(icon: Icon(Icons.person), label: 'Profile'),
      ],
    );
  }
}
