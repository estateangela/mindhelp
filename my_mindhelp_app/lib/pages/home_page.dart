import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/primary_button.dart';
import '../widgets/input_field.dart';

class HomePage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      body: SafeArea(
        child: Column(
          children: [
            SizedBox(height: 16),
            Text('mindhelp', style: Theme.of(context).textTheme.headline1),
            SizedBox(height: 24),

            // 四格按鈕
            Expanded(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16),
                child: GridView.count(
                  crossAxisCount: 2,
                  mainAxisSpacing: 16,
                  crossAxisSpacing: 16,
                  children: [
                    _buildTile(context, Icons.map, '尋找附近\n醫療資源', () {}),
                    _buildTile(context, Icons.menu_book, '心理師專欄', () {}),
                    _buildTile(context, Icons.chat_bubble_outline, 'AI諮詢', () {
                      Navigator.pushNamed(context, '/chat');
                    }),
                    _buildTile(context, Icons.favorite_border, '心理測驗', () {}),
                  ],
                ),
              ),
            ),

            // 下面兩個輸入欄 (依 Figma 顯示)
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 32),
              child: Column(
                children: [
                  InputField(
                      controller: TextEditingController(),
                      label: '',
                      prefixIcon: null),
                  SizedBox(height: 16),
                  InputField(
                      controller: TextEditingController(),
                      label: '',
                      prefixIcon: null),
                ],
              ),
            ),
            SizedBox(height: 24),
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
              SizedBox(height: 8),
              Text(
                label,
                textAlign: TextAlign.center,
                style:
                    TextStyle(color: Colors.white, fontWeight: FontWeight.bold),
              ),
            ],
          ),
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
        // TODO: handle navigation
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
