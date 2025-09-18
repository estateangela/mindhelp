import 'package:flutter/material.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import '../core/theme.dart';
import '../services/location_service.dart';
import '../models/map_item.dart';
import 'package:geocoding/geocoding.dart';

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
      final locations = await _locationService.getLocations();
      for (var item in locations) {
        try {
          // TODO: 考慮將 Geocoding 移至後端，避免前端 API 限制
          final geocodingInstance = GeocodingPlatform.instance;
          if (geocodingInstance != null) {
            final location = await geocodingInstance.locationFromAddress(item.address);
            if (location.isNotEmpty) {
              final LatLng position = LatLng(location.first.latitude, location.first.longitude);
              _markers.add(
                Marker(
                  markerId: MarkerId(item.name),
                  position: position,
                  infoWindow: InfoWindow(
                    title: item.name,
                    snippet: item.address,
                  ),
                ),
              );
            }
          }
        } catch (e) {
          // 替換 print 為一個更適合生產環境的日誌框架
          debugPrint('Error geocoding address for ${item.name}: $e');
        }
      }
    } catch (e) {
      // 替換 print 為一個更適合生產環境的日誌框架
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
