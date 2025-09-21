import 'package:flutter/material.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import '../core/theme.dart';
import '../services/location_service.dart';
import '../models/resource.dart'; // 導入新的 Resource 模型

class MapsPage extends StatefulWidget {
  const MapsPage({super.key});

  @override
  State<MapsPage> createState() => _MapsPageState();
}

class _MapsPageState extends State<MapsPage> {
  final LocationService _locationService = LocationService();
  final Set<Marker> _markers = {};
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadLocations();
  }

  Future<void> _loadLocations() async {
    try {
      // 模擬取得使用者的當前位置
      const LatLng userLocation = LatLng(25.033, 121.565);

      final resources = await _locationService.searchResources(
        lat: userLocation.latitude,
        lon: userLocation.longitude,
      );

      _markers.clear(); // 清除舊的標記點
      for (var resource in resources) {
        // 從後端回傳的資料中直接取得經緯度
        final LatLng position = LatLng(resource.latitude, resource.longitude);
        _markers.add(
          Marker(
            markerId: MarkerId(resource.id), // 使用 ID 作為唯一標識
            position: position,
            infoWindow: InfoWindow(
              title: resource.name,
              snippet: resource.address,
            ),
          ),
        );
      }
    } catch (e) {
      debugPrint('Error loading locations: $e');
      // 顯示錯誤訊息給使用者
    } finally {
      setState(() {
        _isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('地圖', style: TextStyle(color: Colors.white)),
        backgroundColor: AppColors.accent,
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : GoogleMap(
              onMapCreated: (controller) {},
              initialCameraPosition: const CameraPosition(
                target: LatLng(25.033, 121.565), // 預設位置
                zoom: 12,
              ),
              markers: _markers,
            ),
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: 1, // 地圖頁的索引
        selectedItemColor: AppColors.accent,
        unselectedItemColor: AppColors.textBody,
        onTap: (idx) {
          switch (idx) {
            case 0:
              Navigator.pushReplacementNamed(context, '/home');
              break;
            case 1:
              // 已在地圖頁面，不做動作
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
