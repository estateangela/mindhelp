// lib/widgets/custom_app_bar.dart

import 'package:flutter/material.dart';
import '../core/theme.dart';

/// 自定义的 AppBar，可以选择显示「返回按钮」和「右侧图标」
///
/// 用法示例：
///   CustomAppBar(
///     showBackButton: true,            // 是否展示左侧返回箭头
///     titleWidget: Text('標題內容'),
///     rightIcon: IconButton(
///       icon: Icon(Icons.notifications, color: AppColors.textHigh),
///       onPressed: () {
///         // 點擊右側通知按鈕要執行的動作
///       },
///     ),
///   );
///
class CustomAppBar extends StatelessWidget implements PreferredSizeWidget {
  /// 如果 true，就在最左邊顯示一個返回按鈕 (使用 Navigator.pop)
  final bool showBackButton;

  /// 放在中間位置的 Widget (通常是 Text、Image 等)
  final Widget titleWidget;

  /// 可選的右側按鈕，一般用 IconButton，若為 null 則不顯示
  final Widget? rightIcon;

  const CustomAppBar({
    Key? key,
    this.showBackButton = false,
    required this.titleWidget,
    this.rightIcon,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Container(
        color: Colors.transparent,
        padding: const EdgeInsets.only(left: 0, right: 0, top: 8, bottom: 8),
        child: Row(
          children: [
            // 左側：如果需要返回按鈕，就顯示 IconButton，否則保留一個占位 SizedBox
            if (showBackButton)
              IconButton(
                icon: Icon(Icons.arrow_back_ios, color: AppColors.textHigh),
                onPressed: () {
                  Navigator.maybePop(context);
                },
              )
            else
              const SizedBox(width: 48), // 為了讓標題置中，占位寬度和 icon 一致

            // 中間：標題區塊（可自訂為文字、圖片、Logo……）
            Expanded(
              child: Center(child: titleWidget),
            ),

            // 右側：如果傳了 rightIcon，就顯示；沒有就放一個與左側等寬的空白 SizedBox
            if (rightIcon != null)
              SizedBox(
                width: 48, // 統一 IconButton 的寬度
                height: 48,
                child: rightIcon,
              )
            else
              const SizedBox(width: 48),
          ],
        ),
      ),
    );
  }

  /// 這裡用固定的高度 56，和 Flutter 預設 AppBar 一樣
  @override
  Size get preferredSize => const Size.fromHeight(56);
}
