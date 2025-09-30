import 'package:flutter/material.dart';
import 'package:flutter_markdown/flutter_markdown.dart'; // 引入 Markdown 套件
import '../core/theme.dart';
import '../models/article.dart';

class ArticleDetailPage extends StatelessWidget {
  final Article article;

  const ArticleDetailPage({super.key, required this.article});

  // 這個方法會返回一個適合放在 Dialog 裡的 Widget
  @override
  Widget build(BuildContext context) {
    // 硬編碼的文章內容，用來模擬後端回傳的完整內容
    final String hardcodedContent = _getArticleContent(article.id);

    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        mainAxisSize: MainAxisSize.min,
        children: [
          // 圖片
          ClipRRect(
            borderRadius: BorderRadius.circular(12),
            child: Image.asset(
              article.imageUrl, // 使用本地圖片路徑
              height: 150,
              fit: BoxFit.cover,
            ),
          ),
          const SizedBox(height: 16),

          // 標題
          Text(
            article.title,
            style: Theme.of(context).textTheme.titleLarge?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: AppColors.textHigh,
                ),
          ),
          const SizedBox(height: 8),

          // 作者
          Text(
            '作者：${article.author}',
            style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                  color: AppColors.textBody,
                ),
          ),
          const SizedBox(height: 16),

          // 文章內容 (使用 MarkdownBody 渲染)
          // 注意：Dialog 內部的空間限制，如果文章太長，MarkdownBody 搭配 SingleChildScrollView 效果最好
          MarkdownBody(
            data: hardcodedContent,
            styleSheet: MarkdownStyleSheet(
              // 定義標題樣式
              h3: Theme.of(context).textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: AppColors.accent,
                  ),
              // 定義內文樣式
              p: Theme.of(context).textTheme.bodyLarge,
              // 定義列表樣式
              listBullet: Theme.of(context).textTheme.bodyLarge,
            ),
            shrinkWrap: true, // 讓 MarkdownBody 配合 Dialog 內容大小
          ),
        ],
      ),
    );
  }

  // 根據 ID 獲取硬編碼的文章內容
  String _getArticleContent(String id) {
    switch (id) {
      case '1':
        return """
        職場壓力是現代人普遍的困擾，學會有效應對至關重要。

        ### 辨識壓力源
        首先，你必須準確找出壓力的來源。是工作量過大？還是人際關係緊張？一旦確定了壓力源，才能對症下藥。

        ### 正念練習
        每天花 5-10 分鐘進行正念呼吸，可以幫助你將注意力拉回當下，減少焦慮。這不需要複雜的技巧，只需專注於每一次呼吸的進出。

        ### 時間管理
        使用番茄工作法（Pomodoro Technique）或其他時間管理工具，將大任務分解成小塊，確保工作和休息的平衡。

        ### 尋求支持
        不要獨自承受壓力。與信任的家人、朋友或專業心理師傾訴，獲得社會支持是緩解壓力的關鍵步驟。
        """;
      case '2':
        return """
        走出情緒低谷的七個步驟：
        
        1. **允許自己感受**：不要壓抑情緒，承認自己的感受。
        2. **規律生活**：保持固定的作息時間。
        3. **適度運動**：運動是天然的抗抑鬱劑。
        4. **健康飲食**：確保攝取足夠的維生素和礦物質。
        5. **與人連結**：保持社交聯繫，即使只是簡短的聊天。
        6. **設定小目標**：從完成一件小事開始，建立成就感。
        7. **尋求專業幫助**：如果情緒持續低落，請諮詢專業人士。
        """;
      case '3':
        return """
        親密關係中的有效溝通的要素：
        
        * **積極聆聽**：在對方說話時，專注於聆聽而不是思考如何回應。
        * **使用「我」的訊息**：用「我覺得...」來表達感受，而不是用「你總是...」來指責對方。
        * **設定界線**：清楚地表達自己的需求和極限。
        * **避免爭吵中的逃避**：即使感到不舒服，也應保持溝通管道暢通。
        """;
      default:
        return "找不到文章內容。";
    }
  }
}
