import 'package:http/http.dart' as http;
import 'dart:convert';
import '../models/map_item.dart';

class LocationService {
  // 替換為你的 Google Sheet ID
  static const String _sheetUrl =
      'https://docs.google.com/spreadsheets/d/1OjZT5iVkj09gOoY_uJDMVXl-xTbcbF7-IUb1gBArkJc/gviz/tq?tqx=out:json';

  Future<List<MapItem>> getLocations() async {
    try {
      final response = await http.get(Uri.parse(_sheetUrl));
      if (response.statusCode == 200) {
        final String jsonString =
            response.body.substring(47, response.body.length - 2);
        final data = json.decode(jsonString);

        List<MapItem> locations = [];
        for (var row in data['table']['rows']) {
          final cells = row['c'];
          if (cells.length >= 3 &&
              cells[0] != null &&
              cells[1] != null &&
              cells[2] != null) {
            final name = cells[0]['v'];
            final address = cells[1]['v'];
            final description = cells[2]['v'];
            locations.add(MapItem(
                name: name, address: address, description: description));
          }
        }

        // 新增的除錯訊息
        print('成功讀取 ${locations.length} 筆資料');

        return locations;
      } else {
        print('從 Google Sheet 載入資料失敗: ${response.statusCode}');
        throw Exception('Failed to load data from Google Sheet');
      }
    } catch (e) {
      print('錯誤: $e');
      throw Exception('Error: $e');
    }
  }
}
