import 'package:flutter/material.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:geolocator/geolocator.dart';
import '../services/location_service.dart';
import '../core/theme.dart';
import '../widgets/primary_button.dart';
import '../widgets/custom_app_bar.dart';
import '../models/resource.dart';

class MapsPage extends StatefulWidget {
  const MapsPage({super.key});

  @override
  State<MapsPage> createState() => _MapsPageState();
}

class _MapsPageState extends State<MapsPage> {
  late GoogleMapController mapController;
  final LocationService _locationService = LocationService();
  LatLng _currentLocation = const LatLng(25.0487, 121.5175);
  bool _isLoading = true;
  final Set<Marker> _markers = {};
  String _mapStatus = '正在載入地圖...';

  @override
  void initState() {
    super.initState();
    _getLocationAndLoadResources();
  }

  Future<void> _getLocationAndLoadResources() async {
    try {
      setState(() {
        _mapStatus = '正在獲取您的位置...';
      });
      // 新增超時設定
      Position position =
          await _determinePosition().timeout(const Duration(seconds: 10));
      if (mounted) {
        setState(() {
          _currentLocation = LatLng(position.latitude, position.longitude);
        });
      }
      mapController.animateCamera(CameraUpdate.newLatLng(_currentLocation));
      await _loadNearbyClinics();
    } catch (e) {
      if (mounted) {
        _mapStatus = '無法獲取位置，使用預設座標。錯誤：${e.toString()}';
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(_mapStatus)),
        );
        setState(() {}); // 刷新 UI
        await _loadNearbyClinics(); // 即使定位失敗也載入預設位置的資料
      }
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  Future<Position> _determinePosition() async {
    bool serviceEnabled;
    LocationPermission permission;

    serviceEnabled = await Geolocator.isLocationServiceEnabled();
    if (!serviceEnabled) {
      return Future.error('定位服務已停用。');
    }

    permission = await Geolocator.checkPermission();
    if (permission == LocationPermission.denied) {
      permission = await Geolocator.requestPermission();
      if (permission == LocationPermission.denied) {
        return Future.error('定位權限被拒絕。');
      }
    }

    if (permission == LocationPermission.deniedForever) {
      return Future.error('定位權限永久被拒絕，請在設定中啟用。');
    }

    return await Geolocator.getCurrentPosition();
  }

  Future<void> _loadNearbyClinics() async {
    setState(() {
      _mapStatus = '正在搜尋附近診所...';
    });
    try {
      final googleAddresses = await _locationService.getGoogleAddresses(
        lat: _currentLocation.latitude,
        lon: _currentLocation.longitude,
        radius: 5000,
      );
      _markers.clear();
      _markers.add(
        Marker(
          markerId: const MarkerId('current_location'),
          position: _currentLocation,
          infoWindow: const InfoWindow(title: '我的位置'),
          icon: BitmapDescriptor.defaultMarkerWithHue(BitmapDescriptor.hueBlue),
        ),
      );
      for (var address in googleAddresses) {
        _markers.add(
          Marker(
            markerId: MarkerId(address['name']),
            position: LatLng(address['geometry']['location']['lat'],
                address['geometry']['location']['lng']),
            infoWindow: InfoWindow(
              title: address['name'],
              snippet: address['formatted_address'],
            ),
          ),
        );
      }
      _mapStatus = '已找到 ${googleAddresses.length} 間診所';
    } catch (e) {
      _mapStatus = '無法載入診所資訊：${e.toString()}';
    }
  }

  void _onMapCreated(GoogleMapController controller) {
    mapController = controller;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: CustomAppBar(
        titleWidget: const Text(
          '尋找診所',
          style: TextStyle(fontSize: 24, color: AppColors.textHigh),
        ),
        rightIcon: IconButton(
          icon: const Icon(Icons.notifications, color: AppColors.textHigh),
          onPressed: () => Navigator.pushNamed(context, '/notify'),
        ),
      ),
      body: Stack(
        children: [
          GoogleMap(
            onMapCreated: _onMapCreated,
            initialCameraPosition: CameraPosition(
              target: _currentLocation,
              zoom: 14.0,
            ),
            markers: _markers,
            myLocationEnabled: true,
            myLocationButtonEnabled: false,
          ),
        ],
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
