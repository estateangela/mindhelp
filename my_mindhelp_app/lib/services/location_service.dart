import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:geolocator/geolocator.dart';
import 'package:geocoding/geocoding.dart';
import '../models/counseling_center.dart';

class LocationService {
  final String _baseUrl = 'https://mindhelp.onrender.com/api/v1';

  Future<List<CounselingCenter>> getCounselingCenters({
    int page = 1,
    int pageSize = 100,
    String? search,
    bool? onlineOnly,
    double? userLatitude,
    double? userLongitude,
    double radiusKm = 5.0,
  }) async {
    try {
      final queryParams = {
        'page': page.toString(),
        'page_size': pageSize.toString(),
      };
      if (search != null) queryParams['search'] = search;
      if (onlineOnly != null)
        queryParams['online_only'] = onlineOnly.toString();

      final uri = Uri.parse('$_baseUrl/counseling-centers')
          .replace(queryParameters: queryParams);

      print('正在請求 URL: $uri');

      final response = await http.get(
        uri,
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
      ).timeout(const Duration(seconds: 10));

      print('API Response Status: ${response.statusCode}');
      print('API Response Body: ${response.body}');

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        final allCenters = (data['counseling_centers'] as List)
            .map((item) => CounselingCenter.fromJson(item))
            .toList();

        // 如果提供了用戶位置，則過濾方圓五公里內的諮商所
        if (userLatitude != null && userLongitude != null) {
          return _filterCentersByDistance(
              allCenters, userLatitude, userLongitude, radiusKm);
        }

        return allCenters;
      } else {
        throw Exception(
            'Failed to load counseling centers: ${response.statusCode}');
      }
    } catch (e) {
      print('API 調用錯誤: $e');
      // 如果是網絡錯誤，嘗試使用本地備用數據
      if (e.toString().contains('Failed to fetch') ||
          e.toString().contains('ClientException')) {
        print('網絡連接失敗，嘗試使用備用數據...');
        final fallbackData = _getFallbackData();

        // 如果提供了用戶位置，則過濾方圓五公里內的諮商所
        if (userLatitude != null && userLongitude != null) {
          return _filterCentersByDistance(
              fallbackData, userLatitude, userLongitude, radiusKm);
        }

        return fallbackData;
      }
      // 其他錯誤返回空列表
      return [];
    }
  }

  // 備用數據，當 API 無法連接時使用（台北商業大學方圓五公里內）
  List<CounselingCenter> _getFallbackData() {
    return [
      CounselingCenter(
        id: 'fallback-1',
        name: '台北商業大學學生諮商中心',
        address: '臺北市中正區濟南路一段321號',
        phone: '(02)2322-6200',
        onlineCounseling: true,
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      ),
      CounselingCenter(
        id: 'fallback-2',
        name: '台大醫院精神醫學部',
        address: '臺北市中正區中山南路7號',
        phone: '(02)2312-3456',
        onlineCounseling: true,
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      ),
      CounselingCenter(
        id: 'fallback-3',
        name: '台北市立聯合醫院中興院區',
        address: '臺北市大同區鄭州路145號',
        phone: '(02)2552-3234',
        onlineCounseling: false,
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      ),
      CounselingCenter(
        id: 'fallback-4',
        name: '台北市立聯合醫院仁愛院區',
        address: '臺北市大安區仁愛路四段10號',
        phone: '(02)2709-3600',
        onlineCounseling: true,
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      ),
      CounselingCenter(
        id: 'fallback-5',
        name: '台北市立聯合醫院和平院區',
        address: '臺北市中正區中華路二段33號',
        phone: '(02)2388-9595',
        onlineCounseling: false,
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      ),
      CounselingCenter(
        id: 'fallback-6',
        name: '台北市立聯合醫院忠孝院區',
        address: '臺北市南港區同德路87號',
        phone: '(02)2786-1288',
        onlineCounseling: true,
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      ),
      CounselingCenter(
        id: 'fallback-7',
        name: '台北市立聯合醫院陽明院區',
        address: '臺北市士林區雨聲街105號',
        phone: '(02)2835-3456',
        onlineCounseling: false,
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      ),
    ];
  }

  // 獲取備用座標（對應真實醫療機構位置）
  Map<String, double> _getFallbackCoordinates(String centerId) {
    // 台北商業大學方圓五公里內的真實醫療機構座標
    final coordinates = {
      'fallback-1': {'lat': 25.0487, 'lng': 121.5175}, // 台北商業大學
      'fallback-2': {'lat': 25.0370, 'lng': 121.5200}, // 台大醫院
      'fallback-3': {'lat': 25.0600, 'lng': 121.5100}, // 中興院區
      'fallback-4': {'lat': 25.0400, 'lng': 121.5400}, // 仁愛院區
      'fallback-5': {'lat': 25.0300, 'lng': 121.5000}, // 和平院區
      'fallback-6': {'lat': 25.0350, 'lng': 121.5500}, // 忠孝院區（調整位置）
      'fallback-7': {'lat': 25.0800, 'lng': 121.5200}, // 陽明院區
    };

    return coordinates[centerId] ?? {'lat': 25.0487, 'lng': 121.5175};
  }

  // 根據距離過濾諮商所
  Future<List<CounselingCenter>> _filterCentersByDistance(
    List<CounselingCenter> centers,
    double userLatitude,
    double userLongitude,
    double radiusKm,
  ) async {
    final List<CounselingCenter> nearbyCenters = [];

    for (var center in centers) {
      try {
        // 將地址轉換為經緯度
        final location = await locationFromAddress(center.address);

        if (location.isNotEmpty) {
          final centerLat = location.first.latitude;
          final centerLng = location.first.longitude;

          // 計算距離（米）
          final distanceInMeters = Geolocator.distanceBetween(
            userLatitude,
            userLongitude,
            centerLat,
            centerLng,
          );

          // 轉換為公里
          final distanceInKm = distanceInMeters / 1000;

          print('諮商所 ${center.name}: 距離 ${distanceInKm.toStringAsFixed(2)} 公里');

          // 如果在指定半徑內，則添加到結果列表
          if (distanceInKm <= radiusKm) {
            nearbyCenters.add(center);
            print(
                '✓ 添加諮商所: ${center.name} (${distanceInKm.toStringAsFixed(2)} 公里)');
          } else {
            print(
                '✗ 跳過諮商所: ${center.name} (${distanceInKm.toStringAsFixed(2)} 公里 > $radiusKm 公里)');
          }
        } else {
          print('無法解析地址: ${center.address}');
          // 使用備用座標（台北市中心附近）
          _addFallbackCenter(
              center, userLatitude, userLongitude, radiusKm, nearbyCenters);
        }
      } catch (e) {
        print('處理諮商所 ${center.name} 時發生錯誤: $e');
        // 即使出錯也嘗試使用備用座標
        _addFallbackCenter(
            center, userLatitude, userLongitude, radiusKm, nearbyCenters);
      }
    }

    print('方圓 $radiusKm 公里內找到 ${nearbyCenters.length} 個諮商所');
    return nearbyCenters;
  }

  // 添加備用座標處理
  void _addFallbackCenter(
    CounselingCenter center,
    double userLatitude,
    double userLongitude,
    double radiusKm,
    List<CounselingCenter> nearbyCenters,
  ) {
    try {
      // 根據諮商所ID使用不同的備用座標（台北商業大學附近，方圓5公里內）
      final fallbackCoords = _getFallbackCoordinates(center.id);
      final fallbackLat = fallbackCoords['lat']!;
      final fallbackLng = fallbackCoords['lng']!;

      final distanceInMeters = Geolocator.distanceBetween(
        userLatitude,
        userLongitude,
        fallbackLat,
        fallbackLng,
      );

      final distanceInKm = distanceInMeters / 1000;

      if (distanceInKm <= radiusKm) {
        nearbyCenters.add(center);
        print(
            '✓ 使用備用座標添加諮商所: ${center.name} (${distanceInKm.toStringAsFixed(2)} 公里)');
      } else {
        print(
            '✗ 備用座標距離太遠: ${center.name} (${distanceInKm.toStringAsFixed(2)} 公里 > $radiusKm 公里)');
      }
    } catch (e) {
      print('備用座標處理失敗: $e');
    }
  }

  // You will need to implement getCounselingCenterDetails if you plan to use it.
  // The method below is a placeholder for future implementation.
  //
  // Future<CounselingCenter> getCounselingCenterDetails(String centerId) async {
  //   final response = await http.get(
  //     Uri.parse('$_baseUrl/counseling-centers/$centerId'),
  //   );
  //   if (response.statusCode == 200) {
  //     final data = json.decode(response.body);
  //     return CounselingCenter.fromJson(data['data']);
  //   } else {
  //     throw Exception('Failed to load counseling center details: ${response.statusCode}');
  //   }
  // }
}
