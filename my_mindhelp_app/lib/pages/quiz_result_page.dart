import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/primary_button.dart';

class QuizResultPage extends StatelessWidget {
  final int totalScore;
  final String quizTitle;

  const QuizResultPage({
    super.key,
    required this.totalScore,
    required this.quizTitle,
  });

  // 根據分數範圍生成結果和建議
  Map<String, String> _getResult(int score) {
    if (score >= 9) {
      return {
        'title': '中重度困擾',
        'message': '您的得分顯示可能正經歷**較嚴重的**情緒困擾。建議您尋求專業心理諮商或身心科醫師的協助，不需獨自面對。',
        'color': '#D32F2F', // 紅色
      };
    } else if (score >= 5) {
      return {
        'title': '輕至中度困擾',
        'message':
            '您的得分顯示可能存在**輕微或中度**的情緒困擾。建議您可以多利用 App 內資源地圖尋找附近諮商所，或與 AI 聊天助手傾訴。',
        'color': '#FFB300', // 橙色
      };
    } else {
      return {
        'title': '情緒健康',
        'message': '您的分數處於正常範圍。請繼續保持健康的生活習慣和心態，並善用 App 內文章來維持心理健康。',
        'color': '#388E3C', // 綠色
      };
    }
  }

  @override
  Widget build(BuildContext context) {
    final result = _getResult(totalScore);
    final color = Color(
        int.parse(result['color']!.substring(1, 7), radix: 16) + 0xFF000000);

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: Text(quizTitle, style: Theme.of(context).textTheme.titleLarge),
        backgroundColor: color,
        automaticallyImplyLeading: false, // 移除返回按鈕
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text(
              '您的測驗結果',
              style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: AppColors.textHigh,
                  ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 30),

            // 顯示得分區塊
            Container(
              padding: const EdgeInsets.all(24),
              decoration: BoxDecoration(
                color: color.withOpacity(0.1),
                borderRadius: BorderRadius.circular(12),
                border: Border.all(color: color, width: 2),
              ),
              child: Column(
                children: [
                  Text('總得分', style: Theme.of(context).textTheme.titleMedium),
                  const SizedBox(height: 8),
                  Text(
                    totalScore.toString(),
                    style: TextStyle(
                      fontSize: 60,
                      fontWeight: FontWeight.bold,
                      color: color,
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 24),

            // 結果標題
            Text(
              result['title']!,
              style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: color,
                  ),
            ),
            const SizedBox(height: 12),

            // 詳細建議
            Text(
              result['message']!,
              style: Theme.of(context).textTheme.bodyLarge,
            ),
            const SizedBox(height: 48),

            // 行動呼籲按鈕
            PrimaryButton(
              text: '完成並回首頁',
              onPressed: () {
                Navigator.popUntil(context, ModalRoute.withName('/home'));
              },
            ),
          ],
        ),
      ),
    );
  }
}
