import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/primary_button.dart';
import '../models/quiz.dart';
import 'quiz_result_page.dart'; // 導入結果頁面

class QuizQuestionsPage extends StatefulWidget {
  final Quiz quiz;
  const QuizQuestionsPage({super.key, required this.quiz});

  @override
  State<QuizQuestionsPage> createState() => _QuizQuestionsPageState();
}

class _QuizQuestionsPageState extends State<QuizQuestionsPage> {
  int _currentQuestionIndex = 0;
  final Map<int, int> _answers =
      {}; // 儲存答案: {questionIndex: selectedOptionIndex}
  bool _isLoading = false;

  // GAD-7 焦慮量表（示範）
  final List<Map<String, dynamic>> _questions = [
    {
      'question': '感到緊張、焦慮或坐立不安。',
      'options': ['完全沒有 (0分)', '幾天 (1分)', '一半以上天數 (2分)', '幾乎每天 (3分)'],
      'score_values': [0, 1, 2, 3],
    },
    {
      'question': '不能停止或控制煩惱。',
      'options': ['完全沒有 (0分)', '幾天 (1分)', '一半以上天數 (2分)', '幾乎每天 (3分)'],
      'score_values': [0, 1, 2, 3],
    },
    {
      'question': '對各種事情煩惱過度。',
      'options': ['完全沒有 (0分)', '幾天 (1分)', '一半以上天數 (2分)', '幾乎每天 (3分)'],
      'score_values': [0, 1, 2, 3],
    },
    {
      'question': '難以放鬆。',
      'options': ['完全沒有 (0分)', '幾天 (1分)', '一半以上天數 (2分)', '幾乎每天 (3分)'],
      'score_values': [0, 1, 2, 3],
    },
  ];

  void _answerQuestion(int optionIndex) {
    setState(() {
      _answers[_currentQuestionIndex] = optionIndex;
    });
  }

  void _nextQuestion() {
    // 只有在回答當前題目後才能進入下一題
    if (_answers[_currentQuestionIndex] == null) {
      _showSnackBar('請選擇一個答案。');
      return;
    }

    if (_currentQuestionIndex < _questions.length - 1) {
      setState(() {
        _currentQuestionIndex++;
      });
    } else {
      _submitQuiz();
    }
  }

  int _calculateTotalScore() {
    int totalScore = 0;
    _answers.forEach((questionIndex, optionIndex) {
      if (_questions[questionIndex].containsKey('score_values')) {
        final scoreValues =
            _questions[questionIndex]['score_values'] as List<int>;
        if (optionIndex >= 0 && optionIndex < scoreValues.length) {
          totalScore += scoreValues[optionIndex];
        }
      }
    });
    return totalScore;
  }

  Future<void> _submitQuiz() async {
    // 檢查是否有題目未回答
    if (_answers.length != _questions.length) {
      _showSnackBar('請完成所有題目。');
      return;
    }

    setState(() {
      _isLoading = true;
    });

    // 模擬前端處理邏輯 (不呼叫 API)
    await Future.delayed(const Duration(seconds: 1));

    // 1. 計算最終分數
    final finalScore = _calculateTotalScore();

    // 2. 導航到結果頁面
    if (mounted) {
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(
          builder: (context) => QuizResultPage(
            totalScore: finalScore,
            quizTitle: widget.quiz.title,
          ),
        ),
      );
    }

    setState(() {
      _isLoading = false;
    });
  }

  void _showSnackBar(String message) {
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(message)),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final currentQuestion = _questions[_currentQuestionIndex];
    final isLastQuestion = _currentQuestionIndex == _questions.length - 1;
    final selectedOption = _answers[_currentQuestionIndex];

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: Text(widget.quiz.title),
        backgroundColor: AppColors.accent,
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // 進度條
            LinearProgressIndicator(
              value: (_currentQuestionIndex + 1) / _questions.length,
              backgroundColor: AppColors.textBody.withOpacity(0.2),
              valueColor: AlwaysStoppedAnimation<Color>(AppColors.accent),
            ),
            const SizedBox(height: 24),

            // 題目內容
            Text(
              '${_currentQuestionIndex + 1}. ${currentQuestion['question']}',
              style: Theme.of(context).textTheme.headlineMedium,
            ),
            const SizedBox(height: 16),

            // 選項列表
            ...List.generate(
              currentQuestion['options'].length,
              (index) {
                final optionText = currentQuestion['options'][index];
                final isSelected = selectedOption == index;

                // 動態調整按鈕樣式
                final ButtonStyle buttonStyle = ButtonStyle(
                  backgroundColor: WidgetStateProperty.all<Color>(
                    isSelected ? AppColors.accent : Colors.white,
                  ),
                  foregroundColor: WidgetStateProperty.all<Color>(
                    isSelected ? Colors.white : AppColors.textHigh,
                  ),
                  elevation: WidgetStateProperty.all<double>(0),
                );

                return Padding(
                  padding: const EdgeInsets.symmetric(vertical: 8.0),
                  // 使用 TextButton 模擬 PrimaryButton 的點擊效果，並直接套用樣式
                  child: TextButton(
                    style: buttonStyle,
                    child: Text(optionText,
                        style: Theme.of(context).textTheme.bodyLarge),
                    onPressed: () => _answerQuestion(index),
                  ),
                );
              },
            ),
            const Spacer(),

            // 確認/下一題按鈕
            PrimaryButton(
              text: isLastQuestion ? (_isLoading ? '提交中...' : '完成測驗') : '下一題',
              onPressed:
                  _isLoading || selectedOption == null ? () {} : _nextQuestion,
            ),
            const SizedBox(height: 16),
          ],
        ),
      ),
    );
  }
}
