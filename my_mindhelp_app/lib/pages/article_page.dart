import 'package:flutter/material.dart';
import '../core/theme.dart';
import '../models/article.dart';
import '../services/article_service.dart';

class ArticlePage extends StatefulWidget {
  ArticlePage({super.key});

  @override
  State<ArticlePage> createState() => _ArticlePageState();
}

class _ArticlePageState extends State<ArticlePage> {
  final ArticleService _articleService = ArticleService();
  List<Article> _articles = [];
  bool _isLoading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadArticles();
  }

  Future<void> _loadArticles() async {
    try {
      final response = await _articleService.getArticles();
      setState(() {
        _articles = response.articles;
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
    }
  }

  // 備用靜態文章資料
  final List<Article> _staticArticles = [
    Article(
      id: '1',
      title: '如何應對職場壓力？',
      author: Author(name: '張心理師', title: '臨床心理師'),
      summary: '學會辨識壓力源，並透過正念練習、時間管理等技巧來有效緩解工作帶來的焦慮與疲憊。',
      content: '職場壓力是現代人常見的問題...',
      publishDate: '2025-01-01T00:00:00Z',
      tags: ['壓力管理', '職場心理'],
      isBookmarked: false,
      viewCount: 150,
    ),
    Article(
      id: '2',
      title: '走出情緒低谷的七個步驟',
      author: Author(name: '李心理師', title: '諮商心理師'),
      summary: '情緒低落是正常的，但當它持續影響生活時，不妨嘗試這七個實用步驟，幫助你重新找回內心的平靜與力量。',
      content: '情緒低落是人生中不可避免的經歷...',
      publishDate: '2025-01-02T00:00:00Z',
      tags: ['情緒管理', '心理健康'],
      isBookmarked: false,
      viewCount: 200,
    ),
    Article(
      id: '3',
      title: '親密關係中的有效溝通',
      author: Author(name: '王心理師', title: '婚姻家庭治療師'),
      summary: '溝通是維繫關係的橋樑。本專欄將探討如何在與伴侶、家人或朋友的互動中，建立健康且有建設性的溝通模式。',
      content: '良好的溝通是健康關係的基石...',
      publishDate: '2025-01-03T00:00:00Z',
      tags: ['人際關係', '溝通技巧'],
      isBookmarked: false,
      viewCount: 180,
    ),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: Text('專家專欄', style: Theme.of(context).textTheme.headlineLarge),
        centerTitle: true,
        backgroundColor: Colors.transparent,
        elevation: 0,
        actions: [
          IconButton(
            icon: const Icon(Icons.notifications, color: AppColors.textHigh),
            onPressed: () => Navigator.pushNamed(context, '/notify'),
          ),
        ],
      ),
      body: _buildBody(context),
    );
  }

  Widget _buildBody(BuildContext context) {
    if (_isLoading) {
      return const Center(
        child: CircularProgressIndicator(),
      );
    }

    if (_error != null) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.error_outline,
              size: 64,
              color: AppColors.textBody,
            ),
            const SizedBox(height: 16),
            Text(
              '載入文章失敗',
              style: Theme.of(context).textTheme.headlineSmall,
            ),
            const SizedBox(height: 8),
            Text(
              _error!,
              style: Theme.of(context).textTheme.bodyMedium,
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: _loadArticles,
              child: const Text('重試'),
            ),
          ],
        ),
      );
    }

    final articlesToShow = _articles.isNotEmpty ? _articles : _staticArticles;

    return RefreshIndicator(
      onRefresh: _loadArticles,
      child: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: articlesToShow.map((article) {
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
              // 已在 Article Page，不做動作
              break;
            case 1:
              Navigator.pushReplacementNamed(context, '/maps');
              break;
            case 2:
              Navigator.pushReplacementNamed(context, '/chat');
              break;
            case 3:
              Navigator.pushReplacementNamed(context, '/profile');
              break;
          }
        },
        items: const [
          BottomNavigationBarItem(icon: Icon(Icons.article), label: '專欄'),
          BottomNavigationBarItem(icon: Icon(Icons.location_on), label: '地圖'),
          BottomNavigationBarItem(icon: Icon(Icons.chat_bubble), label: '聊天'),
          BottomNavigationBarItem(icon: Icon(Icons.person), label: '我的'),
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
        onTap: () {
          // TODO: 點擊後導航到文章詳情頁面
          print('點擊文章：${article.title}');
        },
        borderRadius: BorderRadius.circular(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            SizedBox(
              height: 150,
              child: ClipRRect(
                borderRadius:
                    const BorderRadius.vertical(top: Radius.circular(12)),
                child: Container(
                  height: 150,
                  color: Colors.grey[200],
                  child: const Center(
                    child: Icon(Icons.article, size: 64, color: Colors.grey),
                  ),
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
                    '作者：${article.author.name}',
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
