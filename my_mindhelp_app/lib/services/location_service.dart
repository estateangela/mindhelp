import 'package:http/http.dart' as http;
import 'dart:convert';
import '../models/counseling_center.dart';

class LocationService {
  // 將基底 URL 替換為你的後端 API 網址
  final String _baseUrl = 'https://mindhelp.onrender.com/api/v1';

  // 獲取諮商所列表
  Future<List<CounselingCenter>> getCounselingCenters({
    int page = 1,
    int pageSize = 10,
    String? search,
    bool? onlineOnly,
  }) async {
    final queryParams = {
      'page': page.toString(),
      'page_size': pageSize.toString(),
    };
    if (search != null) queryParams['search'] = search;
    if (onlineOnly != null) queryParams['online_only'] = onlineOnly.toString();

    final uri = Uri.parse('$_baseUrl/counseling-centers')
        .replace(queryParameters: queryParams);

    final response = await http.get(uri);

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      return (data['counseling_centers'] as List)
          .map((item) => CounselingCenter.fromJson(item))
          .toList();
    } else {
      throw Exception(
          'Failed to load counseling centers: ${response.statusCode}');
    }
  }

  // 獲取單一諮商所詳情
  Future<CounselingCenter> getCounselingCenterDetails(String centerId) async {
    final response = await http.get(
      Uri.parse('$_baseUrl/counseling-centers/$centerId'),
    );

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      return CounselingCenter.fromJson(data['data']);
    } else {
      throw Exception(
          'Failed to load counseling center details: ${response.statusCode}');
    }
  }
}
