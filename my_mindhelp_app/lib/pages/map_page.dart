// lib/pages/maps_page.dart

import 'package:flutter/material.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import '../core/theme.dart';
import '../widgets/custom_app_bar.dart';

class MapsPage extends StatefulWidget {
  @override
  _MapsPageState createState() => _MapsPageState();
}

class _MapsPageState extends State<MapsPage> {
  late GoogleMapController _controller;

  // 初始鏡頭位置：雙北中心座標 (可自行調整)
  static const _initialPosition = CameraPosition(
    target: LatLng(25.0330, 121.5654), // 台北 101 座標
    zoom: 12,
  );

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: const CustomAppBar(
        titleWidget: Text(
          '地圖',
          style: TextStyle(fontSize: 24, color: AppColors.textHigh),
        ),
        showBackButton: true,
      ),
      body: GoogleMap(
        initialCameraPosition: _initialPosition,
        myLocationEnabled: true,
        myLocationButtonEnabled: true,
        onMapCreated: (controller) {
          _controller = controller;
        },
        markers: {
          Marker(
            markerId: MarkerId('center'),
            position: LatLng(25.0330, 121.5654),
            infoWindow: InfoWindow(title: '台北 101'),
          ),
        },
      ),
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: 1,
        selectedItemColor: AppColors.accent,
        unselectedItemColor: AppColors.textBody,
        onTap: (idx) {/* 跳轉邏輯 */},
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
