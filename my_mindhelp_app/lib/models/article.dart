import 'package:json_annotation/json_annotation.dart';

part 'article.g.dart';

@JsonSerializable()
class Article {
  final String id;
  final String title;
  final Author author;
  final String summary;
  final String content;
  final String publishDate;
  final List<String> tags;
  final bool isBookmarked;
  final int viewCount;

  Article({
    required this.id,
    required this.title,
    required this.author,
    required this.summary,
    required this.content,
    required this.publishDate,
    required this.tags,
    required this.isBookmarked,
    required this.viewCount,
  });

  factory Article.fromJson(Map<String, dynamic> json) => _$ArticleFromJson(json);
  Map<String, dynamic> toJson() => _$ArticleToJson(this);
}

@JsonSerializable()
class Author {
  final String name;
  final String title;

  Author({
    required this.name,
    required this.title,
  });

  factory Author.fromJson(Map<String, dynamic> json) => _$AuthorFromJson(json);
  Map<String, dynamic> toJson() => _$AuthorToJson(this);
}

@JsonSerializable()
class ArticleListResponse {
  final List<Article> articles;
  final int total;
  final int page;
  final int pageSize;

  ArticleListResponse({
    required this.articles,
    required this.total,
    required this.page,
    required this.pageSize,
  });

  factory ArticleListResponse.fromJson(Map<String, dynamic> json) => 
      _$ArticleListResponseFromJson(json);
  Map<String, dynamic> toJson() => _$ArticleListResponseToJson(this);
}
