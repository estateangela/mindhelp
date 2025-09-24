import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../widgets/custom_app_bar.dart';

class QuizPage extends StatefulWidget {
  const QuizPage({super.key});

  @override
  State<QuizPage> createState() => _QuizPageState();
}

class _QuizPageState extends State<QuizPage> {
  final List<Quiz> _quizzes = [
    Quiz(
      id: '1',
      title: 'GAD-7 焦慮自評量表',
      summary: '用來初步篩檢廣泛性焦慮症的工具，幫助你了解過去兩週內的焦慮程度。',
      imageUrl: 'assets/images/1.jpg',
    ),
    Quiz(
      id: '2',
      title: 'PHQ-9 憂鬱症篩檢量表',
      summary: '九個問題，幫助你快速評估自己是否可能有憂鬱症狀，是初步篩檢的常用工具。',
      imageUrl: 'assets/images/3.jpg',
    ),
    Quiz(
      id: '3',
      title: '壓力自我評估量表',
      summary: '從多個角度評估你當前的壓力水平，幫助你辨識壓力來源並採取應對措施。',
      imageUrl: 'assets/images/4.jpg',
    ),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: CustomAppBar(
        showBackButton: false,
        titleWidget: const Image(
          image: AssetImage('assets/images/mindhelp.png'),
          width: 200,
          fit: BoxFit.contain,
        ),
        rightIcon: IconButton(
          icon: const Icon(Icons.notifications, color: AppColors.textHigh),
          onPressed: () => Navigator.pushNamed(context, '/notify'),
        ),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: _quizzes.map((quiz) {
            return _buildQuizCard(context, quiz);
          }).toList(),
        ),
      ),
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: 0,
        selectedItemColor: AppColors.accent,
        unselectedItemColor: AppColors.textBody,
        onTap: (idx) {
          switch (idx) {
            case 0:
              Navigator.pushReplacementNamed(context, '/home');
              break;

            case 1:
              Navigator.pushReplacementNamed(context, '/maps');
              break;
            case 2:
              Navigator.pushReplacementNamed(context, '/chat');
              break;
          }
        },
        items: const [
          BottomNavigationBarItem(icon: Icon(Icons.home), label: 'Home'),
          BottomNavigationBarItem(icon: Icon(Icons.location_on), label: 'Maps'),
          BottomNavigationBarItem(icon: Icon(Icons.chat_bubble), label: 'Chat'),
        ],
      ),
    );
  }

  Widget _buildQuizCard(BuildContext context, Quiz quiz) {
    return Container(
      margin: const EdgeInsets.only(bottom: 16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.grey.withOpacity(0.1),
            spreadRadius: 1,
            blurRadius: 5,
            offset: const Offset(0, 3),
          ),
        ],
      ),
      child: InkWell(
        onTap: () {
          // TODO: 點擊後導航到測驗題目頁面
          print('點擊測驗：${quiz.title}');
        },
        borderRadius: BorderRadius.circular(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            ClipRRect(
              borderRadius:
                  const BorderRadius.vertical(top: Radius.circular(12)),
              child: Image.network(
                quiz.imageUrl,
                height: 150,
                fit: BoxFit.cover,
                errorBuilder: (context, error, stackTrace) {
                  return Container(
                    height: 150,
                    color: Colors.grey[200],
                    child: const Center(child: Icon(Icons.image_not_supported)),
                  );
                },
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    quiz.title,
                    style: Theme.of(context).textTheme.titleLarge?.copyWith(
                          fontWeight: FontWeight.bold,
                          color: AppColors.textHigh,
                        ),
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                  const SizedBox(height: 8),
                  Text(
                    quiz.summary,
                    style: Theme.of(context).textTheme.bodyMedium,
                    maxLines: 3,
                    overflow: TextOverflow.ellipsis,
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class Quiz {
  final String id;
  final String title;
  final String summary;
  final String imageUrl;

  Quiz({
    required this.id,
    required this.title,
    required this.summary,
    required this.imageUrl,
  });
}
