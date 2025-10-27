import 'package:flutter/material.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:geolocator/geolocator.dart';
import '../services/location_service.dart';
import '../core/theme.dart';
import '../widgets/custom_app_bar.dart';
import '../models/counseling_center.dart';

class MapsPage extends StatefulWidget {
  const MapsPage({super.key});

  @override
  State<MapsPage> createState() => _MapsPageState();
}

class _MapsPageState extends State<MapsPage> {
  late GoogleMapController mapController;
  final LocationService _locationService = LocationService();
  LatLng _currentLocation = const LatLng(25.0487, 121.5175); // 台北商業大學
  bool _isLoading = true;
  final Set<Marker> _markers = {};
  String _mapStatus = '正在載入地圖...';
  BitmapDescriptor? _userLocationIcon;

  @override
  void initState() {
    super.initState();
    // 添加一些測試標記來驗證地圖功能
    _addTestMarkers();
    // 創建自定義標記圖標
    _createCustomMarkerIcons();
  }

  void _addTestMarkers() {
    // 添加一些固定的測試標記
    _markers.addAll([
      Marker(
        markerId: const MarkerId('test1'),
        position: const LatLng(25.0330, 121.5654), // 台北101
        infoWindow: const InfoWindow(title: '測試標記 1', snippet: '台北101'),
        icon: BitmapDescriptor.defaultMarkerWithHue(BitmapDescriptor.hueRed),
      ),
      Marker(
        markerId: const MarkerId('test2'),
        position: const LatLng(25.0370, 121.5700), // 台北車站附近
        infoWindow: const InfoWindow(title: '測試標記 2', snippet: '台北車站'),
        icon: BitmapDescriptor.defaultMarkerWithHue(BitmapDescriptor.hueGreen),
      ),
      Marker(
        markerId: const MarkerId('test3'),
        position: const LatLng(25.0300, 121.5600), // 中正紀念堂附近
        infoWindow: const InfoWindow(title: '測試標記 3', snippet: '中正紀念堂'),
        icon: BitmapDescriptor.defaultMarkerWithHue(BitmapDescriptor.hueBlue),
      ),
      Marker(
        markerId: const MarkerId('test4'),
        position: const LatLng(25.0480, 121.5170), // 北商附近
        infoWindow: const InfoWindow(title: '測試標記 4', snippet: '台北商業大學'),
        icon: BitmapDescriptor.defaultMarkerWithHue(BitmapDescriptor.hueAzure),
      ),
      Marker(
        markerId: const MarkerId('test5'),
        position: const LatLng(25.0400, 121.5400), // 仁愛院區附近
        infoWindow: const InfoWindow(title: '測試標記 5', snippet: '仁愛院區'),
        icon: BitmapDescriptor.defaultMarkerWithHue(BitmapDescriptor.hueViolet),
      ),
      Marker(
        markerId: const MarkerId('test6'),
        position: const LatLng(25.0600, 121.5100), // 中興院區附近
        infoWindow: const InfoWindow(title: '測試標記 6', snippet: '中興院區'),
        icon: BitmapDescriptor.defaultMarkerWithHue(BitmapDescriptor.hueYellow),
      ),
    ]);
    print('已添加 ${_markers.length} 個測試標記');
  }

  // 創建自定義標記圖標
  Future<void> _createCustomMarkerIcons() async {
    try {
      // 為使用者位置創建特殊的標記圖標 - 使用橙色表示使用者位置
      _userLocationIcon =
          BitmapDescriptor.defaultMarkerWithHue(BitmapDescriptor.hueOrange);
      print('自定義標記圖標創建完成 - 使用者位置使用橙色標記');
    } catch (e) {
      print('創建自定義標記圖標失敗: $e');
    }
  }

  // 獲取備用座標（對應真實醫療機構位置）
  Map<String, double> _getFallbackCoordinates(String centerId) {
    // 台北商業大學方圓五公里內的真實醫療機構座標
    final coordinates = {
      'fallback-1': {'lat': 25.0487, 'lng': 121.5175}, // 台北商業大學
      'fallback-2': {'lat': 25.0370, 'lng': 121.5200}, // 台大醫院
      'fallback-3': {'lat': 25.0600, 'lng': 121.5100}, // 中興院區
      'fallback-4': {'lat': 25.0400, 'lng': 121.5400}, // 仁愛院區
      'fallback-5': {'lat': 25.0300, 'lng': 121.5000}, // 和平院區
      'fallback-6': {'lat': 25.0350, 'lng': 121.5500}, // 忠孝院區（調整位置）
      'fallback-7': {'lat': 25.0800, 'lng': 121.5200}, // 陽明院區
    };

    return coordinates[centerId] ?? {'lat': 25.0487, 'lng': 121.5175};
  }

  void _addFallbackMarker(CounselingCenter center) {
    try {
      // 根據諮商所ID使用不同的備用座標（台北商業大學附近，方圓5公里內）
      final fallbackCoords = _getFallbackCoordinates(center.id);
      final LatLng fallbackPosition =
          LatLng(fallbackCoords['lat']!, fallbackCoords['lng']!);

      // 計算距離
      final double distanceInMeters = Geolocator.distanceBetween(
        _currentLocation.latitude,
        _currentLocation.longitude,
        fallbackPosition.latitude,
        fallbackPosition.longitude,
      );
      final double distanceInKm = distanceInMeters / 1000;

      _markers.add(
        Marker(
          markerId: MarkerId('${center.id}_fallback'),
          position: fallbackPosition,
          infoWindow: InfoWindow(
            title: center.name,
            snippet:
                '${center.phone}\n地址：${center.address}\n距離：${distanceInKm.toStringAsFixed(2)} 公里（備用座標）',
          ),
          icon: BitmapDescriptor.defaultMarkerWithHue(
            center.onlineCounseling
                ? BitmapDescriptor.hueGreen
                : BitmapDescriptor.hueRed,
          ),
        ),
      );
      print('已添加備用標記：${center.name} (${distanceInKm.toStringAsFixed(2)} 公里)');
    } catch (e) {
      print('添加備用標記失敗：$e');
    }
  }

  // 將多個諮商所轉為地圖標記（優先使用座標，缺失則回退備用座標）
  Future<int> _addCenterMarkers(List<CounselingCenter> centers) async {
    int successCount = 0;
    for (var center in centers) {
      try {
        print('正在處理諮商所：${center.name} - ${center.address}');

        LatLng position;
        bool usedFallback = false;
        if (center.latitude != null && center.longitude != null) {
          position = LatLng(center.latitude!, center.longitude!);
          print('使用座標：${position.latitude}, ${position.longitude}');
        } else {
          final fallback = _getFallbackCoordinates(center.id);
          position = LatLng(fallback['lat']!, fallback['lng']!);
          usedFallback = true;
          print('缺少座標，使用備用座標 ${position.latitude}, ${position.longitude}');
        }

        final double distanceInMeters = Geolocator.distanceBetween(
          _currentLocation.latitude,
          _currentLocation.longitude,
          position.latitude,
          position.longitude,
        );
        final double distanceInKm = distanceInMeters / 1000;

        _markers.add(
          Marker(
            markerId: MarkerId(center.id),
            position: position,
            infoWindow: InfoWindow(
              title: center.name,
              snippet: usedFallback
                  ? '${center.phone}\n地址：${center.address}\n距離：${distanceInKm.toStringAsFixed(2)} 公里（備用座標）'
                  : '${center.phone}\n地址：${center.address}\n距離：${distanceInKm.toStringAsFixed(2)} 公里',
            ),
            icon: BitmapDescriptor.defaultMarkerWithHue(
              center.onlineCounseling
                  ? BitmapDescriptor.hueGreen
                  : BitmapDescriptor.hueRed,
            ),
          ),
        );
        successCount++;
      } catch (e) {
        print('添加標記時發生錯誤：${e.toString()}');
        _addFallbackMarker(center);
      }
    }
    return successCount;
  }

  Future<void> _loadAllData() async {
    try {
      setState(() {
        _mapStatus = '正在獲取您的位置...';
        _isLoading = true;
      });

      Position position =
          await _determinePosition().timeout(const Duration(seconds: 10));

      if (mounted) {
        setState(() {
          _currentLocation = LatLng(position.latitude, position.longitude);
        });
        mapController.animateCamera(CameraUpdate.newLatLng(_currentLocation));
      }
    } catch (e) {
      if (mounted) {
        _mapStatus = '無法獲取位置，使用預設座標。錯誤：${e.toString()}';
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(_mapStatus)),
        );
      }
    } finally {
      await _loadNearbyClinics();
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
      _isLoading = true;
    });

    try {
      // 不清除現有的測試標記，只添加當前位置標記（如果還沒有的話）
      bool hasCurrentLocation =
          _markers.any((marker) => marker.markerId.value == 'current_location');
      if (!hasCurrentLocation) {
        _markers.add(
          Marker(
            markerId: const MarkerId('current_location'),
            position: _currentLocation,
            infoWindow: const InfoWindow(
              title: '我的位置',
              snippet: '台北商業大學',
            ),
            icon: _userLocationIcon ??
                BitmapDescriptor.defaultMarkerWithHue(
                    BitmapDescriptor.hueOrange), // 使用橙色標記
          ),
        );
      }

      print('開始調用 API...');
      final List<CounselingCenter> counselingCenters =
          await _locationService.getCounselingCenters(
        userLatitude: _currentLocation.latitude,
        userLongitude: _currentLocation.longitude,
        radiusKm: 5.0, // 方圓五公里
      );

      print('API 調用完成，獲取到 ${counselingCenters.length} 個諮商所');

      if (counselingCenters.isEmpty) {
        print('沒有獲取到諮商所數據，改用備用清單');
        final fallbackCenters = _locationService.getFallbackCenters();
        final count = await _addCenterMarkers(fallbackCenters);
        if (mounted) {
          setState(() {
            _mapStatus = '顯示 $count 間備用諮商所標記（API 無資料）';
          });
        }
        return;
      }

      final successCount = await _addCenterMarkers(counselingCenters);

      if (mounted) {
        setState(() {
          _mapStatus = '已找到 $successCount 間諮商所（方圓 5 公里內）';
        });
        print('成功處理 $successCount 個諮商所');
        print('最終標記數量：${_markers.length}');
        print('標記列表：${_markers.map((m) => m.markerId.value).toList()}');
      }
    } catch (e) {
      print('載入診所資訊時發生錯誤：${e.toString()}');
      _mapStatus = '無法載入診所資訊：${e.toString()}';
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  void _onMapCreated(GoogleMapController controller) {
    mapController = controller;
    _loadAllData(); // 在這裡呼叫 _loadAllData()，確保 mapController 已初始化
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: CustomAppBar(
        showBackButton: false,
        titleWidget: const Image(
          image: AssetImage('assets/images/mindhelp.png'),
          width: 200,
          fit: BoxFit.contain,
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
          if (_isLoading)
            Center(
              child: CircularProgressIndicator(
                valueColor: AlwaysStoppedAnimation(AppColors.accent),
              ),
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
          }
        },
        items: const [
          BottomNavigationBarItem(icon: Icon(Icons.home), label: 'Home'),
          BottomNavigationBarItem(icon: Icon(Icons.location_on), label: 'Maps'),
          BottomNavigationBarItem(icon: Icon(Icons.chat_bubble), label: 'Chat'),
        ],
      ),
    );
  }
}
