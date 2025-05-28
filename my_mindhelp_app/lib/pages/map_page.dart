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

  @override
  void initState() {
    super.initState();
    _setupMap();
  }

  Future<void> _setupMap() async {
    // 1. 取得使用者位置
    Position pos = await Geolocator.getCurrentPosition(
      desiredAccuracy: LocationAccuracy.high,
    );

    final userLatLng = LatLng(pos.latitude, pos.longitude);

    // 2. 更新鏡頭至使用者位置
    setState(() {
      _initialCamera = CameraPosition(target: userLatLng, zoom: 14);
      // 加上一個標記表示使用者
      _markers.add(Marker(
        markerId: MarkerId('user'),
        position: userLatLng,
        icon: BitmapDescriptor.defaultMarkerWithHue(
            BitmapDescriptor.hueAzure),
        infoWindow: InfoWindow(title: '您的位置'),
      ));
    });

    // 3. 向後端請求資料庫提供的附近醫療資源
    await _fetchNearbyResources(pos.latitude, pos.longitude);
  }

  Future<void> _fetchNearbyResources(double lat, double lng) async {
    // 假設你的 API 長這樣： /api/resources?lat=...&lng=...
    final uri = Uri.parse(
        'https://your-backend.com/api/resources?lat=$lat&lng=$lng');
    final resp = await http.get(uri);

    if (resp.statusCode == 200) {
      final List data = jsonDecode(resp.body);
      // 假設每筆資料長 { "id":1, "name":"XX診所", "latitude":25.04, "longitude":121.56 }
      setState(() {
        for (var item in data) {
          _markers.add(Marker(
            markerId: MarkerId(item['id'].toString()),
            position:
                LatLng(item['latitude'], item['longitude']),
            infoWindow: InfoWindow(title: item['name']),
          ));
        }
      });
    } else {
      // 處理請求錯誤
      debugPrint('Fetch resources failed: ${resp.statusCode}');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: const CustomAppBar(
        showBackButton: true,
        titleWidget: Image(
          image: AssetImage('assets/images/mindhelp.png'),
          width: 200,
          fit: BoxFit.contain,
        ),
      ),

      body: GoogleMap(
        initialCameraPosition: _initialCamera,
        myLocationEnabled: true,
        myLocationButtonEnabled: true,
        onMapCreated: (ctrl) => _mapCtrl = ctrl,
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
          BottomNavigationBarItem(
              icon: Icon(Icons.home), label: 'Home'),
          BottomNavigationBarItem(
              icon: Icon(Icons.location_on), label: 'Maps'),
          BottomNavigationBarItem(
              icon: Icon(Icons.chat_bubble), label: 'Chat'),
          BottomNavigationBarItem(
              icon: Icon(Icons.person), label: 'Profile'),
        ],
      ),
    );
  }
}
