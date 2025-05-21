import 'package:flutter/material.dart';

void main() {
  runApp(FlutterApp());
}

class FlutterApp extends StatelessWidget {
  final ValueNotifier<bool> _dark = ValueNotifier<bool>(true);
  final ValueNotifier<double> _widthFactor = ValueNotifier<double>(1.0);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
        home: ValueListenableBuilder<bool>(
            valueListenable: _dark,
            builder: (context, color, child) {
              return ValueListenableBuilder<double>(
                valueListenable: _widthFactor,
                builder: (context, factor, child) {
                  return Scaffold(
                      backgroundColor:
                          _dark.value ? Colors.black : Colors.white,
                      appBar: AppBar(
                        actions: [
                          Switch(
                            value: _dark.value,
                            onChanged: (value) {
                              _dark.value = value;
                            },
                          ),
                          DropdownButton<double>(
                            value: _widthFactor.value,
                            onChanged: (value) {
                              _widthFactor.value = value!;
                            },
                            items: [
                              DropdownMenuItem<double>(
                                value: 0.5,
                                child: Text('Size: 50%'),
                              ),
                              DropdownMenuItem<double>(
                                value: 0.75,
                                child: Text('Size: 75%'),
                              ),
                              DropdownMenuItem<double>(
                                value: 1.0,
                                child: Text('Size: 100%'),
                              ),
                            ],
                          ),
                        ],
                      ),
                      body: Center(
                          child: Container(
                        width: MediaQuery.of(context).size.width *
                            _widthFactor.value,
                        child: Column(
                          mainAxisSize: MainAxisSize.max,
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Group3(),
                          ],
                        ),
                      )));
                },
              );
            }));
  }
}

class Group3 extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Container(
          width: 106,
          height: 57,
          child: Stack(
            children: [
              Positioned(
                left: 0,
                top: 0,
                child: Container(
                  width: 106,
                  height: 57,
                  decoration: ShapeDecoration(
                    color: Color(0xFFFFD590),
                    shape: RoundedRectangleBorder(
                      side: BorderSide(width: 1, color: Color(0xFFE2AD56)),
                      borderRadius: BorderRadius.circular(4),
                    ),
                    shadows: [
                      BoxShadow(
                        color: Color(0x19000000),
                        blurRadius: 4,
                        offset: Offset(2, 4),
                        spreadRadius: 0,
                      )
                    ],
                  ),
                ),
              ),
              Positioned(
                left: 26,
                top: 17,
                child: Text(
                  '登入',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 24,
                    fontFamily: 'jf-openhuninn-2.0',
                    height: 0,
                    letterSpacing: 4.80,
                  ),
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }
}
