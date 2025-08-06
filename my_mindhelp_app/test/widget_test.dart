// test/widget_test.dart
import 'package:flutter_test/flutter_test.dart';
import 'package:my_mindhelp_app/main.dart'; // 确保路径和包名正确

void main() {
  testWidgets('根组件能正常加载', (tester) async {
    await tester.pumpWidget(const MindHelpApp());
    // …你的断言
  });
}
