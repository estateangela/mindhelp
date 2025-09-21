import 'package:json_annotation/json_annotation.dart';

part 'quiz.g.dart';

@JsonSerializable()
class Quiz {
  final String id;
  final String title;
  final String description;
  final List<QuizQuestion> questions;
  final String category;

  Quiz({
    required this.id,
    required this.title,
    required this.description,
    required this.questions,
    required this.category,
  });

  factory Quiz.fromJson(Map<String, dynamic> json) => _$QuizFromJson(json);
  Map<String, dynamic> toJson() => _$QuizToJson(this);
}

@JsonSerializable()
class QuizQuestion {
  final String id;
  final String question;
  final List<QuizOption> options;
  final int order;

  QuizQuestion({
    required this.id,
    required this.question,
    required this.options,
    required this.order,
  });

  factory QuizQuestion.fromJson(Map<String, dynamic> json) => _$QuizQuestionFromJson(json);
  Map<String, dynamic> toJson() => _$QuizQuestionToJson(this);
}

@JsonSerializable()
class QuizOption {
  final String id;
  final String text;
  final int score;

  QuizOption({
    required this.id,
    required this.text,
    required this.score,
  });

  factory QuizOption.fromJson(Map<String, dynamic> json) => _$QuizOptionFromJson(json);
  Map<String, dynamic> toJson() => _$QuizOptionToJson(this);
}

@JsonSerializable()
class QuizSubmission {
  final String quizId;
  final List<QuizAnswer> answers;

  QuizSubmission({
    required this.quizId,
    required this.answers,
  });

  factory QuizSubmission.fromJson(Map<String, dynamic> json) => _$QuizSubmissionFromJson(json);
  Map<String, dynamic> toJson() => _$QuizSubmissionToJson(this);
}

@JsonSerializable()
class QuizAnswer {
  final String questionId;
  final String optionId;

  QuizAnswer({
    required this.questionId,
    required this.optionId,
  });

  factory QuizAnswer.fromJson(Map<String, dynamic> json) => _$QuizAnswerFromJson(json);
  Map<String, dynamic> toJson() => _$QuizAnswerToJson(this);
}

@JsonSerializable()
class QuizResult {
  final String id;
  final String quizId;
  final String quizTitle;
  final int score;
  final String result;
  final String interpretation;
  final String completedAt;

  QuizResult({
    required this.id,
    required this.quizId,
    required this.quizTitle,
    required this.score,
    required this.result,
    required this.interpretation,
    required this.completedAt,
  });

  factory QuizResult.fromJson(Map<String, dynamic> json) => _$QuizResultFromJson(json);
  Map<String, dynamic> toJson() => _$QuizResultToJson(this);
}

@JsonSerializable()
class QuizListResponse {
  final List<Quiz> quizzes;
  final int total;
  final int page;
  final int pageSize;

  QuizListResponse({
    required this.quizzes,
    required this.total,
    required this.page,
    required this.pageSize,
  });

  factory QuizListResponse.fromJson(Map<String, dynamic> json) => 
      _$QuizListResponseFromJson(json);
  Map<String, dynamic> toJson() => _$QuizListResponseToJson(this);
}

@JsonSerializable()
class QuizHistoryResponse {
  final List<QuizResult> results;
  final int total;
  final int page;
  final int pageSize;

  QuizHistoryResponse({
    required this.results,
    required this.total,
    required this.page,
    required this.pageSize,
  });

  factory QuizHistoryResponse.fromJson(Map<String, dynamic> json) => 
      _$QuizHistoryResponseFromJson(json);
  Map<String, dynamic> toJson() => _$QuizHistoryResponseToJson(this);
}
