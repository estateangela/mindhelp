import '../core/api_client.dart';
import '../core/api_config.dart';
import '../models/quiz.dart';

class QuizService {
  final ApiClient _apiClient = ApiClient();

  /// 獲取測驗列表
  Future<QuizListResponse> getQuizzes({
    String? category,
    int page = 1,
    int limit = 10,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'limit': limit,
    };

    if (category != null && category.isNotEmpty) {
      queryParams['category'] = category;
    }

    final response = await _apiClient.dio.get(
      ApiConfig.quizzes,
      queryParameters: queryParams,
    );

    if (response.statusCode == 200) {
      return QuizListResponse.fromJson(response.data['data']);
    }
    throw Exception('獲取測驗列表失敗');
  }

  /// 獲取單一測驗詳情
  Future<Quiz> getQuiz(String quizId) async {
    final response = await _apiClient.dio.get('${ApiConfig.quizzes}/$quizId');

    if (response.statusCode == 200) {
      return Quiz.fromJson(response.data['data']);
    }
    throw Exception('獲取測驗詳情失敗');
  }

  /// 提交測驗答案
  Future<QuizResult> submitQuiz({
    required String quizId,
    required List<QuizAnswer> answers,
  }) async {
    final submission = QuizSubmission(
      quizId: quizId,
      answers: answers,
    );

    final response = await _apiClient.dio.post(
      '${ApiConfig.quizzes}/$quizId/submit',
      data: submission.toJson(),
    );

    if (response.statusCode == 200) {
      return QuizResult.fromJson(response.data['data']);
    }
    throw Exception('提交測驗失敗');
  }

  /// 獲取使用者測驗歷史
  Future<QuizHistoryResponse> getQuizHistory({
    int page = 1,
    int limit = 10,
  }) async {
    final response = await _apiClient.dio.get(
      '${ApiConfig.users}/me/quiz_history',
      queryParameters: {
        'page': page,
        'limit': limit,
      },
    );

    if (response.statusCode == 200) {
      return QuizHistoryResponse.fromJson(response.data['data']);
    }
    throw Exception('獲取測驗歷史失敗');
  }
}
