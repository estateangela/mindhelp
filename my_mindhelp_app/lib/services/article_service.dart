import '../core/api_client.dart';
import '../core/api_config.dart';
import '../models/article.dart';

class ArticleService {
  final ApiClient _apiClient = ApiClient();

  /// 獲取文章列表
  Future<ArticleListResponse> getArticles({
    String? search,
    String? tag,
    String? sortBy,
    int page = 1,
    int limit = 10,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'limit': limit,
    };

    if (search != null && search.isNotEmpty) {
      queryParams['q'] = search;
    }
    if (tag != null && tag.isNotEmpty) {
      queryParams['tag'] = tag;
    }
    if (sortBy != null && sortBy.isNotEmpty) {
      queryParams['sort_by'] = sortBy;
    }

    final response = await _apiClient.dio.get(
      ApiConfig.articles,
      queryParameters: queryParams,
    );

    if (response.statusCode == 200) {
      return ArticleListResponse.fromJson(response.data['data']);
    }
    throw Exception('獲取文章列表失敗');
  }

  /// 獲取單篇文章詳情
  Future<Article> getArticle(String articleId) async {
    final response = await _apiClient.dio.get('${ApiConfig.articles}/$articleId');

    if (response.statusCode == 200) {
      return Article.fromJson(response.data['data']);
    }
    throw Exception('獲取文章詳情失敗');
  }

  /// 收藏文章
  Future<void> bookmarkArticle(String articleId) async {
    final response = await _apiClient.dio.post(
      '${ApiConfig.articles}/$articleId/bookmark',
    );

    if (response.statusCode != 204 && response.statusCode != 200) {
      throw Exception('收藏文章失敗');
    }
  }

  /// 取消收藏文章
  Future<void> unbookmarkArticle(String articleId) async {
    final response = await _apiClient.dio.delete(
      '${ApiConfig.articles}/$articleId/bookmark',
    );

    if (response.statusCode != 204 && response.statusCode != 200) {
      throw Exception('取消收藏文章失敗');
    }
  }

  /// 獲取使用者收藏的文章列表
  Future<ArticleListResponse> getBookmarkedArticles({
    int page = 1,
    int limit = 10,
  }) async {
    final response = await _apiClient.dio.get(
      '${ApiConfig.users}/me/bookmarks/articles',
      queryParameters: {
        'page': page,
        'limit': limit,
      },
    );

    if (response.statusCode == 200) {
      return ArticleListResponse.fromJson(response.data['data']);
    }
    throw Exception('獲取收藏文章列表失敗');
  }
}
