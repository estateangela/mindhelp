import 'dart:convert';
import 'package:http/http.dart' as http;

class OSMGeoResult {
  final double latitude;
  final double longitude;

  const OSMGeoResult({required this.latitude, required this.longitude});
}

class OSMGeocodingService {
  static const String _baseUrl = 'https://nominatim.openstreetmap.org/search';
  final http.Client _client;

  // 簡單節流（約 1 req/sec）
  static DateTime? _lastRequestAt;
  static const Duration _minInterval = Duration(milliseconds: 1100);

  // 簡單快取（記憶體）
  static final Map<String, OSMGeoResult> _cache = {};

  OSMGeocodingService({http.Client? client})
      : _client = client ?? http.Client();

  Future<OSMGeoResult?> geocodeAddress(String address,
      {Duration timeout = const Duration(seconds: 10)}) async {
    final base = address.trim();
    if (base.isEmpty) return null;

    final tried = <String>{};

    Future<OSMGeoResult?> attempt(String q) async {
      final key = _normalize(q);
      if (_cache.containsKey(key)) {
        return _cache[key];
      }

      // 節流
      final now = DateTime.now();
      if (_lastRequestAt != null) {
        final elapsed = now.difference(_lastRequestAt!);
        if (elapsed < _minInterval) {
          final wait = _minInterval - elapsed;
          await Future.delayed(wait);
        }
      }

      final uri = Uri.parse(_baseUrl).replace(queryParameters: {
        'q': q,
        'format': 'json',
        'limit': '1',
        'countrycodes': 'tw',
        'addressdetails': '0',
      });

      print('[OSM] Query: $q');
      final response = await _client.get(
        uri,
        headers: {
          // 請設定成你的實際聯絡方式/專案頁以符合 Nominatim 政策
          'User-Agent': 'mindhelp-app/1.0 (mailto:dev@mindhelp.local)',
          'Accept-Language': 'zh-TW,zh;q=0.9,en;q=0.8',
          'Accept': 'application/json',
        },
      ).timeout(timeout);

      _lastRequestAt = DateTime.now();

      if (response.statusCode != 200) {
        print('[OSM] HTTP ${response.statusCode} for "$q"');
        return null;
      }

      final decoded = json.decode(response.body);
      if (decoded is List && decoded.isNotEmpty) {
        final first = decoded.first;
        final latStr = first['lat']?.toString();
        final lonStr = first['lon']?.toString();
        final lat = latStr != null ? double.tryParse(latStr) : null;
        final lon = lonStr != null ? double.tryParse(lonStr) : null;
        if (lat != null && lon != null) {
          final result = OSMGeoResult(latitude: lat, longitude: lon);
          _cache[key] = result;
          return result;
        }
      }

      print('[OSM] Empty result for "$q"');
      return null;
    }

    // 多策略查詢順序
    final candidates = <String>[
      base,
      base + ' 台灣',
      base + ' 臺灣',
      _replaceTaiToTaiwan(base),
      _toTraditionalTaiwan(base),
      base + ', Taiwan',
    ];

    for (final q in candidates) {
      if (tried.add(q)) {
        final r = await attempt(q);
        if (r != null) return r;
      }
    }

    return null;
  }

  static String _normalize(String s) {
    return s.toLowerCase().replaceAll(RegExp(r'\s+'), ' ').trim();
  }

  static String _replaceTaiToTaiwan(String s) {
    // 若地址未明確國別，嘗試附加 Taiwan
    if (!s.contains('台灣') &&
        !s.contains('臺灣') &&
        !s.toLowerCase().contains('taiwan')) {
      return s + ' Taiwan';
    }
    return s;
  }

  static String _toTraditionalTaiwan(String s) {
    // 將常見的「台」改為「臺」以提升在 OSM 上的匹配率
    return s
        .replaceAll('台北', '臺北')
        .replaceAll('台中', '臺中')
        .replaceAll('台南', '臺南')
        .replaceAll('台東', '臺東')
        .replaceAll('台灣', '臺灣');
  }
}
