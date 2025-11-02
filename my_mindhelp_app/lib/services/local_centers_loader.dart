import 'package:csv/csv.dart';
import 'package:flutter/services.dart' show rootBundle;
import '../models/counseling_center.dart';

class LocalCentersLoader {
  static const String defaultAssetPath =
      'assets/data/geocoded_centers_google.csv';
  static const String nameMappingAssetPath =
      'assets/data/centers_name_mapping.csv';

  /// 從 assets CSV 載入諮商所資料，CSV 欄位：
  /// id,address,latitude,longitude,formatted_address,status
  static Future<List<CounselingCenter>> loadFromCsv({String? assetPath}) async {
    final path = assetPath ?? defaultAssetPath;
    final csvString = await rootBundle.loadString(path);

    // 使用 CSV 套件解析，避免逗號與引號的陷阱
    final rows = const CsvToListConverter(
      eol: '\n',
      shouldParseNumbers: false,
    ).convert(csvString);

    if (rows.isEmpty) return [];

    // 解析表頭索引
    final header = rows.first.map((e) => e.toString().trim()).toList();
    int idxId = header.indexOf('id');
    int idxAddress = header.indexOf('address');
    int idxLat = header.indexOf('latitude');
    int idxLng = header.indexOf('longitude');
    int idxFormatted = header.indexOf('formatted_address');
    int idxName = header.indexOf('name');

    // 載入地址→名稱映射（若存在）
    final Map<String, String> addressToName = await _loadNameMapping();

    final nowIso = DateTime.now().toUtc().toIso8601String();

    final List<CounselingCenter> centers = [];
    for (int i = 1; i < rows.length; i++) {
      final row = rows[i];
      if (row.isEmpty) continue;

      String id = _safeCell(row, idxId);
      String address = _safeCell(row, idxAddress);
      String formatted = _safeCell(row, idxFormatted);
      String latStr = _safeCell(row, idxLat);
      String lngStr = _safeCell(row, idxLng);
      String explicitName = idxName >= 0 ? _safeCell(row, idxName) : '';

      final double? lat = double.tryParse(latStr);
      final double? lng = double.tryParse(lngStr);

      // 名稱優先級：CSV name 欄位 > 映射檔 > formatted_address > address
      final String normAddr =
          _normalizeAddress(address.isNotEmpty ? address : formatted);
      final String mappedName = addressToName[normAddr] ?? '';
      final String name = (explicitName.isNotEmpty
          ? explicitName
          : (mappedName.isNotEmpty
              ? mappedName
              : (formatted.isNotEmpty
                  ? formatted
                  : (address.isNotEmpty ? address : '未命名諮商所'))));

      centers.add(
        CounselingCenter(
          id: id.isNotEmpty ? id : 'csv-$i',
          name: name,
          address: address.isNotEmpty
              ? address
              : (formatted.isNotEmpty ? formatted : ''),
          phone: 'N/A',
          onlineCounseling: false,
          latitude: lat,
          longitude: lng,
          createdAt: nowIso,
          updatedAt: nowIso,
        ),
      );
    }

    return centers;
  }

  static String _safeCell(List<dynamic> row, int idx) {
    if (idx < 0 || idx >= row.length) return '';
    final val = row[idx];
    return val == null ? '' : val.toString().trim();
  }

  static Future<Map<String, String>> _loadNameMapping() async {
    try {
      final content = await rootBundle.loadString(nameMappingAssetPath);
      final rows =
          const CsvToListConverter(eol: '\n', shouldParseNumbers: false)
              .convert(content);
      if (rows.isEmpty) return {};
      final header = rows.first.map((e) => e.toString().trim()).toList();
      final int idxAddr = header.indexOf('address');
      final int idxName = header.indexOf('name');
      if (idxAddr < 0 || idxName < 0) return {};
      final Map<String, String> map = {};
      for (int i = 1; i < rows.length; i++) {
        final row = rows[i];
        if (row.isEmpty) continue;
        final addr = _safeCell(row, idxAddr);
        final name = _safeCell(row, idxName);
        if (addr.isEmpty || name.isEmpty) continue;
        map[_normalizeAddress(addr)] = name;
      }
      return map;
    } catch (_) {
      return {};
    }
  }

  static String _normalizeAddress(String input) {
    final s = input
        .replaceAll('台灣', '')
        .replaceAll('臺灣', '')
        .replaceAll('台北市', '臺北市')
        .replaceAll(RegExp(r"\s+"), '')
        .replaceAll('之', '之') // 保持但去空白
        .replaceAll('，', ',')
        .trim();
    return s;
  }
}
