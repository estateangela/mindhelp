import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../models/article.dart';
import '../widgets/custom_app_bar.dart';
//import 'article_detail_page.dart';

class ArticlePage extends StatelessWidget {
  ArticlePage({super.key});

  final List<Article> _articles = [
    // 使用本地圖片路徑
    Article(
      id: '1',
      title: '如何應對職場壓力？',
      author: '張心理師',
      summary: '學會辨識壓力源，並透過正念練習、時間管理等技巧來有效緩解工作帶來的焦慮與疲憊。',
      imageUrl: 'assets/images/1.jpg',
    ),
    Article(
      id: '2',
      title: '走出情緒低谷的七個步驟',
      author: '李心理師',
      summary: '情緒低落是正常的，但當它持續影響生活時，不妨嘗試這七個實用步驟，幫助你重新找回內心的平靜與力量。',
      imageUrl: 'assets/images/3.jpg',
    ),
    Article(
      id: '3',
      title: '親密關係中的有效溝通',
      author: '王心理師',
      summary: '溝通是維繫關係的橋樑。本專欄將探討如何在與伴侶、家人或朋友的互動中，建立健康且有建設性的溝通模式。',
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
          children: _articles.map((article) {
            return _buildArticleCard(context, article);
          }).toList(),
        ),
      ),
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: 0, // 假設此頁為首頁
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

  Widget _buildArticleCard(BuildContext context, Article article) {
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
        // onTap: () {
        // Navigator.push(
        // context,
        //MaterialPageRoute(
        //builder: (context) => ArticleDetailPage(article: article),
        //),
        //);
        //},
        borderRadius: BorderRadius.circular(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            SizedBox(
              height: 150,
              child: ClipRRect(
                borderRadius:
                    const BorderRadius.vertical(top: Radius.circular(12)),
                // 將 Image.network 替換為 Image.asset
                child: Image.asset(
                  article.imageUrl,
                  fit: BoxFit.cover,
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    article.title,
                    style: Theme.of(context).textTheme.titleLarge?.copyWith(
                          fontWeight: FontWeight.bold,
                          color: AppColors.textHigh,
                        ),
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                  const SizedBox(height: 4),
                  Text(
                    '作者：${article.author}',
                    style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          color: AppColors.textBody,
                        ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    article.summary,
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

class Article {
  final String id;
  final String title;
  final String author;
  final String summary;
  final String imageUrl;

  Article({
    required this.id,
    required this.title,
    required this.author,
    required this.summary,
    required this.imageUrl,
  });
}
