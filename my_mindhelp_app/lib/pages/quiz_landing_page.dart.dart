import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/custom_app_bar.dart';
import '../widgets/primary_button.dart';
import '../models/quiz.dart';
import 'quiz_questions_page.dart'; // 導入測驗題目頁面

class QuizLandingPage extends StatelessWidget {
  QuizLandingPage({super.key});

  // 移除欄位初始化時的 const 關鍵字，只保留 final
  final Quiz defaultQuiz = Quiz(
    id: '1',
    title: 'GAD-7 焦慮自評量表',
    summary: '這是一份快速評估您近期情緒狀態的工具，協助您了解過去兩週內的焦慮程度。\n\n**本測驗結果僅供參考，不能作為臨床診斷依據。**',
    imageUrl: 'assets/images/1.jpg',
  );

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: CustomAppBar(
        showBackButton: false,
        titleWidget: Text(
          defaultQuiz.title,
          style: TextStyle(fontSize: 24, color: AppColors.textHigh),
        ),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text(
              '自我評估',
              style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: AppColors.textHigh,
                  ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),

            // 測驗描述區塊
            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(12),
                border: Border.all(color: AppColors.accent, width: 1.5),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    '測驗目的',
                    style: Theme.of(context).textTheme.titleMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                          color: AppColors.accent,
                        ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    defaultQuiz.summary,
                    style: Theme.of(context).textTheme.bodyLarge,
                  ),
                  const SizedBox(height: 16),
                  Text(
                    '完成後您將會得到一個初步分數，幫助您更了解自己的情緒狀態。',
                    style: Theme.of(context)
                        .textTheme
                        .bodyMedium
                        ?.copyWith(fontStyle: FontStyle.italic),
                  ),
                ],
              ),
            ),

            const SizedBox(height: 48),

            // 開始測驗按鈕
            PrimaryButton(
              text: '開始測驗',
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(
                    // 這裡的 QuizQuestionsPage 移除了 const
                    builder: (context) => QuizQuestionsPage(quiz: defaultQuiz),
                  ),
                );
              },
            ),
            const SizedBox(height: 16),

            // 返回按鈕
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: Text(
                '返回主頁',
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                      decoration: TextDecoration.underline,
                      color: AppColors.textBody,
                    ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
