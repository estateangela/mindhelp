import 'package:flutter/material.dart';
import 'package:flutter_markdown/flutter_markdown.dart'; // 引入 Markdown 套件
import '../core/theme.dart';
import '../widgets/primary_button.dart';
import '../models/quiz.dart'; // 導入 Quiz 模型
import 'quiz_questions_page.dart'; // 導入題目頁面

class QuizDetailPage extends StatelessWidget {
  final Quiz quiz;

  const QuizDetailPage({super.key, required this.quiz});

  // 根據 ID 獲取硬編碼的測驗內容 (問題列表)
  String _getQuizQuestions(String id) {
    switch (id) {
      case '1':
        // GAD-7 焦慮自評量表
        return """
**GAD-7 焦慮自評量表**（過去兩週，您感到煩惱的頻率為？）：
---
1. 感到緊張、焦慮或坐立不安。
2. 不能停止或控制煩惱。
3. 對各種事情煩惱過度。
4. 難以放鬆。
5. 坐立不安，以至於很難靜坐。
6. 變得容易煩躁或發脾氣。
7. 好像有什麼可怕的事情即將發生。
""";
      case '2':
        // PHQ-9 憂鬱症篩檢量表
        return """
**PHQ-9 憂鬱症篩檢量表**（過去兩週，您感到煩惱的頻率為？）：
---
1. 做事時失去興趣或樂趣。
2. 感到心情低落、憂鬱或絕望。
3. 難以入睡、睡不安穩或睡得太多。
4. 感到疲倦或沒有活力。
5. 食慾不振或吃得太多。
6. 覺得自己很糟，或覺得自己是個失敗者，或讓自己和家人失望。
7. 難以集中精神，例如看報紙或看電視時。
8. 動作或說話很慢，以至於別人可能已經察覺。或者正好相反—煩躁或坐立不安，動來動去。
9. 有不如死掉或用某種方式傷害自己的念頭。
""";
      case '3':
        // 壓力自我評估量表
        return """
**壓力自我評估量表**：
---
1. **工作壓力**：我是否覺得工作量已經超出我能承受的範圍？
2. **人際關係**：我是否在與家人或朋友的溝通中感到持續的緊張和疲憊？
3. **睡眠品質**：我是否經常難以入睡，或醒來後感到沒有充分休息？
4. **情緒穩定**：我是否經常感到煩躁、易怒或情緒波動大？
5. **休閒時間**：我是否缺乏足夠的休閒時間或興趣愛好來放鬆身心？
""";
      default:
        return "找不到測驗內容。";
    }
  }

  @override
  Widget build(BuildContext context) {
    final String hardcodedContent = _getQuizQuestions(quiz.id);

    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        mainAxisSize: MainAxisSize.min,
        children: [
          // 圖片
          ClipRRect(
            borderRadius: BorderRadius.circular(12),
            child: Image.asset(
              quiz.imageUrl,
              height: 150,
              fit: BoxFit.cover,
            ),
          ),
          const SizedBox(height: 16),

          // 測驗名稱 (作為內容的一部分)
          Text(
            quiz.title,
            style: Theme.of(context).textTheme.titleLarge?.copyWith(
              fontWeight: FontWeight.bold,
              color: AppColors.textHigh,
            ),
          ),
          const SizedBox(height: 8),

          // 測驗簡介
          Text(
            '測驗簡介:',
            style: Theme.of(context).textTheme.titleMedium?.copyWith(
              fontWeight: FontWeight.bold,
              color: AppColors.accent,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            quiz.summary,
            style: Theme.of(context).textTheme.bodyMedium,
          ),
          const SizedBox(height: 16),

          // 測驗內容 (MarkdownBody 渲染問題列表)
          MarkdownBody(
            data: hardcodedContent,
            styleSheet: MarkdownStyleSheet(
              p: Theme.of(context).textTheme.bodyLarge,
              strong: Theme.of(context).textTheme.bodyLarge?.copyWith(fontWeight: FontWeight.bold),
              listBullet: Theme.of(context).textTheme.bodyLarge,
            ),
            shrinkWrap: true,
          ),
          const SizedBox(height: 16),

          // 提示開始按鈕
          PrimaryButton(
            text: '開始測驗',
            onPressed: () {
              Navigator.of(context).pop(); // 先關閉對話框
              // 導航到測驗題目頁面，並傳遞 Quiz 物件
              Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (context) => QuizQuestionsPage(quiz: quiz),
                ),
              );
            },
          ),
        ],
      ),
    );
  }
}
