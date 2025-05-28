// lib/pages/maps_page.dart

import 'dart:convert';
import 'package:flutter/material.dart';
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
      // 1. 检查并请求定位权限
      var permission = await Geolocator.checkPermission();
      if (permission == LocationPermission.denied) {
        permission = await Geolocator.requestPermission();
        if (permission == LocationPermission.denied) {
          setState(() {
            _errorMessage = '需要位置权限才能使用地图功能';
            _isLoading = false;
          });
          return;
        }
      }

      // 2. 获取用户当前位置
      final pos = await Geolocator.getCurrentPosition(
        desiredAccuracy: LocationAccuracy.high,
      );
      final userLatLng = LatLng(pos.latitude, pos.longitude);

      // 3. 更新摄像头和用户标记
      setState(() {
        _initialCamera = CameraPosition(target: userLatLng, zoom: 14);
        _markers.add(Marker(
          markerId: MarkerId('user'),
          position: userLatLng,
          icon: BitmapDescriptor.defaultMarkerWithHue(
            BitmapDescriptor.hueAzure,
          ),
          infoWindow: InfoWindow(title: '您的位置'),
        ));
      });

      // 4. 拉取后端资源并加标记
      await _fetchNearbyResources(pos.latitude, pos.longitude);
    } catch (e) {
      setState(() {
        _errorMessage = '获取位置信息失败: $e';
        _isLoading = false;
      });
    }
  }

  Future<void> _fetchNearbyResources(double lat, double lng) async {
    try {
      final uri =
          Uri.parse('https://your-backend.com/api/resources?lat=$lat&lng=$lng');
      final resp = await http.get(uri);

      if (resp.statusCode == 200) {
        final List data = jsonDecode(resp.body);
        setState(() {
          for (var item in data) {
            _markers.add(Marker(
              markerId: MarkerId(item['id'].toString()),
              position: LatLng(item['latitude'], item['longitude']),
              infoWindow: InfoWindow(
                title: item['name'],
                snippet: item['address'] ?? '',
              ),
            ));
          }
          _isLoading = false;
        });
      } else {
        setState(() {
          _errorMessage = '获取资源信息失败: ${resp.statusCode}';
          _isLoading = false;
        });
      }
    } catch (e) {
      setState(() {
        _errorMessage = '获取资源信息失败: $e';
        _isLoading = false;
      });
    }
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
          : _errorMessage != null
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        _errorMessage!,
                        style: const TextStyle(color: Colors.red),
                      ),
                      const SizedBox(height: 16),
                      ElevatedButton(
                        onPressed: _setupMap,
                        child: const Text('重試'),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: AppColors.accent,
                        ),
                      ),
                    ],
                  ),
                )
              : Column(
                  children: [
                    // 地图区域
                    Expanded(
                      child: GoogleMap(
                        initialCameraPosition: _initialCamera,
                        myLocationEnabled: true,
                        myLocationButtonEnabled: true,
                        onMapCreated: (ctrl) => _mapCtrl = ctrl,
                        markers: _markers,
                        zoomControlsEnabled: true,
                        mapToolbarEnabled: true,
                      ),
                    ),
                  ],
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
              break; // 当前页
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
