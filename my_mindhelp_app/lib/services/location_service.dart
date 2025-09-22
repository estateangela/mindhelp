import 'package:http/http.dart' as http;
import 'dart:convert';
import '../models/resource.dart';

class LocationService {
  // 將基底 URL 替換為你的後端 API 網址
  final String _baseUrl = 'https://mindhelp.onrender.com/v1';
  // TODO: 請替換為您實際的 JWT Token
  final String _authHeader = 'Bearer <YOUR_JWT_TOKEN>';

  // 搜尋附近的醫療資源
  Future<List<Resource>> searchResources({
    required double lat,
    required double lon,
    int radius = 5000,
    String? type,
    String? specialty,
  }) async {
    // 構建查詢參數
    final Map<String, dynamic> queryParams = {
      'latitude': lat.toString(),
      'longitude': lon.toString(),
      'radius': (radius / 1000).toString(), // API 參數是公里，前端是公尺
    };
    if (type != null) queryParams['type'] = type;
    if (specialty != null) queryParams['specialty'] = specialty;

    final uri = Uri.parse('$_baseUrl/locations/search')
        .replace(queryParameters: queryParams);

    final response = await http.get(
      uri,
      headers: {'Authorization': _authHeader},
    );

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      return (data['data']['locations'] as List)
          .map((item) => Resource.fromJson(item))
          .toList();
    } else {
      throw Exception('Failed to load resources: ${response.statusCode}');
    }
  }

  // 獲取單一資源詳情
  Future<Resource> getResourceDetails(String resourceId) async {
    final response = await http.get(
      Uri.parse('$_baseUrl/locations/$resourceId'),
      headers: {'Authorization': _authHeader},
    );

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      return Resource.fromJson(data['data']);
    } else {
      throw Exception(
          'Failed to load resource details: ${response.statusCode}');
    }
  }

  // 獲取 Google Maps 地址資訊
  Future<List<Map<String, dynamic>>> getGoogleAddresses({
    required double lat,
    required double lon,
    int radius = 5000,
  }) async {
    final Map<String, dynamic> queryParams = {
      'latitude': lat.toString(),
      'longitude': lon.toString(),
      'radius': radius.toString(),
    };

    final uri = Uri.parse('$_baseUrl/maps/google-addresses')
        .replace(queryParameters: queryParams);

    final response = await http.get(uri);

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      return (data['data'] as List).cast<Map<String, dynamic>>();
    } else {
      throw Exception(
          'Failed to load Google addresses: ${response.statusCode}');
    }
  }
}
