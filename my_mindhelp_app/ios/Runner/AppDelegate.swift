import UIKit
import Flutter
import GoogleMaps   // ← 引入

@UIApplicationMain
@objc class AppDelegate: FlutterAppDelegate {
  override func application(
      _ application: UIApplication,
      didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?
  ) -> Bool {
    GMSServices.provideAPIKey("AIzaSyBXHc-rjUBntG7oppaYSpY322EE3q-T4ng")  // ← 這裡設你的 Key
    GeneratedPluginRegistrant.register(with: self)
    return super.application(application, didFinishLaunchingWithOptions: launchOptions)
  }
}
