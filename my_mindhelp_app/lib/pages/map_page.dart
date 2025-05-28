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
      CameraPosition(target: LatLng(25.0330, 121.5654), zoom: 12);
  final Set<Marker> _markers = {};
  bool _isLoading = true;
  String? _errorMessage;

  @override
  void initState() {
    super.initState();
    _setupMap();
  }

  Future<void> _setupMap() async {
    try {
      // 请求权限并定位
      var perm = await Geolocator.checkPermission();
      if (perm == LocationPermission.denied) {
        perm = await Geolocator.requestPermission();
        if (perm == LocationPermission.denied) {
          throw '需要位置权限';
        }
      }
      final pos = await Geolocator.getCurrentPosition();
      final userLatLng = LatLng(pos.latitude, pos.longitude);

      setState(() {
        _initialCamera = CameraPosition(target: userLatLng, zoom: 14);
        _markers.add(Marker(
          markerId: MarkerId('user'),
          position: userLatLng,
          icon:
              BitmapDescriptor.defaultMarkerWithHue(BitmapDescriptor.hueAzure),
          infoWindow: InfoWindow(title: '您的位置'),
        ));
      });
    } catch (e) {
      // 定位失败，但我们依然加载 mock
      debugPrint('📍 定位失败: $e');
    }

    // 不管定位成功与否，都去加载资源
    await _fetchNearbyResources(0, 0);
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
        throw 'HTTP ${resp.statusCode}';
      }
    } catch (e) {
      debugPrint('❗️ 拉取线上资源失败，使用 Mock: $e');
      final jsonStr = await rootBundle.loadString('assets/mock_resources.json');
      data = jsonDecode(jsonStr) as List;
    }

    setState(() {
      _markers
        ..clear()
        ..addAll(data.map((item) => Marker(
              markerId: MarkerId(item['id'].toString()),
              position: LatLng(item['latitude'], item['longitude']),
              infoWindow: InfoWindow(
                title: item['name'],
                snippet: item['address'] ?? '',
              ),
            )));
      _isLoading = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: const CustomAppBar(
        showBackButton: true,
        titleWidget: Text(
          '地圖',
          style: TextStyle(fontSize: 24, color: AppColors.textHigh),
        ),
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : GoogleMap(
              initialCameraPosition: _initialCamera,
              myLocationEnabled: true,
              onMapCreated: (c) => _mapCtrl = c,
              markers: _markers,
            ),
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: 1,
        selectedItemColor: AppColors.accent,
        unselectedItemColor: AppColors.textBody,
        onTap: (i) {
          switch (i) {
            case 0:
              Navigator.pushReplacementNamed(context, '/home');
              break;
            case 1:
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
          BottomNavigationBarItem(icon: Icon(Icons.location_on), label: 'Maps'),
          BottomNavigationBarItem(icon: Icon(Icons.chat_bubble), label: 'Chat'),
          BottomNavigationBarItem(icon: Icon(Icons.person), label: 'Profile'),
        ],
      ),
    );
  }
}
