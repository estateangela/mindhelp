import UIKit
import Flutter
import GoogleMaps

@main
@objc class AppDelegate: FlutterAppDelegate {
  override func application(
      _ application: UIApplication,
      didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?
  ) -> Bool {
    // 這裡從 Info.plist 讀取 GOOGLE_MAPS_API_KEY
    if let mapsApiKey = Bundle.main.infoDictionary?["GOOGLE_MAPS_API_KEY"] as? String {
        GMSServices.provideAPIKey(mapsApiKey)
    } else {
        // 如果金鑰未設定，拋出致命錯誤，提醒開發者檢查配置
        fatalError("GOOGLE_MAPS_API_KEY is not set in Info.plist. Please check your .xcconfig file.")
    }

    GeneratedPluginRegistrant.register(with: self)
    return super.application(application, didFinishLaunchingWithOptions: launchOptions)
  }
}
