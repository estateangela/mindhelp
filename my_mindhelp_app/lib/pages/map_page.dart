// lib/pages/maps_page.dart

import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart' show rootBundle;
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:geolocator/geolocator.dart';
import 'package:http/http.dart' as http;
import '../core/theme.dart';
import '../widgets/custom_app_bar.dart';

class MapsPage extends StatefulWidget {
  @override
  _MapsPageState createState() => _MapsPageState();
}

class _MapsPageState extends State<MapsPage> {
  late GoogleMapController _mapCtrl;
  CameraPosition _initialCamera =
      const CameraPosition(target: LatLng(25.0330, 121.5654), zoom: 12);
  final Set<Marker> _markers = {};
  bool _isLoading = true;
  String? _errorMessage;

  @override
  void initState() {
    super.initState();
    _setupMap();
  }

  Future<void> _setupMap() async {
    double lat = 25.0330, lng = 121.5654;

    try {
      // 请求权限并取得使用者位置
      var perm = await Geolocator.checkPermission();
      if (perm == LocationPermission.denied) {
        perm = await Geolocator.requestPermission();
        if (perm == LocationPermission.denied) {
          throw '用户拒绝位置权限';
        }
      }
      final pos = await Geolocator.getCurrentPosition();
      lat = pos.latitude;
      lng = pos.longitude;

      setState(() {
        _initialCamera = CameraPosition(target: LatLng(lat, lng), zoom: 14);
        // 先把用户标记加入
        _markers.add(Marker(
          markerId: const MarkerId('user'),
          position: LatLng(lat, lng),
          icon:
              BitmapDescriptor.defaultMarkerWithHue(BitmapDescriptor.hueAzure),
          infoWindow: const InfoWindow(title: '您的位置'),
        ));
      });
    } catch (e) {
      // 定位失败，记录错误但继续拉资源
      debugPrint('📍 定位失败: $e');
      _errorMessage = e.toString();
    }

    // 使用定位结果（或默认中心）去拉附近资源
    await _fetchNearbyResources(lat, lng);
  }

  Future<void> _fetchNearbyResources(double lat, double lng) async {
    List data;
    final uri =
        Uri.parse('https://your-backend.com/api/resources?lat=$lat&lng=$lng');

    try {
      final resp = await http.get(uri);
      debugPrint('🔗 请求 URL: $uri  状态: ${resp.statusCode}');
      if (resp.statusCode == 200) {
        data = jsonDecode(resp.body) as List;
      } else {
        throw 'HTTP 错误 ${resp.statusCode}';
      }
    } catch (e) {
      debugPrint('❗️ 拉取线上资源失败: $e ，改用 Mock');
      // 从 assets/mock_resources.json 载入 mock 数据
      final jsonStr = await rootBundle.loadString('assets/mock_resources.json');
      data = jsonDecode(jsonStr) as List;
    }

    setState(() {
      // 只 clear 资源标记，保留用户标记
      _markers.removeWhere((m) => m.markerId.value != 'user');

      for (var item in data) {
        _markers.add(Marker(
          markerId: MarkerId(item['id'].toString()),
          position:
              LatLng(item['latitude'] as double, item['longitude'] as double),
          infoWindow: InfoWindow(
            title: item['name'] as String,
            snippet: item['address'] as String? ?? '',
          ),
        ));
      }

      _isLoading = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,

      // 顶部统一 Logo + 返回 + 通知
      appBar: const CustomAppBar(
        showBackButton: true,
        titleWidget: Image(
          image: AssetImage('assets/images/mindhelp.png'),
          width: 200,
          fit: BoxFit.contain,
        ),
      ),

      body: Stack(
        children: [
          // 加载中或错误提示
          if (_isLoading)
            const Center(child: CircularProgressIndicator())
          else if (_errorMessage != null)
            Center(child: Text(_errorMessage!, style: TextStyle(color: Colors.red))),
          // 地图
          if (!_isLoading && _errorMessage == null)
            GoogleMap(
              initialCameraPosition: _initialCamera,
              myLocationEnabled: true,
              onMapCreated: (c) => _mapCtrl = c,
              markers: _markers,
            ),
        ],
      ),

      bottomNavigationBar: BottomNavigationBar(
        currentIndex: 1, // Maps 在索引 1
        selectedItemColor: AppColors.accent,
        unselectedItemColor: AppColors.textBody,
        onTap: (i) {
          switch (i) {
            case 0:
              Navigator.pushReplacementNamed(context, '/home');
              break;
            case 1:
              // already here
              break;
            case 2:
              Navigator.pushReplacementNamed(context, '/chat');
              break;
            case 3:
              Navigator.pushReplacementNamed(context, '/profile');
              break;
          }
        },
        items: const [
          BottomNavigationBarItem(icon: Icon(Icons.home), label: 'Home'),
          BottomNavigationBarItem(
              icon: Icon(Icons.location_on), label: 'Maps'),
          BottomNavigationBarItem(
              icon: Icon(Icons.chat_bubble), label: 'Chat'),
          BottomNavigationBarItem(icon: Icon(Icons.person), label: 'Profile'),
        ],
      ),
    );
  }
}
