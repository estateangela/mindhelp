import 'package:http/http.dart' as http;
import 'dart:convert';
import '../models/resource.dart';

class LocationService {
  // TODO: 將 'YOUR_SHEET_ID' 替換為你的 Google Sheet ID
  // TODO: 將 'YOUR_SHEET_NAME' 替換為你的工作表名稱
  static const String _sheetUrl =
      'https://docs.google.com/spreadsheets/d/1OjZT5iVkj09gOoY_uJDMVXl-xTbcbF7-IUb1gBArkJc/gviz/tq?tqx=out:json';

  Future<List<Resource>> getResources() async {
    try {
      final response = await http.get(Uri.parse(_sheetUrl));
      if (response.statusCode == 200) {
        final String jsonString =
            response.body.substring(47, response.body.length - 2);
        final data = json.decode(jsonString);

        List<Resource> resources = [];
        for (var row in data['table']['rows']) {
          final cells = row['c'];
          // 確保欄位不是空的
          if (cells.length >= 7 &&
              cells[0] != null &&
              cells[1] != null &&
              cells[2] != null) {
            resources.add(Resource(
              id: cells[0]['v'].toString(),
              name: cells[1]['v'],
              type: cells[2]['v'],
              address: cells[3]['v'],
              phone: cells[4]['v'],
              website: cells[5]?['v'] ?? '',
              description: cells[6]['v'],
            ));
          }
        }
        return resources;
      } else {
        throw Exception('Failed to load data from Google Sheet');
      }
    } catch (e) {
      throw Exception('Error: $e');
    }
  }
}
