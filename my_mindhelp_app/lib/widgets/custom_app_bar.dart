import 'package:flutter/material.dart';
import '../core/theme.dart';

class CustomAppBar extends StatelessWidget implements PreferredSizeWidget {
  final Widget titleWidget;
  final bool showBackButton;
  final IconButton? leftIcon;
  final IconButton? rightIcon;

  const CustomAppBar({
    super.key,
    required this.titleWidget,
    this.showBackButton = true,
    this.leftIcon,
    this.rightIcon,
  });

  @override
  Widget build(BuildContext context) {
    return AppBar(
      title: titleWidget,
      centerTitle: true,
      backgroundColor: Colors.transparent,
      elevation: 0,
      leading: showBackButton
          ? IconButton(
              icon: const Icon(Icons.arrow_back, color: AppColors.textHigh),
              onPressed: () => Navigator.of(context).pop(),
            )
          : leftIcon,
      actions: [
        if (rightIcon != null) rightIcon!,
      ],
    );
  }

  @override
  Size get preferredSize => const Size.fromHeight(kToolbarHeight);
}
