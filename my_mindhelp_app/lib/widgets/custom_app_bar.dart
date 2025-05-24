// lib/widgets/custom_app_bar.dart

import 'package:flutter/material.dart';
import '../core/theme.dart';

/// 一个共用的 AppBar：
/// - 左边：设置/返回按钮
/// - 中间：可传入文字或图片
/// - 右边：通知按钮
class CustomAppBar extends StatelessWidget implements PreferredSizeWidget {
  /// 如果是子页面需要一个返回箭头，传 true
  final bool showBackButton;

  /// 中间可以传入任意 Widget（Text, Image.asset, ...）
  final Widget titleWidget;

  const CustomAppBar({
    Key? key,
    this.showBackButton = false,
    required this.titleWidget,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return AppBar(
      backgroundColor: Colors.transparent,
      elevation: 0,
      leading: showBackButton
          ? BackButton(color: AppColors.textHigh)
          : IconButton(
              icon: Icon(Icons.settings, color: AppColors.textHigh),
              onPressed: () => Navigator.pushNamed(context, '/settings'),
            ),
      title: titleWidget,
      centerTitle: true,
      actions: [
        IconButton(
          icon: Icon(Icons.notifications, color: AppColors.textHigh),
          onPressed: () => Navigator.pushNamed(context, '/notify'),
        ),
      ],
    );
  }

  @override
  Size get preferredSize => Size.fromHeight(kToolbarHeight);
}
